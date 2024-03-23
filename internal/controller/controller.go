package controller

import "github.com/eurofurence/reg-room-service/internal/repository/database"

// Controller is the service interface, which defines
// the functions in the service layer of this application
//
// A type implementing this interface provides functionality
// to interact between the web layer and the data layer.
//
// Deprecated: controller package will be removed in the future.
type Controller interface {
}

type serviceController struct {
	DB database.Repository
}
