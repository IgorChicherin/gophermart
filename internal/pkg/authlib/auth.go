package authlib

type AuthService interface {
	TokenValidator
	Hasher
}

type TokenValidator interface {
	DecodeToken(token string) (string, string, error)
	EncodeToken(login, pwdHash string) string
	AuthValidator
}

type AuthValidator interface {
	Validate(hash, pwd string) error
}

type Hasher interface {
	GetHash(data string) string
	Equals(hash, pwd string) bool
}
