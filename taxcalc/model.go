// Package taxcalc implements tax and totals calculation from ERPNext.
// Migrated from: erpnext/controllers/taxes_and_totals.py
//
// This package demonstrates extraction of complex business logic:
// - Line item calculations (rate, discount, amount)
// - Tax calculations with multiple charge types
// - Totals aggregation
package taxcalc

import (
	"encoding/json"
	"math"
)

// ChargeType defines how tax is calculated.
// Maps to: tax.charge_type in ERPNext
type ChargeType string

const (
	// Actual - Fixed amount distributed proportionally across items
	Actual ChargeType = "Actual"
	// OnNetTotal - Percentage applied to item's net amount
	OnNetTotal ChargeType = "On Net Total"
	// OnPreviousRowAmount - Percentage of previous tax row's amount
	OnPreviousRowAmount ChargeType = "On Previous Row Amount"
	// OnPreviousRowTotal - Percentage of previous tax row's running total
	OnPreviousRowTotal ChargeType = "On Previous Row Total"
	// OnItemQuantity - Fixed amount per unit quantity
	OnItemQuantity ChargeType = "On Item Quantity"
)

// TaxCategory defines whether tax is added or deducted.
type TaxCategory string

const (
	Total     TaxCategory = "Total"
	Valuation TaxCategory = "Valuation"
)

// AddDeduct defines whether tax is added or deducted.
type AddDeduct string

const (
	Add    AddDeduct = "Add"
	Deduct AddDeduct = "Deduct"
)

// LineItem represents a single item in an invoice/order.
// Maps to: Sales Invoice Item, Purchase Invoice Item, etc.
type LineItem struct {
	ItemCode    string  // Item identifier
	Description string  // Item description
	Qty         float64 // Quantity
	UOM         string  // Unit of measure

	// Pricing
	PriceListRate      float64 // Original price from price list
	DiscountPercentage float64 // Discount as percentage (0-100)
	DiscountAmount     float64 // Calculated discount amount
	Rate               float64 // Final rate after discount
	Amount             float64 // Rate * Qty

	// Net values (after exclusive tax adjustment)
	NetRate   float64 // Rate excluding taxes (for inclusive pricing)
	NetAmount float64 // NetRate * Qty

	// Base currency values (company currency)
	BaseRate      float64
	BaseAmount    float64
	BaseNetRate   float64
	BaseNetAmount float64

	// Tax info
	ItemTaxRate   string  // JSON map of account -> rate
	ItemTaxAmount float64 // Total tax for this item
}

// TaxRow represents a single tax/charge line.
// Maps to: Sales Taxes and Charges, Purchase Taxes and Charges
type TaxRow struct {
	AccountHead string     // Tax account
	Description string     // Tax description
	ChargeType  ChargeType // How tax is calculated
	Rate        float64    // Tax rate (percentage or fixed amount)
	RowID       int        // Reference to previous row (1-indexed, for OnPreviousRow*)
	Category    TaxCategory
	AddDeductTax AddDeduct

	// Calculated values
	TaxAmount                     float64 // Total tax amount
	TaxAmountAfterDiscountAmount  float64 // Tax after document discount
	Total                         float64 // Running total (net + cumulative tax)
	NetAmount                     float64 // Applicable net amount for this tax

	// Per-item tracking (used during calculation)
	TaxAmountForCurrentItem      float64
	GrandTotalForCurrentItem     float64
	TaxFractionForCurrentItem    float64
	GrandTotalFractionForCurrentItem float64

	// Base currency
	BaseTaxAmount                    float64
	BaseTaxAmountAfterDiscountAmount float64
	BaseTotal                        float64
}

// Document represents an invoice or order with items and taxes.
// Maps to: Sales Invoice, Purchase Invoice, Sales Order, etc.
type Document struct {
	// Currency
	Currency       string  // Transaction currency
	ConversionRate float64 // Exchange rate to company currency

	// Items
	Items []*LineItem

	// Taxes
	Taxes []*TaxRow

	// Discount
	DiscountAmount               float64
	AdditionalDiscountPercentage float64
	ApplyDiscountOn              string // "Net Total" or "Grand Total"

	// Totals
	TotalQty     float64 // Sum of item quantities
	Total        float64 // Sum of item amounts
	BaseTotal    float64
	NetTotal     float64 // Sum of item net amounts
	BaseNetTotal float64
	GrandTotal   float64 // Net total + taxes
	BaseGrandTotal float64

	// Rounding
	RoundingAdjustment     float64
	BaseRoundingAdjustment float64
	RoundedTotal           float64
	BaseRoundedTotal       float64
}

// PrecisionProvider defines precision settings for calculations.
// This abstracts the frappe field precision system.
type PrecisionProvider interface {
	// GetPrecision returns decimal places for a field
	GetPrecision(fieldName string) int
}

// DefaultPrecision provides standard precision (2 decimal places).
type DefaultPrecision struct{}

func (d DefaultPrecision) GetPrecision(fieldName string) int {
	switch fieldName {
	case "rate", "amount", "net_rate", "net_amount", "tax_amount", "total", "grand_total":
		return 2
	case "qty":
		return 3
	case "discount_percentage":
		return 2
	default:
		return 2
	}
}

// ParseItemTaxRate parses the JSON item tax rate map.
// Maps to: json.loads(item_tax_rate) in Python
func ParseItemTaxRate(itemTaxRate string) (map[string]float64, error) {
	if itemTaxRate == "" {
		return make(map[string]float64), nil
	}
	var result map[string]float64
	err := json.Unmarshal([]byte(itemTaxRate), &result)
	return result, err
}

// Round rounds a value to the specified precision.
func Round(value float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(value*multiplier) / multiplier
}

// Flt converts to float and optionally rounds.
// Maps to: frappe.utils.flt() in Python
func Flt(value float64, precision ...int) float64 {
	if len(precision) > 0 {
		return Round(value, precision[0])
	}
	return value
}
