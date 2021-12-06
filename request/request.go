package request

import (
	"errors"
	"net/http"
)

const (
	// userIdHeaderKey Ключ заголовка запроса, в котором хранится идентификатор пользователя
	userIdHeaderKey = "X-User-Id"
	// userRolesHeaderKey Ключ заголовка запроса, в котором хранятся роли пользователя
	userRolesHeaderKey = "X-User-Roles"
)

// ErrorEmptyUserID Ошибка, которая возвращается в случае, если в заголовке нет идентификатора пользователя
var ErrorEmptyUserID = errors.New("заголовок с идентификатором пользователя пуст")

// ParseUserID Возвращает идентификатор пользователя из заголовка
func ParseUserID(r *http.Request) (string, error) {
	userID := r.Header.Get(userIdHeaderKey)
	if userID == "" {
		return "", ErrorEmptyUserID
	}

	return userID, nil
}