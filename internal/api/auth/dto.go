package auth

import "net/http"

type LoginRedirect struct {
	Location string `header:"Location"`
}

type CallbackRequest struct {
	Code  string `query:"code"`
	Error string `query:"error"`
}

type CallbackRedirect struct {
	Location   string         `header:"Location"`
	SetCookies []*http.Cookie `cookie:"Set-Cookie"`
}

type MeResponse struct {
	User map[string]any `json:"user"`
}
