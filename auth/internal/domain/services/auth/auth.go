package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/074yara/AuthGrpc/auth/internal/domain/entities"
	"github.com/074yara/AuthGrpc/auth/internal/lib/jwt"
	"github.com/074yara/AuthGrpc/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	// ErrInvalidCredentials is error for invalid credentials
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidID          = errors.New("invalid user id")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, hashPass []byte) (uid uint, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (entities.User, error)
	IsAdmin(ctx context.Context, userId uint) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId uint) (entities.App, error)
}

// New returns new Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user exists and password is correct. If it's so Login will return token string
// returns error if user doesn't exist, returns error if user exist, but password is incorrect
func (a *Auth) Login(ctx context.Context, email, password string, appId uint) (token string, err error) {
	const op = "Auth.Login"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Warn("incorrect password")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to compare password", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		log.Error("failed to get app", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged in successfully")
	tokenString, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to create token", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("token created")
	return tokenString, nil

}

// RegisterNewUser registers new user
func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (uint, error) {
	const op = "Auth.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := a.userSaver.SaveUser(ctx, email, hashPass)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed to save user", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user successfully registered", slog.Uint64("id", uint64(id)))
	return id, err
}

// IsAdmin checks if user is admin
func (a *Auth) IsAdmin(ctx context.Context, id uint) (bool, error) {
	const op = "Auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.Uint64("id", uint64(id)),
	)

	log.Info("checking if user is admin")
	isAdmin, err := a.userProvider.IsAdmin(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.Uint64("id", uint64(id)))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidID)
		}
		log.Warn("failed to check if user is admin", err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}
