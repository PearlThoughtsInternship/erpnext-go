package repository

import (
	"testing"
	"time"
)

func TestInMemoryStore_Save(t *testing.T) {
	store := NewInMemoryStore()

	entry := &GLEntry{
		Name:        "GL-00001",
		Account:     "Sales - ACME",
		Credit:      10000,
		VoucherType: "Sales Invoice",
		VoucherNo:   "SINV-2024-00001",
		PostingDate: time.Now(),
	}

	// Save should not return an error
	err := store.Save(entry)
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Entry should be stored
	all := store.GetAll()
	if len(all) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(all))
	}

	// Verify the stored entry
	if all[0].Account != "Sales - ACME" {
		t.Errorf("Stored entry account = %q, want %q", all[0].Account, "Sales - ACME")
	}
}

func TestInMemoryStore_SaveMultiple(t *testing.T) {
	store := NewInMemoryStore()

	entries := []*GLEntry{
		{Account: "Debtors - ACME", Debit: 11800, VoucherType: "Sales Invoice", VoucherNo: "SINV-2024-00001"},
		{Account: "Sales - ACME", Credit: 10000, VoucherType: "Sales Invoice", VoucherNo: "SINV-2024-00001"},
		{Account: "CGST - ACME", Credit: 900, VoucherType: "Sales Invoice", VoucherNo: "SINV-2024-00001"},
		{Account: "SGST - ACME", Credit: 900, VoucherType: "Sales Invoice", VoucherNo: "SINV-2024-00001"},
	}

	for _, entry := range entries {
		if err := store.Save(entry); err != nil {
			t.Errorf("Save() error = %v", err)
		}
	}

	all := store.GetAll()
	if len(all) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(all))
	}
}

func TestInMemoryStore_GetByVoucher(t *testing.T) {
	store := NewInMemoryStore()

	// Save entries for two different vouchers
	store.Save(&GLEntry{Account: "Debtors", Debit: 100, VoucherType: "Sales Invoice", VoucherNo: "SINV-001"})
	store.Save(&GLEntry{Account: "Sales", Credit: 100, VoucherType: "Sales Invoice", VoucherNo: "SINV-001"})
	store.Save(&GLEntry{Account: "Cash", Debit: 100, VoucherType: "Payment Entry", VoucherNo: "PAY-001"})
	store.Save(&GLEntry{Account: "Debtors", Credit: 100, VoucherType: "Payment Entry", VoucherNo: "PAY-001"})

	tests := []struct {
		name        string
		voucherType string
		voucherNo   string
		wantCount   int
	}{
		{
			name:        "find_sales_invoice",
			voucherType: "Sales Invoice",
			voucherNo:   "SINV-001",
			wantCount:   2,
		},
		{
			name:        "find_payment_entry",
			voucherType: "Payment Entry",
			voucherNo:   "PAY-001",
			wantCount:   2,
		},
		{
			name:        "not_found",
			voucherType: "Journal Entry",
			voucherNo:   "JE-001",
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := store.GetByVoucher(tt.voucherType, tt.voucherNo)
			if err != nil {
				t.Errorf("GetByVoucher() error = %v", err)
			}
			if len(entries) != tt.wantCount {
				t.Errorf("GetByVoucher() returned %d entries, want %d", len(entries), tt.wantCount)
			}
		})
	}
}

func TestInMemoryStore_GetAll(t *testing.T) {
	store := NewInMemoryStore()

	// Empty store
	if len(store.GetAll()) != 0 {
		t.Error("Empty store should return empty slice")
	}

	// Add entries
	store.Save(&GLEntry{Account: "Test1"})
	store.Save(&GLEntry{Account: "Test2"})

	all := store.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(all))
	}
}

func TestInMemoryStore_Clear(t *testing.T) {
	store := NewInMemoryStore()

	// Add some entries
	store.Save(&GLEntry{Account: "Test1"})
	store.Save(&GLEntry{Account: "Test2"})

	if len(store.GetAll()) != 2 {
		t.Fatal("Setup failed - should have 2 entries")
	}

	// Clear
	store.Clear()

	// Verify empty
	if len(store.GetAll()) != 0 {
		t.Error("Clear() should remove all entries")
	}

	// Can still add after clear
	store.Save(&GLEntry{Account: "Test3"})
	if len(store.GetAll()) != 1 {
		t.Error("Should be able to add entries after Clear()")
	}
}

func TestInMemoryStore_ImplementsInterface(t *testing.T) {
	// This test verifies that InMemoryStore implements GLEntryStore
	var store GLEntryStore = NewInMemoryStore()

	// If this compiles, the interface is implemented correctly
	_ = store
	t.Log("✅ InMemoryStore correctly implements GLEntryStore interface")
}

// TestRealisticUsage shows how the repository is used in practice
func TestRealisticUsage(t *testing.T) {
	store := NewInMemoryStore()

	// Simulate posting a Sales Invoice
	invoiceNo := "SINV-2024-00001"
	glEntries := []*GLEntry{
		{Account: "Debtors - ACME", Debit: 11800, VoucherType: "Sales Invoice", VoucherNo: invoiceNo},
		{Account: "Sales - ACME", Credit: 10000, VoucherType: "Sales Invoice", VoucherNo: invoiceNo},
		{Account: "CGST Payable - ACME", Credit: 900, VoucherType: "Sales Invoice", VoucherNo: invoiceNo},
		{Account: "SGST Payable - ACME", Credit: 900, VoucherType: "Sales Invoice", VoucherNo: invoiceNo},
	}

	// Save all entries
	for _, entry := range glEntries {
		if err := store.Save(entry); err != nil {
			t.Fatalf("Failed to save entry: %v", err)
		}
	}

	// Retrieve by voucher
	retrieved, err := store.GetByVoucher("Sales Invoice", invoiceNo)
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}

	// Verify
	if len(retrieved) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(retrieved))
	}

	// Calculate totals
	var totalDebit, totalCredit float64
	for _, e := range retrieved {
		totalDebit += e.Debit
		totalCredit += e.Credit
	}

	if totalDebit != totalCredit {
		t.Errorf("Entries don't balance: Debit=%.2f, Credit=%.2f", totalDebit, totalCredit)
	}

	t.Logf("✅ Posted invoice %s: Dr %.2f = Cr %.2f", invoiceNo, totalDebit, totalCredit)
}
