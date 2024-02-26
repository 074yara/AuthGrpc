package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/074yara/AuthGrpc/auth/internal/domain/entities"
	"github.com/074yara/AuthGrpc/auth/internal/storage"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, hashPass []byte) (uid uint, err error) {
	const op = "storage.sqlite.SaveUser"
	stmt, err := s.db.Prepare(`INSERT INTO users(email, pass_hash) VALUES (?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, email, hashPass)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return uint(id), nil
}

func (s *Storage) User(ctx context.Context, email string) (entities.User, error) {
	const op = "storage.sqlite.User"
	stmt, err := s.db.Prepare(`SELECT id, email, pass_hash FROM users WHERE email = ?`)
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, email)
	var user entities.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, storage.ErrUserNotFound
		}
		return entities.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, err
}

func (s *Storage) IsAdmin(ctx context.Context, userId uint) (bool, error) {
	const op = "storage.sqlite.IsAdmin"
	stmt, err := s.db.Prepare(`SELECT is_admin FROM users WHERE id = ?`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, int(userId))
	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, storage.ErrUserNotFound
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appId uint) (entities.App, error) {
	const op = "storage.sqlite.App"
	stmt, err := s.db.Prepare(`SELECT id, name, secret FROM apps WHERE id = ?`)
	if err != nil {
		return entities.App{}, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, int(appId))
	var app entities.App
	err = row.Scan(&app.ID, &app.Secret, &app.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.App{}, storage.ErrAppNotFound
		}
		return entities.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}