package port

type PasswordHasher interface {
	Hash(plainPassword string) (string, error)

	Compare(
		passwordHash string,
		plainPassword string,
	) error
}
