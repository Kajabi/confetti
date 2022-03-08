package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/codegangsta/envy/autoload"
	envy "github.com/codegangsta/envy/lib"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kajabi/developer-platform/services/app-manager/api"
	"github.com/unrolled/render"
)

var r *render.Render

func main() {
	r = render.New(render.Options{
		Layout:     "layout",
		Directory:  "assets/templates",
		Extensions: []string{".html"},
	})

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// Site routes
	router.Route("/sites", func(router chi.Router) {
		router.Use(authorizeMiddleware)

		router.Get("/{id}", Edit)
		router.Post("/{id}", Update)
	})

	// auth routes
	router.Get("/authorize", Authorize)
	router.Get("/callback", AuthCallback)

	port := envy.MustGet("PORT")
	fmt.Println("Listening on port :" + port)
	http.ListenAndServe(":"+port, router)
}

func Edit(w http.ResponseWriter, req *http.Request) {
	siteId, err := strconv.Atoi(chi.URLParam(req, "id"))
	handleErr(err)

	client, err := apiClient(req)
	handleErr(err)

	resp, err := client.ListScriptTagsWithResponse(req.Context(), int64(siteId))
	handleErr(err)

	r.HTML(w, http.StatusOK, "edit", resp.ScriptTagList.ScriptTags)
}

func Update(w http.ResponseWriter, req *http.Request) {
	siteId, err := strconv.Atoi(chi.URLParam(req, "id"))
	handleErr(err)

	client, err := apiClient(req)
	handleErr(err)

	_, err = client.AddScriptTagWithResponse(req.Context(), int64(siteId), &api.ScriptTagCreateUpdate{
		SourceUrl: "https://s3.amazonaws.com/sandbox-integrations-scripts-development.kajabi.dev/scripts/confetti.js",
	})
	handleErr(err)

	http.Redirect(w, req, req.URL.String(), http.StatusTemporaryRedirect)
}

func authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/authorize?site_id="+chi.URLParam(r, "id"), http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func apiClient(req *http.Request) (*api.ClientWithResponses, error) {
	token := req.Context().Value("token").(string)
	return api.NewClientWithResponses("https://api-development.kajabi.dev", api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}))
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
