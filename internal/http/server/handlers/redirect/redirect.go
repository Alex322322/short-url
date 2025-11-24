package redirect

import (
	"log/slog"
	"net/http"

	resp "github.com/Alex322322/short-url/internal/lib/api/response"
	"github.com/Alex322322/short-url/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(logger *slog.Logger, urlRedirector URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Error("Alias parameter is missing")
			render.JSON(w, r, resp.ErrorResponse("Alias parameter is required"))
			return
		}

		resURL, err := urlRedirector.GetURL(alias)
		if err != nil {
			if err == storage.ErrUrlNotFound {
				logger.Info("URL with that alias not found", slog.String("alias", alias))
				render.JSON(w, r, resp.ErrorResponse("URL with that alias not found"))
				return
			}
			logger.Error("Failed to get URL", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ErrorResponse("internal server error"))
			return
		}

		logger.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
