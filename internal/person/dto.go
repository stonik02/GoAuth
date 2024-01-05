package person

type ResponseUserDto struct {
	Id    string
	Name  string
	Email string
}

type AuthDto struct {
	Password string
	Email    string
}
