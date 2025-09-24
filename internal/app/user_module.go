package app

import (
	"user-management-api/internal/handler"
	"user-management-api/internal/repository"
	"user-management-api/internal/routes"
	v1routes "user-management-api/internal/routes/v1"
	v1service "user-management-api/internal/service/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewInMemoryUserRepository()
	userService := v1service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := v1routes.NewUserRoutes(userHandler)
	return &UserModule{routes: userRoutes}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
