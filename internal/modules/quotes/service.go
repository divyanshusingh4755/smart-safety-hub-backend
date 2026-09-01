package quotes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/smart-safety-hub/backend/shared/mail"
	"go.uber.org/zap"
)

type QuoteService struct {
	logger     *zap.Logger
	repo       *QuoteRepo
	mailer     *mail.Mailer
	adminEmail string
}

func NewQuoteService(logger *zap.Logger, repo *QuoteRepo, mailer *mail.Mailer, adminEmail string) *QuoteService {
	return &QuoteService{
		logger:     logger,
		repo:       repo,
		mailer:     mailer,
		adminEmail: adminEmail,
	}
}

func (s *QuoteService) CreateQuote(ctx context.Context, request CreateQuoteRequestDTO) (*CreateQuoteResponseDTO, error) {
	request.FullName = strings.TrimSpace(request.FullName)
	request.CompanyName = strings.TrimSpace(request.CompanyName)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Phone = strings.TrimSpace(request.Phone)
	request.DeliveryLocation = strings.TrimSpace(request.DeliveryLocation)

	var requiredBy *time.Time
	if request.RequiredBy != nil {
		value := strings.TrimSpace(*request.RequiredBy)
		if value != "" {
			parsedDate, err := time.Parse("2006-01-02", value)
			if err != nil {
				return nil, fmt.Errorf("invalid required by date")
			}

			requiredBy = &parsedDate
		}
	}

	quote, items, err := s.repo.SaveQuote(ctx, request, *requiredBy)
	if err != nil {
		s.logger.Error("failed to create quote request", zap.Error(err))
		return nil, fmt.Errorf("error saving quote request: %v", err)
	}
	// Email notifications
	// Admin RFQ Notification
	adminHTML := CreateAdminQuoteEmail(*quote, items)
	adminMessage := mail.Message{
		To:      []string{s.adminEmail},
		Subject: fmt.Sprintf("New RFQ %s - %s", quote.QuoteNumber, quote.CompanyName),
		HTML:    adminHTML,
	}

	// Clicking Reply in admin email should reply directly to customer
	if strings.TrimSpace(quote.Email) != "" {
		adminMessage.ReplyTo = quote.Email
	}

	if err := s.mailer.Send(adminMessage); err != nil {
		s.logger.Error("failed to send admin quote notification", zap.String("quote_id", quote.ID), zap.String("quote_number", quote.QuoteNumber), zap.Error(err))
	}

	// Customer RFQ Confirmation
	if strings.TrimSpace(quote.Email) != "" {
		customerHTML := CreateCustomerQuoteEmail(*quote, items)
		if err := s.mailer.Send(mail.Message{
			To:      []string{quote.Email},
			Subject: fmt.Sprintf("Quote Request Received - %s | Smart Safety Hub", quote.QuoteNumber),
			HTML:    customerHTML,
		}); err != nil {
			s.logger.Error("failed to send customer quote confirmation", zap.String("quote_id", quote.ID), zap.String("quote_number", quote.QuoteNumber), zap.String("email", quote.Email), zap.Error(err))
		}
	}

	return &CreateQuoteResponseDTO{
		Status:      "success",
		Message:     "Quote request submitted successfully",
		ID:          quote.ID,
		QuoteNumber: quote.QuoteNumber,
	}, nil
}

func mapQuoteResponse(quote Quote, items []QuoteItem) QuoteResponseDTO {
	itemResponses := make([]QuoteItemResponseDTO, len(items))

	for i, item := range items {
		itemResponses[i] = QuoteItemResponseDTO{
			ID:         item.ID,
			ProductRef: item.ProductRef,
			Name:       item.Name,
			SKU:        item.SKU,
			Image:      item.Image,
			Quantity:   item.Quantity,
		}
	}

	return QuoteResponseDTO{
		ID:                     quote.ID,
		QuoteNumber:            quote.QuoteNumber,
		FullName:               quote.FullName,
		CompanyName:            quote.CompanyName,
		Email:                  quote.Email,
		Phone:                  quote.Phone,
		DeliveryLocation:       quote.DeliveryLocation,
		RequiredBy:             quote.RequiredBy,
		AdditionalRequirements: quote.AdditionalRequirements,
		Status:                 quote.Status,
		Items:                  itemResponses,
		CreatedAt:              quote.CreatedAt,
		UpdatedAt:              quote.UpdatedAt,
	}
}

func (s *QuoteService) GetQuoteByID(ctx context.Context, quoteID string) (*QuoteResponseDTO, error) {
	quote, items, err := s.repo.GetQuoteByID(ctx, quoteID)
	if err != nil {
		return nil, fmt.Errorf("error fetching quote request: %v", err)
	}

	response := mapQuoteResponse(*quote, items)
	return &response, nil
}

func (s *QuoteService) UpdateQuoteStatus(ctx context.Context, quoteID string, request UpdateQuoteStatusDTO) (*GenericResponseDTO, error) {
	if err := s.repo.UpdateQuoteStatus(ctx, quoteID, request.Status); err != nil {
		return nil, fmt.Errorf("error updating quote status: %v", err)
	}

	return &GenericResponseDTO{Status: "success", Message: "Quote status updated successfully"}, nil
}

func (s *QuoteService) GetAllQuotes(ctx context.Context, request GetQuotesQueryDTO) (*GetAllQuotesResponseDTO, error) {
	quotes, total, err := s.repo.GetAllQuotes(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error fetching quote requests: %v", err)
	}

	responseList := make([]QuoteResponseDTO, len(quotes))

	for i, quote := range quotes {
		responseList[i] = mapQuoteResponse(quote, nil)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + request.Limit - 1) / request.Limit
	}

	return &GetAllQuotesResponseDTO{
		Quotes: responseList,
		Pagination: QuotePaginationDTO{
			Page:       request.Page,
			Limit:      request.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
