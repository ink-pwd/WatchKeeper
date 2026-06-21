package utils

import (
	"context"
	"net/http"
	"net/url"
)

func GetHostName(req string) (string, error) {
	var (
		u   *url.URL
		err error
	)
	u, err = url.ParseRequestURI(req)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func IsServerAlive(ctx context.Context, addr string) bool {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	/*
		Создаем новый GET запрос к указанному сайту используя контекст.
		Тела запроса нет, так как запрос GET.
	*/
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return false
	}
	/*
		Отправляем запрос
	*/
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	/*
		Если хоть какой-то ответ был, значит сервер жив
	*/
	return true
}
