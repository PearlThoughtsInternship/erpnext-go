package taxcalc

import (
	"errors"
	"fmt"
)

// Calculator errors
var (
	ErrNoItems           = errors.New("no items to calculate")
	ErrInvalidRowID      = errors.New("invalid row reference in tax calculation")
	ErrZeroNetTotal      = errors.New("net total is zero, cannot distribute actual tax")
	ErrNegativeQuantity  = errors.New("quantity cannot be negative")
	ErrInvalidDiscount   = errors.New("discount percentage must be between 0 and 100")
	ErrInvalidConversion = errors.New("conversion rate must be greater than zero")
)

// Calculator performs tax and totals calculations.
// Migrated from: class calculate_taxes_and_totals in taxes_and_totals.py
type Calculator struct {
	doc       *Document
	precision PrecisionProvider
}

// NewCalculator creates a new calculator for a document.
func NewCalculator(doc *Document, precision PrecisionProvider) *Calculator {
	if precision == nil {
		precision = DefaultPrecision{}
	}
	return &Calculator{
		doc:       doc,
		precision: precision,
	}
}

// Calculate performs all calculations on the document.
// Maps to: calculate() method in Python
//
// Python equivalent:
//   def calculate(self):
//       self._calculate()
//       self.set_discount_amount()
//       self.apply_discount_amount()
func (c *Calculator) Calculate() error {
	if len(c.doc.Items) == 0 {
		return ErrNoItems
	}

	// Validate conversion rate
	if err := c.validateConversionRate(); err != nil {
		return err
	}

	// Calculate item values (rate, amount, net_amount)
	if err := c.calculateItemValues(); err != nil {
		return err
	}

	// Initialize taxes
	c.initializeTaxes()

	// Calculate net total
	c.calculateNetTotal()

	// Calculate taxes
	if err := c.calculateTaxes(); err != nil {
		return err
	}

	// Calculate final totals
	c.calculateTotals()

	return nil
}

// validateConversionRate ensures conversion rate is valid.
// Maps to: validate_conversion_rate() in Python
//
// Python equivalent:
//   def validate_conversion_rate(self):
//       if not self.doc.currency or self.doc.currency == company_currency:
//           self.doc.conversion_rate = 1.0
//       self.doc.conversion_rate = flt(self.doc.conversion_rate)
func (c *Calculator) validateConversionRate() error {
	if c.doc.ConversionRate <= 0 {
		c.doc.ConversionRate = 1.0
	}
	return nil
}

// calculateItemValues calculates rate, amount, and net values for each item.
// Maps to: calculate_item_values() in Python (lines 161-236)
//
// Python equivalent:
//   def calculate_item_values(self):
//       for item in self.doc.items:
//           if item.discount_percentage == 100:
//               item.rate = 0.0
//           elif item.price_list_rate:
//               item.rate = flt(item.price_list_rate * (1.0 - (item.discount_percentage / 100.0)))
//               item.discount_amount = item.price_list_rate * (item.discount_percentage / 100.0)
//           item.amount = flt(item.rate * item.qty)
//           item.net_amount = item.amount
func (c *Calculator) calculateItemValues() error {
	ratePrecision := c.precision.GetPrecision("rate")
	amountPrecision := c.precision.GetPrecision("amount")

	for _, item := range c.doc.Items {
		// Validate inputs
		if item.Qty < 0 {
			return fmt.Errorf("%w: item %s has qty %.2f", ErrNegativeQuantity, item.ItemCode, item.Qty)
		}
		if item.DiscountPercentage < 0 || item.DiscountPercentage > 100 {
			return fmt.Errorf("%w: item %s has discount %.2f%%", ErrInvalidDiscount, item.ItemCode, item.DiscountPercentage)
		}

		// Calculate rate from price list rate and discount
		if item.DiscountPercentage == 100 {
			item.Rate = 0.0
			item.DiscountAmount = item.PriceListRate
		} else if item.PriceListRate > 0 {
			// Apply discount percentage
			discountMultiplier := 1.0 - (item.DiscountPercentage / 100.0)
			item.Rate = Flt(item.PriceListRate*discountMultiplier, ratePrecision)
			item.DiscountAmount = Flt(item.PriceListRate*(item.DiscountPercentage/100.0), ratePrecision)
		}

		// If rate not set from price list, use existing rate
		if item.Rate == 0 && item.PriceListRate == 0 {
			// Rate should already be set directly
		}

		// Calculate amount
		item.Amount = Flt(item.Rate*item.Qty, amountPrecision)

		// Net values (before tax adjustments for inclusive pricing)
		item.NetRate = item.Rate
		item.NetAmount = item.Amount

		// Convert to base currency
		c.setInCompanyCurrency(item)

		// Initialize tax amount
		item.ItemTaxAmount = 0.0
	}

	return nil
}

// setInCompanyCurrency converts item values to base currency.
// Maps to: _set_in_company_currency() in Python (lines 237-243)
//
// Python equivalent:
//   def _set_in_company_currency(self, doc, fields):
//       for f in fields:
//           val = flt(flt(doc.get(f)) * self.doc.conversion_rate)
//           doc.set("base_" + f, val)
func (c *Calculator) setInCompanyCurrency(item *LineItem) {
	precision := c.precision.GetPrecision("amount")
	rate := c.doc.ConversionRate

	item.BaseRate = Flt(item.Rate*rate, precision)
	item.BaseAmount = Flt(item.Amount*rate, precision)
	item.BaseNetRate = Flt(item.NetRate*rate, precision)
	item.BaseNetAmount = Flt(item.NetAmount*rate, precision)
}

// initializeTaxes resets tax row values before calculation.
// Maps to: initialize_taxes() in Python (lines 245-269)
//
// Python equivalent:
//   def initialize_taxes(self):
//       for tax in self.doc.get("taxes"):
//           tax_fields = ["total", "tax_amount", "tax_amount_for_current_item", ...]
//           for fieldname in tax_fields:
//               tax.set(fieldname, 0.0)
func (c *Calculator) initializeTaxes() {
	for _, tax := range c.doc.Taxes {
		tax.TaxAmount = 0.0
		tax.TaxAmountAfterDiscountAmount = 0.0
		tax.Total = 0.0
		tax.NetAmount = 0.0
		tax.TaxAmountForCurrentItem = 0.0
		tax.GrandTotalForCurrentItem = 0.0
		tax.TaxFractionForCurrentItem = 0.0
		tax.GrandTotalFractionForCurrentItem = 0.0
		tax.BaseTaxAmount = 0.0
		tax.BaseTaxAmountAfterDiscountAmount = 0.0
		tax.BaseTotal = 0.0
	}
}

// calculateNetTotal sums up item amounts.
// Maps to: calculate_net_total() in Python (lines 369-381)
//
// Python equivalent:
//   def calculate_net_total(self):
//       self.doc.total_qty = self.doc.total = self.doc.net_total = 0.0
//       for item in self._items:
//           self.doc.total += item.amount
//           self.doc.total_qty += item.qty
//           self.doc.net_total += item.net_amount
func (c *Calculator) calculateNetTotal() {
	c.doc.TotalQty = 0.0
	c.doc.Total = 0.0
	c.doc.BaseTotal = 0.0
	c.doc.NetTotal = 0.0
	c.doc.BaseNetTotal = 0.0

	for _, item := range c.doc.Items {
		c.doc.TotalQty += item.Qty
		c.doc.Total += item.Amount
		c.doc.BaseTotal += item.BaseAmount
		c.doc.NetTotal += item.NetAmount
		c.doc.BaseNetTotal += item.BaseNetAmount
	}

	// Round totals
	precision := c.precision.GetPrecision("total")
	c.doc.Total = Flt(c.doc.Total, precision)
	c.doc.BaseTotal = Flt(c.doc.BaseTotal, precision)
	c.doc.NetTotal = Flt(c.doc.NetTotal, precision)
	c.doc.BaseNetTotal = Flt(c.doc.BaseNetTotal, precision)
}

// calculateTaxes calculates tax amounts for each tax row.
// Maps to: calculate_taxes() in Python (lines 394-488)
//
// Python equivalent:
//   def calculate_taxes(self):
//       for n, item in enumerate(self._items):
//           for i, tax in enumerate(doc.taxes):
//               current_tax_amount = self.get_current_tax_amount(item, tax, item_tax_map)
//               tax.tax_amount += current_tax_amount
func (c *Calculator) calculateTaxes() error {
	if len(c.doc.Taxes) == 0 {
		return nil
	}

	taxPrecision := c.precision.GetPrecision("tax_amount")

	// Track actual tax amounts for proportional distribution
	actualTaxAmounts := make(map[int]float64)
	for i, tax := range c.doc.Taxes {
		if tax.ChargeType == Actual {
			actualTaxAmounts[i] = Flt(tax.Rate, taxPrecision) // Rate holds the actual amount
		}
	}

	// Process each item
	for itemIdx, item := range c.doc.Items {
		itemTaxMap, _ := ParseItemTaxRate(item.ItemTaxRate)

		for taxIdx, tax := range c.doc.Taxes {
			// Calculate tax amount for this item
			currentTaxAmount, err := c.getCurrentTaxAmount(item, tax, taxIdx, itemTaxMap)
			if err != nil {
				return err
			}

			// Adjust for actual tax distribution
			if tax.ChargeType == Actual {
				actualTaxAmounts[taxIdx] -= currentTaxAmount
				// Add remainder to last item
				if itemIdx == len(c.doc.Items)-1 {
					currentTaxAmount += actualTaxAmounts[taxIdx]
				}
			}

			// Accumulate tax amount
			tax.TaxAmount += currentTaxAmount
			tax.TaxAmountAfterDiscountAmount += currentTaxAmount

			// Track for current item (used by OnPreviousRow*)
			tax.TaxAmountForCurrentItem = currentTaxAmount

			// Calculate running total for current item
			adjustedTaxAmount := c.getAdjustedTaxAmount(currentTaxAmount, tax)
			if taxIdx == 0 {
				tax.GrandTotalForCurrentItem = item.NetAmount + adjustedTaxAmount
			} else {
				tax.GrandTotalForCurrentItem = c.doc.Taxes[taxIdx-1].GrandTotalForCurrentItem + adjustedTaxAmount
			}
		}
	}

	// Round and calculate cumulative totals
	for taxIdx, tax := range c.doc.Taxes {
		tax.TaxAmount = Flt(tax.TaxAmount, taxPrecision)
		tax.TaxAmountAfterDiscountAmount = Flt(tax.TaxAmountAfterDiscountAmount, taxPrecision)

		// Set cumulative total
		c.setCumulativeTotal(taxIdx, tax)

		// Convert to base currency
		rate := c.doc.ConversionRate
		tax.BaseTaxAmount = Flt(tax.TaxAmount*rate, taxPrecision)
		tax.BaseTaxAmountAfterDiscountAmount = Flt(tax.TaxAmountAfterDiscountAmount*rate, taxPrecision)
		tax.BaseTotal = Flt(tax.Total*rate, taxPrecision)
	}

	return nil
}

// getCurrentTaxAmount calculates tax for a single item.
// Maps to: get_current_tax_amount() in Python (lines 566-594)
//
// Python equivalent:
//   def get_current_tax_amount(self, item, tax, item_tax_map):
//       tax_rate = self._get_tax_rate(tax, item_tax_map)
//       if tax.charge_type == "Actual":
//           current_tax_amount = item.net_amount * actual / self.doc.net_total
//       elif tax.charge_type == "On Net Total":
//           current_tax_amount = (tax_rate / 100.0) * item.net_amount
//       elif tax.charge_type == "On Previous Row Amount":
//           current_tax_amount = (tax_rate / 100.0) * prev_row.tax_amount_for_current_item
//       elif tax.charge_type == "On Previous Row Total":
//           current_tax_amount = (tax_rate / 100.0) * prev_row.grand_total_for_current_item
//       elif tax.charge_type == "On Item Quantity":
//           current_tax_amount = tax_rate * item.qty
func (c *Calculator) getCurrentTaxAmount(item *LineItem, tax *TaxRow, taxIdx int, itemTaxMap map[string]float64) (float64, error) {
	// Get applicable tax rate (item-specific or default)
	taxRate := c.getTaxRate(tax, itemTaxMap)

	var currentTaxAmount float64

	switch tax.ChargeType {
	case Actual:
		// Distribute actual amount proportionally by net amount
		if c.doc.NetTotal == 0 {
			currentTaxAmount = 0.0
		} else {
			actualAmount := tax.Rate // For Actual type, Rate holds the fixed amount
			currentTaxAmount = (item.NetAmount * actualAmount) / c.doc.NetTotal
		}

	case OnNetTotal:
		// Percentage of item's net amount
		currentTaxAmount = (taxRate / 100.0) * item.NetAmount

	case OnPreviousRowAmount:
		// Percentage of previous tax row's tax amount
		if tax.RowID < 1 || tax.RowID > len(c.doc.Taxes) {
			return 0, fmt.Errorf("%w: row_id %d for tax %s", ErrInvalidRowID, tax.RowID, tax.AccountHead)
		}
		prevTax := c.doc.Taxes[tax.RowID-1]
		currentTaxAmount = (taxRate / 100.0) * prevTax.TaxAmountForCurrentItem

	case OnPreviousRowTotal:
		// Percentage of previous tax row's running total
		if tax.RowID < 1 || tax.RowID > len(c.doc.Taxes) {
			return 0, fmt.Errorf("%w: row_id %d for tax %s", ErrInvalidRowID, tax.RowID, tax.AccountHead)
		}
		prevTax := c.doc.Taxes[tax.RowID-1]
		currentTaxAmount = (taxRate / 100.0) * prevTax.GrandTotalForCurrentItem

	case OnItemQuantity:
		// Fixed amount per unit
		currentTaxAmount = taxRate * item.Qty

	default:
		currentTaxAmount = 0.0
	}

	return currentTaxAmount, nil
}

// getTaxRate returns the applicable tax rate for an item.
// Maps to: _get_tax_rate() in Python (lines 363-367)
//
// Python equivalent:
//   def _get_tax_rate(self, tax, item_tax_map):
//       if tax.account_head in item_tax_map:
//           return flt(item_tax_map.get(tax.account_head))
//       else:
//           return tax.rate
func (c *Calculator) getTaxRate(tax *TaxRow, itemTaxMap map[string]float64) float64 {
	if rate, ok := itemTaxMap[tax.AccountHead]; ok {
		return rate
	}
	return tax.Rate
}

// getAdjustedTaxAmount adjusts tax for valuation or deduction.
// Maps to: get_tax_amount_if_for_valuation_or_deduction() in Python (lines 543-555)
func (c *Calculator) getAdjustedTaxAmount(taxAmount float64, tax *TaxRow) float64 {
	// Valuation taxes don't add to total
	if tax.Category == Valuation {
		return 0.0
	}

	// Deduction taxes are subtracted
	if tax.AddDeductTax == Deduct {
		return -taxAmount
	}

	return taxAmount
}

// setCumulativeTotal sets the running total for a tax row.
// Maps to: set_cumulative_total() in Python (lines 557-564)
//
// Python equivalent:
//   def set_cumulative_total(self, row_idx, tax):
//       if row_idx == 0:
//           tax.total = flt(self.doc.net_total + tax_amount)
//       else:
//           tax.total = flt(self.doc.get("taxes")[row_idx - 1].total + tax_amount)
func (c *Calculator) setCumulativeTotal(taxIdx int, tax *TaxRow) {
	precision := c.precision.GetPrecision("total")
	taxAmount := c.getAdjustedTaxAmount(tax.TaxAmountAfterDiscountAmount, tax)

	if taxIdx == 0 {
		tax.Total = Flt(c.doc.NetTotal+taxAmount, precision)
	} else {
		tax.Total = Flt(c.doc.Taxes[taxIdx-1].Total+taxAmount, precision)
	}
}

// calculateTotals calculates final grand total.
// Maps to: calculate_totals() in Python
func (c *Calculator) calculateTotals() {
	precision := c.precision.GetPrecision("grand_total")

	if len(c.doc.Taxes) > 0 {
		lastTax := c.doc.Taxes[len(c.doc.Taxes)-1]
		c.doc.GrandTotal = Flt(lastTax.Total, precision)
		c.doc.BaseGrandTotal = Flt(lastTax.BaseTotal, precision)
	} else {
		c.doc.GrandTotal = Flt(c.doc.NetTotal, precision)
		c.doc.BaseGrandTotal = Flt(c.doc.BaseNetTotal, precision)
	}
}

// GetTaxBreakup returns tax amounts by account for display.
func (c *Calculator) GetTaxBreakup() map[string]float64 {
	breakup := make(map[string]float64)
	for _, tax := range c.doc.Taxes {
		breakup[tax.AccountHead] = tax.TaxAmount
	}
	return breakup
}
