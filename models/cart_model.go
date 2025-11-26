package models

var cartStore = make(map[string]map[string]CartItem)

type CartItem struct {
	JewelryID string
	Name      string
	Material  string
	Price     float64
	Quantity  int
	Stock     int
}

func GetUserCart(username string) map[string]CartItem {
	if cart, exists := cartStore[username]; exists {
		return cart
	}
	return make(map[string]CartItem)
}

func SaveUserCart(username string, cart map[string]CartItem) {
	cartStore[username] = cart
}

func ClearUserCart(username string) {
	delete(cartStore, username)
}