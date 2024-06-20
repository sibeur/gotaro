package driver

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
)

type DriverManagerUseCase interface {
	GetDriver(driverName string) DriverClientUseCase
	AddDriver(driver DriverClientUseCase) error
	RemoveDriver(driverName string) error
	GetAllDriverNames() []string
}

type DriverManager struct {
	driverClients map[string]DriverClientUseCase
}

func NewDriverManager() *DriverManager {
	return &DriverManager{
		driverClients: make(map[string]DriverClientUseCase),
	}
}

func (dm *DriverManager) GetDriver(driverName string) DriverClientUseCase {
	return dm.driverClients[driverName]
}

func (dm *DriverManager) AddDriver(driver DriverClientUseCase) error {
	log.Info(fmt.Sprintf("Loading driver %v, type: %v", driver.GetName(), driver.GetTypeString()))
	dm.driverClients[driver.GetName()] = driver
	return nil
}

func (dm *DriverManager) RemoveDriver(driverName string) error {
	delete(dm.driverClients, driverName)
	return nil
}

func (dm *DriverManager) GetAllDriverNames() []string {
	var driverNames []string
	for name := range dm.driverClients {
		driverNames = append(driverNames, name)
	}
	return driverNames
}
