package taoweb

import (
	"net/http"

	"github.com/anyongjitiger/photo-backup-server/web/action"
	"github.com/anyongjitiger/photo-backup-server/web/album"
	"github.com/anyongjitiger/photo-backup-server/web/auth"
	"github.com/anyongjitiger/photo-backup-server/web/common"
	"github.com/anyongjitiger/photo-backup-server/web/upload"
	mux "github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

type Route struct {
	Method string
	Path   string
	Handle mux.Handle // httprouter package as mux
}

var ()
type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/testHtml",
		TestHtml,
	},
	Route{
		"GET",
		"/",
		Index,
	},
	Route{
		"POST",
		"/login",
		auth.LoginHandler,
	},
	Route{
		"POST",
		"/resource",
		/* wrapHandler(negroni.New(
			negroni.HandlerFunc(auth.ValidateTokenMiddleware),
			negroni.Wrap(http.HandlerFunc(auth.ProtectedHandler),
		))), */
		wrapHandlers(auth.ProtectedHandler, auth.ValidateTokenMiddleware),
	},
	Route{
		"POST",
		"/resource2/:name",
		wrapHandlers(auth.ProtectedHandler2, auth.ValidateTokenMiddleware),
	},
	Route{
		"GET",
		"/posts",
		PostIndex,
	},
	Route{
		Method: "GET",
		Path:   "/test",
		Handle: action.TestTaodb,
	},
	Route{
		Method: "GET",
		Path:   "/albums/:prePath",
		// Handle: album.List,
		Handle: wrapHandlers(album.List, auth.ValidateTokenMiddleware),
	},
	Route{
		Method: "GET",
		Path:   "/show/:filePath/:fileName",
		Handle: album.Show,
	},
	Route{
		Method: "POST",
		Path:   "/upload/",
		Handle: upload.Controller{}.Upload,
	},
	Route{
		Method: "POST",
		Path:   "/uploadCheck/",
		Handle: wrapHandlers(upload.CheckUploaded, auth.ValidateTokenMiddleware),
	},
}

func wrapHandler(h http.Handler) mux.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps mux.Params) {
	//   context.Set(r, "params", ps)
	  h.ServeHTTP(w, r)
	}
}

//h: handler function; hfs: middleware function slice
func wrapHandlers(h http.HandlerFunc, hfs ...negroni.HandlerFunc) mux.Handle {
	handler := negroni.New();
	for _, hf := range hfs {
		handler.Use(hf)
	}
	handler.UseHandlerFunc(h)
	return func(w http.ResponseWriter, r *http.Request, ps mux.Params) {
		// ctx := context.WithValue(r.Context(), contextKeyParams, ps)
		ctx := common.WithParams(r.Context(), ps)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

// NewRouter  make router
func NewRouter() *mux.Router {
	router := mux.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handle)
	}
	return router
}