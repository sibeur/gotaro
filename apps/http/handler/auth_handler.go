package handler

import (
	"github.com/sibeur/gotaro/apps/http/handler/dto"
	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandlerV1 struct {
	fiberInstance *fiber.App
	svc           *service.Service
}

func NewAuthHandlerV1(fiberInstance *fiber.App, svc *service.Service) *AuthHandlerV1 {
	return &AuthHandlerV1{
		fiberInstance: fiberInstance,
		svc:           svc,
	}
}

func (h *AuthHandlerV1) Router() {
	auth := h.fiberInstance.Group("/v1").Group("/auth")
	auth.Post("/login", h.login)
	auth.Get("/refresh-token", h.refreshToken)
}
func (h *AuthHandlerV1) login(c *fiber.Ctx) error {

	authData := new(dto.AuthDTO)

	if err := c.BodyParser(authData); err != nil {
		return common.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)
	}

	fValidator := common.NewFiberValidator()

	if errs := fValidator.Validate(authData); len(errs) > 0 {
		return common.ErrorResponse(c, fiber.StatusBadRequest, "Validation error", errs, nil)
	}

	response, err := h.svc.Auth.Login(authData.APIKey, authData.SecretKey)
	if err != nil {
		if err.Error() == common.ErrAuthenticationFailedMsg {
			return common.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil, nil)
		}
		return common.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return common.SuccessResponse(c, "Berhasil login", response, nil)
}

func (h *AuthHandlerV1) refreshToken(c *fiber.Ctx) error {
	// get refresh token from header
	refreshToken := c.Get("Authorization")

	// remove "Bearer " from token
	refreshToken = refreshToken[7:]

	response, err := h.svc.Auth.RefreshToken(refreshToken)
	if err != nil {
		if err.Error() == common.ErrJWTTokenInvalidMsg {
			return common.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil, nil)
		}
		return common.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}
	return common.SuccessResponse(c, "Berhasil refresh token", response, nil)
}
