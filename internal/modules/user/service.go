package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UserService struct {
	logger *zap.Logger
	repo   *UserRepo
	jwt    *JwtManager
}

func NewUserService(logger *zap.Logger, repo *UserRepo, jwt *JwtManager) *UserService {
	return &UserService{
		logger: logger,
		jwt:    jwt,
		repo:   repo,
	}
}

func (u *UserService) Register(ctx context.Context, req *RegisterDTO) (*ResponseDTO, error) {
	user := User{
		FullName:    req.FullName,
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	var userResponse *User
	var rolesPermissions *RolesPermissions

	// Use the tranasaction
	err := u.repo.ExecuteTransaction(ctx, func(repo *UserRepo) error {
		var err error
		// Save user
		userResponse, err = repo.SaveUser(ctx, &user)
		if err != nil {
			return fmt.Errorf("Error came while saving user to DB: %v", err)
		}

		// Get Roles
		role, err := repo.GetRole(ctx, req.UserType)
		if err != nil {
			return fmt.Errorf("Error came getting role from DB: %v", err)
		}

		// Save roles of that user in DB
		if err := repo.SaveUserRoles(ctx, userResponse.ID, role.ID); err != nil {
			return fmt.Errorf("Error while saving user role in DB: %v", err)
		}

		// Get roles and permission and map the data with role ID and Permission ID
		rolesPermissions, err = repo.GetRolesPermissions(ctx, role.ID)
		if err != nil {
			return fmt.Errorf("Error came while getting data from DB: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error came while registering user: %v", err)
	}

	userInfo := UserInformation{
		UserId:      userResponse.ID,
		Email:       userResponse.Email,
		FullName:    userResponse.FullName,
		Roles:       rolesPermissions.Role,
		Permissions: string(rolesPermissions.Permissions),
	}

	response := &ResponseDTO{
		UserInfo: userInfo,
		Status:   "success",
		Message:  "User registered successfully",
	}

	return response, nil
}

func (u *UserService) Login(ctx context.Context, req *LoginDTO) (*ResponseDTO, error) {
	// Get Login Details
	userResponse, err := u.repo.GetUser(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error came while getting login data from DB: %v", err)
	}

	// Check whether password is same or not
	if err := userResponse.ComparePassword(req.Password); err != nil {
		return nil, err
	}

	// Get user Roles
	userRoles, err := u.repo.GetUserRoles(ctx, userResponse.ID)
	if err != nil {
		return nil, fmt.Errorf("error came while getting role data from DB: %v", err)
	}

	// Map user roles with permissions
	rolesPermissions, err := u.repo.GetRolesPermissions(ctx, userRoles.RoleId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting user roles permissions data from DB: %v", err)
	}

	// 3. Generate New Tokens
	accessToken, err := u.jwt.GenerateToken(userResponse.ID, rolesPermissions, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("Error came while generating access token: %v", err)
	}

	newRefreshToken, err := u.jwt.GenerateToken(userResponse.ID, nil, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("Error while generating refresh token: %v", err)
	}

	if err := u.repo.SaveRefreshToken(ctx, userResponse.ID, newRefreshToken, 30*24*time.Hour); err != nil {
		return nil, fmt.Errorf("Error came while saving refresh token to DB: %v", err)
	}

	response := &ResponseDTO{
		UserInfo: UserInformation{
			UserId:      userResponse.ID,
			Email:       userResponse.Email,
			FullName:    userResponse.FullName,
			Roles:       rolesPermissions.Role,
			Permissions: string(rolesPermissions.Permissions),
		},
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    time.Now().Add(15 * time.Minute),
		Status:       "success",
		Message:      "User Logged In successfully",
	}
	return response, nil
}

func (u *UserService) ForgotPassword(ctx context.Context, req *ForgotPasswordDTO) (*GenericResponseDTO, error) {
	userResponse, err := u.repo.GetUser(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("Invalid email: %v", err)
	}

	// Generate Access Token with time 24hours and send mail
	token, err := u.jwt.GenerateToken(userResponse.ID, nil, 900*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error came while generating access token: %v", err)
	}

	// TODO: send mail
	fmt.Println("token", token)

	return &GenericResponseDTO{
		Status:  "success",
		Message: "Email sent successfully",
	}, nil
}

func (u *UserService) ResetPassword(ctx context.Context, req *ResetPasswordDTO) (*GenericResponseDTO, error) {
	var userID any
	var user User

	token, err := u.jwt.Verify(req.Token)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong. Try again: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("Token expired please resend email")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims["sub"]
	}

	if s, ok := userID.(string); ok {
		user = User{
			ID:       s,
			Password: req.Password,
		}
	} else {
		return nil, fmt.Errorf("Something went wrong. Please try again: %v", err)
	}

	if err = user.HashPassword(); err != nil {
		return nil, fmt.Errorf("Error came while hashing password: %v", err)
	}

	if err = u.repo.UpdatePassword(ctx, user.Password, user.ID); err != nil {
		return nil, fmt.Errorf("Error while updating password: %v", err)
	}

	response := &GenericResponseDTO{
		Status:  "success",
		Message: "Password Changed successfully",
	}

	return response, nil
}

func (u *UserService) Logout(ctx context.Context, request LogoutDTO) (*GenericResponseDTO, error) {
	if err := u.repo.RevokeRefreshToken(ctx, request.UserId); err != nil {
		return nil, fmt.Errorf("Error came while revoking old refresh token: %v", err)
	}

	return &GenericResponseDTO{
		Status:  "success",
		Message: "User logged out successfully",
	}, nil
}

func (u *UserService) RefreshToken(ctx context.Context, request RefreshTokenDTO) (*ResponseDTO, error) {
	fmt.Println(request.RefreshToken)
	storedToken, err := u.repo.GetRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, errors.New("Invalid or expired token")
	}

	// Get User
	user, err := u.repo.GetUserById(ctx, storedToken.UserId)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// Get user Roles
	userRoles, err := u.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error came while getting role data from DB: %v", err)
	}

	// Map user roles with permissions
	rolesPermissions, err := u.repo.GetRolesPermissions(ctx, userRoles.RoleId)
	if err != nil {
		return nil, fmt.Errorf("Error came while getting user roles permissions data from DB: %v", err)
	}

	// Generate Access and Refresh Token
	accessToken, err := u.jwt.GenerateToken(user.ID, rolesPermissions, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("Error came while generating access token: %v", err)
	}

	newRefreshToken, err := u.jwt.GenerateToken(user.ID, nil, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("Error while generating refresh token: %v", err)
	}

	err = u.repo.ExecuteTransaction(ctx, func(repo *UserRepo) error {
		repo.UpdateRefreshToken(ctx, request.RefreshToken)
		return repo.SaveRefreshToken(ctx, user.ID, newRefreshToken, 30*24*time.Hour)
	})

	response := &ResponseDTO{
		UserInfo: UserInformation{
			UserId:      user.ID,
			Email:       user.Email,
			FullName:    user.FullName,
			Roles:       rolesPermissions.Role,
			Permissions: string(rolesPermissions.Permissions),
		},
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    time.Now().Add(15 * time.Minute),
		Status:       "success",
		Message:      "Refresh Token Changed Successfully",
	}

	return response, nil
}
