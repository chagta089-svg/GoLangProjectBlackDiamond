package controllers

import (
	"net/http"
	"text/template"
	"webapp/models"
)

type IndexData struct {
	Types    []models.ProductType
	Username string
	Role     string
}

func Index(w http.ResponseWriter, req *http.Request) {
	types, err := models.GetAllProductTypes()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	username := req.URL.Query().Get("username")
	var role string
	if username != "" {
		role, _ = models.GetUserRole(username)
	}

	data := IndexData{
		Types:    types,
		Username: username,
		Role:     role,
	}

	tmpl, _ := template.ParseFiles("views/index.html")
	tmpl.Execute(w, data)
}

type ListProductData struct {
	Products []models.Product
	Username string
}

func ListProduct(w http.ResponseWriter, req *http.Request) {
	typeID := req.URL.Query().Get("typeid")
	username := req.URL.Query().Get("username")

	products, err := models.GetProductsByType(typeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ListProductData{
		Products: products,
		Username: username,
	}

	tmpl, _ := template.ParseFiles("views/productList.html")
	tmpl.Execute(w, data)
}
