package posintf

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . POS
type POS interface {
	CreditPayment(amt float64) bool
	DebitPayment(amt float64) bool
}
