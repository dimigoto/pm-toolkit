package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func ValidationErrors(w http.ResponseWriter, errs interface{}) {
	errorsResponse := map[string]interface{}{
		"errors": errs,
	}

	write(
		w,
		http.StatusUnprocessableEntity,
		errorsResponse,
	)
}

// Error Записывает ошибку с кодом в структуру http.ResponseWriter
func Error(w http.ResponseWriter, code int, err error) {
	write(
		w, code, map[string]string{
			"error": err.Error(),
		},
	)
}

// Success Записывает успех в структуру http.ResponseWriter
func Success(w http.ResponseWriter, data interface{}) {
	code := http.StatusOK

	if data == nil {
		code = http.StatusNoContent
	}

	write(w, code, data)
}

// write Форматирует результат и записывает в структуру http.ResponseWriter
func write(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
		}
	}
}
