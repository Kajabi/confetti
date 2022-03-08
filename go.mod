module github.com/kajabi/confetti

go 1.16

require (
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4
	github.com/go-chi/chi/v5 v5.0.7
	github.com/kajabi/developer-platform/services/app-manager v0.0.0-00010101000000-000000000000
	github.com/mholt/binding v0.3.0 // indirect
	github.com/nirasan/go-oauth-pkce-code-verifier v0.0.0-20170819232839-0fbfe93532da
	github.com/unrolled/render v1.2.0
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b
)

replace github.com/kajabi/developer-platform/services/app-manager => ../developer-platform/services/app-manager
