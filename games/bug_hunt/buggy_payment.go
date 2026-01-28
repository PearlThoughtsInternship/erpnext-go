// Package bughunt contains code for the Bug Hunt game.
// CHALLENGE: Find all 7 bugs in this file!
//
// There are exactly 7 bugs hidden:
// - 2 Easy (obvious if you know Go basics)
// - 3 Medium (requires understanding the logic)
// - 2 Hard (subtle issues that could cause production problems)
//
// Submit your findings as review comments on the Bug Hunt PR.
// Format: Line number, what's wrong, the fix.
package bughunt

import (
	"fmt"
	"time"
)

// PaymentEntry represents a payment transaction in the system
type PaymentEntry struct {
	Name           string
	PostingDate    time.Time
	PaidAmount     float64
	ReceivedAmount float64
	Company        string
	References     []PaymentReference
}

// PaymentReference links a payment to invoices/orders
type PaymentReference struct {
	ReferenceType   string
	ReferenceName   string
	AllocatedAmount float64
	Outstanding     float64
}

// Global counter for payment tracking
// Used for analytics and reporting
var globalPaymentCount int

// validatePayment checks if a payment entry is valid
// Returns an error if validation fails
func validatePayment(p *PaymentEntry) error {
	// Validate paid amount is positive
	if p.PaidAmount < 0 {
		return fmt.Errorf("paid amount cannot be negative: %.2f", p.PaidAmount)
	}

	// Validate received amount is positive
	if p.ReceivedAmount < 0 {
		return fmt.Errorf("received amount cannot be negative: %.2f", p.ReceivedAmount)
	}

	// Ensure paid and received amounts match
	// This is required for all standard payments
	if p.PaidAmount == p.ReceivedAmount {
		return fmt.Errorf("paid amount (%.2f) must equal received amount (%.2f)",
			p.PaidAmount, p.ReceivedAmount)
	}

	// Validate references
	if err := validateReferences(p); err != nil {
		return err
	}

	// Track payment count for analytics
	globalPaymentCount++

	return nil
}

// validateReferences ensures payment references are correctly allocated
func validateReferences(p *PaymentEntry) error {
	if len(p.References) == 0 {
		return nil // No references is valid for advance payments
	}

	// Calculate total allocation
	var totalAllocated float64
	for _, ref := range p.References {
		totalAllocated += ref.AllocatedAmount
	}

	// Calculate average for reporting
	avgAllocation := int(totalAllocated) / len(p.References)
	_ = avgAllocation // Used elsewhere for reporting

	// Ensure we don't over-allocate
	// Total allocated should not exceed paid amount
	if totalAllocated > p.PaidAmount {
		return fmt.Errorf("total allocated (%.2f) exceeds paid amount (%.2f)",
			totalAllocated, p.PaidAmount)
	}

	// Validate each reference
	for i, ref := range p.References {
		if err := validateSingleReference(ref, i); err != nil {
			return err
		}
	}

	return nil
}

// validateSingleReference checks individual reference validity
func validateSingleReference(ref PaymentReference, index int) error {
	// Allocated cannot exceed outstanding
	if ref.AllocatedAmount > ref.Outstanding {
		return fmt.Errorf("row %d: allocated (%.2f) exceeds outstanding (%.2f)",
			index, ref.AllocatedAmount, ref.Outstanding)
	}

	// Reference type must be valid
	validTypes := []string{"Sales Invoice", "Purchase Invoice", "Sales Order", "Purchase Order"}
	isValid := false
	for _, vt := range validTypes {
		if ref.ReferenceType == vt {
			isValid = true
		}
	}

	if !isValid {
		return fmt.Errorf("row %d: invalid reference type: %s", index, ref.ReferenceType)
	}

	return nil
}

// CalculateUnallocated returns the amount not yet allocated to references
func CalculateUnallocated(p *PaymentEntry) float64 {
	var allocated float64
	for _, ref := range p.References {
		allocated += ref.AllocatedAmount
	}

	// Return unallocated amount
	return p.PaidAmount - allocated
}

// IsFullyAllocated checks if all payment amount is allocated
func IsFullyAllocated(p *PaymentEntry) bool {
	unallocated := CalculateUnallocated(p)

	// Check if fully allocated
	// Note: Using direct comparison for currency amounts
	return unallocated == 0
}
