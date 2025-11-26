package controllers

import (
	"net/http"
	"strconv"
	"text/template"
	"webapp/models"
)

func AdminDashboard(w http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	message := req.URL.Query().Get("message")
	messageType := req.URL.Query().Get("type")

	if username == "" {
		http.Redirect(w, req, "/Login", http.StatusSeeOther)
		return
	}

	
	role, err := models.GetUserRole(username)
	if err != nil || role != "admin" {
		http.Redirect(w, req, "/?username="+username, http.StatusSeeOther)
		return
	}

	
	products, _ := models.GetAllProducts()
	users, _ := models.GetAllUsers()

	
	searchCustomer := req.URL.Query().Get("search_customer")
	searchOrder := req.URL.Query().Get("search_order")
	orderStatus := req.URL.Query().Get("order_status")

	
	var orderResults []models.OrderResult
	hasSearch := searchCustomer != "" || searchOrder != "" || orderStatus != ""

	if hasSearch {
		
		orderResults, err = models.SearchOrders(searchCustomer, searchOrder, orderStatus)
		if err != nil {
			
			orderResults = []models.OrderResult{}
		}
	}

	type DashboardData struct {
		Username       string
		Products       []models.Product
		Users          []models.User
		OrderResults   []models.OrderResult
		SearchCustomer string
		SearchOrder    string
		OrderStatus    string
		HasSearch      bool
		Message        string
		MessageType    string
	}

	data := DashboardData{
		Username:       username,
		Products:       products,
		Users:          users,
		OrderResults:   orderResults,
		SearchCustomer: searchCustomer,
		SearchOrder:    searchOrder,
		OrderStatus:    orderStatus,
		HasSearch:      hasSearch,
		Message:        message,
		MessageType:    messageType,
	}

	tmpl, err := template.ParseFiles("views/admin_dashboard.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}


func UpdateStock(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/admin", http.StatusSeeOther)
		return
	}

	username := req.FormValue("username")
	jewelryID := req.FormValue("jewelry_id")
	stock := req.FormValue("stock")

	role, err := models.GetUserRole(username)
	if err != nil || role != "admin" {
		http.Redirect(w, req, "/?username="+username, http.StatusSeeOther)
		return
	}

	stockInt, _ := strconv.Atoi(stock)
	err = models.UpdateProductStock(jewelryID, stockInt)
	if err != nil {
		http.Redirect(w, req, "/admin?username="+username+"&message=เกิดข้อผิดพลาด&type=error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, req, "/admin?username="+username+"&message=อัพเดทจำนวนสต็อกเรียบร้อย&type=success", http.StatusSeeOther)
}

func UpdateUserRole(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/admin", http.StatusSeeOther)
		return
	}

	username := req.FormValue("username")
	targetUsername := req.FormValue("target_username")
	newRole := req.FormValue("role")

	role, err := models.GetUserRole(username)
	if err != nil || role != "admin" {
		http.Redirect(w, req, "/?username="+username, http.StatusSeeOther)
		return
	}

	err = models.UpdateUserRole(targetUsername, newRole)
	if err != nil {
		http.Redirect(w, req, "/admin?username="+username+"&message=เกิดข้อผิดพลาด&type=error", http.StatusSeeOther)
		return
	}

	http.Redirect(w, req, "/admin?username="+username+"&message=เปลี่ยนบทบาทผู้ใช้เรียบร้อย&type=success", http.StatusSeeOther)
}
