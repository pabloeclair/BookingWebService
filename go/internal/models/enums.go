package models

//todo: enum by

type Role int16

const (
	Role_ADMIN      Role = 0
	Role_USER       Role = 1
	Role_MAIN_ADMIN Role = 2
)

var (
	Role_name = map[Role]string{
		0: "ADMIN",
		1: "USER",
		2: "MAIN_ADMIN",
	}
	Role_array = []string{
		Role_USER.String(),
		Role_ADMIN.String(),
		Role_MAIN_ADMIN.String(),
	}
)

func (r Role) String() string {
	if name, ok := Role_name[r]; ok {
		return name
	}
	return "UNKNOWN"
}
