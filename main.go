package main

import (
	"container/heap"
	"fmt"

	"github.com/shopspring/decimal"
)

type net_balance struct {
	user_id int64
	balance decimal.Decimal
}

type minNetBalanceHeap []*net_balance

func (h minNetBalanceHeap) Len() int           { return len(h) }
func (h minNetBalanceHeap) Less(i, j int) bool { return h[i].balance.LessThan(h[j].balance) }
func (h minNetBalanceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minNetBalanceHeap) Push(x any) {
	*h = append(*h, x.(*net_balance))
}

func (h *minNetBalanceHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type maxNetBalanceHeap []*net_balance

func (h maxNetBalanceHeap) Len() int           { return len(h) }
func (h maxNetBalanceHeap) Less(i, j int) bool { return h[i].balance.GreaterThan(h[j].balance) }
func (h maxNetBalanceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxNetBalanceHeap) Push(x any) {
	*h = append(*h, x.(*net_balance))
}

func (h *maxNetBalanceHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func printMinNetBalanceHeap(h minNetBalanceHeap) {
	fmt.Println("minNetBalannceHeap:")
	for i, nb := range h {
		fmt.Printf("  [%d] user_id: %d, balance: %s\n", i, nb.user_id, nb.balance.String())
	}
}

func printMaxNetBalanceHeap(h maxNetBalanceHeap) {
	fmt.Println("maxNetBalannceHeap:")
	for i, nb := range h {
		fmt.Printf("  [%d] user_id: %d, balance: %s\n", i, nb.user_id, nb.balance.String())
	}
}

type payment struct {
	from_id int64
	to_id   int64
	amount  decimal.Decimal
}

func printPayment(p payment) {
	fmt.Printf("Payment: from_id: %d -> to_id: %d, amount: %s\n", p.from_id, p.to_id, p.amount.String())
}

func printPayments(payments []payment) {
	fmt.Println("Payments:")
	for _, p := range payments {
		printPayment(p)
	}
}

func sumBalances(balances ...[]*net_balance) decimal.Decimal {
	var sum decimal.Decimal
	for _, v := range balances {
		for _, b := range v {
			sum = sum.Add(b.balance)
		}
	}
	return sum

}

func main() {

	payerNetBalances := []*net_balance{
		{user_id: 1, balance: decimal.NewFromInt(5)},
		{user_id: 2, balance: decimal.NewFromInt(5)},
		{user_id: 3, balance: decimal.NewFromInt(10)},
	}

	payeeNetBalances := []*net_balance{
		{user_id: 4, balance: decimal.NewFromInt(-10)},
		{user_id: 5, balance: decimal.NewFromInt(-10)},
	}

	if !sumBalances(payerNetBalances, payeeNetBalances).IsZero() {
		// in real application, this check would return an error
		panic("Balances do not net to zero, function aborted")
	}

	payerHeap := (*maxNetBalanceHeap)(&payerNetBalances)

	// Initialize (and sort) payerHeap
	heap.Init(payerHeap)
	printMaxNetBalanceHeap(*payerHeap)

	payeeHeap := (*minNetBalanceHeap)(&payeeNetBalances)

	heap.Init(payeeHeap)
	printMinNetBalanceHeap(*payeeHeap)

	// Create a payment slice for the maximum amount of payments
	payments := make([]payment, payerHeap.Len())
	var paymentCount int

	for payerHeap.Len() > 0 {
		payer := heap.Pop(payerHeap).(*net_balance)
		payee := heap.Pop(payeeHeap).(*net_balance)
		delta := payee.balance.Add(payer.balance)

		var pa decimal.Decimal

		if delta.IsZero() {
			pa = payer.balance
		}
		if delta.IsNegative() {
			payeeHeap.Push(&net_balance{user_id: payee.user_id, balance: delta})
			pa = delta.Abs()
		}
		if delta.IsPositive() {
			payerHeap.Push(&net_balance{user_id: payer.user_id, balance: delta})
			pa = delta
		}
		pmt := payment{from_id: payer.user_id, to_id: payee.user_id, amount: pa}
		payments = append(payments, pmt)
		paymentCount++
		printPayment(pmt)
		printMaxNetBalanceHeap(*payerHeap)
		printMinNetBalanceHeap(*payeeHeap)

	}

	payments = payments[paymentCount:]
	printPayments(payments)

}
