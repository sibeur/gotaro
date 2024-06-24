package handler

import (
	"io"
	"log"
	"os"

	"github.com/sibeur/gotaro/apps/http/handler/dto"
	"github.com/sibeur/gotaro/apps/http/handler/middleware"
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
	medias := h.fiberInstance.Group("/v1").Group("/medias", middleware.VerifyAuth(h.svc))
	medias.Post("/get-batch", h.getMediaBatch)
	medias.Post("/:slug", middleware.VerifyAuthAudiences([]string{common.APIClientSuperAdminScope, common.APIClientUploaderScope}), h.uploadMedia)
	medias.Get("/:slug/:fileAliasName", middleware.VerifyAuthAudiences([]string{common.APIClientSuperAdminScope, common.APIClientUploaderScope}), h.getMedia)
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
		log.Printf("Error uploading file: %v", err)
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	isCommit := false
	isCommitString := c.FormValue("commit")

	if isCommitString != "" && isCommitString == "true" {
		isCommit = true
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer src.Close()

	// Create a destination file to save the uploaded file
	tempFilePath := common.TemporaryFolder + "/" + file.Filename
	dst, err := os.Create(tempFilePath)
	if err != nil {
		log.Printf("Error creating destination file: %v", err)
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}
	defer dst.Close()

	defer func() {
		err := os.Remove(tempFilePath)
		if err != nil {
			log.Printf("Failed to delete temp file %v", tempFilePath)
		}
	}()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("Error copying file: %v", err)
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	ruleSlug := c.Params("slug")

	media, err := h.svc.Media.Upload(ruleSlug, file.Filename, isCommit)
	if err != nil {
		log.Printf("Error uploading media: %v", err)
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", media.ToMediaResult(), nil)
}

func (h *MediaHandler) getMedia(c *fiber.Ctx) error {
	ruleSlug := c.Params("slug")
	fileAliasName := c.Params("fileAliasName")

	media, err := h.svc.Media.FindMedia(ruleSlug, fileAliasName)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	if media == nil {
		return errorResponse(c, fiber.StatusNotFound, common.ErrMediaNotFoundMsg, nil, nil)
	}

	return successResponse(c, "", media.ToMediaResult(), nil)
}

func (h *MediaHandler) getMediaBatch(c *fiber.Ctx) error {
	mediaData := new(dto.GetMediaBatchDTO)

	if err := c.BodyParser(mediaData); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)
	}

	fValidator := common.NewFiberValidator()

	if errs := fValidator.Validate(mediaData); len(errs) > 0 {
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, errs, nil)
	}

	medias := h.svc.Media.FindMediaBatch(mediaData.Files)
	return successResponse(c, "", medias, nil)
}
