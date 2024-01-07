package token

type Repository interface {
	CreateJWTAccessToken(userID string) (string, error)
	CreateJWTRefreshToken(userID string) (string, error)
	TokenVrification(token string) PersonDataInToken
}
