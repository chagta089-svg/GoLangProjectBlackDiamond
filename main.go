package main

import (
	"log"
	"net/http"
	"webapp/config"
	"webapp/controllers"
	_ "webapp/models"
)

func main() {

	err := config.InitDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/Login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/checkLogin", controllers.CheckLogin)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/createAccount", controllers.CreateAccount)

	http.HandleFunc("/listProduct", controllers.ListProduct)

	http.HandleFunc("/cart", controllers.Cart)
	http.HandleFunc("/addToCart", controllers.AddToCart)
	http.HandleFunc("/updateCartItem", controllers.UpdateCartItem)
	http.HandleFunc("/removeFromCart", controllers.RemoveFromCart)
	http.HandleFunc("/checkout", controllers.Checkout)

	http.HandleFunc("/orders", controllers.Orders)

	http.HandleFunc("/admin", controllers.AdminDashboard)
	http.HandleFunc("/admin/updateStock", controllers.UpdateStock)
	http.HandleFunc("/admin/updateUserRole", controllers.UpdateUserRole)

	log.Println("Server starting on :8011")
	log.Fatal(http.ListenAndServe(":8011", nil))
}
