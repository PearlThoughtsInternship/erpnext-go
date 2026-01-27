package taxcalc

import (
	"errors"
	"math"
	"testing"
)

// almostEqual checks if two floats are approximately equal.
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// --- Test Calculate Item Values ---

func TestCalculateItemValues(t *testing.T) {
	tests := []struct {
		name              string
		items             []*LineItem
		conversionRate    float64
		wantErr           error
		checkItem         int // index to check
		expectedRate      float64
		expectedAmount    float64
		expectedDiscount  float64
		expectedBaseAmount float64
	}{
		{
			name: "simple item - no discount",
			items: []*LineItem{
				{ItemCode: "ITEM-001", PriceListRate: 100.0, Qty: 5},
			},
			conversionRate:    1.0,
			wantErr:           nil,
			checkItem:         0,
			expectedRate:      100.0,
			expectedAmount:    500.0,
			expectedDiscount:  0.0,
			expectedBaseAmount: 500.0,
		},
		{
			name: "item with 10% discount",
			items: []*LineItem{
				{ItemCode: "ITEM-002", PriceListRate: 100.0, DiscountPercentage: 10, Qty: 2},
			},
			conversionRate:    1.0,
			wantErr:           nil,
			checkItem:         0,
			expectedRate:      90.0,
			expectedAmount:    180.0,
			expectedDiscount:  10.0,
			expectedBaseAmount: 180.0,
		},
		{
			name: "item with 100% discount (free)",
			items: []*LineItem{
				{ItemCode: "ITEM-003", PriceListRate: 50.0, DiscountPercentage: 100, Qty: 3},
			},
			conversionRate:    1.0,
			wantErr:           nil,
			checkItem:         0,
			expectedRate:      0.0,
			expectedAmount:    0.0,
			expectedDiscount:  50.0,
			expectedBaseAmount: 0.0,
		},
		{
			name: "item with currency conversion",
			items: []*LineItem{
				{ItemCode: "ITEM-004", PriceListRate: 100.0, Qty: 1},
			},
			conversionRate:    1.5, // 1 USD = 1.5 base currency
			wantErr:           nil,
			checkItem:         0,
			expectedRate:      100.0,
			expectedAmount:    100.0,
			expectedDiscount:  0.0,
			expectedBaseAmount: 150.0,
		},
		{
			name: "item with fractional quantity",
			items: []*LineItem{
				{ItemCode: "ITEM-005", PriceListRate: 10.0, Qty: 2.5},
			},
			conversionRate:    1.0,
			wantErr:           nil,
			checkItem:         0,
			expectedRate:      10.0,
			expectedAmount:    25.0,
			expectedDiscount:  0.0,
			expectedBaseAmount: 25.0,
		},
		{
			name: "negative quantity - error",
			items: []*LineItem{
				{ItemCode: "ITEM-006", PriceListRate: 100.0, Qty: -1},
			},
			conversionRate: 1.0,
			wantErr:        ErrNegativeQuantity,
		},
		{
			name: "invalid discount (over 100%) - error",
			items: []*LineItem{
				{ItemCode: "ITEM-007", PriceListRate: 100.0, DiscountPercentage: 150, Qty: 1},
			},
			conversionRate: 1.0,
			wantErr:        ErrInvalidDiscount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{
				Items:          tt.items,
				ConversionRate: tt.conversionRate,
			}
			calc := NewCalculator(doc, nil)

			// Need to validate conversion rate first
			calc.validateConversionRate()

			err := calc.calculateItemValues()

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			item := doc.Items[tt.checkItem]
			if !almostEqual(item.Rate, tt.expectedRate, 0.01) {
				t.Errorf("rate: got %.2f, want %.2f", item.Rate, tt.expectedRate)
			}
			if !almostEqual(item.Amount, tt.expectedAmount, 0.01) {
				t.Errorf("amount: got %.2f, want %.2f", item.Amount, tt.expectedAmount)
			}
			if !almostEqual(item.DiscountAmount, tt.expectedDiscount, 0.01) {
				t.Errorf("discount: got %.2f, want %.2f", item.DiscountAmount, tt.expectedDiscount)
			}
			if !almostEqual(item.BaseAmount, tt.expectedBaseAmount, 0.01) {
				t.Errorf("base_amount: got %.2f, want %.2f", item.BaseAmount, tt.expectedBaseAmount)
			}
		})
	}
}

// --- Test Calculate Net Total ---

func TestCalculateNetTotal(t *testing.T) {
	tests := []struct {
		name            string
		items           []*LineItem
		expectedTotalQty float64
		expectedTotal    float64
		expectedNetTotal float64
	}{
		{
			name: "single item",
			items: []*LineItem{
				{Qty: 5, Amount: 500, NetAmount: 500, BaseAmount: 500, BaseNetAmount: 500},
			},
			expectedTotalQty: 5,
			expectedTotal:    500,
			expectedNetTotal: 500,
		},
		{
			name: "multiple items",
			items: []*LineItem{
				{Qty: 2, Amount: 200, NetAmount: 200, BaseAmount: 200, BaseNetAmount: 200},
				{Qty: 3, Amount: 150, NetAmount: 150, BaseAmount: 150, BaseNetAmount: 150},
				{Qty: 1, Amount: 50, NetAmount: 50, BaseAmount: 50, BaseNetAmount: 50},
			},
			expectedTotalQty: 6,
			expectedTotal:    400,
			expectedNetTotal: 400,
		},
		{
			name:            "empty items",
			items:           []*LineItem{},
			expectedTotalQty: 0,
			expectedTotal:    0,
			expectedNetTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &Document{Items: tt.items, ConversionRate: 1.0}
			calc := NewCalculator(doc, nil)
			calc.calculateNetTotal()

			if !almostEqual(doc.TotalQty, tt.expectedTotalQty, 0.01) {
				t.Errorf("total_qty: got %.2f, want %.2f", doc.TotalQty, tt.expectedTotalQty)
			}
			if !almostEqual(doc.Total, tt.expectedTotal, 0.01) {
				t.Errorf("total: got %.2f, want %.2f", doc.Total, tt.expectedTotal)
			}
			if !almostEqual(doc.NetTotal, tt.expectedNetTotal, 0.01) {
				t.Errorf("net_total: got %.2f, want %.2f", doc.NetTotal, tt.expectedNetTotal)
			}
		})
	}
}

// --- Test Tax Calculations ---

func TestCalculateTaxes_OnNetTotal(t *testing.T) {
	// Test: 10% tax on net total
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", Qty: 2, Rate: 100, Amount: 200, NetAmount: 200, BaseNetAmount: 200},
		},
		Taxes: []*TaxRow{
			{AccountHead: "GST", ChargeType: OnNetTotal, Rate: 10},
		},
	}

	calc := NewCalculator(doc, nil)
	calc.calculateNetTotal()
	err := calc.calculateTaxes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tax := doc.Taxes[0]
	expectedTaxAmount := 20.0 // 10% of 200

	if !almostEqual(tax.TaxAmount, expectedTaxAmount, 0.01) {
		t.Errorf("tax_amount: got %.2f, want %.2f", tax.TaxAmount, expectedTaxAmount)
	}
	if !almostEqual(tax.Total, 220.0, 0.01) {
		t.Errorf("total: got %.2f, want %.2f", tax.Total, 220.0)
	}
}

func TestCalculateTaxes_OnPreviousRowAmount(t *testing.T) {
	// Test: CGST 9% on net total, SGST 9% on CGST (cascading)
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", Qty: 1, Rate: 1000, Amount: 1000, NetAmount: 1000, BaseNetAmount: 1000},
		},
		Taxes: []*TaxRow{
			{AccountHead: "CGST", ChargeType: OnNetTotal, Rate: 9},
			{AccountHead: "SGST", ChargeType: OnPreviousRowAmount, Rate: 100, RowID: 1}, // 100% of CGST
		},
	}

	calc := NewCalculator(doc, nil)
	calc.calculateNetTotal()
	err := calc.calculateTaxes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cgst := doc.Taxes[0]
	sgst := doc.Taxes[1]

	// CGST: 9% of 1000 = 90
	if !almostEqual(cgst.TaxAmount, 90.0, 0.01) {
		t.Errorf("CGST tax_amount: got %.2f, want %.2f", cgst.TaxAmount, 90.0)
	}

	// SGST: 100% of CGST (90) = 90
	if !almostEqual(sgst.TaxAmount, 90.0, 0.01) {
		t.Errorf("SGST tax_amount: got %.2f, want %.2f", sgst.TaxAmount, 90.0)
	}

	// Final total: 1000 + 90 + 90 = 1180
	if !almostEqual(sgst.Total, 1180.0, 0.01) {
		t.Errorf("final total: got %.2f, want %.2f", sgst.Total, 1180.0)
	}
}

func TestCalculateTaxes_OnPreviousRowTotal(t *testing.T) {
	// Test: Tax on running total (compound tax)
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", Qty: 1, Rate: 100, Amount: 100, NetAmount: 100, BaseNetAmount: 100},
		},
		Taxes: []*TaxRow{
			{AccountHead: "Tax1", ChargeType: OnNetTotal, Rate: 10},         // 10% on 100 = 10
			{AccountHead: "Tax2", ChargeType: OnPreviousRowTotal, Rate: 5, RowID: 1}, // 5% on 110 = 5.50
		},
	}

	calc := NewCalculator(doc, nil)
	calc.calculateNetTotal()
	err := calc.calculateTaxes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tax1 := doc.Taxes[0]
	tax2 := doc.Taxes[1]

	// Tax1: 10% of 100 = 10, Total = 110
	if !almostEqual(tax1.TaxAmount, 10.0, 0.01) {
		t.Errorf("Tax1 amount: got %.2f, want %.2f", tax1.TaxAmount, 10.0)
	}
	if !almostEqual(tax1.Total, 110.0, 0.01) {
		t.Errorf("Tax1 total: got %.2f, want %.2f", tax1.Total, 110.0)
	}

	// Tax2: 5% of 110 = 5.50
	if !almostEqual(tax2.TaxAmount, 5.50, 0.01) {
		t.Errorf("Tax2 amount: got %.2f, want %.2f", tax2.TaxAmount, 5.50)
	}
	if !almostEqual(tax2.Total, 115.50, 0.01) {
		t.Errorf("Tax2 total: got %.2f, want %.2f", tax2.Total, 115.50)
	}
}

func TestCalculateTaxes_Actual(t *testing.T) {
	// Test: Fixed amount tax distributed proportionally
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", Qty: 1, Rate: 300, Amount: 300, NetAmount: 300, BaseNetAmount: 300},
			{ItemCode: "ITEM-002", Qty: 1, Rate: 200, Amount: 200, NetAmount: 200, BaseNetAmount: 200},
		},
		Taxes: []*TaxRow{
			{AccountHead: "Shipping", ChargeType: Actual, Rate: 50}, // Fixed 50 shipping
		},
	}

	calc := NewCalculator(doc, nil)
	calc.calculateNetTotal() // NetTotal = 500
	err := calc.calculateTaxes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tax := doc.Taxes[0]

	// Total shipping = 50, distributed: Item1 gets 30 (60%), Item2 gets 20 (40%)
	if !almostEqual(tax.TaxAmount, 50.0, 0.01) {
		t.Errorf("tax_amount: got %.2f, want %.2f", tax.TaxAmount, 50.0)
	}
	if !almostEqual(tax.Total, 550.0, 0.01) {
		t.Errorf("total: got %.2f, want %.2f", tax.Total, 550.0)
	}
}

func TestCalculateTaxes_OnItemQuantity(t *testing.T) {
	// Test: Tax per unit quantity
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", Qty: 5, Rate: 100, Amount: 500, NetAmount: 500, BaseNetAmount: 500},
		},
		Taxes: []*TaxRow{
			{AccountHead: "PerUnit", ChargeType: OnItemQuantity, Rate: 2}, // $2 per unit
		},
	}

	calc := NewCalculator(doc, nil)
	calc.calculateNetTotal()
	err := calc.calculateTaxes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tax := doc.Taxes[0]
	// 5 units * $2 = $10
	if !almostEqual(tax.TaxAmount, 10.0, 0.01) {
		t.Errorf("tax_amount: got %.2f, want %.2f", tax.TaxAmount, 10.0)
	}
}

func TestCalculateTaxes_DeductTax(t *testing.T) {
	// Test: Tax with deduction (discount-like)
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", Qty: 1, Rate: 100, Amount: 100, NetAmount: 100, BaseNetAmount: 100},
		},
		Taxes: []*TaxRow{
			{AccountHead: "Discount", ChargeType: OnNetTotal, Rate: 10, AddDeductTax: Deduct},
		},
	}

	calc := NewCalculator(doc, nil)
	calc.calculateNetTotal()
	err := calc.calculateTaxes()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tax := doc.Taxes[0]
	// 10% of 100 = 10, but deducted
	if !almostEqual(tax.TaxAmount, 10.0, 0.01) {
		t.Errorf("tax_amount: got %.2f, want %.2f", tax.TaxAmount, 10.0)
	}
	// Total should be 100 - 10 = 90
	if !almostEqual(tax.Total, 90.0, 0.01) {
		t.Errorf("total: got %.2f, want %.2f", tax.Total, 90.0)
	}
}

func TestCalculateTaxes_InvalidRowID(t *testing.T) {
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", Qty: 1, Rate: 100, Amount: 100, NetAmount: 100, BaseNetAmount: 100},
		},
		Taxes: []*TaxRow{
			{AccountHead: "Bad", ChargeType: OnPreviousRowAmount, Rate: 10, RowID: 5}, // Invalid!
		},
	}

	calc := NewCalculator(doc, nil)
	calc.calculateNetTotal()
	err := calc.calculateTaxes()

	if !errors.Is(err, ErrInvalidRowID) {
		t.Errorf("expected error %v, got %v", ErrInvalidRowID, err)
	}
}

// --- Test Full Calculation ---

func TestCalculate_FullInvoice(t *testing.T) {
	// Simulate a real invoice:
	// - 2 items with discounts
	// - GST 18% (split CGST 9% + SGST 9%)
	// - Shipping charge $50
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "LAPTOP", PriceListRate: 1000, DiscountPercentage: 10, Qty: 2},
			{ItemCode: "MOUSE", PriceListRate: 50, DiscountPercentage: 0, Qty: 5},
		},
		Taxes: []*TaxRow{
			{AccountHead: "CGST", ChargeType: OnNetTotal, Rate: 9},
			{AccountHead: "SGST", ChargeType: OnNetTotal, Rate: 9},
			{AccountHead: "Shipping", ChargeType: Actual, Rate: 50},
		},
	}

	calc := NewCalculator(doc, nil)
	err := calc.Calculate()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Item calculations:
	// LAPTOP: 1000 * 0.9 * 2 = 1800
	// MOUSE: 50 * 1 * 5 = 250
	// Net Total = 2050

	if !almostEqual(doc.NetTotal, 2050.0, 0.01) {
		t.Errorf("net_total: got %.2f, want %.2f", doc.NetTotal, 2050.0)
	}

	// Tax calculations:
	// CGST: 9% of 2050 = 184.50
	// SGST: 9% of 2050 = 184.50
	// Shipping: 50 (actual)
	// Total Tax = 419

	cgst := doc.Taxes[0]
	sgst := doc.Taxes[1]
	shipping := doc.Taxes[2]

	if !almostEqual(cgst.TaxAmount, 184.50, 0.01) {
		t.Errorf("CGST: got %.2f, want %.2f", cgst.TaxAmount, 184.50)
	}
	if !almostEqual(sgst.TaxAmount, 184.50, 0.01) {
		t.Errorf("SGST: got %.2f, want %.2f", sgst.TaxAmount, 184.50)
	}
	if !almostEqual(shipping.TaxAmount, 50.0, 0.01) {
		t.Errorf("Shipping: got %.2f, want %.2f", shipping.TaxAmount, 50.0)
	}

	// Grand Total = 2050 + 184.50 + 184.50 + 50 = 2469
	if !almostEqual(doc.GrandTotal, 2469.0, 0.01) {
		t.Errorf("grand_total: got %.2f, want %.2f", doc.GrandTotal, 2469.0)
	}
}

func TestCalculate_WithCurrencyConversion(t *testing.T) {
	// Test multi-currency: USD to INR (rate 83)
	doc := &Document{
		Currency:       "USD",
		ConversionRate: 83.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", PriceListRate: 100, Qty: 1}, // $100
		},
		Taxes: []*TaxRow{
			{AccountHead: "GST", ChargeType: OnNetTotal, Rate: 18},
		},
	}

	calc := NewCalculator(doc, nil)
	err := calc.Calculate()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Net Total: $100 (INR 8300)
	if !almostEqual(doc.NetTotal, 100.0, 0.01) {
		t.Errorf("net_total: got %.2f, want %.2f", doc.NetTotal, 100.0)
	}
	if !almostEqual(doc.BaseNetTotal, 8300.0, 0.01) {
		t.Errorf("base_net_total: got %.2f, want %.2f", doc.BaseNetTotal, 8300.0)
	}

	// GST: 18% of $100 = $18 (INR 1494)
	gst := doc.Taxes[0]
	if !almostEqual(gst.TaxAmount, 18.0, 0.01) {
		t.Errorf("tax_amount: got %.2f, want %.2f", gst.TaxAmount, 18.0)
	}
	if !almostEqual(gst.BaseTaxAmount, 1494.0, 0.01) {
		t.Errorf("base_tax_amount: got %.2f, want %.2f", gst.BaseTaxAmount, 1494.0)
	}

	// Grand Total: $118 (INR 9794)
	if !almostEqual(doc.GrandTotal, 118.0, 0.01) {
		t.Errorf("grand_total: got %.2f, want %.2f", doc.GrandTotal, 118.0)
	}
	if !almostEqual(doc.BaseGrandTotal, 9794.0, 0.01) {
		t.Errorf("base_grand_total: got %.2f, want %.2f", doc.BaseGrandTotal, 9794.0)
	}
}

func TestCalculate_NoItems(t *testing.T) {
	doc := &Document{
		ConversionRate: 1.0,
		Items:          []*LineItem{},
	}

	calc := NewCalculator(doc, nil)
	err := calc.Calculate()

	if !errors.Is(err, ErrNoItems) {
		t.Errorf("expected error %v, got %v", ErrNoItems, err)
	}
}

func TestCalculate_NoTaxes(t *testing.T) {
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "ITEM-001", PriceListRate: 100, Qty: 1},
		},
		Taxes: []*TaxRow{}, // No taxes
	}

	calc := NewCalculator(doc, nil)
	err := calc.Calculate()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Grand Total should equal Net Total when no taxes
	if !almostEqual(doc.GrandTotal, 100.0, 0.01) {
		t.Errorf("grand_total: got %.2f, want %.2f", doc.GrandTotal, 100.0)
	}
}

// --- Test Item Tax Rate Override ---

func TestCalculate_ItemSpecificTaxRate(t *testing.T) {
	// Test: Item with different tax rate (tax-exempt item)
	doc := &Document{
		ConversionRate: 1.0,
		Items: []*LineItem{
			{ItemCode: "TAXABLE", PriceListRate: 100, Qty: 1, ItemTaxRate: ""},                      // Normal tax
			{ItemCode: "EXEMPT", PriceListRate: 100, Qty: 1, ItemTaxRate: `{"GST Account": 0}`},    // 0% tax
		},
		Taxes: []*TaxRow{
			{AccountHead: "GST Account", ChargeType: OnNetTotal, Rate: 18},
		},
	}

	calc := NewCalculator(doc, nil)
	err := calc.Calculate()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only TAXABLE item gets 18% = 18
	// EXEMPT item gets 0% = 0
	// Total tax = 18
	gst := doc.Taxes[0]
	if !almostEqual(gst.TaxAmount, 18.0, 0.01) {
		t.Errorf("tax_amount: got %.2f, want %.2f", gst.TaxAmount, 18.0)
	}
}
