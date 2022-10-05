package store

import (
	"davesstore/pkg/posintf/posintffakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store application", func() {

	var (
		s         store
		o         Order
		name      = "Aruna"
		rec       string
		err       error
		posclient posintffakes.FakePOS
	)

	BeforeEach(func() {
		posclient = posintffakes.FakePOS{}
		s = store{
			Name:        name,
			OpeningTime: 9,
			ClosingTime: 17,
			Inventory: map[string]int{
				"apples": 3,
			},
			Prices: map[string]float64{
				"apples": 3.00,
			},
			PosClient: &posclient,
		}
		o = Order{
			OrderTime: 10,
			Items: map[string]int{
				"apples": 2,
			},
		}

	})

	Context("When the customer is welcomed", func() {
		It("should include the store's name in the welcome message", func() {
			Expect(s.GetWelcome()).To(ContainSubstring(name))
		})
	})

	Context("When we say goodbye to the customer", func() {
		It("should include the store's name in the goodbye message", func() {
			Expect(s.GetGoodbye()).To(ContainSubstring(name))
		})
	})

	Context("When we process an order", func() {

		JustBeforeEach(func() {
			rec, err = s.ProcessOrder(o)
		})

		Context("but the store is not open", func() {

			BeforeEach(func() {
				o.OrderTime = 2
			})

			It("should return an error", func() {
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorStoreClosed))
			})
		})

		Context("but the store does not have the inventory", func() {
			BeforeEach(func() {
				o.Items = map[string]int{
					"bananas": 3,
				}
			})

			It("should return an error", func() {
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorOutOfStock))
			})
		})

		Context("and the customer pays in cash", func() {
			BeforeEach(func() {
				o.PaymentMethod = CashPayment
			})

			It("should print a receipt with items we expect", func() {
				Expect(err).To(BeNil())
				Expect(rec).To(ContainSubstring("2 of apples"))
				Expect(rec).To(ContainSubstring(CashPayment))
				Expect(rec).To(ContainSubstring("Total: 6.00"))
			})
		})

		Context("and the customer pays in credit", func() {
			BeforeEach(func() {
				posclient.CreditPaymentReturns(true)
				o.PaymentMethod = CreditPayment
			})

			It("should print a receipt with items we expect", func() {
				Expect(err).To(BeNil())
				Expect(rec).To(ContainSubstring("2 of apples"))
				Expect(rec).To(ContainSubstring(CreditPayment))
				Expect(rec).To(ContainSubstring("Total: 6.00"))
			})
		})

		Context("and the customer pays in credit", func() {
			BeforeEach(func() {
				posclient.CreditPaymentReturns(false)
				o.PaymentMethod = CreditPayment
			})

			It("should print a receipt with items we expect", func() {
				Expect(err).To(Equal(ErrorProcessingPayment))
			})
		})

	})
})
