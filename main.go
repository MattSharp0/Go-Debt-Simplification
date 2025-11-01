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

	// debtorNetBalances := []*net_balance{
	// 	{user_id: 1, balance: decimal.NewFromInt(5)},
	// 	{user_id: 2, balance: decimal.NewFromInt(5)},
	// 	{user_id: 3, balance: decimal.NewFromInt(10)},
	// }

	// creditorNetBalances := []*net_balance{
	// 	{user_id: 4, balance: decimal.NewFromInt(-10)},
	// 	{user_id: 5, balance: decimal.NewFromInt(-10)},
	// }

	debtorNetBalances := []*net_balance{
		{user_id: 3, balance: decimal.NewFromInt(35)},
		{user_id: 4, balance: decimal.NewFromInt(225)},
	}

	creditorNetBalances := []*net_balance{
		{user_id: 2, balance: decimal.NewFromInt(-260)},
	}

	if !sumBalances(debtorNetBalances, creditorNetBalances).IsZero() {
		// in real application, this check would return an error
		panic("Balances do not net to zero, function aborted")
	}

	debtorHeap := (*maxNetBalanceHeap)(&debtorNetBalances)

	// Initialize (and sort) debtorHeap
	heap.Init(debtorHeap)
	printMaxNetBalanceHeap(*debtorHeap)

	creditorHeap := (*minNetBalanceHeap)(&creditorNetBalances)

	heap.Init(creditorHeap)
	printMinNetBalanceHeap(*creditorHeap)

	// Create a payment slice for the maximum amount of payments
	payments := make([]payment, debtorHeap.Len())
	var paymentCount int

	fmt.Printf("Debtor heap length: %d\n", debtorHeap.Len())

	for debtorHeap.Len() > 0 {
		fmt.Printf("Debtor heap length: %d\n", debtorHeap.Len())
		debtor := heap.Pop(debtorHeap).(*net_balance)
		fmt.Printf("debtor: %v, Balance: %d\n", debtor.user_id, debtor.balance.IntPart())

		creditor := heap.Pop(creditorHeap).(*net_balance)
		fmt.Printf("creditor: %v, Balance: %d\n", creditor.user_id, creditor.balance.IntPart())
		delta := creditor.balance.Add(debtor.balance)

		var pa decimal.Decimal

		if delta.IsZero() {
			pa = debtor.balance
		}
		if delta.IsNegative() {
			creditorHeap.Push(&net_balance{user_id: creditor.user_id, balance: delta})
			pa = decimal.Min(creditor.balance.Abs(), debtor.balance)
		}
		if delta.IsPositive() {
			debtorHeap.Push(&net_balance{user_id: debtor.user_id, balance: delta})
			pa = decimal.Min(creditor.balance.Abs(), debtor.balance)
		}
		pmt := payment{from_id: debtor.user_id, to_id: creditor.user_id, amount: pa}
		payments = append(payments, pmt)
		paymentCount++
		printPayment(pmt)
		printMaxNetBalanceHeap(*debtorHeap)
		printMinNetBalanceHeap(*creditorHeap)

	}

	payments = payments[paymentCount:]

	printPayments(payments)

}
