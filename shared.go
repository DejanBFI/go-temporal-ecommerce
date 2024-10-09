package app

const CartTaskQueue = "CART_TASK_QUEUE"
const CartMessagesSignal = "cartMessages"

var QueryTypes = struct {
	GET_CART string
}{
	GET_CART: "getCart",
}

var RouteTypes = struct {
	ADD_TO_CART      string
	REMOVE_FROM_CART string
}{
	ADD_TO_CART:      "add_to_cart",
	REMOVE_FROM_CART: "remove_from_cart",
}

type RouteSignal struct {
	Route string
}

type AddToCartSignal struct {
	Route string
	Item  CartItem
}

type RemoveFromCartSignal struct {
	Route string
	Item  CartItem
}
