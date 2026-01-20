package app

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/smart-safety-hub/backend/internal/modules/aws"
	"github.com/smart-safety-hub/backend/internal/modules/brand"
	"github.com/smart-safety-hub/backend/internal/modules/categories"
	"github.com/smart-safety-hub/backend/internal/modules/products"
	"github.com/smart-safety-hub/backend/internal/modules/user"
	"github.com/smart-safety-hub/backend/shared"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	GrpcAddr   string
	HTTPAddr   string
	DBURL      string
	PrivateKey string
	PublicKey  string
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
	s3Client, err := shared.NewS3Client()
	if err != nil {
		log.Fatalf("Falied to init s3: %v", err)
	}

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
	brandService := brand.NewBrandService(l, brandRepo)
	brandRestHandler := brand.NewRestHandler(brandService, v)

	// Category
	categoryRepo := categories.NewCategoryRepo(sqlxDB)
	categoryService := categories.NewCategoryService(l, categoryRepo)
	categoryRestHandler := categories.NewRestHandler(categoryService, v)

	// Product
	productRepo := products.NewProductRepo(sqlxDB)
	productService := products.NewProductService(l, productRepo)
	productRestHandler := products.NewRestHandler(productService, v)

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
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
		v1.Get("/get-brand/{id}", brandRestHandler.GetBrandByID)
		v1.Get("/get-all-brands", brandRestHandler.GetAllBrand)

		// Category
		v1.Get("/get-category/{id}", categoryRestHandler.GetCategoryByID)
		v1.Get("/get-all-category", categoryRestHandler.GetAllCategory)

		// Product
		v1.Get("/get-product/id/{id}", productRestHandler.GetProductByID)
		v1.Get("/get-product/slug/{slug}", productRestHandler.GetProductBySlug)
		v1.Get("/get-all-products", productRestHandler.GetAllProducts)

		// Product Attribute
		v1.Get("/get-product-attribute/{id}", productRestHandler.GetProductAttributeByID)

		// Product Media
		v1.Get("/get-product-media/{id}", productRestHandler.GetProductMedia)

		// Product Variant
		v1.Get("/get-product-variants/{id}", productRestHandler.GetProductVariants)

		// Product SEO
		v1.Get("/get-product-seo/{id}", productRestHandler.GetProductSEO)
		v1.Group(func(r chi.Router) {
			r.Use(jwtMiddleware)
			// Protected Routes
			// Brands
			r.With(shared.HasScope("catalog:create")).Post("/upload-brand-image", uploadHandler.UploadImage)
			r.With(shared.HasScope("catalog:create")).Post("/create-brand", brandRestHandler.CreateBrand)
			r.With(shared.HasScope("catalog:update")).Patch("/update-brand/{id}", brandRestHandler.UpdateBrand)
			r.With(shared.HasScope("catalog:delete")).Delete("/delete-brand/{id}", brandRestHandler.DeleteBrand)

			// Categories
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
	}

	return container, cleanup
}
