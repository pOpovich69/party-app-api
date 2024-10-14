package http

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go-rest-api/config/container"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/controllers"
	"go-rest-api/internal/infra/http/middlewares"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CreateRouter(con container.Container) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/ping", func(apiRouter chi.Router) {
			apiRouter.Get("/", pingHandler())
			apiRouter.Handle("/*", notFoundJson())
		})
		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, con.SessionController, con.AuthMw)
				})
				apiRouter.Route("/user", func(apiRouter chi.Router) {
					apiRouter.Use(con.AuthMw)
					UserRouter(apiRouter, con.UserController)
				})
				apiRouter.Route("/", func(apiRouter chi.Router) {
					apiRouter.Use(con.AuthMw)
					PartyRouter(apiRouter, con.PartyController, con)
				})
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "file_storage"))
		requestCtx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(requestCtx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func AuthRouter(r chi.Router, sc controllers.SessionController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			sc.Register(),
		)
		apiRouter.Post(
			"/login",
			sc.Login(),
		)
		apiRouter.With(amw).Delete(
			"/logout",
			sc.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/me",
			uc.FindMe(),
		)
		apiRouter.Post(
			"/me/balance",
			uc.UpdateMyBalance(),
		)
	})
}

func PartyRouter(r chi.Router, pc controllers.PartyController, con container.Container) {
	pathObjMw := middlewares.PathObjectMiddleware(con.PartyService)
	isOwnerMw := middlewares.IsOwnerMiddleware[domain.Party]()
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/parties",
			pc.GetParties(),
		)
		apiRouter.Get(
			"/parties/creator/{creatorId}",
			pc.FindByCreatorId(),
		)
		apiRouter.Get(
			"/party/{partyId}",
			pc.FindById(),
		)
		apiRouter.Post(
			"/party",
			pc.Save(),
		)
		apiRouter.With(pathObjMw).With(isOwnerMw).Put(
			"/party/{partyId}",
			pc.Update(),
		)
		apiRouter.With(pathObjMw).With(isOwnerMw).Delete(
			"/party/{partyId}",
			pc.Delete(),
		)

	})
}

func notFoundJson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.NotFound(w, errors.New("resource Not Found"))
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.Ok(w)
	}
}
