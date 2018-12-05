package redirect

import (
	"net/http"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
)

// SecureRedirectMiddleware redirects the client to the identical URL served via HTTPS
type SecureRedirectMiddleware struct {
	WhiteListedPaths map[string]struct{}
}

func NewSecureRedirectMiddleware(paths ...string) SecureRedirectMiddleware {
	uniquePaths := make(map[string]struct{})
	for _, p := range paths {
		uniquePaths[p] = struct{}{}
	}
	return SecureRedirectMiddleware{
		WhiteListedPaths: uniquePaths,
	}
}

func (srm SecureRedirectMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		_, whiteListed := srm.WhiteListedPaths[r.URL.Path]
		if strings.ToLower(r.Header.Get("X-Forwarded-Proto")) == "http" && !whiteListed {
			redirectURL := r.URL
			redirectURL.Host = r.Host
			redirectURL.Scheme = "https"
			http.Redirect(w.(http.ResponseWriter), r.Request, redirectURL.String(), http.StatusMovedPermanently)
			return
		}
		h(w, r)
	}
}
