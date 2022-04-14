package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	envy "github.com/codegangsta/envy/lib"
	"github.com/kajabi/confetti/go/auth"
)

var authClient = auth.NewClient(&auth.Options{
	ClientID:     envy.MustGet("CLIENT_ID"),
	ClientSecret: envy.MustGet("CLIENT_SECRET"),
	RedirectURL:  envy.MustGet("REDIRECT_URL"),
	Scopes:       strings.Split(envy.MustGet("SCOPES"), ","),
	PKCE:         false,
	AuthDomain:   envy.MustGet("AUTH_DOMAIN"),
	Audience:     envy.MustGet("AUTH_AUDIENCE"),
})

func Authorize(w http.ResponseWriter, r *http.Request) {
	raw := authClient.URL("test") // TODO generate state better

	url, err := url.Parse(raw)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	query := url.Query()
	query.Add("site_id", r.URL.Query().Get("site_id"))

	redirect := query.Get("redirect_uri") + fmt.Sprintf("?site_id=%s", r.URL.Query().Get("site_id"))
	query.Set("redirect_uri", redirect)

	url.RawQuery = query.Encode()

	http.Redirect(w, r, url.String(), http.StatusFound)
}

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	token, err := authClient.GetAccessToken("test", r) // TODO verify state better
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token.AccessToken,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, fmt.Sprintf("/sites/%s", r.URL.Query().Get("site_id")), http.StatusFound)
}
