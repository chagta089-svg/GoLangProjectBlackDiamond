package models

const roleAdmin string = "admin"

func GetRoleAdmin() string {
	return roleAdmin
}
func IsAdmin(role string) bool {
	return role == roleAdmin
}
