package integrationtest

import (
	"math"
	"testing"
	"time"
)

// GLEntry represents a General Ledger Entry
type GLEntry struct {
	Account     string
	Debit       float64
	Credit      float64
	PartyType   string
	Party       string
	VoucherType string
	VoucherNo   string
	PostingDate time.Time
}

// GLMap is a slice of GL entries with helper methods
type GLMap []GLEntry

// TotalDebit returns the sum of all debit amounts
func (m GLMap) TotalDebit() float64 {
	var total float64
	for _, e := range m {
		total += e.Debit
	}
	return total
}

// TotalCredit returns the sum of all credit amounts
func (m GLMap) TotalCredit() float64 {
	var total float64
	for _, e := range m {
		total += e.Credit
	}
	return total
}

// IsBalanced checks if total debits equal total credits
func (m GLMap) IsBalanced() bool {
	return math.Abs(m.TotalDebit()-m.TotalCredit()) < 0.0001
}

// TestPaymentEntryFlow tests a payment entry scenario
// This is the main test you need to complete!
//
// YOUR TASK:
// 1. Create GL entries for a payment
// 2. Verify they balance
// 3. Verify specific amounts match expected values
func TestPaymentEntryFlow(t *testing.T) {
	tests := []struct {
		name           string
		paymentAmount  float64
		writeOffAmount float64
		bankAccount    string
		customerName   string
		// TODO: Add expected values
		expectedDebit  float64
		expectedCredit float64
	}{
		{
			name:           "full_payment",
			paymentAmount:  11800.00,
			writeOffAmount: 0,
			bankAccount:    "Bank - ACME",
			customerName:   "Acme Corporation",
			expectedDebit:  11800.00,
			expectedCredit: 11800.00,
		},
		{
			name:           "partial_payment",
			paymentAmount:  5000.00,
			writeOffAmount: 0,
			bankAccount:    "Bank - ACME",
			customerName:   "Acme Corporation",
			expectedDebit:  5000.00,
			expectedCredit: 5000.00,
		},
		// TODO: Add more test cases
		// Hint: What about a payment with a write-off?
		// {
		//     name:           "payment_with_writeoff",
		//     paymentAmount:  11700.00,
		//     writeOffAmount: 100.00,  // Small amount written off
		//     ...
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Create GL entries for the payment
			//
			// Payment Entry creates these GL entries:
			// 1. Bank/Cash account (Debit) - receiving money
			// 2. Debtors account (Credit) - reducing receivable
			// 3. Write-off account (Debit) if there's a write-off
			//
			// YOUR CODE HERE:
			glEntries := GLMap{
				// TODO: Add the bank entry
				// {
				//     Account: tt.bankAccount,
				//     Debit:   tt.paymentAmount,
				//     ...
				// },

				// TODO: Add the debtors entry
				// {
				//     Account:   "Debtors - ACME",
				//     Credit:    ???, // What should this be?
				//     PartyType: "Customer",
				//     Party:     tt.customerName,
				//     ...
				// },

				// TODO: Add write-off entry if writeOffAmount > 0
			}

			// Verify entries balance
			if !glEntries.IsBalanced() {
				t.Errorf("Entries don't balance: Debit=%.2f, Credit=%.2f",
					glEntries.TotalDebit(), glEntries.TotalCredit())
			}

			// Verify total amounts match expected
			if glEntries.TotalDebit() != tt.expectedDebit {
				t.Errorf("Total debit = %.2f, want %.2f",
					glEntries.TotalDebit(), tt.expectedDebit)
			}

			if glEntries.TotalCredit() != tt.expectedCredit {
				t.Errorf("Total credit = %.2f, want %.2f",
					glEntries.TotalCredit(), tt.expectedCredit)
			}

			t.Logf("✅ %s: Dr %.2f = Cr %.2f",
				tt.name, glEntries.TotalDebit(), glEntries.TotalCredit())
		})
	}
}

// TestPaymentWithWriteOff tests a payment where some amount is written off
// BONUS: Implement this test!
func TestPaymentWithWriteOff(t *testing.T) {
	// Scenario: Customer owes ₹11,800 but pays only ₹11,700
	// The ₹100 difference is written off as bad debt
	//
	// GL Entries:
	// Bank (Dr)          11,700  - money received
	// Write Off (Dr)        100  - expense for write-off
	// Debtors (Cr)       11,800  - full receivable cleared
	//
	// YOUR CODE HERE

	t.Skip("TODO: Implement write-off test")
}

// TestMultiplePaymentsToSameInvoice tests partial payments
// BONUS: Implement this test!
func TestMultiplePaymentsToSameInvoice(t *testing.T) {
	// Scenario: Invoice of ₹10,000 paid in two installments
	// Payment 1: ₹6,000
	// Payment 2: ₹4,000
	//
	// YOUR CODE HERE

	t.Skip("TODO: Implement multiple payments test")
}
