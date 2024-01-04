package roles

type AssignRoleDto struct {
	UserId string `json:"userId"`
	RoleId string `json:"roleId"`
}

type TakeRoleDto struct {
	UserId string `json:"userId"`
	RoleId string `json:"roleId"`
}

// AllUserRolesDto структура для представления общего результата
type AllUserRolesDto struct {
	UserId string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Roles  []Role `json:"roles"`
}
