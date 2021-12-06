package validator

import (
	"context"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	type Model struct {
		RequiredField string `json:"required_field" validate:"required"`
		NotRequiredField string
	}

	validator, _ := New()

	m := &Model{
		NotRequiredField: "123",
	}

	err := validator.Validate(context.Background(), m, "ru")

	if err.Empty() {
		t.Error("Должна вернуться ошибка валидации")
	}

	m = &Model{
		RequiredField: "123",
	}

	err = validator.Validate(context.Background(), m, "ru")

	if !err.Empty() {
		t.Error("Ошибок валидации быть не должно")
	}
}
