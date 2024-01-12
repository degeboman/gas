package verify

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
	Token string `json:"token"`
}

type ProviderVerify interface {
	VerifyToken(signingKey []byte, token string) (interface{}, error)
}

func New(log *slog.Logger, cfg config.Config, providerVerify ProviderVerify) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.verify.New"

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

		_, err = providerVerify.VerifyToken(cfg.SigningKey, req.Token)
		if err != nil {
			log.Error("failed to verify token", sl.Err(err))

			render.JSON(w, r, response.Error("failed to sign up"+sl.Err(err).String()))

			return
		}

		render.Status(r, http.StatusOK)
	}
}
