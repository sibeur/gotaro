package handler

import (
	"github.com/sibeur/gotaro/apps/http/handler/dto"
	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/entity"
	"github.com/sibeur/gotaro/core/service"

	"github.com/gofiber/fiber/v2"
)

type RuleHandler struct {
	fiberInstance *fiber.App
	svc           *service.Service
}

func NewRuleHandler(fiberInstance *fiber.App, svc *service.Service) *RuleHandler {
	return &RuleHandler{
		fiberInstance: fiberInstance,
		svc:           svc,
	}
}

func (h *RuleHandler) Router() {
	rules := h.fiberInstance.Group("/v1").Group("/rules")
	rules.Get("/", h.findAllRules)
	rules.Post("/", h.createRule)
	rules.Get("/:slug", h.findRuleBySlug)
	rules.Put("/:slug", h.updateRule)
	rules.Delete("/:slug", h.deleteRule)
}

func (h *RuleHandler) findAllRules(c *fiber.Ctx) error {
	rules, err := h.svc.Rule.FindAll()
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	response := make([]common.GotaroMap, 0)
	for _, rule := range rules {
		response = append(response, rule.ToJSONSimple())
	}

	return successResponse(c, "", response, nil)
}

func (h *RuleHandler) createRule(c *fiber.Ctx) error {
	ruleData := new(dto.NewRuleDTO)

	if err := c.BodyParser(ruleData); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)
	}

	fValidator := common.NewFiberValidator()

	if errs := fValidator.Validate(ruleData); len(errs) > 0 {
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, errs, nil)
	}

	if !common.IsSlugValid(ruleData.Slug) {
		errSlug := common.NewFiberErrorMessage("Slug", common.ErrSlugInvalidMsg)
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, []common.FiberErrorMessage{errSlug}, nil)
	}

	if !common.IsSlugValid(ruleData.DriverSlug) {
		errSlug := common.NewFiberErrorMessage("DriverSlug", common.ErrSlugInvalidMsg)
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, []common.FiberErrorMessage{errSlug}, nil)
	}

	existingDriver, err := h.svc.Driver.FindBySlug(ruleData.DriverSlug)
	if err != nil || existingDriver == nil {
		return errorResponse(c, fiber.StatusBadRequest, common.ErrDriverNotFoundMsg, nil, nil)
	}

	rule := &entity.Rule{
		Name:     ruleData.Name,
		Slug:     ruleData.Slug,
		MaxSize:  ruleData.MaxSize,
		Mimes:    ruleData.Mimes,
		DriverID: existingDriver.ID,
	}

	err = h.svc.Rule.Create(rule)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", rule, nil)
}

func (h *RuleHandler) findRuleBySlug(c *fiber.Ctx) error {
	ruleSlug := c.Params("slug")

	rule, err := h.svc.Rule.FindBySlug(ruleSlug)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	if rule == nil {
		return errorResponse(c, fiber.StatusNotFound, common.ErrRuleNotFoundMsg, nil, nil)
	}

	return successResponse(c, "", rule.ToJSON(), nil)
}

func (h *RuleHandler) updateRule(c *fiber.Ctx) error {
	ruleSlug := c.Params("slug")

	ruleData := new(dto.EditRuleDTO)

	if err := c.BodyParser(ruleData); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil, nil)
	}

	fValidator := common.NewFiberValidator()

	if errs := fValidator.Validate(ruleData); len(errs) > 0 {
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, errs, nil)
	}

	if !common.IsSlugValid(ruleData.DriverSlug) {
		errSlug := common.NewFiberErrorMessage("DriverSlug", common.ErrSlugInvalidMsg)
		return errorResponse(c, fiber.StatusBadRequest, common.ErrValidationMsg, []common.FiberErrorMessage{errSlug}, nil)
	}

	existingDriver, err := h.svc.Driver.FindBySlug(ruleData.DriverSlug)
	if err != nil || existingDriver == nil {
		return errorResponse(c, fiber.StatusBadRequest, common.ErrDriverNotFoundMsg, nil, nil)
	}

	rule := &entity.Rule{
		Name:     ruleData.Name,
		Slug:     ruleSlug,
		MaxSize:  ruleData.MaxSize,
		Mimes:    ruleData.Mimes,
		DriverID: existingDriver.ID,
	}

	err = h.svc.Rule.Update(rule)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", rule.ToJSONSimple(), nil)
}

func (h *RuleHandler) deleteRule(c *fiber.Ctx) error {
	ruleSlug := c.Params("slug")

	rule, err := h.svc.Rule.FindBySlug(ruleSlug)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	if rule == nil {
		return errorResponse(c, fiber.StatusNotFound, common.ErrRuleNotFoundMsg, nil, nil)
	}

	err = h.svc.Rule.Delete(ruleSlug)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, err.Error(), nil, nil)
	}

	return successResponse(c, "", nil, nil)
}
