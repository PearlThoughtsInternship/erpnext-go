// ports.go defines interfaces for external dependencies.
// Following hexagonal architecture, these ports abstract database access
// and allow easy testing with mock implementations.
package ledger

import "time"

// AccountLookup abstracts queries for Account master data.
// Production implementations query the Account doctype.
//
// Maps to: frappe.get_cached_value("Account", ...) calls in general_ledger.py
type AccountLookup interface {
	// GetAccount returns full account details.
	GetAccount(name string) (*Account, error)

	// GetAccountCurrency returns the account's designated currency.
	GetAccountCurrency(name string) (string, error)

	// IsGroup returns true if the account is a group (parent) account.
	// GL entries cannot be posted to group accounts.
	IsGroup(name string) (bool, error)

	// IsFrozen returns true if the account is frozen for posting.
	IsFrozen(name string) (bool, error)

	// IsDisabled returns true if the account is disabled.
	IsDisabled(name string) (bool, error)

	// GetBalanceMustBe returns "Debit", "Credit", or "" for balance constraint.
	GetBalanceMustBe(name string) (string, error)
}

// Account represents an account from the Chart of Accounts.
// Simplified view for ledger operations.
type Account struct {
	Name            string
	AccountName     string
	Company         string
	AccountCurrency string
	IsGroup         bool
	Disabled        bool
	FreezeAccount   bool
	BalanceMustBe   string // "Debit", "Credit", or ""
	RootType        string // "Asset", "Liability", "Equity", "Income", "Expense"
}

// CompanySettings abstracts company-level accounting configuration.
// Maps to: frappe.get_cached_value("Company", ...) and related calls
type CompanySettings interface {
	// GetDefaultCurrency returns the company's base currency.
	GetDefaultCurrency(company string) (string, error)

	// GetRoundOffAccount returns the round-off account for the company.
	GetRoundOffAccount(company string) (string, error)

	// GetRoundOffCostCenter returns the cost center for round-off entries.
	GetRoundOffCostCenter(company string) (string, error)

	// GetAccountsFrozenTillDate returns the date until which accounts are frozen.
	// Returns nil if no freeze is set.
	GetAccountsFrozenTillDate(company string) (*time.Time, error)

	// GetBookClosingDate returns the date until which books are closed.
	// Returns nil if no closing is set.
	GetBookClosingDate(company string) (*time.Time, error)
}

// AccountingPeriodChecker validates posting against accounting periods.
// ERPNext supports closing specific document types for specific periods.
//
// Maps to: validate_accounting_period() in general_ledger.py
type AccountingPeriodChecker interface {
	// IsDocumentTypeClosed returns true if the document type is closed
	// for the given company and posting date.
	IsDocumentTypeClosed(company, docType string, postingDate time.Time) (bool, error)

	// GetClosedPeriodMessage returns user-friendly error message if closed.
	GetClosedPeriodMessage(company, docType string, postingDate time.Time) (string, error)
}

// FiscalYearLookup resolves fiscal year for a given date.
// Maps to: get_fiscal_year() in accounts/utils.py
type FiscalYearLookup interface {
	// GetFiscalYear returns the fiscal year name for the given date and company.
	GetFiscalYear(date time.Time, company string) (string, error)

	// GetFiscalYearDates returns start and end dates for a fiscal year.
	GetFiscalYearDates(fiscalYear string, company string) (start, end time.Time, err error)
}

// GLEntryStore abstracts GL entry persistence.
// Maps to: save_entries() and related functions in general_ledger.py
type GLEntryStore interface {
	// Save persists a GL entry to the database.
	Save(entry *GLEntry) error

	// SaveBatch persists multiple GL entries in a transaction.
	SaveBatch(entries []GLEntry) error

	// GetByVoucher retrieves all GL entries for a voucher.
	GetByVoucher(voucherType, voucherNo string) ([]GLEntry, error)

	// MarkCancelled marks all entries for a voucher as cancelled.
	MarkCancelled(voucherType, voucherNo string) error
}

// PaymentLedgerStore abstracts payment ledger entry persistence.
// Maps to: create_payment_ledger_entry() in accounts/utils.py
type PaymentLedgerStore interface {
	// Save persists a payment ledger entry.
	Save(entry *PaymentLedgerEntry) error

	// SaveBatch persists multiple payment ledger entries.
	SaveBatch(entries []PaymentLedgerEntry) error

	// GetByVoucher retrieves payment ledger entries for a voucher.
	GetByVoucher(voucherType, voucherNo string) ([]PaymentLedgerEntry, error)

	// Delink marks payment ledger entries as delinked (for cancellation).
	Delink(voucherType, voucherNo string) error
}

// BudgetValidator validates GL entries against budgets.
// Maps to: BudgetValidation class in controllers/budget_controller.py
type BudgetValidator interface {
	// Validate checks if GL entries violate any budget constraints.
	// Returns nil if validation passes, error with details if budget exceeded.
	Validate(entries []GLEntry) error
}

// AccountingDimensionProvider retrieves accounting dimensions for offsetting.
// Maps to: get_accounting_dimensions_for_offsetting_entry() in general_ledger.py
type AccountingDimensionProvider interface {
	// GetDimensionsForOffsetting returns dimensions that need offsetting entries.
	GetDimensionsForOffsetting(glMap []GLEntry, company string) ([]AccountingDimension, error)
}

// AccountingDimension represents a dimension that requires offsetting entries.
type AccountingDimension struct {
	Fieldname        string // Field name in GL entry (e.g., "cost_center")
	Name             string // Dimension name
	OffsettingAccount string // Account for offsetting entries
	AccountCurrency  string // Currency of the offsetting account
}

// Engine combines all ports needed for GL posting.
// This is the main dependency injection point for the ledger engine.
type Engine struct {
	Accounts          AccountLookup
	Company           CompanySettings
	Periods           AccountingPeriodChecker
	FiscalYears       FiscalYearLookup
	GLStore           GLEntryStore
	PaymentStore      PaymentLedgerStore
	Budget            BudgetValidator
	Dimensions        AccountingDimensionProvider
}

// NewEngine creates a new ledger engine with all dependencies.
func NewEngine(
	accounts AccountLookup,
	company CompanySettings,
	periods AccountingPeriodChecker,
	fiscalYears FiscalYearLookup,
	glStore GLEntryStore,
	paymentStore PaymentLedgerStore,
	budget BudgetValidator,
	dimensions AccountingDimensionProvider,
) *Engine {
	return &Engine{
		Accounts:     accounts,
		Company:      company,
		Periods:      periods,
		FiscalYears:  fiscalYears,
		GLStore:      glStore,
		PaymentStore: paymentStore,
		Budget:       budget,
		Dimensions:   dimensions,
	}
}
