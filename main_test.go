package main

import (
	"fmt"
	"testing"
)

func TestCheckout(t *testing.T) {
	type testCase struct {
		name      string
		items     []string
		expected  int
		errString string
	}

	pricingRules := map[string]PricingRule{
		"A": {UnitPrice: 50, SpecialPrice: struct {
			Quantity int
			Price    int
		}{Quantity: 3, Price: 130}},
		"B": {UnitPrice: 30, SpecialPrice: struct {
			Quantity int
			Price    int
		}{Quantity: 2, Price: 45}},
		"C": {UnitPrice: 20},
		"D": {UnitPrice: 15},
	}

	tests := []testCase{
		{
			name:      "Single item",
			items:     []string{"A"},
			expected:  50,
			errString: "",
		},
		{
			name:      "Multiple items without special pricing",
			items:     []string{"A", "B", "C", "D"},
			expected:  115,
			errString: "",
		},
		{
			name:      "Special pricing for A",
			items:     []string{"A", "A", "A"},
			expected:  130,
			errString: "",
		},
		{
			name:      "Special pricing for B",
			items:     []string{"B", "B"},
			expected:  45,
			errString: "",
		},
		{
			name:      "Mixed items with special pricing",
			items:     []string{"A", "A", "B", "B", "A", "C", "D"},
			expected:  130 + 45 + 20 + 15,
			errString: "",
		},
		{
			name:      "Invalid item",
			items:     []string{"A", "E"},
			expected:  0,
			errString: "invalid SKU: E",
		},
	}

	passCount := 0
	failCount := 0

	for _, test := range tests {
		checkout := NewCheckout(pricingRules)
		var err error

		for _, item := range test.items {
			err = checkout.Scan(item)
			if err != nil {
				break
			}
		}

		var output int
		if err == nil {
			output, err = checkout.GetTotalPrice()
		}

		if test.errString != "" && err == nil {
			failCount++
			t.Errorf(`
---------------------------------
Test Failed: %s
 items: %v
 expected err: %v
 actual err: none
`, test.name, test.items, test.errString)
		} else if test.errString == "" && err != nil {
			failCount++
			t.Errorf(`
---------------------------------
Test Failed: %s
 items: %v
 expected err: none
 actual err: %v
`, test.name, test.items, err)
		} else if test.errString != "" && err != nil && err.Error() != test.errString {
			failCount++
			t.Errorf(`
---------------------------------
Test Failed: %s
 items: %v
 expected err: %v
 actual err: %v
`, test.name, test.items, test.errString, err)
		} else if output != test.expected {
			failCount++
			t.Errorf(`
---------------------------------
Test Failed: %s
 items: %v
 expected: %d
 actual: %d
`, test.name, test.items, test.expected, output)
		} else {
			passCount++
			fmt.Printf(`
---------------------------------
Test Passed: %s
 items: %v
 expected: %d
 actual: %d
`, test.name, test.items, test.expected, output)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}
