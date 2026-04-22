package apimiddleware

import (
	"strings"

	"github.com/PopularVote/internal/infrastructure/auth"
	"github.com/gofiber/fiber/v3"
)

const (
	LocalUserID = "api_user_id"
	LocalEmail  = "api_email"
)

// JWTProtected validates the Bearer token and stores claims in locals.
func JWTProtected(jwtSvc *auth.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtSvc.ValidateAccessToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals(LocalUserID, claims.UserID)
		c.Locals(LocalEmail, claims.Email)

		return c.Next()
	}
}
