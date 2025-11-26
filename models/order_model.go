package models

import (
	"fmt"
	"webapp/config"
)

type OrderItem struct {
	OrderID   string
	JewelryID string
	Name      string
	Quantity  int
	Price     float64
	Total     float64
	Status    string
}

type OrderGroup struct {
	OrderID   string
	OrderDate string
	Status    string
	Items     []OrderItem
	Total     float64
}

type OrderResult struct {
	OrderID      string
	Username     string
	CustomerName string
	Email        string
	JewelryID    string
	ProductName  string
	Quantity     int
	Price        float64
	Total        float64
	OrderDate    string
	Status       string
}

func CreateOrder(username string, cart map[string]CartItem) (string, error) {
	db := config.GetDB()

	// สร้าง orderid
	var orderCount int
	db.QueryRow("SELECT COUNT(*) FROM Orders").Scan(&orderCount)
	orderID := fmt.Sprintf("O%03d", orderCount+1)

	_, err := db.Exec("INSERT INTO Orders (order_id, username, status) VALUES (?, ?, 'pending')", orderID, username)
	if err != nil {
		return "", err
	}

	for jewelryID, item := range cart {

		_, err := db.Exec("INSERT INTO OrderDetail (order_id, jewelry_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderID, jewelryID, item.Quantity, item.Price)
		if err != nil {
			return "", err
		}

		_, err = db.Exec("UPDATE Jewelry SET stock = stock - ? WHERE jewelry_id = ?", item.Quantity, jewelryID)
		if err != nil {
			return "", err
		}
	}

	return orderID, nil
}

func GetUserOrders(username string) ([]OrderGroup, error) {
	db := config.GetDB()

	rows, err := db.Query(`
		SELECT o.order_id, o.order_date, o.status, 
			   od.jewelry_id, j.name, od.quantity, od.price,
			   (od.quantity * od.price) as item_total
		FROM Orders o
		JOIN OrderDetail od ON o.order_id = od.order_id
		JOIN Jewelry j ON od.jewelry_id = j.jewelry_id
		WHERE o.username = ?
		ORDER BY o.order_date DESC, o.order_id
	`, username)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []OrderGroup
	var currentOrderID string
	var currentGroup *OrderGroup

	for rows.Next() {
		var item OrderItem
		var orderDate string

		err := rows.Scan(
			&item.OrderID,
			&orderDate,
			&item.Status,
			&item.JewelryID,
			&item.Name,
			&item.Quantity,
			&item.Price,
			&item.Total,
		)

		if err != nil {
			return nil, err
		}

		if item.OrderID != currentOrderID {
			if currentGroup != nil {
				orders = append(orders, *currentGroup)
			}
			currentOrderID = item.OrderID
			currentGroup = &OrderGroup{
				OrderID:   item.OrderID,
				OrderDate: orderDate,
				Status:    item.Status,
				Items:     []OrderItem{item},
				Total:     item.Total,
			}
		} else {
			currentGroup.Items = append(currentGroup.Items, item)
			currentGroup.Total += item.Total
		}
	}

	if currentGroup != nil {
		orders = append(orders, *currentGroup)
	}

	return orders, nil
}

func SearchOrders(searchCustomer, searchOrder, orderStatus string) ([]OrderResult, error) {
	db := config.GetDB()

	query := `
        SELECT 
            o.order_id,
            u.username,
            u.fullname,
            u.email,
            od.jewelry_id,
            j.name,
            od.quantity,
            od.price,
            (od.quantity * od.price) as total,
            DATE(o.order_date) as order_date,
            o.status
        FROM Orders o
        JOIN user u ON o.username = u.username
        JOIN OrderDetail od ON o.order_id = od.order_id
        JOIN Jewelry j ON od.jewelry_id = j.jewelry_id
        WHERE 1=1
    `

	var args []interface{}

	if searchCustomer != "" {
		query += " AND (u.username LIKE ? OR u.email LIKE ? OR u.fullname LIKE ?)"
		searchPattern := "%" + searchCustomer + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	if searchOrder != "" {
		query += " AND o.order_id LIKE ?"
		args = append(args, "%"+searchOrder+"%")
	}

	if orderStatus != "" {
		query += " AND o.status = ?"
		args = append(args, orderStatus)
	}

	query += " ORDER BY o.order_date DESC, o.order_id"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []OrderResult
	for rows.Next() {
		var result OrderResult
		var orderDate string

		err := rows.Scan(
			&result.OrderID,
			&result.Username,
			&result.CustomerName,
			&result.Email,
			&result.JewelryID,
			&result.ProductName,
			&result.Quantity,
			&result.Price,
			&result.Total,
			&orderDate,
			&result.Status,
		)
		if err != nil {
			return nil, err
		}

		result.OrderDate = orderDate
		results = append(results, result)
	}

	return results, nil
}
