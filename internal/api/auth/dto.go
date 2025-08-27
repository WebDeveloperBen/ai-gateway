package auth

import "net/http"

type LoginRedirect struct {
	Location string `header:"Location"`
}

type CallbackRequest struct {
	Code  string `query:"code" doc:"OAuth code"`
	Error string `query:"error" doc:"OAuth error"`
	State string `query:"state" doc:"CSRF state"`
}

type CallbackRedirect struct {
	Location  string      `header:"Location"`
	SetCookie http.Cookie `header:"Set-Cookie"`
}

type MeResponseBody struct {
	Email string `json:"email"`
	Sub   string `json:"sub"`

	Name              string   `json:"name,omitempty"`
	GivenName         string   `json:"given_name,omitempty"`
	FamilyName        string   `json:"family_name,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	Roles             []string `json:"roles,omitempty"`
	Groups            []string `json:"groups,omitempty"`
}

type MeResponse struct {
	Body MeResponseBody `json:"body"`
}
