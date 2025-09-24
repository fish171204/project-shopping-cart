package app

import (
	"user-management-api/internal/handler"
	"user-management-api/internal/repository"
	"user-management-api/internal/routes"
	"user-management-api/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	//Initialize repository
	userRepo := repository.NewInMemoryUserRepository()

	//Initialize service
	userService := service.NewUserService(userRepo)

	//Initialize handler
	userHandler := handler.NewUserHandler(userService)

	//Initialize routes
	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{routes: userRoutes}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
