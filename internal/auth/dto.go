package auth

type RegisterDto struct {
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type AuthDto struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type AuthResponseDto struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type RefreshDto struct {
	Refresh string `json:"Refresh"`
}

type RefreshResponseDto struct {
	Access string
}
