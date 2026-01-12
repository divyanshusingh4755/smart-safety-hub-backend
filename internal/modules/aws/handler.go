package aws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
)

type UploadHandler struct {
	Service   *UploadService
	Validator *validator.Validate
}

func NewUploadHandler(service *UploadService, validator *validator.Validate) *UploadHandler {
	return &UploadHandler{
		Service:   service,
		Validator: validator,
	}
}

func (u *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		fmt.Println("er", err)
		http.Error(w, "Multipart parse error", http.StatusBadRequest)
		return
	}

	bucketName := r.PostFormValue("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name required", http.StatusBadRequest)
		return
	}

	headers := r.MultipartForm.File["file"]
	if len(headers) == 0 {
		http.Error(w, "No files provided", http.StatusBadRequest)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	responses := make([]interface{}, len(headers))

	for i, header := range headers {
		// Create local copies for the closure
		i, header := i, header

		g.Go(func() error {
			file, err := header.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			// Upload via Service
			resp, err := u.Service.UploadImage(ctx, file, header, bucketName)
			if err != nil {
				fmt.Println("errr", err)
				return err
			}

			responses[i] = resp
			return nil
		})
	}

	// 4. Wait for all uploads to finish
	if err := g.Wait(); err != nil {
		http.Error(w, "One or more uploads failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(responses),
		"data":    responses,
	})
}
