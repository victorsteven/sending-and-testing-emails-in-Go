package helpers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type WelcomeMail struct {
	FromAdmin string
	Temp      string
	Name      string
	Email     string
	Emails    []string
	Subject   string
}

type EmailResponse struct {
	Status   int
	RespBody string
}

type WelcomeModel struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func ValidateInputs(dataSet interface{}) error {

	var validate *validator.Validate

	validate = validator.New()

	err := validate.Struct(dataSet)

	if err != nil {
		//Validation syntax is invalid
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("validation syntax is invalid")
		}

		reflected := reflect.ValueOf(dataSet)

		for _, err := range err.(validator.ValidationErrors) {

			// Attempt to find field by name and get json tag name
			field, _ := reflected.Type().FieldByName(err.StructField())

			//If json tag doesn't exist, use lower case of name
			name := field.Tag.Get("json")
			if name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				return errors.New(fmt.Sprintf("The %s is required", name))
			case "email":
				return errors.New(fmt.Sprintf("The %s should be a valid email", name))
			default:
				return errors.New(fmt.Sprintf("The %s is invalid", name))
			}
		}
	}

	return nil
}
