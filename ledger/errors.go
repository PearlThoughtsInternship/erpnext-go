// errors.go defines typed errors for the ledger package.
// Using sentinel errors allows callers to check error types with errors.Is()
package ledger

import (
	"errors"
	"fmt"
)

// Sentinel errors for GL posting validation.
// These map to frappe.throw() calls in general_ledger.py
var (
	// Account validation errors
	ErrAccountDisabled = errors.New("account is disabled")
	ErrAccountFrozen   = errors.New("account is frozen")
	ErrAccountIsGroup  = errors.New("cannot post to group account")

	// Balance validation errors
	ErrDebitCreditMismatch = errors.New("debit and credit amounts do not balance")
	ErrInsufficientEntries = errors.New("incorrect number of GL entries")

	// Period validation errors
	ErrPeriodClosed        = errors.New("accounting period is closed")
	ErrFiscalYearNotFound  = errors.New("fiscal year not found for date")
	ErrAccountsFrozenTill  = errors.New("accounts frozen till date")
	ErrBooksClosedTill     = errors.New("books closed till date")

	// Budget validation errors
	ErrBudgetExceeded = errors.New("budget exceeded")

	// Currency validation errors
	ErrInvalidAccountCurrency = errors.New("invalid account currency")
	ErrCurrencyMismatch       = errors.New("currency mismatch")

	// Voucher validation errors
	ErrVoucherNotFound    = errors.New("voucher not found")
	ErrVoucherAlreadyPosted = errors.New("voucher already has GL entries")
)

// ValidationError wraps a sentinel error with additional context.
// This provides both programmatic error checking and user-friendly messages.
type ValidationError struct {
	Err     error  // Underlying sentinel error
	Account string // Account involved (if applicable)
	Details string // User-friendly explanation
}

func (e *ValidationError) Error() string {
	if e.Account != "" {
		return fmt.Sprintf("%s: %s - %s", e.Err.Error(), e.Account, e.Details)
	}
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Err.Error(), e.Details)
	}
	return e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a validation error with context.
func NewValidationError(err error, account, details string) *ValidationError {
	return &ValidationError{
		Err:     err,
		Account: account,
		Details: details,
	}
}

// DisabledAccountsError lists multiple disabled accounts.
// Maps to the error format in validate_disabled_accounts()
type DisabledAccountsError struct {
	Accounts []string
}

func (e *DisabledAccountsError) Error() string {
	return fmt.Sprintf("cannot create accounting entries against disabled accounts: %v", e.Accounts)
}

func (e *DisabledAccountsError) Unwrap() error {
	return ErrAccountDisabled
}

// PeriodClosedError provides details about a closed accounting period.
type PeriodClosedError struct {
	Company     string
	DocType     string
	PostingDate string
	PeriodName  string
}

func (e *PeriodClosedError) Error() string {
	return fmt.Sprintf(
		"accounting period is closed for %s in %s on %s (Period: %s)",
		e.DocType, e.Company, e.PostingDate, e.PeriodName,
	)
}

func (e *PeriodClosedError) Unwrap() error {
	return ErrPeriodClosed
}

// BudgetExceededError provides details about budget violation.
type BudgetExceededError struct {
	Account    string
	CostCenter string
	Budget     float64
	Actual     float64
	Variance   float64
}

func (e *BudgetExceededError) Error() string {
	return fmt.Sprintf(
		"budget exceeded for %s in %s: budget %.2f, actual %.2f, over by %.2f",
		e.Account, e.CostCenter, e.Budget, e.Actual, e.Variance,
	)
}

func (e *BudgetExceededError) Unwrap() error {
	return ErrBudgetExceeded
}

// GLEntryCountError indicates wrong number of GL entries.
// A valid transaction needs at least 2 entries (debit and credit sides).
type GLEntryCountError struct {
	Expected int
	Actual   int
	Message  string
}

func (e *GLEntryCountError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("expected at least %d GL entries, got %d", e.Expected, e.Actual)
}

func (e *GLEntryCountError) Unwrap() error {
	return ErrInsufficientEntries
}
