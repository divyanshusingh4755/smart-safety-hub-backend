package cache

import "fmt"

func ProductById(id string) string {
	return fmt.Sprintf("product:id:%s", id)
}

func ProductBySlug(slug string) string {
	return fmt.Sprintf("product:slug:%s", slug)
}

func ProductAttributes(id string) string {
	return fmt.Sprintf("product:attributes:%s", id)
}

func ProductMedia(id string) string {
	return fmt.Sprintf("product:media:%s", id)
}

func ProductVariants(id string) string {
	return fmt.Sprintf("product:variants:%s", id)
}

func ProductSEO(id string) string {
	return fmt.Sprintf("product:seo:%s", id)
}

func ProductList(filters string) string {
	return fmt.Sprintf("products:list:%s", filters)
}

func BrandByID(id string) string {
	return fmt.Sprintf("brand:id:%s", id)
}

func BrandList(page, limit int) string {
	return fmt.Sprintf("brand:list:p%d:l%d", page, limit)
}

func CategoryByID(id string) string {
	return fmt.Sprintf("category:id:%s", id)
}

func CategoryBySlug(id string) string {
	return fmt.Sprintf("category:slug:%s", id)
}

func CategoryList() string {
	return "category:list:all"
}
