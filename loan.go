package main

import "time"

// creating a struct Loan
type Loan struct {
	PrincipalLoan int        // Load which will disbursed to the customer
	PayableLoan   int        // How much amount customer need to repay
	PaidLoan      int        // How much amount which customer already paid
	ApprovedDate  time.Time  // date when the Loan got approved , Day 1 of loan schedule
	Term          int        // Term how many customer need to pay
	Billings      []*Billing // List of billings, status about repayment of the loan
}

// creating struct of Billing
type Billing struct {
	amount  int       // amount of billing
	dueDate time.Time // due date when the payment need to do
	state   string    // paid or pending
}

func (l *Loan) GetOutstanding() int {
	return l.PayableLoan - l.PaidLoan
}

// createBilling when initialize from loan
func (l *Loan) CreateBilling() {
	// find the amount per billing
	amountBilling := l.PayableLoan / l.Term
	dueDate := l.ApprovedDate.AddDate(0, 0, 7) // adding the due date 7 days, since we have a repaymeent on weekly basis

	//create a billing
	var billings []*Billing
	for i := 0; i < l.Term; i++ {
		billings = append(billings, &Billing{amount: amountBilling, dueDate: dueDate, state: "pending"}) // assign the initial state is pending
		dueDate = dueDate.AddDate(0, 0, 7)
	}
	l.Billings = billings
}

// function to get the billing which is already pass the due date
func (l *Loan) GetDueDateBillings(currentTime time.Time) (dueDateBilling []*Billing) {
	for _, billing := range l.Billings {
		// just find the pending state have been pass the due date.
		// assumtion when customer haven't paid until repayment schedule will be marked as pass duedate.
		if billing.state == "pending" && (billing.dueDate.Before(currentTime) || billing.dueDate.Equal(currentTime)) {
			dueDateBilling = append(dueDateBilling, billing)
		}
	}

	return dueDateBilling
}

// function to get the next pending billing from current time
func (l *Loan) GetNextPendingBilling(currentTime time.Time) (nextBilling *Billing) {
	for _, billing := range l.Billings {
		if billing.dueDate.After(currentTime) && billing.state == "pending" {
			nextBilling = billing
			break
		}
	}

	return nextBilling
}

// funtion to get the status wether the status of loan already delinquent or not
func (l *Loan) IsDelinquent(currentTime time.Time) bool {
	dueDateBillings := l.GetDueDateBillings(currentTime)
	if len(dueDateBillings) > 2 {
		return true
	} else {
		return false
	}
}

// making payment, since the making payment is sequential, just find the most oldest billing have not been paid
func (l *Loan) MakePayment() {
	for _, billing := range l.Billings {
		if billing.state == "pending" {
			billing.state = "paid"
			l.PaidLoan += billing.amount
			break
		}
	}
}

// function to get list billing with offset and limit parameters
func (l *Loan) ListBillings(offset int, limit int) (b []*Billing) {
	end := min(offset+limit, len(l.Billings))
	return l.Billings[offset:end]
}
