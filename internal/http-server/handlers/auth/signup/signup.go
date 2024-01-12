package signup

import (
	"errors"
	"github.com/degeboman/gas/internal/lib/api/response"
	"github.com/degeboman/gas/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	AccessToken string `json:"access_token"`
}

// TODO maybe add userinfo interface
// TODO check user info interface

type Request struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	UserInfo interface{} `json:"user_info"`
}

type UserCreator interface {
	CreateUser(email, password string, userInfo interface{}) (string, error)
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.signup.New"

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

		_, err = userCreator.CreateUser(req.Email, req.Password, req.UserInfo)
		if err != nil {
			log.Error("failed to sign up", sl.Err(err))

			render.JSON(w, r, response.Error("failed to sign up"+sl.Err(err).String()))

			return
		}

		responseOK(r)
	}
}

func responseOK(r *http.Request) {
	render.Status(r, http.StatusOK)
}
