package refresh

import (
	"errors"
	"github.com/degeboman/gas/internal/config"
	"github.com/degeboman/gas/internal/lib/api/response"
	"github.com/degeboman/gas/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	AccessToken string `json:"access_token"`
}

type providerRefresh interface {
	RefreshToken(cfg config.Config, refreshToken string) (string, error)
}

func New(log *slog.Logger, cfg config.Config, providerRefresh providerRefresh) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.refresh.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		access, err := providerRefresh.RefreshToken(cfg, req.RefreshToken)

		if err != nil {
			log.Error("failed to sign in", sl.Err(err))

			render.JSON(w, r, response.Error("failed to sign up"+sl.Err(err).String()))

			return
		}

		responseOK(w, r, access)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, accessToken string) {
	render.JSON(w, r, Response{
		AccessToken: accessToken,
	})
}
