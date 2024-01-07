package token

type Repository interface {
	CreateJWTAccessToken(person PersonDataInToken) (string, error)
	CreateJWTRefreshToken(person PersonDataInToken) (string, error)
	TokenVrification(tokenString string, key string) (PersonDataInToken, error)
}
