package v1dto

import "user-management-api/internal/models"

// Data Transfer Object (DTO): chuyển dữ liệu ra response
type UserDTO struct {
	UUID   string `json:"uuid"`
	Name   string `json:"full_name"`
	Email  string `json:"email_address"`
	Age    int    `json:"age"`
	Status string `json:"status"`
	Level  string `json:"level"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,gt=18"`
	Password string `json:"password" binding:"required,min=8,password_strong"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}

type UpdateUserInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Age      int    `json:"age" binding:"omitempty,gt=18"`
	Password string `json:"password" binding:"omitempty,min=8,password_strong"`
	Status   int    `json:"status" binding:"omitempty,oneof=1 2"`
	Level    int    `json:"level" binding:"omitempty,oneof=1 2"`
}

// Request
func (input *CreateUserInput) MapCreateInputToModel() {

}

func (input *UpdateUserInput) MapUpdateInputToModel() {

}

// Response
func MapUserToDTO(user models.User) *UserDTO {
	return &UserDTO{}
}

func MapUsersToDTO(users []models.User) []UserDTO {
	dtos := make([]UserDTO, 0, len(users))

	for _, user := range users {
		dto := MapUserToDTO(user)
		dtos = append(dtos, *dto)
	}

	return dtos
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Show"
	case 2:
		return "Hide"
	default:
		return "None"
	}
}

func mapLevelText(status int) string {
	switch status {
	case 1:
		return "Admin"
	case 2:
		return "Member"
	default:
		return "None"
	}
}
