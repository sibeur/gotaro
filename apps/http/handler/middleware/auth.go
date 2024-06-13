package middleware

import (
	"log"

	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/service"

	"github.com/gofiber/fiber/v2"
)

func VerifyAuth(svc *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		bearerToken := c.Get("Authorization")
		if bearerToken == "" {
			log.Printf("[VerifyAuth] Bearer token not found")
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}
		bearerToken = bearerToken[7:]

		// Verify token
		token, err := svc.Auth.ValidateToken(bearerToken)
		if err != nil {
			log.Printf("[VerifyAuth] ValidateToken Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		// Get subject from token
		claims := token.Claims
		subject, err := claims.GetSubject()
		if err != nil {
			log.Printf("[VerifyAuth] GetSubject Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		audiences, err := claims.GetAudience()
		if err != nil {
			log.Printf("[VerifyAuth] GetAudience Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		// check issuer
		issuer, err := claims.GetIssuer()
		if err != nil {
			log.Printf("[VerifyAuth] GetIssuer Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}
		if issuer != common.JWTIssuerAccessToken {
			log.Printf("[VerifyAuth] Issuer Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		// Set user_id to locals
		c.Locals("user_id", subject)
		// set role_ids to locals
		c.Locals("role_ids", []string(audiences))

		// Continue stack
		return c.Next()
	}
}

func VerifyAuthWithUserData(svc *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		bearerToken := c.Get("Authorization")
		if bearerToken == "" {
			log.Printf("[VerifyAuthWithUserData] Bearer token not found")
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}
		bearerToken = bearerToken[7:]

		// Verify token
		token, err := svc.Auth.ValidateToken(bearerToken)
		if err != nil {
			log.Printf("[VerifyAuthWithUserData] ValidateToken Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		// Get subject from token
		claims := token.Claims
		subject, err := claims.GetSubject()
		if err != nil {
			log.Printf("[VerifyAuthWithUserData] GetSubject Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		audiences, err := claims.GetAudience()
		if err != nil {
			log.Printf("[VerifyAuthWithUserData] GetAudience Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		// Check issuer
		issuer, err := claims.GetIssuer()
		if err != nil {
			log.Printf("[VerifyAuthWithUserData] GetIssuer Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}
		if issuer != common.JWTIssuerAccessToken {
			log.Printf("[VerifyAuthWithUserData] Invalid issuer")
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}
		// Get user data
		user, err := svc.ApiClient.FindByID(subject)
		if err != nil {
			log.Printf("[VerifyAuthWithUserData] FindByID Error: %v", err)
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		if user == nil {
			log.Printf("[VerifyAuthWithUserData] User not found")
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		// Set user id to locals
		c.Locals("user_id", user.ID)

		// set role_ids to locals
		c.Locals("role_ids", []string(audiences))

		// Set user data to locals
		c.Locals("user", user)

		// Continue stack
		return c.Next()
	}
}

func VerifyAuthAudiences(audiences []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get role ids from locals
		roleIDs := c.Locals("role_ids").([]string)

		// throw error if role ids not found
		if len(roleIDs) == 0 {
			return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
		}

		// Check role ids
		for _, roleID := range roleIDs {
			for _, audience := range audiences {
				if audience == roleID {
					return c.Next()
				}
			}
		}
		return common.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorizedMsg, nil, nil)
	}
}
