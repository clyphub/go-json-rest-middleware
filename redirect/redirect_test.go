package redirect

import (
	"net/http"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
)

func TestSecureRedirectMiddleware(t *testing.T) {
	api := rest.NewApi()
	api.Use(&SecureRedirectMiddleware{})
	api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
		w.WriteJson(map[string]string{"Id": "123"})
	}))
	handler := api.MakeHandler()

	req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	req.Header.Set("X-Forwarded-Proto", "HtTp")
	recorded := test.RunRequest(t, handler, req)
	recorded.CodeIs(http.StatusMovedPermanently)
	recorded.HeaderIs("Location", "https://localhost/")

	req = test.MakeSimpleRequest("GET", "http://localhost/", nil)
	recorded = test.RunRequest(t, handler, req)
	recorded.CodeIs(http.StatusOK)
	recorded.BodyIs(`{"Id":"123"}`)
}
