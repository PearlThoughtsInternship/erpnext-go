// Package ledger implements the General Ledger posting engine from ERPNext.
// Migrated from: erpnext/accounts/general_ledger.py
//
// This package is the foundation of all accounting in ERPNext - every transaction
// (Sales Invoice, Purchase Invoice, Journal Entry, Payment Entry) posts through
// this engine to create GL entries for double-entry bookkeeping.
package ledger

import (
	"time"
)

// IsOpeningEntry defines whether a GL entry is an opening balance entry.
type IsOpeningEntry string

const (
	IsOpeningNo  IsOpeningEntry = "No"
	IsOpeningYes IsOpeningEntry = "Yes"
)

// IsAdvanceEntry defines whether a GL entry is an advance payment entry.
type IsAdvanceEntry string

const (
	IsAdvanceNo  IsAdvanceEntry = "No"
	IsAdvanceYes IsAdvanceEntry = "Yes"
)

// GLEntry represents a single General Ledger entry.
// This is the atomic unit of accounting - every financial transaction
// creates one or more GL entries that must balance (total debit = total credit).
//
// Maps to: erpnext/accounts/doctype/gl_entry/gl_entry.json
// ERPNext naming: ACC-GLE-.YYYY.-.#####
type GLEntry struct {
	// Identity
	Name string // Document name (auto-generated)

	// Dates
	PostingDate     time.Time  // Date when entry is recorded
	TransactionDate time.Time  // Original transaction date
	DueDate         *time.Time // Payment due date (for receivables/payables)

	// Account Details
	Account         string // Ledger account (e.g., "Sales - ABC")
	AccountCurrency string // Account's designated currency

	// Party (Customer/Supplier)
	PartyType string // "Customer", "Supplier", "Employee", etc.
	Party     string // Party name

	// Counter account for journal display
	Against string // Comma-separated list of counter accounts

	// Voucher (Source Document)
	VoucherType    string // "Sales Invoice", "Journal Entry", etc.
	VoucherNo      string // Document number
	VoucherSubtype string // Additional classification
	VoucherDetailNo string // Line item reference

	// Against Voucher (for AR/AP matching)
	AgainstVoucherType string // Linked voucher type
	AgainstVoucher     string // Linked voucher number

	// Amounts in Company Currency (reporting currency)
	Debit  float64 // Debit in company currency
	Credit float64 // Credit in company currency

	// Amounts in Account Currency
	DebitInAccountCurrency  float64
	CreditInAccountCurrency float64

	// Amounts in Transaction Currency (customer/supplier sees this)
	TransactionCurrency            string
	TransactionExchangeRate        float64
	DebitInTransactionCurrency     float64
	CreditInTransactionCurrency    float64

	// Amounts in Reporting Currency (for multi-currency consolidation)
	ReportingCurrencyExchangeRate float64
	DebitInReportingCurrency      float64
	CreditInReportingCurrency     float64

	// Dimensions
	CostCenter string // Cost center for expense/revenue analysis
	Project    string // Project for project accounting

	// Classification
	Company     string         // Company this entry belongs to
	FiscalYear  string         // Fiscal year reference
	FinanceBook string         // Finance book for parallel ledgers
	IsOpening   IsOpeningEntry // Opening balance entry flag
	IsAdvance   IsAdvanceEntry // Advance payment flag
	IsCancelled bool           // Cancellation flag
	Remarks     string         // Free-text remarks

	// Internal flags (not persisted, used during processing)
	ToRename bool // Flag for renaming logic
}

// VoucherRef identifies a source document for GL operations.
type VoucherRef struct {
	VoucherType string
	VoucherNo   string
	Company     string
}

// PostingOptions controls GL posting behavior.
// Maps to: function parameters in make_gl_entries()
type PostingOptions struct {
	Cancel            bool   // True if cancelling/reversing entries
	AdvAdj            bool   // Advance adjustment flag
	MergeEntries      bool   // Merge similar GL entries
	UpdateOutstanding string // "Yes" or "No" - update AR/AP outstanding
	FromRepost        bool   // True if reposting (e.g., valuation change)
}

// DefaultPostingOptions returns standard posting options.
func DefaultPostingOptions() PostingOptions {
	return PostingOptions{
		Cancel:            false,
		AdvAdj:            false,
		MergeEntries:      true,
		UpdateOutstanding: "Yes",
		FromRepost:        false,
	}
}

// PaymentLedgerEntry represents an AR/AP ledger entry.
// This is a separate ledger from GL that tracks customer/supplier balances
// with more granular linking for payment allocation.
//
// Maps to: erpnext/accounts/doctype/payment_ledger_entry/
type PaymentLedgerEntry struct {
	Name string

	PostingDate time.Time
	Company     string
	Account     string

	// Party information
	PartyType string
	Party     string

	// Voucher references
	VoucherType    string
	VoucherNo      string
	VoucherDetailNo string

	// Against voucher (for matching)
	AgainstVoucherType string
	AgainstVoucherNo   string

	// Amounts
	AccountCurrency        string
	Amount                 float64 // In company currency
	AmountInAccountCurrency float64
	DueDate                *time.Time

	// Finance book
	FinanceBook string
	Delinked    bool
}

// GLMap is a slice of GL entries that form a complete transaction.
// In double-entry bookkeeping, the sum of debits must equal sum of credits.
type GLMap []GLEntry

// TotalDebit returns the sum of all debit amounts.
func (m GLMap) TotalDebit() float64 {
	var total float64
	for _, e := range m {
		total += e.Debit
	}
	return total
}

// TotalCredit returns the sum of all credit amounts.
func (m GLMap) TotalCredit() float64 {
	var total float64
	for _, e := range m {
		total += e.Credit
	}
	return total
}

// IsBalanced returns true if total debits equal total credits.
func (m GLMap) IsBalanced() bool {
	return Flt(m.TotalDebit()-m.TotalCredit(), 2) == 0
}

// Copy creates a deep copy of a GL entry.
func (e *GLEntry) Copy() GLEntry {
	copy := *e
	if e.DueDate != nil {
		d := *e.DueDate
		copy.DueDate = &d
	}
	return copy
}

// Flt converts to float and optionally rounds.
// Maps to: frappe.utils.flt() in Python
func Flt(value float64, precision ...int) float64 {
	if len(precision) > 0 {
		return Round(value, precision[0])
	}
	return value
}

// Round rounds a value to the specified precision.
func Round(value float64, precision int) float64 {
	if precision < 0 {
		return value
	}
	multiplier := 1.0
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}
	return float64(int64(value*multiplier+0.5)) / multiplier
}
