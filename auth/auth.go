// package auth is an abstraction of our authentication and authorization layer
// for Kajabi apps. If you wish to create a third-party, first party or internal app that utilizes Kajabi id, this is the package to use.
package auth

import (
	"errors"
	"net/http"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"golang.org/x/oauth2"
)

type Client struct {
	config   *oauth2.Config
	verifier *cv.CodeVerifier
	audience string
}

type Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	PKCE         bool
	AuthDomain   string
	Audience     string
}

func NewClient(options *Options) *Client {
	client := &Client{}

	config := &oauth2.Config{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		RedirectURL:  options.RedirectURL,
		Scopes:       options.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  options.AuthDomain + "/authorize",
			TokenURL: options.AuthDomain + "/oauth/token",
		},
	}
	client.config = config
	client.audience = options.Audience

	if options.PKCE {
		// Ignoring this error as this call will use the default length and not
		// surface any errors
		client.verifier, _ = cv.CreateCodeVerifier()
	}

	return client
}

func (c *Client) URL(state string) string {
	params := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("audience", c.audience),
	}

	// PKCE parameters
	if c.verifier != nil {
		params = append(params,
			oauth2.SetAuthURLParam("code_challenge", c.verifier.CodeChallengeS256()),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		)
	}

	return c.config.AuthCodeURL(state, params...)
}

func (c *Client) GetAccessToken(state string, r *http.Request) (*oauth2.Token, error) {
	// Verify state
	if r.URL.Query().Get("state") != state {
		return nil, errors.New("Invalid state parameter")
	}

	params := []oauth2.AuthCodeOption{}

	// PKCE parameters
	if c.verifier != nil {
		params = append(params, oauth2.SetAuthURLParam("code_verifier", c.verifier.String()))
	}

	return c.config.Exchange(r.Context(), r.URL.Query().Get("code"), params...)
}
