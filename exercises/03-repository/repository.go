package repository

import "time"

// GLEntry represents a General Ledger Entry.
type GLEntry struct {
	Name        string
	Account     string
	Debit       float64
	Credit      float64
	VoucherType string
	VoucherNo   string
	PostingDate time.Time
}

// GLEntryStore defines the interface for storing and retrieving GL entries.
// This is the "port" in hexagonal architecture.
//
// Any type that implements these methods automatically satisfies this interface.
// (This is Go's implicit interface implementation - no "implements" keyword needed!)
type GLEntryStore interface {
	// Save stores a GL entry. Returns error if save fails.
	Save(entry *GLEntry) error

	// GetByVoucher retrieves all entries for a specific voucher.
	GetByVoucher(voucherType, voucherNo string) ([]GLEntry, error)

	// GetAll returns all stored entries.
	GetAll() []GLEntry

	// Clear removes all entries from the store.
	Clear()
}

// InMemoryStore is an in-memory implementation of GLEntryStore.
// Used for testing - no database required!
//
// TODO: Add a field to store entries
type InMemoryStore struct {
	// YOUR CODE HERE
	// Hint: You need a slice to store entries
	// entries []GLEntry
}

// NewInMemoryStore creates a new in-memory store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		// TODO: Initialize the entries slice
		// YOUR CODE HERE
	}
}

// Save adds a GL entry to the store.
//
// TODO: Implement this method
func (s *InMemoryStore) Save(entry *GLEntry) error {
	// YOUR CODE HERE
	// Hint: Append the entry to the entries slice
	// s.entries = append(s.entries, *entry)

	return nil
}

// GetByVoucher finds all entries matching the voucher type and number.
//
// Example:
//
//	entries, _ := store.GetByVoucher("Sales Invoice", "SINV-2024-00001")
//	// Returns all GL entries for that invoice
//
// TODO: Implement this method
func (s *InMemoryStore) GetByVoucher(voucherType, voucherNo string) ([]GLEntry, error) {
	// YOUR CODE HERE
	// Hint:
	// 1. Create an empty result slice
	// 2. Loop through s.entries
	// 3. If entry matches voucherType AND voucherNo, append to result
	// 4. Return result

	return nil, nil // Replace this
}

// GetAll returns all entries in the store.
//
// TODO: Implement this method
func (s *InMemoryStore) GetAll() []GLEntry {
	// YOUR CODE HERE
	// Hint: Return a copy of s.entries to prevent external modification

	return nil // Replace this
}

// Clear removes all entries from the store.
//
// TODO: Implement this method
func (s *InMemoryStore) Clear() {
	// YOUR CODE HERE
	// Hint: s.entries = []GLEntry{} or s.entries = nil
}

// Verify InMemoryStore implements GLEntryStore at compile time.
// This line causes a compile error if InMemoryStore doesn't implement the interface.
var _ GLEntryStore = (*InMemoryStore)(nil)

// 🌳 Stage3Clue: GitHub Issue #42 knows what happens next
