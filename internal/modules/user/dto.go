package user

import "time"

type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"user_type" validate:"required"`
}

type RegisterDTO struct {
	FullName    string `json:"full_name" validate:"required,min=3"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=12,max=72"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	UserType    string `json:"user_type" validate:"required"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutDTO struct {
	UserId string `json:"user_id" validate:"required"`
}

type ResetPasswordDTO struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GenericResponseDTO struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseDTO struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	UserInfo     UserInformation `json:"user_info"`
	ExpiresIn    time.Time       `json:"expires_in"`
	Status       string          `json:"status"`
	Message      string          `json:"message"`
}

type UserInformation struct {
	UserId      string  `json:"user_id"`
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	CompanyId   *string `json:"company_id"`
	Roles       string  `json:"roles"`
	Permissions string  `json:"permissions"`
}
