// engine.go implements the GL posting engine.
// Migrated from: erpnext/accounts/general_ledger.py
//
// This is the heart of ERPNext's accounting system. Every financial
// transaction flows through MakeGLEntries to create double-entry
// bookkeeping records.
package ledger

import (
	"fmt"
	"strings"
)

// MakeGLEntries is the main entry point for posting GL entries.
// All accounting transactions (Sales Invoice, Purchase Invoice,
// Journal Entry, Payment Entry, etc.) call this function.
//
// Maps to: make_gl_entries() in general_ledger.py (lines 28-67)
//
// Python equivalent:
//
//	def make_gl_entries(gl_map, cancel=False, adv_adj=False,
//	    merge_entries=True, update_outstanding="Yes", from_repost=False):
//	    if gl_map:
//	        if not cancel:
//	            make_acc_dimensions_offsetting_entry(gl_map)
//	            validate_accounting_period(gl_map)
//	            validate_disabled_accounts(gl_map)
//	            gl_map = process_gl_map(gl_map, merge_entries)
//	            ...
//	            save_entries(gl_map, ...)
//	        else:
//	            make_reverse_gl_entries(gl_map, ...)
func (e *Engine) MakeGLEntries(glMap []GLEntry, opts PostingOptions) error {
	if len(glMap) == 0 {
		return nil
	}

	// Budget validation (if enabled)
	if e.Budget != nil && glMap[0].VoucherType != "Period Closing Voucher" {
		if err := e.Budget.Validate(glMap); err != nil {
			return err
		}
	}

	if !opts.Cancel {
		// Add accounting dimension offsetting entries
		if e.Dimensions != nil {
			if err := e.makeAccDimensionsOffsettingEntry(&glMap); err != nil {
				return err
			}
		}

		// Validate accounting period
		if e.Periods != nil {
			if err := e.validateAccountingPeriod(glMap); err != nil {
				return err
			}
		}

		// Validate disabled accounts
		if err := e.validateDisabledAccounts(glMap); err != nil {
			return err
		}

		// Process GL map (distribute, merge, toggle)
		processedMap, err := e.ProcessGLMap(glMap, opts.MergeEntries, opts.FromRepost)
		if err != nil {
			return err
		}

		// Validate we have enough entries
		if len(processedMap) < 2 {
			return &GLEntryCountError{
				Expected: 2,
				Actual:   len(processedMap),
				Message:  "Incorrect number of General Ledger Entries found. You might have selected a wrong Account in the transaction.",
			}
		}

		// Create payment ledger entries (for AR/AP tracking)
		if e.PaymentStore != nil && glMap[0].VoucherType != "Period Closing Voucher" {
			if err := e.createPaymentLedgerEntries(processedMap, opts); err != nil {
				return err
			}
		}

		// Save GL entries
		if err := e.saveEntries(processedMap, opts); err != nil {
			return err
		}
	} else {
		// Cancellation - create reverse entries
		if err := e.makeReverseGLEntries(glMap, opts); err != nil {
			return err
		}
	}

	return nil
}

// ProcessGLMap processes GL entries: distributes by cost center, merges
// similar entries, and normalizes negative amounts.
//
// Maps to: process_gl_map() in general_ledger.py (lines 188-200)
//
// Python equivalent:
//
//	def process_gl_map(gl_map, merge_entries=True, precision=None, from_repost=False):
//	    if gl_map[0].voucher_type != "Period Closing Voucher":
//	        gl_map = distribute_gl_based_on_cost_center_allocation(gl_map)
//	    if merge_entries:
//	        gl_map = merge_similar_entries(gl_map)
//	    gl_map = toggle_debit_credit_if_negative(gl_map)
//	    return gl_map
func (e *Engine) ProcessGLMap(glMap []GLEntry, mergeEntries bool, fromRepost bool) ([]GLEntry, error) {
	if len(glMap) == 0 {
		return []GLEntry{}, nil
	}

	result := make([]GLEntry, len(glMap))
	copy(result, glMap)

	// Cost center allocation distribution (skip for Period Closing Voucher)
	// Note: This is a complex feature - simplified for initial implementation
	// Full implementation would use CostCenterAllocationProvider interface

	// Merge similar entries
	if mergeEntries {
		result = MergeSimilarEntries(result)
	}

	// Toggle debit/credit if negative
	result = ToggleDebitCreditIfNegative(result)

	return result, nil
}

// MergeSimilarEntries combines GL entries with the same merge key.
// This reduces the number of GL entries by consolidating entries
// to the same account/party/cost center/etc.
//
// Maps to: merge_similar_entries() in general_ledger.py (lines 273-326)
//
// Python equivalent:
//
//	def merge_similar_entries(gl_map, precision=None):
//	    merged_gl_map = []
//	    for entry in gl_map:
//	        entry.merge_key = get_merge_key(entry, merge_properties)
//	        same_head = check_if_in_list(entry, merged_gl_map)
//	        if same_head:
//	            same_head.debit += entry.debit
//	            same_head.credit += entry.credit
//	            ...
//	        else:
//	            merged_gl_map.append(entry)
//	    # Filter zero entries
//	    return merged_gl_map
func MergeSimilarEntries(glMap []GLEntry) []GLEntry {
	if len(glMap) == 0 {
		return glMap
	}

	merged := make([]GLEntry, 0, len(glMap))
	keyIndex := make(map[string]int) // merge key -> index in merged

	for _, entry := range glMap {
		key := getMergeKey(entry)

		if idx, exists := keyIndex[key]; exists {
			// Add to existing entry
			merged[idx].Debit += entry.Debit
			merged[idx].DebitInAccountCurrency += entry.DebitInAccountCurrency
			merged[idx].DebitInTransactionCurrency += entry.DebitInTransactionCurrency
			merged[idx].Credit += entry.Credit
			merged[idx].CreditInAccountCurrency += entry.CreditInAccountCurrency
			merged[idx].CreditInTransactionCurrency += entry.CreditInTransactionCurrency
		} else {
			// Add new entry
			keyIndex[key] = len(merged)
			merged = append(merged, entry)
		}
	}

	// Filter zero entries (but keep Exchange Gain Or Loss journal entries)
	result := make([]GLEntry, 0, len(merged))
	for _, entry := range merged {
		if Flt(entry.Debit, 2) != 0 || Flt(entry.Credit, 2) != 0 {
			result = append(result, entry)
		}
		// Note: In full implementation, also keep Exchange Gain Or Loss entries
	}

	return result
}

// getMergeKey creates a unique key for merging GL entries.
// Entries with the same key can be consolidated.
//
// Maps to: get_merge_key() in general_ledger.py (lines 349-354)
func getMergeKey(entry GLEntry) string {
	// Key fields that must match for merging
	parts := []string{
		entry.Account,
		entry.CostCenter,
		entry.Party,
		entry.PartyType,
		entry.VoucherDetailNo,
		entry.AgainstVoucher,
		entry.AgainstVoucherType,
		entry.Project,
		entry.FinanceBook,
		entry.VoucherNo,
	}
	return strings.Join(parts, "|")
}

// ToggleDebitCreditIfNegative normalizes negative amounts.
// If debit is negative, it's moved to credit (and vice versa).
// This ensures all GL entries have non-negative debit and credit.
//
// Maps to: toggle_debit_credit_if_negative() in general_ledger.py (lines 363-403)
//
// Python equivalent:
//
//	def toggle_debit_credit_if_negative(gl_map):
//	    for entry in gl_map:
//	        if debit < 0:
//	            credit = credit - debit
//	            debit = 0.0
//	        if credit < 0:
//	            debit = debit - credit
//	            credit = 0.0
func ToggleDebitCreditIfNegative(glMap []GLEntry) []GLEntry {
	for i := range glMap {
		entry := &glMap[i]

		// Process each debit/credit pair
		togglePair(&entry.Debit, &entry.Credit)
		togglePair(&entry.DebitInAccountCurrency, &entry.CreditInAccountCurrency)
		togglePair(&entry.DebitInTransactionCurrency, &entry.CreditInTransactionCurrency)
	}
	return glMap
}

// togglePair normalizes a debit/credit pair to non-negative values.
func togglePair(debit, credit *float64) {
	d := *debit
	c := *credit

	// Handle both negative case
	if d < 0 && c < 0 && d == c {
		d *= -1
		c *= -1
	}

	// Move negative debit to credit
	if d < 0 {
		c = c - d
		d = 0
	}

	// Move negative credit to debit
	if c < 0 {
		d = d - c
		c = 0
	}

	*debit = d
	*credit = c
}

// validateDisabledAccounts checks that no GL entries use disabled accounts.
//
// Maps to: validate_disabled_accounts() in general_ledger.py (lines 134-150)
func (e *Engine) validateDisabledAccounts(glMap []GLEntry) error {
	if e.Accounts == nil {
		return nil
	}

	var disabledAccounts []string
	checked := make(map[string]bool)

	for _, entry := range glMap {
		if entry.Account == "" || checked[entry.Account] {
			continue
		}
		checked[entry.Account] = true

		disabled, err := e.Accounts.IsDisabled(entry.Account)
		if err != nil {
			return err
		}
		if disabled {
			disabledAccounts = append(disabledAccounts, entry.Account)
		}
	}

	if len(disabledAccounts) > 0 {
		return &DisabledAccountsError{Accounts: disabledAccounts}
	}

	return nil
}

// validateAccountingPeriod checks that posting is allowed for the date.
//
// Maps to: validate_accounting_period() in general_ledger.py (lines 153-185)
func (e *Engine) validateAccountingPeriod(glMap []GLEntry) error {
	if e.Periods == nil || len(glMap) == 0 {
		return nil
	}

	entry := glMap[0]
	closed, err := e.Periods.IsDocumentTypeClosed(
		entry.Company,
		entry.VoucherType,
		entry.PostingDate,
	)
	if err != nil {
		return err
	}

	if closed {
		msg, _ := e.Periods.GetClosedPeriodMessage(
			entry.Company,
			entry.VoucherType,
			entry.PostingDate,
		)
		return &PeriodClosedError{
			Company:     entry.Company,
			DocType:     entry.VoucherType,
			PostingDate: entry.PostingDate.Format("2006-01-02"),
			PeriodName:  msg,
		}
	}

	return nil
}

// makeAccDimensionsOffsettingEntry creates offsetting entries for
// accounting dimensions when entries span multiple dimension values.
//
// Maps to: make_acc_dimensions_offsetting_entry() in general_ledger.py (lines 70-103)
func (e *Engine) makeAccDimensionsOffsettingEntry(glMap *[]GLEntry) error {
	if e.Dimensions == nil || len(*glMap) == 0 {
		return nil
	}

	company := (*glMap)[0].Company
	dimensions, err := e.Dimensions.GetDimensionsForOffsetting(*glMap, company)
	if err != nil {
		return err
	}

	if len(dimensions) == 0 {
		return nil
	}

	numDimensions := float64(len(dimensions))
	var offsettingEntries []GLEntry

	for _, entry := range *glMap {
		for _, dim := range dimensions {
			offsetting := entry.Copy()

			// Swap debit/credit for offsetting
			debit := Flt(entry.Credit, 2) / numDimensions
			credit := Flt(entry.Debit, 2) / numDimensions

			offsetting.Account = dim.OffsettingAccount
			offsetting.Debit = debit
			offsetting.Credit = credit
			offsetting.DebitInAccountCurrency = debit
			offsetting.CreditInAccountCurrency = credit
			offsetting.Remarks = fmt.Sprintf("Offsetting for Accounting Dimension - %s", dim.Name)
			offsetting.AccountCurrency = dim.AccountCurrency
			offsetting.AgainstVoucher = ""
			offsetting.AgainstVoucherType = ""
			offsetting.PartyType = ""
			offsetting.Party = ""

			offsettingEntries = append(offsettingEntries, offsetting)
		}
	}

	*glMap = append(*glMap, offsettingEntries...)
	return nil
}

// saveEntries validates and persists GL entries.
//
// Maps to: save_entries() in general_ledger.py (lines 406-421)
func (e *Engine) saveEntries(glMap []GLEntry, opts PostingOptions) error {
	if e.GLStore == nil {
		return nil
	}

	// Process debit/credit difference (rounding)
	if err := e.processDebitCreditDifference(&glMap); err != nil {
		return err
	}

	// Validate freezing date
	if err := e.checkFreezingDate(glMap, opts.AdvAdj); err != nil {
		return err
	}

	// Save all entries
	return e.GLStore.SaveBatch(glMap)
}

// processDebitCreditDifference handles rounding differences.
// If total debit != total credit within allowance, creates a round-off entry.
//
// Maps to: process_debit_credit_difference() in general_ledger.py (lines 469-499)
func (e *Engine) processDebitCreditDifference(glMap *[]GLEntry) error {
	if len(*glMap) == 0 {
		return nil
	}

	precision := 2
	diff := getDebitCreditDifference(*glMap, precision)
	allowance := getDebitCreditAllowance((*glMap)[0].VoucherType, precision)

	if absFloat(diff) > allowance {
		return fmt.Errorf(
			"debit and credit not equal for %s #%s. Difference is %.2f",
			(*glMap)[0].VoucherType,
			(*glMap)[0].VoucherNo,
			diff,
		)
	}

	// Create round-off entry if difference is significant but within allowance
	minDiff := 1.0 / pow10(precision)
	if absFloat(diff) >= minDiff {
		if err := e.makeRoundOffGLE(glMap, diff, precision); err != nil {
			return err
		}
	}

	return nil
}

// getDebitCreditDifference calculates total debit - total credit.
//
// Maps to: get_debit_credit_difference() in general_ledger.py (lines 502-520)
func getDebitCreditDifference(glMap []GLEntry, precision int) float64 {
	var diff float64
	for _, entry := range glMap {
		diff += Flt(entry.Debit, precision) - Flt(entry.Credit, precision)
	}
	return Flt(diff, precision)
}

// getDebitCreditAllowance returns the maximum allowed difference.
//
// Maps to: get_debit_credit_allowance() in general_ledger.py (lines 523-529)
func getDebitCreditAllowance(voucherType string, precision int) float64 {
	if voucherType == "Journal Entry" || voucherType == "Payment Entry" {
		return 5.0 / pow10(precision)
	}
	return 0.5
}

// makeRoundOffGLE creates a GL entry to balance rounding differences.
//
// Maps to: make_round_off_gle() in general_ledger.py (lines 547+)
func (e *Engine) makeRoundOffGLE(glMap *[]GLEntry, diff float64, precision int) error {
	if e.Company == nil || len(*glMap) == 0 {
		return nil
	}

	company := (*glMap)[0].Company
	roundOffAccount, err := e.Company.GetRoundOffAccount(company)
	if err != nil || roundOffAccount == "" {
		// No round-off account configured, skip
		return nil
	}

	roundOffCostCenter, _ := e.Company.GetRoundOffCostCenter(company)

	entry := (*glMap)[0].Copy()
	entry.Account = roundOffAccount
	entry.CostCenter = roundOffCostCenter
	entry.Remarks = "Round Off"
	entry.AgainstVoucher = ""
	entry.AgainstVoucherType = ""
	entry.PartyType = ""
	entry.Party = ""

	if diff > 0 {
		entry.Debit = 0
		entry.Credit = Flt(diff, precision)
		entry.DebitInAccountCurrency = 0
		entry.CreditInAccountCurrency = Flt(diff, precision)
	} else {
		entry.Debit = Flt(-diff, precision)
		entry.Credit = 0
		entry.DebitInAccountCurrency = Flt(-diff, precision)
		entry.CreditInAccountCurrency = 0
	}

	*glMap = append(*glMap, entry)
	return nil
}

// checkFreezingDate validates against accounts frozen date.
func (e *Engine) checkFreezingDate(glMap []GLEntry, advAdj bool) error {
	if e.Company == nil || len(glMap) == 0 || advAdj {
		return nil
	}

	company := glMap[0].Company
	postingDate := glMap[0].PostingDate

	frozenDate, err := e.Company.GetAccountsFrozenTillDate(company)
	if err != nil {
		return err
	}

	if frozenDate != nil && postingDate.Before(*frozenDate) {
		return NewValidationError(
			ErrAccountsFrozenTill,
			"",
			fmt.Sprintf("Accounts are frozen till %s", frozenDate.Format("2006-01-02")),
		)
	}

	return nil
}

// makeReverseGLEntries creates reversing entries for cancellation.
//
// Maps to: make_reverse_gl_entries() in general_ledger.py
func (e *Engine) makeReverseGLEntries(glMap []GLEntry, opts PostingOptions) error {
	// Get existing entries for the voucher
	if e.GLStore == nil || len(glMap) == 0 {
		return nil
	}

	voucherType := glMap[0].VoucherType
	voucherNo := glMap[0].VoucherNo

	existingEntries, err := e.GLStore.GetByVoucher(voucherType, voucherNo)
	if err != nil {
		return err
	}

	// Create reversed entries
	reversedEntries := make([]GLEntry, len(existingEntries))
	for i, entry := range existingEntries {
		reversed := entry.Copy()
		// Swap debit and credit
		reversed.Debit, reversed.Credit = entry.Credit, entry.Debit
		reversed.DebitInAccountCurrency, reversed.CreditInAccountCurrency =
			entry.CreditInAccountCurrency, entry.DebitInAccountCurrency
		reversed.DebitInTransactionCurrency, reversed.CreditInTransactionCurrency =
			entry.CreditInTransactionCurrency, entry.DebitInTransactionCurrency
		reversed.Remarks = "Cancelled: " + entry.Remarks
		reversedEntries[i] = reversed
	}

	// Mark original entries as cancelled
	if err := e.GLStore.MarkCancelled(voucherType, voucherNo); err != nil {
		return err
	}

	// Save reversed entries
	return e.GLStore.SaveBatch(reversedEntries)
}

// createPaymentLedgerEntries creates payment ledger entries for AR/AP tracking.
//
// Maps to: create_payment_ledger_entry() in accounts/utils.py
func (e *Engine) createPaymentLedgerEntries(glMap []GLEntry, opts PostingOptions) error {
	if e.PaymentStore == nil {
		return nil
	}

	var entries []PaymentLedgerEntry

	for _, gl := range glMap {
		// Only create payment ledger entries for party-based entries
		if gl.PartyType == "" || gl.Party == "" {
			continue
		}

		entry := PaymentLedgerEntry{
			PostingDate:     gl.PostingDate,
			Company:         gl.Company,
			Account:         gl.Account,
			PartyType:       gl.PartyType,
			Party:           gl.Party,
			VoucherType:     gl.VoucherType,
			VoucherNo:       gl.VoucherNo,
			VoucherDetailNo: gl.VoucherDetailNo,
			AgainstVoucherType: gl.AgainstVoucherType,
			AgainstVoucherNo:   gl.AgainstVoucher,
			AccountCurrency:    gl.AccountCurrency,
			Amount:             gl.Debit - gl.Credit,
			AmountInAccountCurrency: gl.DebitInAccountCurrency - gl.CreditInAccountCurrency,
			DueDate:         gl.DueDate,
			FinanceBook:     gl.FinanceBook,
		}

		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		return e.PaymentStore.SaveBatch(entries)
	}

	return nil
}

// Helper functions

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func pow10(n int) float64 {
	result := 1.0
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}
