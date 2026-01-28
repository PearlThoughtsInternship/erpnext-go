package ledger

import (
	"errors"
	"testing"
	"time"
)

// Mock implementations for testing

type mockAccountLookup struct {
	accounts map[string]*Account
}

func newMockAccountLookup() *mockAccountLookup {
	return &mockAccountLookup{
		accounts: map[string]*Account{
			"Sales - ABC": {
				Name:            "Sales - ABC",
				AccountName:     "Sales",
				Company:         "ABC Company",
				AccountCurrency: "USD",
				IsGroup:         false,
				Disabled:        false,
			},
			"Debtors - ABC": {
				Name:            "Debtors - ABC",
				AccountName:     "Debtors",
				Company:         "ABC Company",
				AccountCurrency: "USD",
				IsGroup:         false,
				Disabled:        false,
			},
			"Cash - ABC": {
				Name:            "Cash - ABC",
				AccountName:     "Cash",
				Company:         "ABC Company",
				AccountCurrency: "USD",
				IsGroup:         false,
				Disabled:        false,
			},
			"Disabled Account - ABC": {
				Name:            "Disabled Account - ABC",
				AccountName:     "Disabled Account",
				Company:         "ABC Company",
				AccountCurrency: "USD",
				IsGroup:         false,
				Disabled:        true,
			},
		},
	}
}

func (m *mockAccountLookup) GetAccount(name string) (*Account, error) {
	if acc, ok := m.accounts[name]; ok {
		return acc, nil
	}
	return nil, errors.New("account not found")
}

func (m *mockAccountLookup) GetAccountCurrency(name string) (string, error) {
	if acc, ok := m.accounts[name]; ok {
		return acc.AccountCurrency, nil
	}
	return "", errors.New("account not found")
}

func (m *mockAccountLookup) IsGroup(name string) (bool, error) {
	if acc, ok := m.accounts[name]; ok {
		return acc.IsGroup, nil
	}
	return false, errors.New("account not found")
}

func (m *mockAccountLookup) IsFrozen(name string) (bool, error) {
	if acc, ok := m.accounts[name]; ok {
		return acc.FreezeAccount, nil
	}
	return false, errors.New("account not found")
}

func (m *mockAccountLookup) IsDisabled(name string) (bool, error) {
	if acc, ok := m.accounts[name]; ok {
		return acc.Disabled, nil
	}
	return false, errors.New("account not found")
}

func (m *mockAccountLookup) GetBalanceMustBe(name string) (string, error) {
	if acc, ok := m.accounts[name]; ok {
		return acc.BalanceMustBe, nil
	}
	return "", errors.New("account not found")
}

type mockGLStore struct {
	entries []GLEntry
}

func (m *mockGLStore) Save(entry *GLEntry) error {
	m.entries = append(m.entries, *entry)
	return nil
}

func (m *mockGLStore) SaveBatch(entries []GLEntry) error {
	m.entries = append(m.entries, entries...)
	return nil
}

func (m *mockGLStore) GetByVoucher(voucherType, voucherNo string) ([]GLEntry, error) {
	var result []GLEntry
	for _, e := range m.entries {
		if e.VoucherType == voucherType && e.VoucherNo == voucherNo {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockGLStore) MarkCancelled(voucherType, voucherNo string) error {
	for i := range m.entries {
		if m.entries[i].VoucherType == voucherType && m.entries[i].VoucherNo == voucherNo {
			m.entries[i].IsCancelled = true
		}
	}
	return nil
}

type mockCompanySettings struct{}

func (m *mockCompanySettings) GetDefaultCurrency(company string) (string, error) {
	return "USD", nil
}

func (m *mockCompanySettings) GetRoundOffAccount(company string) (string, error) {
	return "Round Off - ABC", nil
}

func (m *mockCompanySettings) GetRoundOffCostCenter(company string) (string, error) {
	return "Main - ABC", nil
}

func (m *mockCompanySettings) GetAccountsFrozenTillDate(company string) (*time.Time, error) {
	return nil, nil
}

func (m *mockCompanySettings) GetBookClosingDate(company string) (*time.Time, error) {
	return nil, nil
}

// Test helper functions

func makeTestDate() time.Time {
	return time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
}

func makeTestGLEntry(account string, debit, credit float64) GLEntry {
	return GLEntry{
		PostingDate:     makeTestDate(),
		Account:         account,
		Debit:           debit,
		Credit:          credit,
		DebitInAccountCurrency:  debit,
		CreditInAccountCurrency: credit,
		Company:         "ABC Company",
		VoucherType:     "Sales Invoice",
		VoucherNo:       "SINV-001",
		AccountCurrency: "USD",
	}
}

// Tests for GLMap methods

func TestGLMap_TotalDebit(t *testing.T) {
	tests := []struct {
		name     string
		entries  GLMap
		expected float64
	}{
		{
			name:     "empty map",
			entries:  GLMap{},
			expected: 0,
		},
		{
			name: "single entry",
			entries: GLMap{
				makeTestGLEntry("Sales - ABC", 100, 0),
			},
			expected: 100,
		},
		{
			name: "multiple entries",
			entries: GLMap{
				makeTestGLEntry("Sales - ABC", 100, 0),
				makeTestGLEntry("Debtors - ABC", 50, 0),
				makeTestGLEntry("Cash - ABC", 0, 150),
			},
			expected: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entries.TotalDebit()
			if got != tt.expected {
				t.Errorf("TotalDebit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGLMap_TotalCredit(t *testing.T) {
	tests := []struct {
		name     string
		entries  GLMap
		expected float64
	}{
		{
			name:     "empty map",
			entries:  GLMap{},
			expected: 0,
		},
		{
			name: "single entry",
			entries: GLMap{
				makeTestGLEntry("Sales - ABC", 0, 100),
			},
			expected: 100,
		},
		{
			name: "multiple entries",
			entries: GLMap{
				makeTestGLEntry("Sales - ABC", 0, 100),
				makeTestGLEntry("Debtors - ABC", 0, 50),
				makeTestGLEntry("Cash - ABC", 150, 0),
			},
			expected: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entries.TotalCredit()
			if got != tt.expected {
				t.Errorf("TotalCredit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGLMap_IsBalanced(t *testing.T) {
	tests := []struct {
		name     string
		entries  GLMap
		expected bool
	}{
		{
			name:     "empty map is balanced",
			entries:  GLMap{},
			expected: true,
		},
		{
			name: "balanced entries",
			entries: GLMap{
				makeTestGLEntry("Debtors - ABC", 100, 0),
				makeTestGLEntry("Sales - ABC", 0, 100),
			},
			expected: true,
		},
		{
			name: "unbalanced entries",
			entries: GLMap{
				makeTestGLEntry("Debtors - ABC", 100, 0),
				makeTestGLEntry("Sales - ABC", 0, 90),
			},
			expected: false,
		},
		{
			name: "balanced with rounding",
			entries: GLMap{
				makeTestGLEntry("Debtors - ABC", 100.001, 0),
				makeTestGLEntry("Sales - ABC", 0, 100.005),
			},
			expected: true, // Difference of 0.004 rounds to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entries.IsBalanced()
			if got != tt.expected {
				t.Errorf("IsBalanced() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Tests for MergeSimilarEntries

func TestMergeSimilarEntries(t *testing.T) {
	tests := []struct {
		name           string
		entries        []GLEntry
		expectedCount  int
		expectedDebits []float64
	}{
		{
			name:           "empty entries",
			entries:        []GLEntry{},
			expectedCount:  0,
			expectedDebits: []float64{},
		},
		{
			name: "no merge needed - different accounts",
			entries: []GLEntry{
				makeTestGLEntry("Sales - ABC", 100, 0),
				makeTestGLEntry("Debtors - ABC", 0, 100),
			},
			expectedCount:  2,
			expectedDebits: []float64{100, 0},
		},
		{
			name: "merge same account",
			entries: []GLEntry{
				makeTestGLEntry("Sales - ABC", 50, 0),
				makeTestGLEntry("Sales - ABC", 50, 0),
			},
			expectedCount:  1,
			expectedDebits: []float64{100},
		},
		{
			name: "merge multiple same accounts",
			entries: []GLEntry{
				makeTestGLEntry("Sales - ABC", 30, 0),
				makeTestGLEntry("Debtors - ABC", 0, 100),
				makeTestGLEntry("Sales - ABC", 40, 0),
				makeTestGLEntry("Sales - ABC", 30, 0),
			},
			expectedCount:  2,
			expectedDebits: []float64{100, 0}, // Sales merged, Debtors separate
		},
		{
			name: "filter zero entries",
			entries: []GLEntry{
				makeTestGLEntry("Sales - ABC", 100, 0),
				makeTestGLEntry("Debtors - ABC", 0, 0), // Zero entry
				makeTestGLEntry("Cash - ABC", 0, 100),
			},
			expectedCount:  2,
			expectedDebits: []float64{100, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeSimilarEntries(tt.entries)

			if len(result) != tt.expectedCount {
				t.Errorf("MergeSimilarEntries() count = %d, want %d", len(result), tt.expectedCount)
			}

			for i, expectedDebit := range tt.expectedDebits {
				if i < len(result) && result[i].Debit != expectedDebit {
					t.Errorf("Entry[%d].Debit = %v, want %v", i, result[i].Debit, expectedDebit)
				}
			}
		})
	}
}

// Tests for ToggleDebitCreditIfNegative

func TestToggleDebitCreditIfNegative(t *testing.T) {
	tests := []struct {
		name           string
		debit          float64
		credit         float64
		expectedDebit  float64
		expectedCredit float64
	}{
		{
			name:           "positive values unchanged",
			debit:          100,
			credit:         0,
			expectedDebit:  100,
			expectedCredit: 0,
		},
		{
			name:           "negative debit moved to credit",
			debit:          -100,
			credit:         0,
			expectedDebit:  0,
			expectedCredit: 100,
		},
		{
			name:           "negative credit moved to debit",
			debit:          0,
			credit:         -100,
			expectedDebit:  100,
			expectedCredit: 0,
		},
		{
			name:           "negative debit adds to existing credit",
			debit:          -50,
			credit:         100,
			expectedDebit:  0,
			expectedCredit: 150,
		},
		{
			name:           "both negative and equal",
			debit:          -100,
			credit:         -100,
			expectedDebit:  100,
			expectedCredit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := makeTestGLEntry("Sales - ABC", tt.debit, tt.credit)
			result := ToggleDebitCreditIfNegative([]GLEntry{entry})

			if result[0].Debit != tt.expectedDebit {
				t.Errorf("Debit = %v, want %v", result[0].Debit, tt.expectedDebit)
			}
			if result[0].Credit != tt.expectedCredit {
				t.Errorf("Credit = %v, want %v", result[0].Credit, tt.expectedCredit)
			}
		})
	}
}

// Tests for validateDisabledAccounts

func TestValidateDisabledAccounts(t *testing.T) {
	engine := &Engine{
		Accounts: newMockAccountLookup(),
	}

	tests := []struct {
		name      string
		entries   []GLEntry
		wantError bool
	}{
		{
			name: "valid accounts",
			entries: []GLEntry{
				makeTestGLEntry("Sales - ABC", 100, 0),
				makeTestGLEntry("Debtors - ABC", 0, 100),
			},
			wantError: false,
		},
		{
			name: "disabled account",
			entries: []GLEntry{
				makeTestGLEntry("Disabled Account - ABC", 100, 0),
				makeTestGLEntry("Debtors - ABC", 0, 100),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validateDisabledAccounts(tt.entries)

			if tt.wantError && err == nil {
				t.Error("validateDisabledAccounts() expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("validateDisabledAccounts() unexpected error: %v", err)
			}

			if tt.wantError && err != nil {
				var disabledErr *DisabledAccountsError
				if !errors.As(err, &disabledErr) {
					t.Error("expected DisabledAccountsError")
				}
			}
		})
	}
}

// Tests for ProcessGLMap

func TestProcessGLMap(t *testing.T) {
	engine := &Engine{
		Accounts: newMockAccountLookup(),
	}

	tests := []struct {
		name          string
		entries       []GLEntry
		mergeEntries  bool
		expectedCount int
	}{
		{
			name:          "empty entries",
			entries:       []GLEntry{},
			mergeEntries:  true,
			expectedCount: 0,
		},
		{
			name: "with merge",
			entries: []GLEntry{
				makeTestGLEntry("Sales - ABC", 50, 0),
				makeTestGLEntry("Sales - ABC", 50, 0),
				makeTestGLEntry("Debtors - ABC", 0, 100),
			},
			mergeEntries:  true,
			expectedCount: 2,
		},
		{
			name: "without merge",
			entries: []GLEntry{
				makeTestGLEntry("Sales - ABC", 50, 0),
				makeTestGLEntry("Sales - ABC", 50, 0),
				makeTestGLEntry("Debtors - ABC", 0, 100),
			},
			mergeEntries:  false,
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ProcessGLMap(tt.entries, tt.mergeEntries, false)

			if err != nil {
				t.Errorf("ProcessGLMap() error = %v", err)
				return
			}

			if len(result) != tt.expectedCount {
				t.Errorf("ProcessGLMap() count = %d, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

// Tests for getDebitCreditDifference

func TestGetDebitCreditDifference(t *testing.T) {
	tests := []struct {
		name     string
		entries  []GLEntry
		expected float64
	}{
		{
			name:     "empty entries",
			entries:  []GLEntry{},
			expected: 0,
		},
		{
			name: "balanced entries",
			entries: []GLEntry{
				makeTestGLEntry("Debtors - ABC", 100, 0),
				makeTestGLEntry("Sales - ABC", 0, 100),
			},
			expected: 0,
		},
		{
			name: "debit excess",
			entries: []GLEntry{
				makeTestGLEntry("Debtors - ABC", 100, 0),
				makeTestGLEntry("Sales - ABC", 0, 90),
			},
			expected: 10,
		},
		{
			name: "credit excess",
			entries: []GLEntry{
				makeTestGLEntry("Debtors - ABC", 90, 0),
				makeTestGLEntry("Sales - ABC", 0, 100),
			},
			expected: -10.0, // Note: floating point, may be -9.99 due to precision
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDebitCreditDifference(tt.entries, 2)
			// Compare with tolerance for floating point
			if absFloat(got-tt.expected) > 0.01 {
				t.Errorf("getDebitCreditDifference() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Tests for getDebitCreditAllowance

func TestGetDebitCreditAllowance(t *testing.T) {
	tests := []struct {
		name        string
		voucherType string
		precision   int
		expected    float64
	}{
		{
			name:        "journal entry",
			voucherType: "Journal Entry",
			precision:   2,
			expected:    0.05,
		},
		{
			name:        "payment entry",
			voucherType: "Payment Entry",
			precision:   2,
			expected:    0.05,
		},
		{
			name:        "sales invoice",
			voucherType: "Sales Invoice",
			precision:   2,
			expected:    0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDebitCreditAllowance(tt.voucherType, tt.precision)
			if got != tt.expected {
				t.Errorf("getDebitCreditAllowance() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Tests for Flt and Round

func TestFlt(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		precision []int
		expected  float64
	}{
		{
			name:      "no precision",
			value:     123.456789,
			precision: nil,
			expected:  123.456789,
		},
		{
			name:      "precision 2",
			value:     123.456,
			precision: []int{2},
			expected:  123.46,
		},
		{
			name:      "precision 0",
			value:     123.456,
			precision: []int{0},
			expected:  123,
		},
		{
			name:      "precision 3",
			value:     123.4564,
			precision: []int{3},
			expected:  123.456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flt(tt.value, tt.precision...)
			if got != tt.expected {
				t.Errorf("Flt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		precision int
		expected  float64
	}{
		{
			name:      "round up",
			value:     123.456,
			precision: 2,
			expected:  123.46,
		},
		{
			name:      "round down",
			value:     123.454,
			precision: 2,
			expected:  123.45,
		},
		{
			name:      "round half up",
			value:     123.455,
			precision: 2,
			expected:  123.46,
		},
		{
			name:      "negative precision",
			value:     123.456,
			precision: -1,
			expected:  123.456, // Returns unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Round(tt.value, tt.precision)
			if got != tt.expected {
				t.Errorf("Round() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test error types

func TestValidationError(t *testing.T) {
	err := NewValidationError(ErrAccountDisabled, "Test Account", "Account is disabled for posting")

	if !errors.Is(err, ErrAccountDisabled) {
		t.Error("ValidationError should wrap ErrAccountDisabled")
	}

	expectedMsg := "account is disabled: Test Account - Account is disabled for posting"
	if err.Error() != expectedMsg {
		t.Errorf("Error() = %q, want %q", err.Error(), expectedMsg)
	}
}

func TestDisabledAccountsError(t *testing.T) {
	err := &DisabledAccountsError{
		Accounts: []string{"Account A", "Account B"},
	}

	if !errors.Is(err, ErrAccountDisabled) {
		t.Error("DisabledAccountsError should wrap ErrAccountDisabled")
	}

	if !containsString(err.Error(), "Account A") || !containsString(err.Error(), "Account B") {
		t.Errorf("Error() should contain account names: %s", err.Error())
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Integration test for full GL posting flow

func TestMakeGLEntries_Integration(t *testing.T) {
	glStore := &mockGLStore{}
	engine := &Engine{
		Accounts: newMockAccountLookup(),
		Company:  &mockCompanySettings{},
		GLStore:  glStore,
	}

	entries := []GLEntry{
		makeTestGLEntry("Debtors - ABC", 100, 0),
		makeTestGLEntry("Sales - ABC", 0, 100),
	}

	err := engine.MakeGLEntries(entries, DefaultPostingOptions())

	if err != nil {
		t.Errorf("MakeGLEntries() error = %v", err)
		return
	}

	if len(glStore.entries) != 2 {
		t.Errorf("Expected 2 entries saved, got %d", len(glStore.entries))
	}
}

func TestMakeGLEntries_DisabledAccount(t *testing.T) {
	engine := &Engine{
		Accounts: newMockAccountLookup(),
	}

	entries := []GLEntry{
		makeTestGLEntry("Disabled Account - ABC", 100, 0),
		makeTestGLEntry("Sales - ABC", 0, 100),
	}

	err := engine.MakeGLEntries(entries, DefaultPostingOptions())

	if err == nil {
		t.Error("MakeGLEntries() expected error for disabled account")
	}

	var disabledErr *DisabledAccountsError
	if !errors.As(err, &disabledErr) {
		t.Errorf("Expected DisabledAccountsError, got %T", err)
	}
}

func TestMakeGLEntries_Empty(t *testing.T) {
	engine := &Engine{}

	err := engine.MakeGLEntries([]GLEntry{}, DefaultPostingOptions())

	if err != nil {
		t.Errorf("MakeGLEntries() for empty should not error, got %v", err)
	}
}
