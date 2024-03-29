package beanstalkd

import (
	"github.com/esportsdrafts/esportsdrafts/libs/beanstalkd/models"
)

// CreateBeanstalkdClient creates a simple client wrapper
func CreateBeanstalkdClient(address, port string) *models.Client {
	return &models.Client{
		Address: address,
		Port:    port,
	}
}
