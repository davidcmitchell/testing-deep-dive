package store

/*

	Context("When the customer is welcome", func() {
		It("should include the store's name", func() {
			s := store{
				Name: "David",
			}
			Expect(s.GetWelcome()).To(ContainSubstring("David"))
		})
	})

	Context("When the customer is thanked", func() {
		It("should include the store's name", func() {
			s := store{
				Name: "David",
			}
			Expect(s.GetGoodbye()).To(ContainSubstring("David"))
		})
	})


    ...


		var (
		s store
		name = "David"
	)

	BeforeEach(func() {
		s = store{
			Name: name,
		}
	})

	Context("When the customer is welcome", func() {
		It("should include the store's name", func() {
			Expect(s.GetWelcome()).To(ContainSubstring(name))
		})
	})

	Context("When the customer is thanked", func() {
		It("should include the store's name", func() {
			Expect(s.GetGoodbye()).To(ContainSubstring(name))
		})
	})


   ...


   	Context("When the customer makes an order", func() {

		Context("but the store is not open", func() {
			It("should return an error", func() {
				order := Order{
					OrderTime: 2,
				}
				s.OpeningTime = 9
				s.ClosingTime = 17
				rec, err := s.ProcessOrder(order)
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorStoreClosed))
			})
		})

		Context("but we are out of inventory", func() {
			It("should return an error", func() {
				order := Order{
					OrderTime: 10,
					Items: map[string]int{
						"banana": 10,
					},
				}
				s.OpeningTime = 9
				s.ClosingTime = 17
				rec, err := s.ProcessOrder(order)
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorOutOfStock))

			})
		})


...

package store

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store application", func() {

	var (
		s    store
		o    Order
		name = "David"
		rec  string
		err  error
	)

	BeforeEach(func() {
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
		}
		o = Order{
			OrderTime: 10,
			Items: map[string]int{
				"apples": 2,
			},
		}
	})

	Context("When the customer is welcome", func() {
		It("should include the store's name", func() {
			Expect(s.GetWelcome()).To(ContainSubstring(name))
		})
	})

	Context("When the customer is thanked", func() {
		It("should include the store's name", func() {
			Expect(s.GetGoodbye()).To(ContainSubstring(name))
		})
	})

	Context("When the customer makes an order", func() {
		JustBeforeEach(func() {
			rec, err = s.ProcessOrder(o)
		})

		Context("but the store is not open", func() {
			BeforeEach(func() {
				o.OrderTime = 1
			})

			It("should return an error", func() {
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorStoreClosed))
			})
		})

		Context("but we are out of inventory", func() {
			BeforeEach(func() {
				o.Items = map[string]int{
					"banana": 10,
				}
			})

			It("should return an error", func() {
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorOutOfStock))
			})
		})

		Context("and the payment is cash", func() {
			BeforeEach(func() {
				o.PaymentMethod = CashPayment
			})
			It("should not return an error", func() {
				Expect(err).To(BeNil())
				Expect(rec).To(ContainSubstring("2 of apples"))
				Expect(rec).To(ContainSubstring(CashPayment))
				Expect(rec).To(ContainSubstring("Total: 6.00"))
			})
		})
	})
})

...


package store

import (
	"davesstore/pkg/posintf/posintffakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store application", func() {

	var (
		s          store
		o          Order
		name       = "David"
		rec        string
		err        error
		fakeclient posintffakes.FakePOS
	)

	BeforeEach(func() {
		fakeclient = posintffakes.FakePOS{}
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
			PosClient: &fakeclient,
		}
		o = Order{
			OrderTime: 10,
			Items: map[string]int{
				"apples": 2,
			},
		}
	})

	Context("When the customer is welcome", func() {
		It("should include the store's name", func() {
			Expect(s.GetWelcome()).To(ContainSubstring(name))
		})
	})

	Context("When the customer is thanked", func() {
		It("should include the store's name", func() {
			Expect(s.GetGoodbye()).To(ContainSubstring(name))
		})
	})

	Context("When the customer makes an order", func() {
		JustBeforeEach(func() {
			rec, err = s.ProcessOrder(o)
		})

		Context("but the store is not open", func() {
			BeforeEach(func() {
				o.OrderTime = 1
			})

			It("should return an error", func() {
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorStoreClosed))
			})
		})

		Context("but we are out of inventory", func() {
			BeforeEach(func() {
				o.Items = map[string]int{
					"banana": 10,
				}
			})

			It("should return an error", func() {
				Expect(rec).To(Equal(""))
				Expect(err).To(Equal(ErrorOutOfStock))
			})
		})

		Context("and the payment is cash", func() {
			BeforeEach(func() {
				o.PaymentMethod = CashPayment
			})
			It("should not return an error", func() {
				Expect(err).To(BeNil())
				Expect(rec).To(ContainSubstring("2 of apples"))
				Expect(rec).To(ContainSubstring(CashPayment))
				Expect(rec).To(ContainSubstring("Total: 6.00"))
			})
		})

		Context("and the payment is credit", func() {
			BeforeEach(func() {
				fakeclient.CreditPaymentReturns(true)
				o.PaymentMethod = CreditPayment
			})
			It("should not return an error", func() {
				Expect(err).To(BeNil())
				Expect(rec).To(ContainSubstring("2 of apples"))
				Expect(rec).To(ContainSubstring(CreditPayment))
				Expect(rec).To(ContainSubstring("Total: 6.00"))
			})
		})
	})
})

*/
