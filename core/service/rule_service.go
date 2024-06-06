package service

import (
	"errors"

	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/entity"
	"github.com/sibeur/gotaro/core/repository"
)

type RuleService struct {
	repo *repository.Repository
}

func NewRuleService(repo *repository.Repository) *RuleService {
	return &RuleService{repo: repo}
}

func (u *RuleService) FindAll() ([]*entity.Rule, error) {
	result, err := u.repo.Rule.FindAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *RuleService) Create(rule *entity.Rule) error {
	existingRule, err := u.repo.Rule.FindBySlug(rule.Slug)
	if err == nil && existingRule != nil {
		return errors.New(common.ErrRuleAlreadyExistMsg)
	}
	return u.repo.Rule.Create(rule)
}

func (u *RuleService) Update(rule *entity.Rule) error {
	return u.repo.Rule.Update(rule)
}

func (u *RuleService) Delete(slug string) error {
	return u.repo.Rule.Delete(slug)
}

func (u *RuleService) FindBySlug(slug string) (*entity.Rule, error) {
	return u.repo.Rule.FindBySlug(slug)
}
