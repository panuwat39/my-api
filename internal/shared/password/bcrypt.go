package password

import (
	"fmt"

	"github.com/panuwat39/my-api/internal/user/port"
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	cost int
}

func NewHasher() *Hasher {
	return &Hasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *Hasher) Hash(
	plainPassword string,
) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		h.cost,
	)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(hash), nil
}

func (h *Hasher) Compare(
	passwordHash string,
	plainPassword string,
) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(plainPassword),
	); err != nil {
		return fmt.Errorf("compare password: %w", err)
	}

	return nil
}

var _ port.PasswordHasher = (*Hasher)(nil)
