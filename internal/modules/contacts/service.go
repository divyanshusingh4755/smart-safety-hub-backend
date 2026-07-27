package contacts

import (
	"context"
	"fmt"
	"strings"

	"github.com/smart-safety-hub/backend/shared/mail"
	"go.uber.org/zap"
)

type ContactService struct {
	logger     *zap.Logger
	repo       *ContactRepo
	mailer     *mail.Mailer
	adminEmail string
}

func NewContactService(
	logger *zap.Logger,
	repo *ContactRepo,
	mailer *mail.Mailer,
	adminEmail string,
) *ContactService {

	return &ContactService{
		logger:     logger,
		repo:       repo,
		mailer:     mailer,
		adminEmail: adminEmail,
	}
}

func (s *ContactService) CreateContact(
	ctx context.Context,
	request CreateContactRequestDTO,
) (*CreateContactResponseDTO, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Phone = strings.TrimSpace(request.Phone)

	contact, err := s.repo.SaveContact(
		ctx,
		request,
	)

	if err != nil {
		s.logger.Error(
			"failed to create contact enquiry",
			zap.Error(err),
		)

		return nil, fmt.Errorf(
			"error saving contact enquiry: %v",
			err,
		)
	}

	// Send admin notification
	adminHTML := CreateAdminContactEmail(*contact)
	adminMessage := mail.Message{
		To: []string{
			s.adminEmail,
		},

		Subject: fmt.Sprintf(
			"New Website Enquiry - %s",
			contact.Name,
		),

		HTML: adminHTML,
	}

	// When customer email exists, clicking Reply will reply directly to that customer

	if contact.Email != nil {
		adminMessage.ReplyTo = *contact.Email
	}

	if err := s.mailer.Send(adminMessage); err != nil {
		s.logger.Error("failed to send admin contact notification", zap.String("contact_id", contact.ID), zap.Error(err))
	}

	// Send customer confirmation only if email was provided
	if contact.Email != nil && strings.TrimSpace(*contact.Email) != "" {
		customerHTML := CreateCustomerContactEmail(*contact)

		if err := s.mailer.Send(mail.Message{
			To: []string{
				*contact.Email,
			},

			Subject: "We Received Your Enquiry | Smart Safety Hub",

			HTML: customerHTML,
		}); err != nil {
			s.logger.Error("failed to send customer contact confirmation", zap.String("contact_id", contact.ID), zap.String("email", *contact.Email), zap.Error(err))
		}
	}

	return &CreateContactResponseDTO{
		Status:  "success",
		Message: "Enquiry submitted successfully",
		ID:      contact.ID,
	}, nil
}

func mapContactResponse(
	contact Contact,
) ContactResponseDTO {
	return ContactResponseDTO{
		ID:          contact.ID,
		Name:        contact.Name,
		CompanyName: contact.CompanyName,
		Email:       contact.Email,
		Phone:       contact.Phone,
		Country:     contact.Country,
		InquiryType: contact.InquiryType,
		ProductName: contact.ProductName,
		Quantity:    contact.Quantity,
		Message:     contact.Message,
		Source:      contact.Source,
		PageTitle:   contact.PageTitle,
		PageURL:     contact.PageURL,
		Status:      contact.Status,
		CreatedAt:   contact.CreatedAt,
		UpdatedAt:   contact.UpdatedAt,
	}
}

func (s *ContactService) GetContactByID(
	ctx context.Context,
	contactID string,
) (*ContactResponseDTO, error) {
	contact, err := s.repo.GetContactByID(
		ctx,
		contactID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error fetching contact enquiry: %v",
			err,
		)
	}

	response := mapContactResponse(*contact)
	return &response, nil
}

func (s *ContactService) GetAllContacts(
	ctx context.Context,
	request GetContactsQueryDTO,
) (*GetAllContactsResponseDTO, error) {

	contacts, total, err := s.repo.GetAllContacts(
		ctx,
		request,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"error fetching contact enquiries: %v",
			err,
		)
	}

	responseList := make(
		[]ContactResponseDTO,
		len(contacts),
	)

	for i, contact := range contacts {
		responseList[i] =
			mapContactResponse(contact)
	}

	totalPages := 0

	if total > 0 {
		totalPages =
			(total + request.Limit - 1) /
				request.Limit
	}

	return &GetAllContactsResponseDTO{
		Contacts: responseList,

		Pagination: ContactPaginationDTO{
			Page:       request.Page,
			Limit:      request.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ContactService) UpdateContactStatus(
	ctx context.Context,
	contactID string,
	request UpdateContactStatusDTO,
) (*GenericResponseDTO, error) {
	if err := s.repo.UpdateContactStatus(
		ctx,
		contactID,
		request.Status,
	); err != nil {
		return nil, fmt.Errorf(
			"error updating contact status: %v",
			err,
		)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Contact status updated successfully",
	}, nil
}
