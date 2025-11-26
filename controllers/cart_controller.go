package controllers

import (
	"net/http"
	"strconv"
	"text/template"
	"webapp/models"
)

type CartData struct {
	CartItems   []models.CartItem
	TotalItems  int
	Subtotal    float64
	Total       float64
	Username    string
	Message     string
	MessageType string
}

func Cart(w http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	message := req.URL.Query().Get("message")
	messageType := req.URL.Query().Get("type")

	userCart := models.GetUserCart(username)
	var cartItems []models.CartItem
	totalItems := 0
	subtotal := 0.0

	for _, cartItem := range userCart {
		// ดึงข้อมูลล่าสุดจาก database
		product, err := models.GetProductByID(cartItem.JewelryID)
		if err == nil {
			item := models.CartItem{
				JewelryID: cartItem.JewelryID,
				Name:      product.Name,
				Material:  product.Material,
				Price:     product.Price,
				Quantity:  cartItem.Quantity,
				Stock:     product.Stock,
			}
			cartItems = append(cartItems, item)
			totalItems += item.Quantity
			subtotal += item.Price * float64(item.Quantity)
		}
	}

	data := CartData{
		CartItems:   cartItems,
		TotalItems:  totalItems,
		Subtotal:    subtotal,
		Total:       subtotal,
		Username:    username,
		Message:     message,
		MessageType: messageType,
	}

	tmpl, _ := template.ParseFiles("views/cart.html")
	tmpl.Execute(w, data)
}

func AddToCart(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/cart", http.StatusSeeOther)
		return
	}

	jewelryID := req.FormValue("jewelry_id")
	quantity := req.FormValue("quantity")
	username := req.FormValue("username")

	if username == "" {
		http.Redirect(w, req, "/Login", http.StatusSeeOther)
		return
	}

	qty, err := strconv.Atoi(quantity)
	if err != nil || qty <= 0 {
		http.Redirect(w, req, "/cart?username="+username+"&message=จำนวนไม่ถูกต้อง&type=error", http.StatusSeeOther)
		return
	}

	product, err := models.GetProductByID(jewelryID)
	if err != nil {
		http.Redirect(w, req, "/cart?username="+username+"&message=ไม่พบสินค้า&type=error", http.StatusSeeOther)
		return
	}

	if product.Stock < qty {
		http.Redirect(w, req, "/cart?username="+username+"&message=สต็อกไม่เพียงพอ&type=error", http.StatusSeeOther)
		return
	}

	userCart := models.GetUserCart(username)

	if existingItem, exists := userCart[jewelryID]; exists {
		existingItem.Quantity += qty
		userCart[jewelryID] = existingItem
	} else {
		userCart[jewelryID] = models.CartItem{
			JewelryID: jewelryID,
			Name:      product.Name,
			Material:  product.Material,
			Price:     product.Price,
			Quantity:  qty,
			Stock:     product.Stock,
		}
	}

	models.SaveUserCart(username, userCart)

	http.Redirect(w, req, "/cart?username="+username+"&message=เพิ่มสินค้าในตะกร้าเรียบร้อย&type=success", http.StatusSeeOther)
}

func UpdateCartItem(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/cart", http.StatusSeeOther)
		return
	}

	jewelryID := req.FormValue("jewelry_id")
	quantity := req.FormValue("quantity")
	username := req.FormValue("username")

	qty, _ := strconv.Atoi(quantity)

	userCart := models.GetUserCart(username)
	if item, exists := userCart[jewelryID]; exists {
		item.Quantity = qty
		userCart[jewelryID] = item
		models.SaveUserCart(username, userCart)
	}

	http.Redirect(w, req, "/cart?username="+username+"&message=อัพเดทจำนวนสินค้าเรียบร้อย&type=success", http.StatusSeeOther)
}

func RemoveFromCart(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/cart", http.StatusSeeOther)
		return
	}

	jewelryID := req.FormValue("jewelry_id")
	username := req.FormValue("username")

	userCart := models.GetUserCart(username)
	delete(userCart, jewelryID)
	models.SaveUserCart(username, userCart)

	http.Redirect(w, req, "/cart?username="+username+"&message=ลบสินค้าจากตะกร้าเรียบร้อย&type=success", http.StatusSeeOther)
}
