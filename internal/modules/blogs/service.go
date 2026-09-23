package blogs

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type BlogService struct {
	logger *zap.Logger
	repo   *BlogRepo
}

func NewBlogService(logger *zap.Logger, repo *BlogRepo) *BlogService {
	return &BlogService{
		logger: logger,
		repo:   repo,
	}
}

func (s *BlogService) CreateBlog(ctx context.Context, request BlogRequestDTO) (*GenericResponseDTO, error) {
	if err := s.repo.SaveBlog(ctx, request); err != nil {
		return nil, fmt.Errorf("error while saving blog to DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Blog Created Successfully",
	}, nil
}

func (s *BlogService) UpdateBlog(ctx context.Context, blogID string, request UpdateBlogDTO) (*GenericResponseDTO, error) {
	if err := s.repo.UpdateBlog(ctx, blogID, request); err != nil {
		return nil, fmt.Errorf("error while updating blog in DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Blog Updated Successfully",
	}, nil
}

func (s *BlogService) DeleteBlog(ctx context.Context, blogID string) (*GenericResponseDTO, error) {
	if err := s.repo.DeleteBlog(ctx, blogID); err != nil {
		return nil, fmt.Errorf("error while deleting blog from DB: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Blog Deleted Successfully",
	}, nil
}

func (s *BlogService) GetBlogByID(ctx context.Context, blogID string) (*BlogResponse, error) {
	blog, err := s.repo.GetBlogByID(ctx, blogID)
	if err != nil {
		return nil, fmt.Errorf("error while getting blog from DB: %v", err)
	}

	return mapBlogResponse(blog), nil
}

func (s *BlogService) GetBlogBySlug(ctx context.Context, slug string) (*BlogResponse, error) {
	blog, err := s.repo.getBlogBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("error while getting blog from DB: %v", err)
	}

	return mapBlogResponse(blog), nil
}

func (s *BlogService) GetAllBlogs(ctx context.Context) (*GetAllBlogs, error) {
	blogs, err := s.repo.GetAllBlogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting blogs from DB: %v", err)
	}

	blogList := make([]BlogResponse, len(blogs))
	for i, blog := range blogs {
		blogList[i] = *mapBlogResponse(&blog)
	}

	return &GetAllBlogs{
		Blogs: blogList,
	}, nil
}

func (s *BlogService) GetPublishedBlogs(ctx context.Context) (*GetAllBlogs, error) {

	blogs, err := s.repo.GetPublishedBlogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting published blogs from DB: %v", err)
	}

	response := make([]BlogResponse, 0, len(blogs))

	for _, blog := range blogs {
		response = append(response, *mapBlogResponse(&blog))
	}

	return &GetAllBlogs{Blogs: response}, nil
}

func (s *BlogService) GetPublishedPages(ctx context.Context) (*GetAllBlogs, error) {
	pages, err := s.repo.GetPublishedPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while getting published pages from DB: %v", err)
	}

	response := make([]BlogResponse, 0, len(pages))
	for _, page := range pages {
		response = append(response, *mapBlogResponse(&page))
	}

	return &GetAllBlogs{Blogs: response}, nil
}

func mapBlogResponse(blog *Blog) *BlogResponse {
	return &BlogResponse{
		ID:              blog.ID,
		Title:           blog.Title,
		Slug:            blog.Slug,
		Type:            blog.Type,
		Excerpt:         blog.Excerpt,
		Content:         blog.Content,
		FeaturedImage:   blog.FeaturedImage,
		MetaTitle:       blog.MetaTitle,
		MetaDescription: blog.MetaDescription,
		IsPublished:     blog.IsPublished,
		PublishedAt:     blog.PublishedAt,
		CreatedAt:       blog.CreatedAt,
		UpdatedAt:       blog.UpdatedAt,
	}
}
