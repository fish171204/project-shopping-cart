package validation

import (
	"fmt"
	"strings"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	RegisterCustomValidation(v)

	return nil
}

func HandleValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationError {
			root := strings.Split(e.Namespace(), ".")[0]

			rawPath := strings.TrimPrefix(e.Namespace(), root+".")

			parts := strings.Split(rawPath, ".")

			for i, part := range parts {
				if strings.Contains("part", "[") {
					idx := strings.Index(part, "[")
					base := utils.CamelToSnake(part[:idx]) // => 0 đến trước [
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = utils.CamelToSnake(part)
				}
			}

			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[e.Field()] = fmt.Sprintf("%s phải lớn hơn : %s", e.Field(), e.Param())
			case "lt":
				errors[e.Field()] = fmt.Sprintf("%s phải nhỏ hơn : %s", e.Field(), e.Param())
			case "gte":
				errors[e.Field()] = fmt.Sprintf("%s phải lớn hơn hoặc bằng giá trị tối thiểu là: %s", e.Field(), e.Param())
			case "lte":
				errors[e.Field()] = fmt.Sprintf("%s phải lớn hơn hoặc bằng giá trị tối đa là: %s", e.Field(), e.Param())
			case "min":
				errors[e.Field()] = fmt.Sprintf("%s phải từ %s ký tự", e.Field(), e.Param())
			case "max":
				errors[e.Field()] = fmt.Sprintf("%s phải ít hơn %s ký tự", e.Field(), e.Param())
			case "min_int":
				errors[e.Field()] = fmt.Sprintf("%s phải có giá trị lớn hơn hoặc bằng %s", e.Field(), e.Param())
			case "max_int":
				errors[e.Field()] = fmt.Sprintf("%s phải có giá trị bé hơn hoặc bằng %s", e.Field(), e.Param())
			// users
			case "uuid":
				errors[e.Field()] = e.Field() + " phải là UUID hợp lệ"
			// products
			case "slug":
				errors[e.Field()] = e.Field() + " chỉ được chứa chữ thường, số, dấu gạch ngang hoặc dấu chấm"
			case "required": // case slice struct
				errors[fieldPath] = fieldPath + " là bắt buộc"
			case "search":
				errors[e.Field()] = e.Field() + " chỉ được chứa chữ thường, in hoa ,số và khoảng trắng"
			case "email":
				errors[e.Field()] = e.Field() + " phải đúng định dạng email"
			case "datetime":
				errors[e.Field()] = e.Field() + " phải theo đúng định dạng YYYY-MM-DD"
			case "password_strong":
				errors[e.Field()] = e.Field() + " phải ít nhất 8 ký tự bao gồm (chữ thường, chữ in hoa, số và ký tự đặc biệt)"
			case "file_ext":
				allowedValue := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[e.Field()] = fmt.Sprintf("%s chỉ cho phép những file có extension: %s", e.Field(), allowedValue)
			// category
			case "oneof":
				allowedValue := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[e.Field()] = fmt.Sprintf("%s phải là một trong các giá trị: %s", e.Field(), allowedValue)
			}

		}
		return gin.H{"error": errors}

	}

	return gin.H{"error": "Yêu cầu không hợp lệ " + err.Error()}
}
