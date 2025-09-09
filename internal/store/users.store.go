package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}
	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Match(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users ( name, email, password_hash, role)
	VALUES($1,$2,$3,'user')
	RETURNING id, role, created_at, updated_at
	`
	err := s.db.QueryRow(query, user.Name, user.Email, user.PasswordHash.hash).Scan(&user.ID, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `
	SELECT id, name, email, password_hash, role, created_at, updated_at
	FROM users
	WHERE email = $1
	`
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) GetUserByID(id string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `
	SELECT id, name, email, password_hash, role, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET name = $1, email= $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $4
	RETURNING updated_at
	`
	result, err := s.db.Exec(query, user.Name, user.Email)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresUserStore) GetUserToken(scope, plainTextPassword string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plainTextPassword))
	query := `
	SELECT u.id, u.name, u.email, u.password_hash, u.role, u.created_at, u.updated_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.hash = $1 AND t.scope = $2 and t.expiry > $3
	`

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}
