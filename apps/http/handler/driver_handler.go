package handler

import (
	"encoding/json"

	"github.com/sibeur/gotaro/apps/http/handler/dto"
	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/common/driver"
	"github.com/sibeur/gotaro/core/entity"
	"github.com/sibeur/gotaro/core/service"

	"github.com/gofiber/fiber/v2"
)

type DriverHandler struct {
	fiberInstance *fiber.App
	svc           *service.Service
}

func NewDriverHandler(fiberInstance *fiber.App, svc *service.Service) *DriverHandler {
	return &DriverHandler{
		fiberInstance: fiberInstance,
		svc:           svc,
	}
}

func (h *DriverHandler) Router() {
	driver := h.fiberInstance.Group("/v1").Group("/drivers")
	driver.Get("/", h.findAllDrivers)
	driver.Get("/:slug", h.findDriverBySlug)
	driver.Post("/", h.createDriver)
	driver.Put("/:slug", h.updateDriver)
	driver.Delete("/:slug", h.deleteDriver)
}

func (h *DriverHandler) findAllDrivers(c *fiber.Ctx) error {
	drivers, err := h.svc.Driver.FindAll()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	response := make([]common.GotaroMap, 0)
	for _, driver := range drivers {
		response = append(response, driver.ToJSONSimple())
	}

	return successResponse(c, "Success", response, nil)
}

func (h *DriverHandler) findDriverBySlug(c *fiber.Ctx) error {
	driverSlug := c.Params("slug")

	driver, err := h.svc.Driver.FindBySlug(driverSlug)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	if driver == nil {
		return errorResponse(c, fiber.StatusNotFound, common.ErrDriverNotFoundMsg, nil, nil)
	}
	return successResponse(c, "", driver.ToJSON(), nil)
}

func (h *DriverHandler) createDriver(c *fiber.Ctx) error {
	driverData := new(dto.NewDriverDTO)

	if err := c.BodyParser(driverData); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)
	}

	validatorErrs := make([]common.FiberErrorMessage, 0)
	fValidator := common.NewFiberValidator()

	if errs := fValidator.Validate(driverData); len(errs) > 0 {
		validatorErrs = append(validatorErrs, errs...)
	}

	if !common.IsSlugValid(driverData.Slug) {
		errSlug := common.NewFiberErrorMessage("Slug", common.ErrSlugInvalidMsg)
		validatorErrs = append(validatorErrs, errSlug)
	}

	driverInput := &entity.Driver{
		Slug: driverData.Slug,
		Name: driverData.Name,
		Type: driverData.Type,
	}

	switch driverData.Type {
	case uint32(driver.GCSDriverType):
		errs, err := validateGCSConfig(driverInput, driverData.DriverConfig, c)
		if err != nil {
			return err
		}
		validatorErrs = append(validatorErrs, errs...)

	}

	if len(validatorErrs) > 0 {
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, validatorErrs, nil)
	}
	err := h.svc.Driver.Create(driverInput)
	if err != nil {
		if err.Error() == common.ErrDriverAlreadyExistMsg {
			return errorResponse(c, fiber.StatusConflict, err.Error(), nil, nil)
		}
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", driverInput.ToJSONSimple(), nil)
}

func (h *DriverHandler) updateDriver(c *fiber.Ctx) error {
	driverData := new(dto.EditDriverDTO)

	if err := c.BodyParser(driverData); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)
	}

	validatorErrs := make([]common.FiberErrorMessage, 0)
	fValidator := common.NewFiberValidator()

	if errs := fValidator.Validate(driverData); len(errs) > 0 {
		validatorErrs = append(validatorErrs, errs...)
	}

	driverSlug := c.Params("slug")

	driverInput := &entity.Driver{
		Slug: driverSlug,
		Name: driverData.Name,
		Type: driverData.Type,
	}

	switch driverData.Type {
	case uint32(driver.GCSDriverType):
		errs, err := validateGCSConfig(driverInput, driverData.DriverConfig, c)
		if err != nil {
			return err
		}
		validatorErrs = append(validatorErrs, errs...)
	}

	if len(validatorErrs) > 0 {
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, validatorErrs, nil)
	}
	err := h.svc.Driver.Update(driverInput)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", driverInput.ToJSONSimple(), nil)
}

func (h *DriverHandler) deleteDriver(c *fiber.Ctx) error {
	driverSlug := c.Params("slug")

	driver, err := h.svc.Driver.FindBySlug(driverSlug)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	if driver == nil {
		return errorResponse(c, fiber.StatusNotFound, common.ErrDriverNotFoundMsg, nil, nil)
	}

	err = h.svc.Driver.Delete(driverSlug)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", nil, nil)
}

func validateGCSConfig(driverInput *entity.Driver, driverConfigInput any, c *fiber.Ctx) ([]common.FiberErrorMessage, error) {
	validatorErrs := make([]common.FiberErrorMessage, 0)

	driverConfigJSON, err := json.Marshal(driverConfigInput)
	if err != nil {
		return nil, errorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)
	}
	driverConfig := new(dto.GCSDriverConfigDTO)
	if err := json.Unmarshal(driverConfigJSON, &driverConfig); err != nil {
		return nil, errorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)

	}
	serviceAccountJSON, err := driverConfig.GetServiceAccountJSONBytes()
	if err != nil {
		errServiceAccount := common.NewFiberErrorMessage("ServiceAccount", err.Error())
		validatorErrs = append(validatorErrs, errServiceAccount)
	}

	gcsValidator := common.NewFiberValidator()
	if errs := gcsValidator.Validate(driverConfig); len(errs) > 0 {
		validatorErrs = append(validatorErrs, errs...)
	}
	driverInput.DriverConfig = driver.NewGCSDriverConfig(driverConfig.ProjectID, driverConfig.BucketName, driverConfig.DefaultFolder, serviceAccountJSON)

	return validatorErrs, nil
}
