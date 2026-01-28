// SOLUTION FILE - For mentors only!
// Do not peek unless you're truly stuck.
package glentry

import "time"

type GLEntry struct {
	Name         string
	Account      string
	Company      string
	Debit        float64
	Credit       float64
	VoucherType  string
	VoucherNo    string
	PartyType    string
	Party        string
	PostingDate  time.Time
	CreationDate time.Time
	CostCenter   string
	IsCancelled  bool
	IsOpening    string
}

func NewGLEntry(account string, debit, credit float64) GLEntry {
	return GLEntry{
		Account:      account,
		Debit:        debit,
		Credit:       credit,
		PostingDate:  time.Now(),
		CreationDate: time.Now(),
		IsCancelled:  false,
		IsOpening:    "No",
	}
}

func (e GLEntry) IsValid() bool {
	if e.Account == "" {
		return false
	}

	hasDebit := e.Debit > 0
	hasCredit := e.Credit > 0

	// Entry should have either debit OR credit, not both, not neither
	return (hasDebit && !hasCredit) || (!hasDebit && hasCredit)
}
