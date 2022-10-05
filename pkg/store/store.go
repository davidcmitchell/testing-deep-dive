package store

import (
	"davesstore/pkg/pointofsale"
	"davesstore/pkg/posintf"
	"fmt"
)

const (
	CashPayment   = "cash"
	DebitPayment  = "debit"
	CreditPayment = "credit"
)

var (
	ErrorStoreClosed       = fmt.Errorf("store is closed")
	ErrorOutOfStock        = fmt.Errorf("store is out of stock")
	ErrorProcessingPayment = fmt.Errorf("could not process payment")
)

type Order struct {
	OrderTime     int // presumed 24h clock
	Items         map[string]int
	PaymentMethod string
}

type Store interface {
	GetWelcome() string
	GetGoodbye() string
	ProcessOrder(Order) (string, error)
}

type store struct {
	Name        string
	Inventory   map[string]int
	Prices      map[string]float64
	OpeningTime int // presumed 24h clock
	ClosingTime int // presumed 24h clock
	PosClient   posintf.POS
}

func NewStore(name string, inv map[string]int, prices map[string]float64, open int, close int) Store {
	pos, err := pointofsale.NewPointOfSaleClient()
	if err != nil {
		return nil
	}

	return store{
		Name:        name,
		Inventory:   inv,
		Prices:      prices,
		OpeningTime: open,
		ClosingTime: close,
		PosClient:   pos,
	}
}

func (s store) GetWelcome() string {
	return fmt.Sprintf("Welcome to %s's Store", s.Name)
}

func (s store) GetGoodbye() string {
	return fmt.Sprintf("Thank you for shopping at %s's Store", s.Name)
}

func (s store) ProcessOrder(o Order) (string, error) {

	// is the store open?
	if s.OpeningTime > o.OrderTime || s.ClosingTime < o.OrderTime {
		return "", ErrorStoreClosed
	}

	// does it have the inventory?
	for item, quant := range o.Items {
		stock, ok := s.Inventory[item]
		if !ok || quant > stock {
			return "", ErrorOutOfStock
		}
	}

	// does the payment go through
	if o.PaymentMethod != CashPayment {
		/*pos, err := pointofsale.NewPointOfSaleClient()
		if err != nil {
			return "", err
		}*/

		var ok bool
		if o.PaymentMethod == DebitPayment {
			//ok = pos.DebitPayment(s.getTotal(o))
		}
		if o.PaymentMethod == CreditPayment {
			ok = s.PosClient.CreditPayment(s.getTotal(o))
		}
		if !ok {
			return "", ErrorProcessingPayment
		}

	}

	return s.createReceipt(o), nil
}

func (s store) getTotal(o Order) float64 {
	var total float64
	for item, quant := range o.Items {
		amt := float64(quant) * s.Prices[item]
		total += amt
	}
	return total
}

func (s store) createReceipt(o Order) string {
	var receipt string
	for item, quant := range o.Items {
		amt := float64(quant) * s.Prices[item]
		receipt += fmt.Sprintf("%d of %s: %.2f\n", quant, item, amt)
	}
	receipt += fmt.Sprintf("Total: %.2f\n", s.getTotal(o))
	receipt += fmt.Sprintf("Payment type: %s\n", o.PaymentMethod)
	return receipt
}
