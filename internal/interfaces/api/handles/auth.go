package apihandles

import (
	"github.com/PopularVote/internal/domain/user"
	"github.com/PopularVote/internal/infrastructure/auth"
	sqlrepo "github.com/PopularVote/internal/infrastructure/sql"
	apimiddleware "github.com/PopularVote/internal/interfaces/api/middleware"
	"github.com/gofiber/fiber/v3"
)

// AuthAPIHandler handles JWT-based authentication endpoints.
type AuthAPIHandler struct {
	jwtSvc *auth.JWTService
}

func NewAuthAPIHandler(jwtSvc *auth.JWTService) *AuthAPIHandler {
	return &AuthAPIHandler{
		jwtSvc: jwtSvc,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}

// Register handles user registration.
func (h *AuthAPIHandler) Register(c fiber.Ctx) error {
	p := new(RegisterRequest)
	if err := c.Bind().JSON(p); err != nil {
		return err
	}
	db := c.Locals("db").(*sqlrepo.DB)
	repo := sqlrepo.NewUserRepo(db)
	userService := user.NewService(repo)
	createdUser, err := userService.CreateUser(c.Context(), p.Name, p.Email, p.Password)
	if err != nil {
		return err
	}

	acess, err := h.jwtSvc.GenerateAccessToken(createdUser.PublicID, createdUser.Email)
	if err != nil {
		return err
	}
	refresh, err := h.jwtSvc.GenerateRefreshToken(createdUser.PublicID)
	if err != nil {
		return err
	}
	return c.JSON(tokenResponse{
		AccessToken:  acess,
		RefreshToken: refresh,
		ExpiresIn:    h.jwtSvc.GetAccessTokenTTL(),
	})
}

// Login handles user login and JWT issuance.
func (h *AuthAPIHandler) Login(c fiber.Ctx) {
	p := new(LoginRequest)
	if err := c.Bind().JSON(p); err != nil {
		return
	}
	db := c.Locals("db").(*sqlrepo.DB)
	repo := sqlrepo.NewUserRepo(db)
	userService := user.NewService(repo)
	authenticatedUser, err := userService.AuthenticateUser(c.Context(), p.Email, p.Password)
	if err != nil {
		return
	}

	acess, err := h.jwtSvc.GenerateAccessToken(authenticatedUser.PublicID, authenticatedUser.Email)
	if err != nil {
		return
	}
	refresh, err := h.jwtSvc.GenerateRefreshToken(authenticatedUser.PublicID)
	if err != nil {
		return
	}
	c.JSON(tokenResponse{
		AccessToken:  acess,
		RefreshToken: refresh,
		ExpiresIn:    h.jwtSvc.GetAccessTokenTTL(),
	})

}

// RefreshToken handles access token refreshing using a valid refresh token.
func (h *AuthAPIHandler) RefreshToken(c fiber.Ctx) {
	p := new(struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	})
	if err := c.Bind().JSON(p); err != nil {
		return
	}
	userID, err := h.jwtSvc.ValidateRefreshToken(p.RefreshToken)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired refresh token"})
		return
	}
	// check user exists in database
	db := c.Locals("db").(*sqlrepo.DB)
	repo := sqlrepo.NewUserRepo(db)
	userService := user.NewService(repo)
	_, err = userService.GetUserByID(c.Context(), userID)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		return
	}

	// Generate new access token
	accessToken, err := h.jwtSvc.GenerateAccessToken(userID, "")
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate access token"})
		return
	}
	refresh, err := h.jwtSvc.GenerateRefreshToken(userID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refresh, // Keep the same refresh token
		ExpiresIn:    h.jwtSvc.GetAccessTokenTTL(),
	})
}

func (h *AuthAPIHandler) Logout(c *fiber.Ctx) {
	// For stateless JWT, logout is typically handled client-side by deleting tokens.
	// todo add blacklist for refresh tokens if needed.
}

func (h *AuthAPIHandler) Me(c fiber.Ctx) error {
	// This endpoint can return user info based on the access token.
	return c.JSON(fiber.Map{
		"user_id": c.Locals(apimiddleware.LocalUserID),
		"email":   c.Locals(apimiddleware.LocalEmail),
	})
}
