package models

import (
	"webapp/config"
)

type ProductType struct {
	TypeID   string
	TypeName string
}

type Product struct {
	ID       string
	Name     string
	Price    float64
	Type     string
	Stock    int
	Material string
	TypeName string
}

func GetAllProductTypes() ([]ProductType, error) {
	db := config.GetDB()
	rows, err := db.Query("SELECT typeid, typename FROM type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []ProductType
	for rows.Next() {
		var t ProductType
		if err := rows.Scan(&t.TypeID, &t.TypeName); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func GetProductsByType(typeID string) ([]Product, error) {
	db := config.GetDB()
	rows, err := db.Query(`
		SELECT j.jewelry_id, j.name, j.price, t.typename, j.stock, j.material
		FROM Jewelry j
		JOIN Type t ON j.typeid = t.typeid
		WHERE j.typeid = ?`, typeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.TypeName, &p.Stock, &p.Material); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProductByID(productID string) (*Product, error) {
	db := config.GetDB()
	var product Product
	err := db.QueryRow(`
		SELECT j.jewelry_id, j.name, j.material, j.price, j.stock, t.typename 
		FROM Jewelry j 
		JOIN Type t ON j.typeid = t.typeid 
		WHERE j.jewelry_id = ?`, productID,
	).Scan(&product.ID, &product.Name, &product.Material, &product.Price, &product.Stock, &product.TypeName)
	
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func UpdateProductStock(productID string, stock int) error {
	db := config.GetDB()
	_, err := db.Exec("UPDATE Jewelry SET stock = ? WHERE jewelry_id = ?", stock, productID)
	return err
}

func GetAllProducts() ([]Product, error) {
	db := config.GetDB()
	rows, err := db.Query(`
		SELECT j.jewelry_id, j.name, t.typename, j.price, j.stock
		FROM Jewelry j 
		JOIN Type t ON j.typeid = t.typeid
		ORDER BY j.jewelry_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.TypeName, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}