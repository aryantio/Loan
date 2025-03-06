package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	var strDate string
	var initialDate time.Time
	var isValidDate bool

	// initialize the loan with 5mio principalLoan with 10% interest, we got 5,5 mio for repayment, and put the term is 50 times repayment date.
	objLoan := Loan{
		PrincipalLoan: 5000000,
		PayableLoan:   5500000,
		PaidLoan:      0,
		Term:          50,
	}

	fmt.Println("Selamat datang")
	fmt.Println()

	// We need to input the approved date of loan
	// validate the date format
	for {
		fmt.Print("Masukkan tanggal pinjaman kamu di terima (format DD-MM-YYYY) : ")
		fmt.Scan(&strDate)
		initialDate, isValidDate = validDate(strDate)
		if isValidDate {
			objLoan.ApprovedDate = initialDate
			break
		} else {
			fmt.Println("Not valid date")
		}
	}

	// after we input the balance, we will crate the billing
	objLoan.CreateBilling()

	// show the profile page with given loan object and we put the approved date as the current date.
	profilePage(&objLoan, initialDate)
}

// function to validate the input of users
func validDate(strDate string) (time.Time, bool) {
	validTime, err := time.Parse("02-01-2006", strDate)
	if err != nil {
		return time.Time{}, false
	} else {
		return validTime, true
	}
}

// page to show the profile Page
func profilePage(loan *Loan, currentDate time.Time) {
	var choice int
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("Tanggal hari ini :", currentDate.Format("02-01-2006"))
	fmt.Println("Sisa tagihan kamu : Rp", loan.GetOutstanding())

	// check wether have a payment pass the due date
	dueDateBillings := loan.GetDueDateBillings(currentDate)
	if len(dueDateBillings) == 0 {
		fmt.Println("Pembayaran kamu bagus, tidak ada yang melewati jatuh tempo")
	} else {
		fmt.Println("Tagihan telah melewati jatuh tempo: ")
		for _, dueDateBilling := range dueDateBillings {
			fmt.Println("Tagihan sebesar Rp", dueDateBilling.amount, " telah jatuh tempo pada ", dueDateBilling.dueDate.Format("02-01-2006"), ". Silakan melakukan pembayaran")
		}
	}

	// inform the next billing due date
	nextBilling := loan.GetNextPendingBilling(currentDate)
	fmt.Println("Tagihan kamu sebesar Rp", nextBilling.amount, " akan jatuh tempo pada ", nextBilling.dueDate.Format("02-01-2006"))

	// inform the status of customer's credit, if customer haven't paid more thant 2 weeks, the status will be delinquent
	fmt.Print("Status credit kamu ")
	if loan.IsDelinquent(currentDate) {
		fmt.Println("Delinquent")
	} else {
		fmt.Println("Credit Bagus")
	}

	fmt.Println()
	// this is menu prompt, can show list billings, paying the oldest repayment, can use the easter egg to change the current date amd exit
	fmt.Println("Pilih menu")
	fmt.Println("1. List tagihan kamu")
	fmt.Println("2. Bayar tagihan terbaru")
	fmt.Println("3. Easter egg untuk setting current date format")
	fmt.Println("4. Keluar")
	fmt.Print("Masukkan pilihan anda : ")

	fmt.Scan(&choice)

	switch choice {
	case 1:
		// go to page list tagihan
		billingPage(loan, 0, currentDate)
	case 2:
		// pay the transaction.
		loan.MakePayment()
	case 3:
		// go to change date
		fmt.Println("Masukkan tanggal yang mau di set, harus lebih dari tanggal approval (format DD-MM-YYYY) :")
		var dateStr string
		var err error

		fmt.Scan(&dateStr)
		currentDate, err = time.Parse("02-01-2006", dateStr)
		if err != nil {
			fmt.Println("Gagal set tanggal")
		} else {
			fmt.Println("Berhasil set tanggal")
		}
	case 4:
		// exit the program
		os.Exit(0)
	}

	profilePage(loan, currentDate)
}

// function to show billing list Page
func billingPage(loan *Loan, offset int, currentDate time.Time) {
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("List tagihan :")

	// create the list of billings with pagination
	for _, billing := range loan.ListBillings(offset, 10) {
		fmt.Println("Jumlah Tagihan : Rp", billing.amount, "Status :", billing.state, "Jatuh Tempo: ", billing.dueDate.Format("02-01-2006"))
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("Menu pilihan : ")
	if offset < len(loan.Billings)-10 {
		fmt.Println("1. Next")
	}

	if offset != 0 {
		fmt.Println("2. Back")
	}
	fmt.Println("3. Kembali ke halaman profile")

	fmt.Print("Masukkan pilihan anda : ")
	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		if offset == len(loan.Billings)-10 {
			fmt.Println("Gak boleh yaa")
		} else {
			offset += 10
		}
		billingPage(loan, offset, currentDate)
	case 2:
		if offset == 0 {
			fmt.Println("Gak boleh yaa")
		} else {
			offset -= 10
		}
		billingPage(loan, offset, currentDate)
	case 3:
		profilePage(loan, currentDate)
	}

}
