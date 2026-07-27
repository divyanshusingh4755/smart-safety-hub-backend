package app

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/smart-safety-hub/backend/internal/modules/aws"
	"github.com/smart-safety-hub/backend/internal/modules/brand"
	"github.com/smart-safety-hub/backend/internal/modules/categories"
	"github.com/smart-safety-hub/backend/internal/modules/contacts"
	"github.com/smart-safety-hub/backend/internal/modules/products"
	"github.com/smart-safety-hub/backend/internal/modules/social"
	"github.com/smart-safety-hub/backend/internal/modules/user"
	"github.com/smart-safety-hub/backend/shared"
	"github.com/smart-safety-hub/backend/shared/cache"
	"github.com/smart-safety-hub/backend/shared/mail"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	GrpcAddr   string
	HTTPAddr   string
	DBURL      string
	RedisURL   string
	PrivateKey string
	PublicKey  string

	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string

	ContactAdminEmail string
}

type Container struct {
	HTTPRouter *chi.Mux
	GRPCServer *grpc.Server
	DB         *sqlx.DB
	Logger     *zap.Logger
}

func Bootstrap(cfg Config) (*Container, func()) {
	l := shared.NewLogger()
	sqlxDB := shared.Connect(cfg.DBURL, l)
	redisClient, err := cache.NewRedisClient(cfg.RedisURL, l)
	if err != nil {
		log.Fatalf("Falied to init redis: %v", err)
	}
	s3Client, err := shared.NewS3Client()
	if err != nil {
		log.Fatalf("Falied to init s3: %v", err)
	}

	user.InitSessionStore()
	user.InitOAuthProviders()
	cacheStore := cache.New(redisClient)

	ctx := context.Background()
	if err := cacheStore.InitProductListVersion(ctx); err != nil {
		log.Printf("failed to init product list version: %v", err)
	}

	mailer := mail.NewMailer(
		mail.Config{
			Host:      cfg.SMTPHost,
			Port:      cfg.SMTPPort,
			Username:  cfg.SMTPUsername,
			Password:  cfg.SMTPPassword,
			FromEmail: cfg.SMTPFromEmail,
			FromName:  cfg.SMTPFromName,
		},
	)
	// Create a shared JWT Manager
	jwtManager, _ := shared.NewJWTManager(cfg.PrivateKey, cfg.PublicKey, l)
	jwtMiddleware := shared.JWTMiddleware(jwtManager)

	// Initialize validator
	v := validator.New(validator.WithRequiredStructEnabled())

	// Create Modules
	// User
	userRepo := user.NewUserRepo(sqlxDB)
	userService := user.NewUserService(l, userRepo, jwtManager)
	userRestHandler := user.NewRestHandler(userService, v)

	// upload
	uploadService := aws.NewUploadService(s3Client)
	uploadHandler := aws.NewUploadHandler(uploadService, v)

	// brand
	brandRepo := brand.NewBrandRepo(sqlxDB)
	brandService := brand.NewBrandService(l, brandRepo, cacheStore)
	brandRestHandler := brand.NewRestHandler(brandService, v)

	// Category
	categoryRepo := categories.NewCategoryRepo(sqlxDB)
	categoryService := categories.NewCategoryService(l, categoryRepo, cacheStore)
	categoryRestHandler := categories.NewRestHandler(categoryService, v)

	// Product
	productRepo := products.NewProductRepo(sqlxDB)
	productService := products.NewProductService(l, productRepo, cacheStore)
	productRestHandler := products.NewRestHandler(productService, v)

	// Social Accounts
	socialRepo := social.NewRepository(sqlxDB)
	socialService := social.NewSocialService(l, socialRepo)
	socialRestHandler := social.NewRestHandler(socialService, v)

	// Contact
	contactRepo := contacts.NewContactRepo(sqlxDB)
	contactService := contacts.NewContactService(l, contactRepo, mailer, cfg.ContactAdminEmail)
	contactRestHandler := contacts.NewRestHandler(contactService, v)

	// GRPC
	grpcSrv := grpc.NewServer()

	// Http
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// router.Use(httprate.Limit(
	// 	100,
	// 	1*time.Minute,
	// 	httprate.WithKeyFuncs(httprate.KeyByIP),
	// ))
	// Cors
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"http://localhost:3000", "https://smartsafetyhub.com"},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		OptionsPassthrough: false,
		Debug:              true,
	})
	router.Use(c.Handler)
	router.Route("/v1", func(v1 chi.Router) {
		// Public Routes
		// Auth
		v1.Post("/auth/register", userRestHandler.Register)
		v1.Post("/auth/login", userRestHandler.Login)
		v1.Post("/auth/forgot-password", userRestHandler.ForgotPassword)
		v1.Post("/auth/reset-password", userRestHandler.ResetPassword)
		v1.Post("/auth/logout", userRestHandler.Logout)
		v1.Post("/auth/refresh", userRestHandler.RefreshToken)

		// Brand
		v1.Get("/get-brand/id/{id}", brandRestHandler.GetBrandByID)
		v1.Get("/get-brand/slug/{slug}", brandRestHandler.GetBrandBySlug)
		v1.Get("/get-all-brands", brandRestHandler.GetAllBrand)

		// Category
		v1.Get("/get-category/id/{id}", categoryRestHandler.GetCategoryByID)
		v1.Get("/get-category/slug/{slug}", categoryRestHandler.GetCategoryBySlug)
		v1.Get("/get-all-category", categoryRestHandler.GetAllCategory)

		// Product
		v1.Get("/get-product/id/{id}", productRestHandler.GetProductByID)
		v1.Get("/get-product/slug/{slug}", productRestHandler.GetProductBySlug)
		v1.Get("/get-all-products", productRestHandler.GetAllProducts)
		v1.Get("/products/category/{slug}", productRestHandler.GetProductsByCategory)
		v1.Get("/products/brand/{slug}", productRestHandler.GetProductsByBrand)

		// Product Attribute
		v1.Get("/get-product-attribute/{id}", productRestHandler.GetProductAttributeByID)

		// Product Media
		v1.Get("/get-product-media/{id}", productRestHandler.GetProductMedia)

		// Product Variant
		v1.Get("/get-product-variants/{id}", productRestHandler.GetProductVariants)

		// Product SEO
		v1.Get("/get-product-seo/{id}", productRestHandler.GetProductSEO)

		// Social
		v1.Get("/auth/social/{provider}/callback", socialRestHandler.Callback)
		v1.Get("/auth/social/{provider}", socialRestHandler.BeginAuth)

		// Contact
		v1.Post("/contacts", contactRestHandler.CreateContact)

		// ToDo:
		// Create api of contact and request a quote and integrate it in frontend

		v1.Group(func(r chi.Router) {
			r.Use(jwtMiddleware)
			// Protected Routes
			// Brands
			r.With(shared.HasScope("catalog:create")).Post("/upload-brand-image", uploadHandler.UploadImage)
			r.With(shared.HasScope("catalog:create")).Post("/create-brand", brandRestHandler.CreateBrand)
			r.With(shared.HasScope("catalog:update")).Patch("/update-brand/{id}", brandRestHandler.UpdateBrand)
			r.With(shared.HasScope("catalog:delete")).Delete("/delete-brand/{id}", brandRestHandler.DeleteBrand)

			// Categories
			r.With(shared.HasScope("catalog:create")).Post("/upload-category-image", uploadHandler.UploadImage)
			r.With(shared.HasScope("catalog:create")).Post("/create-category", categoryRestHandler.CreateCategory)
			r.With(shared.HasScope("catalog:update")).Patch("/update-category/{id}", categoryRestHandler.UpdateCategory)
			r.With(shared.HasScope("catalog:delete")).Delete("/delete-category/{id}", categoryRestHandler.DeleteCategory)

			// Products
			r.With(shared.HasScope("catalog:create")).Post("/create-product", productRestHandler.CreateProduct)
			r.With(shared.HasScope("catalog:update")).Patch("/update-product/{id}", productRestHandler.UpdateProduct)
			r.With(shared.HasScope("catalog:delete")).Delete("/delete-product/{id}", productRestHandler.DeleteProduct)

			// Product Attributes
			r.With(shared.HasScope("catalog:update")).Post("/add-product-attribute", productRestHandler.AddProductAttribute)

			// Product Variants
			r.With(shared.HasScope("catalog:update")).Post("/add-product-variants/{id}", productRestHandler.SyncProductVariants)

			// Products Media
			r.With(shared.HasScope("catalog:update")).Post("/add-product-media/{id}", productRestHandler.AddProductMedia)

			// Product SEO
			r.With(shared.HasScope("catalog:update")).Post("/add-product-seo/{id}", productRestHandler.SaveProductSEO)

			// Social
			r.Get("/auth/social/{provider}/prepare", socialRestHandler.PrepareAuth)
			r.Get("/auth/social/connections", socialRestHandler.HandleListConnections)

			// Create post
			r.Post("/posts/create", socialRestHandler.HandleCreatePost)

			// Import Products
			r.Post("/product/import", productRestHandler.ImportProduct)

			// Contact
			r.Get("/contacts", contactRestHandler.GetAllContacts)
			r.Get("/contacts/{id}", contactRestHandler.GetContactByID)
			r.Patch("/contacts/{id}/status", contactRestHandler.UpdateContactStatus)
		})
	})

	container := &Container{
		HTTPRouter: router,
		GRPCServer: grpcSrv,
		DB:         sqlxDB,
		Logger:     l,
	}

	cleanup := func() {
		l.Sync()
		sqlxDB.Close()
		redisClient.Close()
	}

	return container, cleanup
}
