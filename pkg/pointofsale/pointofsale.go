package pointofsale

import (
	"fmt"
	"time"
)

type PointOfSaleClient struct {
}

func NewPointOfSaleClient() (*PointOfSaleClient, error) {
	return nil, fmt.Errorf("unable to initialize")
}

func (posc *PointOfSaleClient) CreditPayment(amt float64) bool {
	time.Sleep(time.Hour)
	return true
}

func (posc *PointOfSaleClient) DebitPayment(amt float64) bool {
	time.Sleep(time.Hour)
	return true
}
