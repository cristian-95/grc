package postgres

import (
	"fmt"

	"github.com/cristian-95/grc"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	//embed type:
	*sqlx.DB
}

// User consulta o banco de dados e retorna a user correspondente ao parametro id.
func (s *UserStore) User(id uuid.UUID) (grc.User, error) {
	var u grc.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return grc.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) Users() ([]grc.User, error) {
	var uu []grc.User
	if err := s.Select(&uu, `SELECT * FROM users`); err != nil {
		return []grc.User{}, fmt.Errorf("error getting users: %w", err)
	}
	return uu, nil
}

func (s *UserStore) CreateUser(u *grc.User) error {
	if err := s.Get(u, `INSERT INTO users VALUES ($1, $2, $3) RETURNING *`,
		u.ID,
		u.Username,
		u.Password); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (s *UserStore) UpdateUser(u *grc.User) error {
	if err := s.Get(u, `UPDATE users SET username = $1, password = $2 WHERE id = $3 RETURNING *`,
		u.Username,
		u.Password,
		u.ID); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) UsersByUsername(username string) (grc.User, error) {
	var u grc.User
	if err := s.Get(&u, `SELECT * FROM users WHERE username = $1`, username); err != nil {
		return grc.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
