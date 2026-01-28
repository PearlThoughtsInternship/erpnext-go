package glentry

import "time"

// GLEntry represents a General Ledger Entry in the accounting system.
// This mirrors ERPNext's GL Entry DocType.
//
// Every financial transaction creates GL entries that must balance:
// Total Debits = Total Credits
//
// YOUR TASK: Complete the missing fields marked with TODO comments.

type GLEntry struct {
	// ============================================
	// SECTION 1: Core Identification Fields
	// ============================================

	// Name is the unique identifier for this entry (e.g., "GL-00001")
	Name string

	// TODO: Add a field called "Account" of type string
	// This is the ledger account name (e.g., "Sales - ACME", "Debtors - ACME")
	// YOUR CODE HERE

	// TODO: Add a field called "Company" of type string
	// This is the company this entry belongs to (e.g., "ACME Corp")
	// YOUR CODE HERE

	// ============================================
	// SECTION 2: Debit and Credit Amounts
	// ============================================

	// TODO: Add a field called "Debit" of type float64
	// This is the debit amount in company currency
	// YOUR CODE HERE

	// TODO: Add a field called "Credit" of type float64
	// This is the credit amount in company currency
	// YOUR CODE HERE

	// ============================================
	// SECTION 3: Voucher Reference
	// ============================================

	// The voucher is the source document (Sales Invoice, Payment, etc.)

	// TODO: Add a field called "VoucherType" of type string
	// Examples: "Sales Invoice", "Payment Entry", "Journal Entry"
	// YOUR CODE HERE

	// TODO: Add a field called "VoucherNo" of type string
	// Examples: "SINV-2024-00001", "PAY-2024-00001"
	// YOUR CODE HERE

	// ============================================
	// SECTION 4: Party Information
	// ============================================

	// Party = Customer or Supplier involved in the transaction

	// TODO: Add a field called "PartyType" of type string
	// Examples: "Customer", "Supplier", "Employee"
	// YOUR CODE HERE

	// TODO: Add a field called "Party" of type string
	// This is the actual party name (e.g., "Acme Corporation")
	// YOUR CODE HERE

	// ============================================
	// SECTION 5: Dates
	// ============================================

	// TODO: Add a field called "PostingDate" of type time.Time
	// This is when the transaction is recorded in the books
	// YOUR CODE HERE

	// TODO: Add a field called "CreationDate" of type time.Time
	// This is when this GL entry record was created
	// YOUR CODE HERE

	// ============================================
	// SECTION 6: Cost Center (for cost accounting)
	// ============================================

	// TODO: Add a field called "CostCenter" of type string
	// Examples: "Main - ACME", "North Region - ACME"
	// YOUR CODE HERE

	// ============================================
	// SECTION 7: Status Flags
	// ============================================

	// TODO: Add a field called "IsCancelled" of type bool
	// True if this entry has been cancelled/reversed
	// YOUR CODE HERE

	// TODO: Add a field called "IsOpening" of type string
	// Values: "Yes", "No" - indicates opening balance entry
	// YOUR CODE HERE
}

// NewGLEntry creates a GLEntry with sensible defaults.
// This is a constructor pattern common in Go.
//
// TODO: Complete this function by setting default values.
func NewGLEntry(account string, debit, credit float64) GLEntry {
	return GLEntry{
		// TODO: Set the Account field to the account parameter
		// YOUR CODE HERE

		// TODO: Set the Debit field to the debit parameter
		// YOUR CODE HERE

		// TODO: Set the Credit field to the credit parameter
		// YOUR CODE HERE

		// TODO: Set PostingDate to time.Now()
		// YOUR CODE HERE

		// TODO: Set CreationDate to time.Now()
		// YOUR CODE HERE

		// TODO: Set IsCancelled to false
		// YOUR CODE HERE

		// TODO: Set IsOpening to "No"
		// YOUR CODE HERE
	}
}

// IsBalanced checks if a single entry has either debit OR credit (not both).
// In double-entry accounting, a single entry should only affect one side.
//
// TODO: Implement this function
func (e GLEntry) IsValid() bool {
	// An entry is valid if:
	// 1. It has a debit OR a credit (not both non-zero)
	// 2. The account is not empty

	// TODO: Return true if the entry is valid, false otherwise
	// Hint: (Debit > 0 && Credit == 0) || (Debit == 0 && Credit > 0)
	// Also check: Account != ""

	// YOUR CODE HERE
	return false // Replace this
}
