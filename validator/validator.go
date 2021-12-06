package validator

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslator "github.com/go-playground/validator/v10/translations/en"
	rutranslator "github.com/go-playground/validator/v10/translations/ru"
)

const (
	LocaleRu = "ru"
	LocaleEn = "en"
)

// ValidationErrors Структура ошибок валидации: ключ - валидируемое поле, значение – текст ошибки
type ValidationErrors map[string] interface{}

// Empty В случае если ошибок нет вернёт true
func (e *ValidationErrors) Empty() bool {
	return len(*e) == 0
}

// Validator Обертка над go-playground/validator для более удобной работы с валидацией запросов. Достаточно вызвать
// метод Validate с валидируемой структурой и локалью
type Validator struct {
	validate *validator.Validate
	uni *ut.UniversalTranslator
}

func New() (*Validator, error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(jsonTagValueAsErrorKey)

	uni, err := registerTranslators(validate, LocaleRu, LocaleEn); if err != nil {
		return nil, err
	}

	return &Validator{
		validate: validate,
		uni: uni,
	}, nil
}

// Validate Валидирует структуру s и возвращает ValidationErrors, где ключ - название поля, значение - текст ошибки,
// локализованный согласно loc
func (v *Validator) Validate(ctx context.Context, s interface{}, loc string) *ValidationErrors {
	result := ValidationErrors{}
	translator, _ := v.uni.GetTranslator(loc)
	err := v.validate.StructCtx(ctx, s)

	if err == nil {
		return &result
	}

	for _, err := range err.(validator.ValidationErrors) {
		result[err.Field()] = err.Translate(translator)
	}

	return &result
}

// registerTranslators Регистрирует трансляторы для локалей locs
func registerTranslators(v *validator.Validate, locs ...string) (*ut.UniversalTranslator, error) {
	var translators []locales.Translator

	for _, loc := range locs {
		if loc == "ru" {
			translators = append(translators, ru.New())

			continue
		}

		if loc == "en" {
			translators = append(translators, en.New())
		}
	}

	uni := ut.New(en.New(), translators...)

	for _, loc := range locs {
		if loc == "ru" {
			trans, found := uni.GetTranslator("ru"); if found {
				err := rutranslator.RegisterDefaultTranslations(v, trans); if err != nil {
					return nil, err
				}
			}

			continue
		}

		if loc == "en" {
			trans, found := uni.GetTranslator("en"); if found {
				err := entranslator.RegisterDefaultTranslations(v, trans); if err != nil {
					return nil, err
				}
			}

			continue
		}
	}

	return uni, nil
}

// jsonTagValueAsErrorKey Регистрирует функцию, которая вытаскивает названия полей структуры из тегов json
func jsonTagValueAsErrorKey(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}