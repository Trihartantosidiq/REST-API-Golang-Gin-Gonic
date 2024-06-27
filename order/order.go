package order

var (
	orderDefaulType = "NORMAL"
	orderMercgant   = "PREMIUM"
)

type Order struct {
	Number int
	state  string
}

func cretaorder() Order {
	order := Order{state: orderDefaulType}
	return order
}
