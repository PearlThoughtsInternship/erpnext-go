// Package modeofpayment implements the Mode of Payment doctype from ERPNext.
// Migrated from: erpnext/accounts/doctype/mode_of_payment/mode_of_payment.py
package modeofpayment

// PaymentType represents the type of payment method.
type PaymentType string

const (
	Cash    PaymentType = "Cash"
	Bank    PaymentType = "Bank"
	General PaymentType = "General"
	Phone   PaymentType = "Phone"
)

// ModeOfPaymentAccount represents a child table row linking
// a company to its default ledger account for this payment mode.
type ModeOfPaymentAccount struct {
	Company        string
	DefaultAccount string
}

// ModeOfPayment represents a payment method master record.
// Maps to ERPNext's Mode of Payment doctype.
type ModeOfPayment struct {
	Name     string
	Type     PaymentType
	Enabled  bool
	Accounts []ModeOfPaymentAccount
}

// AccountLookup abstracts database queries for account information.
// Production implementations query the Account doctype.
type AccountLookup interface {
	// GetAccountCompany returns the company that owns the given account.
	GetAccountCompany(accountName string) (string, error)
}

// POSChecker abstracts database queries for POS profile information.
// Production implementations query Sales Invoice Payment records.
type POSChecker interface {
	// GetPOSProfilesUsingMode returns POS profile names that use this payment mode.
	GetPOSProfilesUsingMode(modeName string) ([]string, error)
}
