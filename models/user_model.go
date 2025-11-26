package models

import (
	"webapp/config"
)

type User struct {
	Username string
	Fullname string
	Email    string
	Password string
	Role     string
}

func GetUserByLogin(usernameOrEmail, password string) (*User, error) {
	db := config.GetDB()
	var user User
	err := db.QueryRow(
		"SELECT username, fullname, email, role FROM user WHERE (username = ? OR email = ?) AND password = ?",
		usernameOrEmail, usernameOrEmail, password,
	).Scan(&user.Username, &user.Fullname, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *User) error {
	db := config.GetDB()
	_, err := db.Exec(
		"INSERT INTO user (username, fullname, email, password, role) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Fullname, user.Email, user.Password, user.Role,
	)
	return err
}

func IsUsernameOrEmailExists(username, email string) (bool, error) {
	db := config.GetDB()
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM user WHERE username = ? OR email = ?",
		username, email,
	).Scan(&count)
	return count > 0, err
}

func GetUserRole(username string) (string, error) {
	db := config.GetDB()
	var role string
	err := db.QueryRow("SELECT role FROM user WHERE username = ?", username).Scan(&role)
	return role, err
}

func GetAllUsers() ([]User, error) {
	db := config.GetDB()
	rows, err := db.Query("SELECT username, fullname, email, role FROM user ORDER BY username")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Fullname, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func UpdateUserRole(username, role string) error {
	db := config.GetDB()
	_, err := db.Exec("UPDATE user SET role = ? WHERE username = ?", role, username)
	return err
}
