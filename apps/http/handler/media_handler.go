package handler

import (
	"io"
	"os"

	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/service"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	fiberInstance *fiber.App
	svc           *service.Service
}

func NewMediaHandler(fiberInstance *fiber.App, svc *service.Service) *MediaHandler {
	return &MediaHandler{
		fiberInstance: fiberInstance,
		svc:           svc,
	}
}

func (h *MediaHandler) Router() {
	medias := h.fiberInstance.Group("/v1").Group("/medias")
	medias.Post("/:slug", h.uploadMedia)
}

func (h *MediaHandler) findAllMedias(c *fiber.Ctx) error {
	medias, err := h.svc.Media.FindAll()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	response := make([]common.GotaroMap, 0)
	for _, media := range medias {
		response = append(response, media.ToJSONSimple())
	}

	return successResponse(c, "", response, nil)
}

func (h *MediaHandler) uploadMedia(c *fiber.Ctx) error {

	file, err := c.FormFile("file")
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer src.Close()

	// Create a destination file to save the uploaded file
	tempFilePath := common.TemporaryFolder + "/" + file.Filename
	dst, err := os.Create(tempFilePath)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	ruleSlug := c.Params("slug")

	url, err := h.svc.Media.Upload(ruleSlug, file.Filename)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", common.GotaroMap{"url": url}, nil)
}
