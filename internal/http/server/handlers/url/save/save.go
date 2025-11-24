package save

import (
	"log/slog"
	"net/http"
	"time"

	resp "github.com/Alex322322/short-url/internal/lib/api/response"
	"github.com/Alex322322/short-url/internal/lib/random"
	"github.com/Alex322322/short-url/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"` // validate tag for validator
	Alias string `json:"alias,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 8

//go:generate go run github.com/vektra/mockery/v3@latest --dir . --name URLSaver --output ./mocks
type URLSaver interface {
	SaveURL(urlToSave string, alias string, timestamp time.Time) (int64, error)
}

func New(logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var request Request
		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			logger.Error("Failed to decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ErrorResponse("Invalid request body"))
			return
		}

		logger.Info("Request body decoded", slog.Any("request", request))

		// validate request
		if err := validator.New().Struct(request); err != nil {
			// error type assertion
			validateErr := err.(validator.ValidationErrors)
			// log validation errors
			logger.Error("Request validation failed", slog.String("error", err.Error()))
			// respond with validation errors
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := request.Alias
		if alias == "" {
			// generate alias if not provided
			alias = random.GenerateAlias(aliasLength)
		}

		id, err := urlSaver.SaveURL(request.URL, alias, request.Timestamp)
		if err != nil {
			if err == storage.ErrUrlExists {
				logger.Info("URL with the same alias already exists", slog.String("url", request.URL), slog.String("alias", alias))
				render.JSON(w, r, resp.ErrorResponse("URL with the same alias already exists"))
				return
			}
			logger.Error("Failed to save URL", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ErrorResponse("Failed to save URL"))
			return
		}

		logger.Info("URL saved successfully", slog.Int64("id", id), slog.String("url", request.URL), slog.String("alias", alias))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OKResponse(),
		Alias:    alias,
	})
}
