package service

import (
	"errors"
	"log"

	"github.com/sibeur/gotaro/core/common"
	driver_lib "github.com/sibeur/gotaro/core/common/driver"
	"github.com/sibeur/gotaro/core/entity"
	"github.com/sibeur/gotaro/core/repository"
)

type DriverService struct {
	repo          *repository.Repository
	DriverManager *driver_lib.DriverManager
}

func NewDriverService(repo *repository.Repository, driverManager *driver_lib.DriverManager) *DriverService {
	return &DriverService{repo: repo, DriverManager: driverManager}
}

func (u *DriverService) FindAll() ([]*entity.Driver, error) {
	result, err := u.repo.Driver.FindAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *DriverService) FindAllSimple() ([]*entity.Driver, error) {
	result, err := u.repo.Driver.FindAllSimple()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *DriverService) Create(driver *entity.Driver) error {
	existingDriver, err := u.repo.Driver.FindBySlug(driver.Slug)
	if err == nil && existingDriver != nil {
		return errors.New(common.ErrDriverAlreadyExistMsg)
	}
	err = u.repo.Driver.Create(driver)
	if err != nil {
		return err
	}
	go u.updateDriverClient(driver.Slug)
	return nil
}

func (u *DriverService) Update(driver *entity.Driver) error {
	err := u.repo.Driver.Update(driver)
	if err != nil {
		return err
	}
	go u.updateDriverClient(driver.Slug)
	return nil
}

func (u *DriverService) updateDriverClient(slug string) {
	driver, err := u.repo.Driver.FindBySlug(slug)
	if err != nil {
		log.Printf("Error updating driver: %v", err)
	}
	driverClient, err := driver_lib.NewDriverClient(driver.Slug, driver_lib.StorageDriverType(driver.Type), driver.GetDriverConfig())
	if err != nil {
		log.Printf("Error updating driver: %v", err)
	}
	u.DriverManager.AddDriver(driverClient)
}

func (u *DriverService) Delete(slug string) error {
	err := u.repo.Driver.Delete(slug)
	if err != nil {
		return err
	}
	u.DriverManager.RemoveDriver(slug)
	return nil
}

func (u *DriverService) FindBySlug(slug string) (*entity.Driver, error) {
	result, err := u.repo.Driver.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *DriverService) LoadDriverManager() error {
	drivers, err := u.FindAll()
	if err != nil {
		return err
	}

	for _, driver := range drivers {
		driverConfigJSON := driver.GetDriverConfig()
		driverClient, err := driver_lib.NewDriverClient(driver.Slug, driver_lib.StorageDriverType(driver.Type), driverConfigJSON)
		if err != nil {
			return err
		}
		u.DriverManager.AddDriver(driverClient)
	}

	return nil
}
