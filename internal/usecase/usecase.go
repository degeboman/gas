package usecase

import (
	"errors"
	"fmt"
	"github.com/degeboman/gas/internal/config"
	"github.com/degeboman/gas/internal/storage"
	"github.com/degeboman/gas/internal/storage/mongodb"
	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
)

//TODO разбить на отдельные файлы как у handler

type Usecase struct {
	storage.Storage
}

func (u Usecase) RefreshToken(cfg config.Config, refreshToken string) (string, error) {
	const op = "usecase.usecase.RefreshToken"

	_, err := parseToken(cfg.SigningKey, refreshToken)
	if err != nil {
		return "", err
	}

	// TODO get user info

	//accessToken, err := generateToken(userInfo, cfg.AccessDuration, cfg.SigningKey)
	//if err != nil {
	//	// TODO handling error with defer
	//	return "", fmt.Errorf("%s: %w", op, err)
	//}

	return "accessToken", nil
}

func (u Usecase) VerifyToken(signingKey []byte, token string) (interface{}, error) {
	return parseToken(signingKey, token)
}

func (u Usecase) Signin(cfg config.Config, email, password string) (access string, refresh string, err error) {
	const op = "usecase.usecase.Signin"

	userInfo, err := u.Storage.UserByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := comparePassword(userInfo.(map[string]interface{})["password"].(string), password); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, errors.New("password or email is not correct"))
	}

	accessToken, err := generateToken(userInfo, cfg.AccessDuration, cfg.SigningKey)
	if err != nil {
		// TODO handling error with defer
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := generateToken(userInfo, cfg.RefreshDuration, cfg.SigningKey)
	if err != nil {
		// TODO handling error with defer
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}

func (u Usecase) CreateUser(email, password string, userInfo interface{}) (string, error) {
	const op = "usecase.usecase.CreateUser"

	if !isValidEmail(email) {
		return "0", fmt.Errorf("%s: %w", op, errors.New("email is not valid"))
	}

	if len(password) < 6 {
		return "0", fmt.Errorf("%s: %w", op, errors.New("password is too short"))
	}

	if err := u.Storage.DoesEmailExist(email); err != nil {
		return "0", fmt.Errorf("%s: %w", op, err)
	}

	passwordHash := hashPassword(password)

	return u.Storage.CreateUser(email, passwordHash, userInfo)
}

func New(storage *mongodb.UsersStorage) Usecase {
	return Usecase{
		Storage: storage,
	}
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hashPassword(s string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hashed)
}

func comparePassword(hashed string, normal string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}

func generateToken(userInfo interface{}, duration time.Duration, signingKey []byte) (token string, err error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(duration * time.Second)),
		},
		UserInfo: userInfo,
	}

	ss := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return ss.SignedString(signingKey)
}

func parseToken(key []byte, token string) (*UserClaims, error) {
	data, err := jwt.ParseWithClaims(token, &UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := data.Claims.(*UserClaims); ok && data.Valid {
		// removing a field containing a password hash
		if _, ok := claims.UserInfo.(map[string]interface{})["password"]; ok {
			delete(claims.UserInfo.(map[string]interface{}), "password")
		}

		return claims, nil
	} else {
		return nil, err
	}
}

type UserClaims struct {
	jwt.StandardClaims
	UserInfo interface{}
}
