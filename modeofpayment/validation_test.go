package modeofpayment

import (
	"errors"
	"testing"
)

// --- Mock implementations for testing ---

// mockAccountLookup simulates database queries for Account records.
type mockAccountLookup struct {
	// accounts maps account name -> company name
	accounts map[string]string
}

func (m *mockAccountLookup) GetAccountCompany(accountName string) (string, error) {
	company, ok := m.accounts[accountName]
	if !ok {
		return "", errors.New("account not found")
	}
	return company, nil
}

// mockPOSChecker simulates database queries for POS profiles.
type mockPOSChecker struct {
	// profilesByMode maps mode name -> list of POS profile names
	profilesByMode map[string][]string
}

func (m *mockPOSChecker) GetPOSProfilesUsingMode(modeName string) ([]string, error) {
	profiles, ok := m.profilesByMode[modeName]
	if !ok {
		return []string{}, nil
	}
	return profiles, nil
}

// --- Tests ---

func TestValidateRepeatingCompanies(t *testing.T) {
	tests := []struct {
		name     string
		accounts []ModeOfPaymentAccount
		wantErr  error
	}{
		{
			name:     "empty accounts - valid",
			accounts: []ModeOfPaymentAccount{},
			wantErr:  nil,
		},
		{
			name: "single company - valid",
			accounts: []ModeOfPaymentAccount{
				{Company: "Company A", DefaultAccount: "Cash - A"},
			},
			wantErr: nil,
		},
		{
			name: "unique companies - valid",
			accounts: []ModeOfPaymentAccount{
				{Company: "Company A", DefaultAccount: "Cash - A"},
				{Company: "Company B", DefaultAccount: "Cash - B"},
				{Company: "Company C", DefaultAccount: "Cash - C"},
			},
			wantErr: nil,
		},
		{
			name: "duplicate companies - error",
			accounts: []ModeOfPaymentAccount{
				{Company: "Company A", DefaultAccount: "Cash - A"},
				{Company: "Company A", DefaultAccount: "Bank - A"},
			},
			wantErr: ErrDuplicateCompany,
		},
		{
			name: "duplicate among many - error",
			accounts: []ModeOfPaymentAccount{
				{Company: "Company A", DefaultAccount: "Cash - A"},
				{Company: "Company B", DefaultAccount: "Cash - B"},
				{Company: "Company A", DefaultAccount: "Bank - A"},
			},
			wantErr: ErrDuplicateCompany,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModeOfPayment{
				Name:     "Test Mode",
				Accounts: tt.accounts,
			}

			err := m.ValidateRepeatingCompanies()

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got: %v", tt.wantErr, err)
				}
			}
		})
	}
}

func TestValidateAccounts(t *testing.T) {
	// Setup mock account data: maps account name to its owning company
	lookup := &mockAccountLookup{
		accounts: map[string]string{
			"Cash - Company A": "Company A",
			"Bank - Company A": "Company A",
			"Cash - Company B": "Company B",
		},
	}

	tests := []struct {
		name     string
		mode     *ModeOfPayment
		wantErr  error
	}{
		{
			name: "empty accounts - valid",
			mode: &ModeOfPayment{
				Name:     "Credit Card",
				Accounts: []ModeOfPaymentAccount{},
			},
			wantErr: nil,
		},
		{
			name: "account matches company - valid",
			mode: &ModeOfPayment{
				Name: "Cash",
				Accounts: []ModeOfPaymentAccount{
					{Company: "Company A", DefaultAccount: "Cash - Company A"},
				},
			},
			wantErr: nil,
		},
		{
			name: "multiple accounts all match - valid",
			mode: &ModeOfPayment{
				Name: "Cash",
				Accounts: []ModeOfPaymentAccount{
					{Company: "Company A", DefaultAccount: "Cash - Company A"},
					{Company: "Company B", DefaultAccount: "Cash - Company B"},
				},
			},
			wantErr: nil,
		},
		{
			name: "account company mismatch - error",
			mode: &ModeOfPayment{
				Name: "Cash",
				Accounts: []ModeOfPaymentAccount{
					{Company: "Company B", DefaultAccount: "Cash - Company A"}, // Account belongs to A, not B
				},
			},
			wantErr: ErrAccountMismatch,
		},
		{
			name: "empty default account - skipped",
			mode: &ModeOfPayment{
				Name: "Cash",
				Accounts: []ModeOfPaymentAccount{
					{Company: "Company A", DefaultAccount: ""}, // Empty, should skip
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mode.ValidateAccounts(lookup)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got: %v", tt.wantErr, err)
				}
			}
		})
	}
}

func TestValidatePOSModeOfPayment(t *testing.T) {
	// Setup mock POS profile data
	checker := &mockPOSChecker{
		profilesByMode: map[string][]string{
			"Cash":        {"Retail POS", "Restaurant POS"},
			"Credit Card": {"Retail POS"},
		},
	}

	tests := []struct {
		name    string
		mode    *ModeOfPayment
		wantErr error
	}{
		{
			name: "enabled mode - always valid",
			mode: &ModeOfPayment{
				Name:    "Cash",
				Enabled: true,
			},
			wantErr: nil,
		},
		{
			name: "disabled, not in POS - valid",
			mode: &ModeOfPayment{
				Name:    "Wire Transfer",
				Enabled: false,
			},
			wantErr: nil,
		},
		{
			name: "disabled, used in POS - error",
			mode: &ModeOfPayment{
				Name:    "Cash",
				Enabled: false,
			},
			wantErr: ErrModeInUse,
		},
		{
			name: "disabled, used in one POS - error",
			mode: &ModeOfPayment{
				Name:    "Credit Card",
				Enabled: false,
			},
			wantErr: ErrModeInUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mode.ValidatePOSModeOfPayment(checker)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got: %v", tt.wantErr, err)
				}
			}
		})
	}
}

func TestValidate_Integration(t *testing.T) {
	// Full validation integration test
	lookup := &mockAccountLookup{
		accounts: map[string]string{
			"Cash - A": "Company A",
			"Cash - B": "Company B",
		},
	}
	checker := &mockPOSChecker{
		profilesByMode: map[string][]string{
			"Cash": {"Retail POS"},
		},
	}

	tests := []struct {
		name    string
		mode    *ModeOfPayment
		wantErr error
	}{
		{
			name: "valid mode - all checks pass",
			mode: &ModeOfPayment{
				Name:    "Credit Card",
				Type:    Bank,
				Enabled: true,
				Accounts: []ModeOfPaymentAccount{
					{Company: "Company A", DefaultAccount: "Cash - A"},
					{Company: "Company B", DefaultAccount: "Cash - B"},
				},
			},
			wantErr: nil,
		},
		{
			name: "fails on duplicate company",
			mode: &ModeOfPayment{
				Name:    "Test",
				Enabled: true,
				Accounts: []ModeOfPaymentAccount{
					{Company: "Company A", DefaultAccount: "Cash - A"},
					{Company: "Company A", DefaultAccount: "Cash - A"},
				},
			},
			wantErr: ErrDuplicateCompany,
		},
		{
			name: "fails on account mismatch",
			mode: &ModeOfPayment{
				Name:    "Test",
				Enabled: true,
				Accounts: []ModeOfPaymentAccount{
					{Company: "Company B", DefaultAccount: "Cash - A"}, // Wrong company
				},
			},
			wantErr: ErrAccountMismatch,
		},
		{
			name: "fails on POS in use",
			mode: &ModeOfPayment{
				Name:     "Cash",
				Enabled:  false,
				Accounts: []ModeOfPaymentAccount{},
			},
			wantErr: ErrModeInUse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mode.Validate(lookup, checker)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got: %v", tt.wantErr, err)
				}
			}
		})
	}
}
