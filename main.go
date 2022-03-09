package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/codegangsta/envy/autoload"
	envy "github.com/codegangsta/envy/lib"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kajabi/developer-platform/services/app-manager/api"
	"github.com/unrolled/render"
)

var r *render.Render

type ConfettiForm struct {
	Enabled bool
	Max     string
	Size    string
}

func (f *ConfettiForm) ScriptURL() (*url.URL, error) {
	url, err := url.Parse("https://s3.amazonaws.com/sandbox-integrations-scripts-development.kajabi.dev/scripts/confetti.js")
	if err != nil {
		return url, err
	}

	q := url.Query()
	q.Add("enabled", strconv.FormatBool(f.Enabled))
	q.Add("max", f.Max)
	q.Add("size", f.Size)
	url.RawQuery = q.Encode()

	return url, nil
}

func main() {
	r = render.New(render.Options{
		Layout:     "layout",
		Directory:  "assets/templates",
		Extensions: []string{".html"},
	})

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

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

	src, err := findScriptTag(int64(siteId), req)
	handleErr(err)

	url, err := url.Parse(src)
	handleErr(err)

	// Map the url to a form
	form := &ConfettiForm{
		Enabled: url.Query().Get("enabled") == "true",
		Max:     queryWithDefault(url, "max", "80"),
		Size:    queryWithDefault(url, "size", "1"),
	}

	r.HTML(w, http.StatusOK, "edit", form)
}

func queryWithDefault(url *url.URL, key string, def string) string {
	if len(url.Query().Get(key)) > 0 {
		return url.Query().Get(key)
	} else {
		return def
	}
}

func Update(w http.ResponseWriter, req *http.Request) {
	siteId, err := strconv.Atoi(chi.URLParam(req, "id"))
	handleErr(err)

	_, err = apiClient(req)
	handleErr(err)

	form := &ConfettiForm{
		Enabled: req.PostFormValue("enabled") == "on",
		Max:     req.PostFormValue("max"),
		Size:    req.PostFormValue("size"),
	}

	url, err := form.ScriptURL()
	handleErr(err)

	err = createOrUpdateScriptTag(int64(siteId), url.String(), req)
	handleErr(err)

	http.Redirect(w, req, req.URL.Path, http.StatusFound)
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
	return api.NewClientWithResponses(envy.MustGet("API_DOMAIN"), api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}))
}

func findScriptTag(siteId int64, r *http.Request) (string, error) {
	client, err := apiClient(r)
	if err != nil {
		return "", err
	}

	resp, err := client.ListScriptTagsWithResponse(r.Context(), siteId)
	if err != nil {
		return "", err
	}

	if len(resp.ScriptTagList.ScriptTags) > 0 {
		return resp.ScriptTagList.ScriptTags[0].SourceUrl, nil
	}

	return "", nil
}

func createOrUpdateScriptTag(siteId int64, source string, r *http.Request) error {
	client, err := apiClient(r)
	if err != nil {
		return err
	}

	resp, err := client.ListScriptTagsWithResponse(r.Context(), siteId)
	if err != nil {
		return err
	}

	if len(resp.ScriptTagList.ScriptTags) > 0 {
		_, err := client.UpdateScriptTag(r.Context(), siteId, resp.ScriptTagList.ScriptTags[0].Id, &api.ScriptTagCreateUpdate{
			SourceUrl: source,
		})

		return err
	}

	_, err = client.AddScriptTag(r.Context(), siteId, &api.ScriptTagCreateUpdate{
		SourceUrl: source,
	})

	return err
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
