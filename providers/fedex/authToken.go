package fedex

import (
	"fmt"
	"net/http"
)

type authToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (me authToken) setRequestAuthorization(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("%s %s", me.TokenType, me.AccessToken))
}
