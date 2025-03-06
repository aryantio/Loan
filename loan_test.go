package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateBilling(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	expectedLoan := Loan{
		Billings: []*Billing{
			{amount: 110000, state: "pending", dueDate: approvedDate.AddDate(0, 0, 7)},
		},
	}

	loan.CreateBilling()
	assert.Equal(t, expectedLoan.Billings[0], loan.Billings[0], "will have the same the first billing")
}

func TestGetDueDateBillings(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	loan.CreateBilling()
	currentTime, _ := time.Parse("02-01-2006", "15-01-2022")
	dueDateLoan := loan.GetDueDateBillings(currentTime)

	dueDate1, _ := time.Parse("02-01-2006", "08-01-2022")
	dueDate2, _ := time.Parse("02-01-2006", "15-01-2022")
	expectedBilling := []*Billing{
		{amount: 110000, state: "pending", dueDate: dueDate1},
		{amount: 110000, state: "pending", dueDate: dueDate2},
	}

	assert.Equal(t, expectedBilling, dueDateLoan, "Will return the same due date")
}

func TestGetNextPendingBilling(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	loan.CreateBilling()
	currentTime, _ := time.Parse("02-01-2006", "01-01-2022")
	nextDueDateLoan := loan.GetNextPendingBilling(currentTime)

	nextDueDate, _ := time.Parse("02-01-2006", "08-01-2022")
	expectedBilling := &Billing{amount: 110000, state: "pending", dueDate: nextDueDate}

	assert.Equal(t, expectedBilling, nextDueDateLoan, "Will return the same due date")
}

func TestIsDelinquentTrue(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	loan.CreateBilling()
	currentTime, _ := time.Parse("02-01-2006", "01-03-2022")
	delinquent := loan.IsDelinquent(currentTime)
	assert.True(t, delinquent, "should be deliquent")
}

func TestIsDelinquentFalse(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	loan.CreateBilling()
	currentTime, _ := time.Parse("02-01-2006", "08-01-2022")
	delinquent := loan.IsDelinquent(currentTime)
	assert.False(t, delinquent, "should not be deliquent")
}

func TestMakePayment(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	loan.CreateBilling()
	loan.MakePayment()

	dueDate, _ := time.Parse("02-01-2006", "08-01-2022")
	expectedBillingPaid := &Billing{amount: 110000, state: "paid", dueDate: dueDate}
	assert.Equal(t, expectedBillingPaid, loan.Billings[0], "First billing should be already paid")
}

func TestListBilling(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	loan.CreateBilling()
	loanBillings := loan.ListBillings(0, 1)

	dueDate, _ := time.Parse("02-01-2006", "08-01-2022")
	expectedBillings := []*Billing{
		{amount: 110000, state: "pending", dueDate: dueDate},
	}

	assert.Equal(t, expectedBillings, loanBillings, "Have the same value")
}

func TestGetOutstanding(t *testing.T) {
	approvedDate, _ := time.Parse("02-01-2006", "01-01-2022")
	loan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		Term:          50,
		PaidLoan:      0,
		ApprovedDate:  approvedDate,
	}

	loan.CreateBilling()
	loan.MakePayment()
	outstandingBalance := loan.GetOutstanding()
	expectedOutstanding := 5390000
	assert.Equal(t, expectedOutstanding, outstandingBalance, "Will have the same outstanding balance")
}
