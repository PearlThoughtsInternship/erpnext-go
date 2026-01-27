package modeofpayment

import (
	"errors"
	"fmt"
	"strings"
)

// Validation errors matching ERPNext's frappe.throw() messages.
var (
	ErrDuplicateCompany = errors.New("same company is entered more than once")
	ErrAccountMismatch  = errors.New("account does not match with company")
	ErrModeInUse        = errors.New("mode of payment is used in POS profiles")
)

// ValidationError provides detailed error information.
type ValidationError struct {
	Err     error
	Details string
}

func (e *ValidationError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Err.Error(), e.Details)
	}
	return e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// ValidateRepeatingCompanies checks that no company appears multiple times
// in the accounts table.
//
// Python equivalent:
//
//	def validate_repeating_companies(self):
//	    accounts_list = [entry.company for entry in self.accounts]
//	    if len(accounts_list) != len(set(accounts_list)):
//	        frappe.throw(_("Same Company is entered more than once"))
func (m *ModeOfPayment) ValidateRepeatingCompanies() error {
	seen := make(map[string]bool)
	for _, account := range m.Accounts {
		if seen[account.Company] {
			return &ValidationError{
				Err:     ErrDuplicateCompany,
				Details: fmt.Sprintf("company '%s' appears multiple times", account.Company),
			}
		}
		seen[account.Company] = true
	}
	return nil
}

// ValidateAccounts verifies that each account's parent company matches
// the company specified in the accounts table.
//
// Python equivalent:
//
//	def validate_accounts(self):
//	    for entry in self.accounts:
//	        if frappe.get_cached_value("Account", entry.default_account, "company") != entry.company:
//	            frappe.throw(_("Account {0} does not match with Company {1}..."))
func (m *ModeOfPayment) ValidateAccounts(lookup AccountLookup) error {
	for _, account := range m.Accounts {
		if account.DefaultAccount == "" {
			continue // Skip empty accounts
		}

		accountCompany, err := lookup.GetAccountCompany(account.DefaultAccount)
		if err != nil {
			return fmt.Errorf("failed to lookup account %s: %w", account.DefaultAccount, err)
		}

		if accountCompany != account.Company {
			return &ValidationError{
				Err: ErrAccountMismatch,
				Details: fmt.Sprintf("account '%s' belongs to '%s', not '%s' in mode '%s'",
					account.DefaultAccount, accountCompany, account.Company, m.Name),
			}
		}
	}
	return nil
}

// ValidatePOSModeOfPayment prevents disabling a payment mode that is
// currently used in POS profiles.
//
// Python equivalent:
//
//	def validate_pos_mode_of_payment(self):
//	    if not self.enabled:
//	        pos_profiles = frappe.db.sql("SELECT ... WHERE mode_of_payment = %s", self.name)
//	        if pos_profiles:
//	            frappe.throw(_("POS Profile {} contains Mode of Payment {}..."))
func (m *ModeOfPayment) ValidatePOSModeOfPayment(checker POSChecker) error {
	// Only check when disabling
	if m.Enabled {
		return nil
	}

	profiles, err := checker.GetPOSProfilesUsingMode(m.Name)
	if err != nil {
		return fmt.Errorf("failed to check POS profiles: %w", err)
	}

	if len(profiles) > 0 {
		return &ValidationError{
			Err:     ErrModeInUse,
			Details: fmt.Sprintf("POS Profile '%s' contains Mode of Payment '%s'. Please remove them to disable this mode", strings.Join(profiles, ", "), m.Name),
		}
	}
	return nil
}

// Validate runs all validation checks on the Mode of Payment.
// This matches ERPNext's validate() method that calls all validation methods.
func (m *ModeOfPayment) Validate(lookup AccountLookup, checker POSChecker) error {
	if err := m.ValidateAccounts(lookup); err != nil {
		return err
	}
	if err := m.ValidateRepeatingCompanies(); err != nil {
		return err
	}
	if err := m.ValidatePOSModeOfPayment(checker); err != nil {
		return err
	}
	return nil
}
