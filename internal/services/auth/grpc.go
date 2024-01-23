package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IskanderSh/taqwa-auth/internal/config"
	"github.com/IskanderSh/taqwa-auth/internal/domain/models"
	"github.com/IskanderSh/taqwa-auth/internal/lib/error/wrapper"
	"github.com/IskanderSh/taqwa-auth/internal/lib/jwt"
	"github.com/IskanderSh/taqwa-auth/internal/lib/logger/sl"
	"github.com/IskanderSh/taqwa-auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	token       *config.Token
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, hashPass []byte) (string, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (*models.User, error)
}

func New(log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	token *config.Token,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		token:       token,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.grpc.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")

			return "", wrapper.Wrap(op, storage.ErrUserNotFound)
		}

		log.Error("failed to get user", sl.Err(err))

		return "", wrapper.Wrap(op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.HashPass, []byte(password)); err != nil {
		log.Info("invalid credentials", sl.Err(err))

		return "", wrapper.Wrap(op, err)
	}

	token, err := jwt.NewToken(user, a.token)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))

		return "", wrapper.Wrap(op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (string, error) {
	const op = "auth.grpc.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil && !errors.Is(err, storage.ErrUserNotFound) {
		log.Error("failed to get user", sl.Err(err))

		return "", wrapper.Wrap(op, err)
	}

	if user != nil {
		log.Warn("user with such email already registered")

		return "", wrapper.Wrap(op, err)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return "", wrapper.Wrap(op, err)
	}

	userID, err := a.usrSaver.SaveUser(ctx, email, hashPass)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return "", wrapper.Wrap(op, storage.ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return "", wrapper.Wrap(op, err)
	}

	log.Info("user successfully registered")
	return userID, nil
}
