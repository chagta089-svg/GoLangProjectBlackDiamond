package controllers

import (
	"net/http"
	"text/template"
	"webapp/models"
)

func Checkout(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/cart", http.StatusSeeOther)
		return
	}

	username := req.FormValue("username")

	userCart := models.GetUserCart(username)
	if len(userCart) == 0 {
		http.Redirect(w, req, "/cart?username="+username+"&message=ตะกร้าว่างเปล่า&type=error", http.StatusSeeOther)
		return
	}

	// สร้าง order ใหม่
	orderID, err := models.CreateOrder(username, userCart)
	if err != nil {
		http.Redirect(w, req, "/cart?username="+username+"&message=เกิดข้อผิดพลาด&type=error", http.StatusSeeOther)
		return
	}

	// ล้างตะกร้า
	models.ClearUserCart(username)

	http.Redirect(w, req, "/cart?username="+username+"&message=สั่งซื้อสำเร็จ รหัสออเดอร์: "+orderID+"&type=success", http.StatusSeeOther)
}

func Orders(w http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	if username == "" {
		http.Redirect(w, req, "/Login", http.StatusSeeOther)
		return
	}

	orders, err := models.GetUserOrders(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	type OrdersData struct {
		Orders   []models.OrderGroup
		Username string
		Count    int
	}

	data := OrdersData{
		Orders:   orders,
		Username: username,
		Count:    len(orders),
	}

	tmpl, _ := template.ParseFiles("views/orders.html")
	tmpl.Execute(w, data)
}
