package controllers

import (
	"net/http"
	"text/template"
	"webapp/models"
)

// Authentication = การพิสูจน์ตัวตน
func Login(w http.ResponseWriter, req *http.Request) {
	errMsg := ""
	if req.URL.Query().Get("error") == "1" {
		errMsg = "❌ Username หรือ Password ไม่ถูกต้อง"
	}
	data := struct {
		Error string
	}{
		Error: errMsg,
	}

	temp, _ := template.ParseFiles("views/login.html")
	temp.Execute(w, data)
}

func CheckLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/Login", http.StatusSeeOther)
		return
	}

	user := req.FormValue("username")
	pass := req.FormValue("password")

	userData, err := models.GetUserByLogin(user, pass)
	if err == nil {
		http.Redirect(w, req, "/?username="+userData.Username, http.StatusSeeOther)
	} else {
		http.Redirect(w, req, "/Login?error=1", http.StatusSeeOther)
	}
}

func Logout(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "authenticated",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, req *http.Request) {
	data := map[string]string{
		"Error":   "",
		"Success": "",
	}
	temp, _ := template.ParseFiles("views/register.html")
	temp.Execute(w, data)
}

func CreateAccount(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/register", http.StatusSeeOther)
		return
	}

	user := models.User{
		Fullname: req.FormValue("fullname"),
		Email:    req.FormValue("email"),
		Username: req.FormValue("username"),
		Password: req.FormValue("password"),
		Role:     "user",
	}

	data := map[string]string{}

	exists, err := models.IsUsernameOrEmailExists(user.Username, user.Email)
	if err != nil {
		data["Error"] = "เกิดข้อผิดพลาดในการตรวจสอบข้อมูล"
	} else if exists {
		data["Error"] = "Username นี้มีคนใช้แล้ว หรือ Email นี้มีคนใช้แล้ว"
	} else {
		err = models.CreateUser(&user)
		if err != nil {
			data["Error"] = "เกิดข้อผิดพลาดในการสมัครสมาชิก"
		} else {
			data["Success"] = "สมัครสมาชิกสำเร็จ"
		}
	}

	temp, _ := template.ParseFiles("views/register.html")
	temp.Execute(w, data)
}
