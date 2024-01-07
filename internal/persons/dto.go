package person

type ResponseUserDto struct {
	Id    string
	Name  string
	Email string
}

type ResponseUserAuthDto struct {
	Id            string
	Name          string
	Email         string
	Hash_Password string
}

type AuthDto struct {
	Password string
	Email    string
}
