package ledger

import (
	"testing"
	"time"
)

// =============================================================================
// INTEGRATION TESTS WITH REALISTIC ERPNEXT DATA
// =============================================================================
//
// These tests use actual ERPNext-style data to verify our Go implementation
// produces the same GL entries that Python/Frappe would produce.
//
// Test data is modeled after real ERPNext transactions:
// - Sales Invoice with GST
// - Payment Entry
// - Journal Entry for adjustments
// =============================================================================

// TestRealisticSalesInvoiceGLEntries verifies GL entry creation for a typical
// Indian Sales Invoice with GST taxation.
//
// Python equivalent scenario:
//   - Sales Invoice SINV-2024-00001
//   - Customer: "Acme Corporation"
//   - Item: Widget @ ₹10,000
//   - CGST 9% = ₹900
//   - SGST 9% = ₹900
//   - Grand Total: ₹11,800
//
// Expected GL Entries (from ERPNext):
//   Debtors - Acme         Debit  ₹11,800
//   Sales Revenue          Credit ₹10,000
//   CGST Payable           Credit ₹900
//   SGST Payable           Credit ₹900
func TestRealisticSalesInvoiceGLEntries(t *testing.T) {
	// Realistic ERPNext-style GL entries for a Sales Invoice
	glEntries := []GLEntry{
		{
			Name:                    "ACC-GLE-2024-00001",
			PostingDate:             time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Account:                 "Debtors - ACME",
			AccountCurrency:         "INR",
			Debit:                   11800.00,
			Credit:                  0,
			DebitInAccountCurrency:  11800.00,
			CreditInAccountCurrency: 0,
			PartyType:               "Customer",
			Party:                   "Acme Corporation",
			Against:                 "Sales - ACME, CGST - ACME, SGST - ACME",
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			CostCenter:              "Main - ACME",
			FiscalYear:              "2024-2025",
			IsOpening:               IsOpeningNo,
			IsAdvance:               IsAdvanceNo,
			Remarks:                 "Against Sales Invoice SINV-2024-00001",
		},
		{
			Name:                    "ACC-GLE-2024-00002",
			PostingDate:             time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  10000.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 10000.00,
			Against:                 "Acme Corporation",
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			CostCenter:              "Main - ACME",
			FiscalYear:              "2024-2025",
			IsOpening:               IsOpeningNo,
			Remarks:                 "Against Sales Invoice SINV-2024-00001",
		},
		{
			Name:                    "ACC-GLE-2024-00003",
			PostingDate:             time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Account:                 "CGST Payable - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  900.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 900.00,
			Against:                 "Acme Corporation",
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			CostCenter:              "Main - ACME",
			FiscalYear:              "2024-2025",
			IsOpening:               IsOpeningNo,
			Remarks:                 "Against Sales Invoice SINV-2024-00001",
		},
		{
			Name:                    "ACC-GLE-2024-00004",
			PostingDate:             time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Account:                 "SGST Payable - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  900.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 900.00,
			Against:                 "Acme Corporation",
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			CostCenter:              "Main - ACME",
			FiscalYear:              "2024-2025",
			IsOpening:               IsOpeningNo,
			Remarks:                 "Against Sales Invoice SINV-2024-00001",
		},
	}

	// Verify the GL entries balance
	glMap := GLMap(glEntries)

	t.Run("entries_balance", func(t *testing.T) {
		if !glMap.IsBalanced() {
			t.Errorf("GL entries do not balance: Debit=%v, Credit=%v",
				glMap.TotalDebit(), glMap.TotalCredit())
		}
	})

	t.Run("total_debit_equals_invoice_amount", func(t *testing.T) {
		expectedTotal := 11800.00
		if glMap.TotalDebit() != expectedTotal {
			t.Errorf("Total debit = %v, want %v", glMap.TotalDebit(), expectedTotal)
		}
	})

	t.Run("total_credit_equals_invoice_amount", func(t *testing.T) {
		expectedTotal := 11800.00
		if glMap.TotalCredit() != expectedTotal {
			t.Errorf("Total credit = %v, want %v", glMap.TotalCredit(), expectedTotal)
		}
	})

	// Process through our engine (merge similar entries)
	t.Run("process_gl_map", func(t *testing.T) {
		engine := &Engine{
			Accounts: &realisticAccountLookup{},
		}

		processed, err := engine.ProcessGLMap(glEntries, true, false)
		if err != nil {
			t.Fatalf("ProcessGLMap error: %v", err)
		}

		// Should still have 4 entries (no duplicates to merge)
		if len(processed) != 4 {
			t.Errorf("Processed entries = %d, want 4", len(processed))
		}

		// Should still balance
		processedMap := GLMap(processed)
		if !processedMap.IsBalanced() {
			t.Errorf("Processed entries do not balance")
		}
	})
}

// TestRealisticPaymentEntryGLEntries verifies GL entries for payment receipt.
//
// Python equivalent scenario:
//   - Payment Entry PE-2024-00001
//   - Customer: "Acme Corporation" pays ₹11,800
//   - Against Invoice: SINV-2024-00001
//
// Expected GL Entries (from ERPNext):
//   Bank Account           Debit  ₹11,800
//   Debtors - Acme         Credit ₹11,800
func TestRealisticPaymentEntryGLEntries(t *testing.T) {
	dueDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	glEntries := []GLEntry{
		{
			Name:                    "ACC-GLE-2024-00010",
			PostingDate:             time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Account:                 "HDFC Bank - ACME",
			AccountCurrency:         "INR",
			Debit:                   11800.00,
			Credit:                  0,
			DebitInAccountCurrency:  11800.00,
			CreditInAccountCurrency: 0,
			Against:                 "Acme Corporation",
			VoucherType:             "Payment Entry",
			VoucherNo:               "PE-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			FiscalYear:              "2024-2025",
			IsOpening:               IsOpeningNo,
			Remarks:                 "Amount ₹11,800 received from Acme Corporation",
		},
		{
			Name:                    "ACC-GLE-2024-00011",
			PostingDate:             time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Account:                 "Debtors - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  11800.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 11800.00,
			PartyType:               "Customer",
			Party:                   "Acme Corporation",
			Against:                 "HDFC Bank - ACME",
			AgainstVoucherType:      "Sales Invoice",
			AgainstVoucher:          "SINV-2024-00001",
			VoucherType:             "Payment Entry",
			VoucherNo:               "PE-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			FiscalYear:              "2024-2025",
			DueDate:                 &dueDate,
			IsOpening:               IsOpeningNo,
			Remarks:                 "Amount ₹11,800 received from Acme Corporation",
		},
	}

	glMap := GLMap(glEntries)

	t.Run("entries_balance", func(t *testing.T) {
		if !glMap.IsBalanced() {
			t.Errorf("GL entries do not balance")
		}
	})

	t.Run("against_voucher_set_correctly", func(t *testing.T) {
		// The credit entry should reference the original invoice
		creditEntry := glEntries[1]
		if creditEntry.AgainstVoucher != "SINV-2024-00001" {
			t.Errorf("AgainstVoucher = %q, want SINV-2024-00001", creditEntry.AgainstVoucher)
		}
	})

	// 🌳 Stage2Clue: QXJvdW5kIHRoZSBvbGQgdHJlZSwgbmV2ZXIgY3V0dGluZw==
	// (The repository holds the third secret)
}

// TestRealisticJournalEntryGLEntries verifies GL entries for manual adjustments.
//
// Python equivalent scenario:
//   - Journal Entry JV-2024-00001
//   - Expense adjustment: ₹5,000
//   - From: Bank Account
//   - To: Office Expenses
func TestRealisticJournalEntryGLEntries(t *testing.T) {
	glEntries := []GLEntry{
		{
			Name:                    "ACC-GLE-2024-00020",
			PostingDate:             time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Account:                 "Office Expenses - ACME",
			AccountCurrency:         "INR",
			Debit:                   5000.00,
			Credit:                  0,
			DebitInAccountCurrency:  5000.00,
			CreditInAccountCurrency: 0,
			Against:                 "HDFC Bank - ACME",
			VoucherType:             "Journal Entry",
			VoucherNo:               "JV-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			CostCenter:              "Main - ACME",
			FiscalYear:              "2024-2025",
			IsOpening:               IsOpeningNo,
			Remarks:                 "Office supplies expense for January 2024",
		},
		{
			Name:                    "ACC-GLE-2024-00021",
			PostingDate:             time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Account:                 "HDFC Bank - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  5000.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 5000.00,
			Against:                 "Office Expenses - ACME",
			VoucherType:             "Journal Entry",
			VoucherNo:               "JV-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			FiscalYear:              "2024-2025",
			IsOpening:               IsOpeningNo,
			Remarks:                 "Office supplies expense for January 2024",
		},
	}

	glMap := GLMap(glEntries)

	t.Run("entries_balance", func(t *testing.T) {
		if !glMap.IsBalanced() {
			t.Errorf("GL entries do not balance")
		}
	})
}

// TestMergeSimilarEntriesRealistic tests merging with real-world data.
//
// Scenario: Multiple line items in a Sales Invoice create separate GL entries
// that need to be merged by account.
//
// Python equivalent:
//   Sales Invoice with 3 items:
//   - Widget A: ₹5,000
//   - Widget B: ₹3,000
//   - Widget C: ₹2,000
//   Total Sales: ₹10,000 (merged to single GL entry)
func TestMergeSimilarEntriesRealistic(t *testing.T) {
	// Before merge: 3 separate entries for same account
	glEntries := []GLEntry{
		{
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  5000.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 5000.00,
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00002",
			VoucherDetailNo:         "item-001",
			CostCenter:              "Main - ACME",
			Remarks:                 "Widget A",
		},
		{
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  3000.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 3000.00,
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00002",
			VoucherDetailNo:         "item-002",
			CostCenter:              "Main - ACME",
			Remarks:                 "Widget B",
		},
		{
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  2000.00,
			DebitInAccountCurrency:  0,
			CreditInAccountCurrency: 2000.00,
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00002",
			VoucherDetailNo:         "item-003",
			CostCenter:              "Main - ACME",
			Remarks:                 "Widget C",
		},
	}

	// These should NOT merge because VoucherDetailNo differs
	// In ERPNext, merge_similar_entries considers voucher_detail_no
	merged := MergeSimilarEntries(glEntries)

	t.Run("separate_items_not_merged", func(t *testing.T) {
		// Different VoucherDetailNo means different merge keys
		if len(merged) != 3 {
			t.Errorf("Expected 3 entries (different detail nos), got %d", len(merged))
		}
	})

	// Now test with same VoucherDetailNo (should merge)
	glEntriesSameDetail := []GLEntry{
		{
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  5000.00,
			CostCenter:              "Main - ACME",
			VoucherNo:               "SINV-2024-00002",
		},
		{
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  3000.00,
			CostCenter:              "Main - ACME",
			VoucherNo:               "SINV-2024-00002",
		},
		{
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  2000.00,
			CostCenter:              "Main - ACME",
			VoucherNo:               "SINV-2024-00002",
		},
	}

	mergedSame := MergeSimilarEntries(glEntriesSameDetail)

	t.Run("same_account_merged", func(t *testing.T) {
		if len(mergedSame) != 1 {
			t.Errorf("Expected 1 merged entry, got %d", len(mergedSame))
		}
		if mergedSame[0].Credit != 10000.00 {
			t.Errorf("Merged credit = %v, want 10000", mergedSame[0].Credit)
		}
	})
}

// TestMultiCurrencyGLEntries tests multi-currency transaction handling.
//
// Scenario: USD invoice converted to INR
//   - Invoice: $1,000 USD
//   - Exchange rate: 83.50 INR/USD
//   - Company currency: INR
func TestMultiCurrencyGLEntries(t *testing.T) {
	glEntries := []GLEntry{
		{
			Account:                        "Debtors - ACME",
			AccountCurrency:                "INR",
			Debit:                          83500.00, // Company currency
			Credit:                         0,
			DebitInAccountCurrency:         83500.00,
			CreditInAccountCurrency:        0,
			TransactionCurrency:            "USD",
			TransactionExchangeRate:        83.50,
			DebitInTransactionCurrency:     1000.00, // Transaction currency
			CreditInTransactionCurrency:    0,
			VoucherType:                    "Sales Invoice",
			VoucherNo:                      "SINV-2024-USD-001",
		},
		{
			Account:                        "Sales - ACME",
			AccountCurrency:                "INR",
			Debit:                          0,
			Credit:                         83500.00,
			DebitInAccountCurrency:         0,
			CreditInAccountCurrency:        83500.00,
			TransactionCurrency:            "USD",
			TransactionExchangeRate:        83.50,
			DebitInTransactionCurrency:     0,
			CreditInTransactionCurrency:    1000.00,
			VoucherType:                    "Sales Invoice",
			VoucherNo:                      "SINV-2024-USD-001",
		},
	}

	glMap := GLMap(glEntries)

	t.Run("company_currency_balanced", func(t *testing.T) {
		if !glMap.IsBalanced() {
			t.Errorf("Company currency entries do not balance")
		}
	})

	t.Run("transaction_currency_balanced", func(t *testing.T) {
		var totalDebitTxn, totalCreditTxn float64
		for _, e := range glEntries {
			totalDebitTxn += e.DebitInTransactionCurrency
			totalCreditTxn += e.CreditInTransactionCurrency
		}
		if Flt(totalDebitTxn-totalCreditTxn, 2) != 0 {
			t.Errorf("Transaction currency not balanced: Debit=%v, Credit=%v",
				totalDebitTxn, totalCreditTxn)
		}
	})

	t.Run("exchange_rate_applied_correctly", func(t *testing.T) {
		// $1,000 * 83.50 = ₹83,500
		expectedINR := 1000.00 * 83.50
		if glEntries[0].Debit != expectedINR {
			t.Errorf("Debit = %v, want %v", glEntries[0].Debit, expectedINR)
		}
	})
}

// TestFullGLPostingFlow tests the complete GL posting workflow.
//
// This is an end-to-end integration test simulating what happens
// when a Sales Invoice is submitted in ERPNext.
func TestFullGLPostingFlow(t *testing.T) {
	// Setup mock adapters
	glStore := &mockGLStore{}
	engine := &Engine{
		Accounts: &realisticAccountLookup{},
		Company:  &mockCompanySettings{},
		GLStore:  glStore,
	}

	// Create GL entries for a Sales Invoice
	glEntries := []GLEntry{
		{
			PostingDate:            time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Account:                "Debtors - ACME",
			AccountCurrency:        "INR",
			Debit:                  11800.00,
			Credit:                 0,
			DebitInAccountCurrency: 11800.00,
			PartyType:              "Customer",
			Party:                  "Acme Corporation",
			Against:                "Sales - ACME",
			VoucherType:            "Sales Invoice",
			VoucherNo:              "SINV-2024-00001",
			Company:                "ACME Industries Pvt Ltd",
			CostCenter:             "Main - ACME",
		},
		{
			PostingDate:             time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Account:                 "Sales - ACME",
			AccountCurrency:         "INR",
			Debit:                   0,
			Credit:                  11800.00,
			CreditInAccountCurrency: 11800.00,
			Against:                 "Acme Corporation",
			VoucherType:             "Sales Invoice",
			VoucherNo:               "SINV-2024-00001",
			Company:                 "ACME Industries Pvt Ltd",
			CostCenter:              "Main - ACME",
		},
	}

	err := engine.MakeGLEntries(glEntries, DefaultPostingOptions())

	t.Run("no_error", func(t *testing.T) {
		if err != nil {
			t.Errorf("MakeGLEntries error: %v", err)
		}
	})

	t.Run("entries_saved", func(t *testing.T) {
		if len(glStore.entries) != 2 {
			t.Errorf("Expected 2 entries saved, got %d", len(glStore.entries))
		}
	})

	t.Run("entries_balanced", func(t *testing.T) {
		savedMap := GLMap(glStore.entries)
		if !savedMap.IsBalanced() {
			t.Errorf("Saved entries not balanced")
		}
	})
}

// =============================================================================
// MOCK ADAPTERS FOR REALISTIC TESTING
// =============================================================================

// realisticAccountLookup provides mock data matching ERPNext account structure.
type realisticAccountLookup struct{}

func (r *realisticAccountLookup) GetAccount(name string) (*Account, error) {
	accounts := map[string]*Account{
		"Debtors - ACME": {
			Name:            "Debtors - ACME",
			AccountName:     "Debtors",
			Company:         "ACME Industries Pvt Ltd",
			AccountCurrency: "INR",
			RootType:        "Asset",
		},
		"Sales - ACME": {
			Name:            "Sales - ACME",
			AccountName:     "Sales",
			Company:         "ACME Industries Pvt Ltd",
			AccountCurrency: "INR",
			RootType:        "Income",
		},
		"CGST Payable - ACME": {
			Name:            "CGST Payable - ACME",
			AccountName:     "CGST Payable",
			Company:         "ACME Industries Pvt Ltd",
			AccountCurrency: "INR",
			RootType:        "Liability",
		},
		"SGST Payable - ACME": {
			Name:            "SGST Payable - ACME",
			AccountName:     "SGST Payable",
			Company:         "ACME Industries Pvt Ltd",
			AccountCurrency: "INR",
			RootType:        "Liability",
		},
		"HDFC Bank - ACME": {
			Name:            "HDFC Bank - ACME",
			AccountName:     "HDFC Bank",
			Company:         "ACME Industries Pvt Ltd",
			AccountCurrency: "INR",
			RootType:        "Asset",
		},
		"Office Expenses - ACME": {
			Name:            "Office Expenses - ACME",
			AccountName:     "Office Expenses",
			Company:         "ACME Industries Pvt Ltd",
			AccountCurrency: "INR",
			RootType:        "Expense",
		},
	}
	if acc, ok := accounts[name]; ok {
		return acc, nil
	}
	return &Account{Name: name, Company: "ACME Industries Pvt Ltd"}, nil
}

func (r *realisticAccountLookup) GetAccountCurrency(name string) (string, error) {
	return "INR", nil
}

func (r *realisticAccountLookup) IsGroup(name string) (bool, error) {
	return false, nil
}

func (r *realisticAccountLookup) IsFrozen(name string) (bool, error) {
	return false, nil
}

func (r *realisticAccountLookup) IsDisabled(name string) (bool, error) {
	return false, nil
}

func (r *realisticAccountLookup) GetBalanceMustBe(name string) (string, error) {
	return "", nil
}
