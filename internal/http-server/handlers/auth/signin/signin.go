package signin

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
	"time"
)

type Response struct {
	AccessToken string `json:"access_token"`
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginProvider interface {
	Signin(cgf config.Config, email, password string) (access string, refresh string, err error)
}

func New(log *slog.Logger, cfg config.Config, loginProvider LoginProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.signin.New"

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

		access, refresh, err := loginProvider.Signin(cfg, req.Email, req.Password)

		http.SetCookie(w, &http.Cookie{
			Name:    "refresh_token",
			Value:   refresh,
			Expires: time.Now().Add(cfg.RefreshDuration),
		})

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
