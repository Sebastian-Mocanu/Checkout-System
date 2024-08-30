package main

import (
	"fmt"
	"testing"
)

func TestCheckout(t *testing.T) {
	type testCase struct {
		name               string
		items              []string
		expectedTotal      int
		expectedPromotions map[string]int
		errString          string
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
			name:               "Single item",
			items:              []string{"A"},
			expectedTotal:      50,
			expectedPromotions: map[string]int{},
			errString:          "",
		},
		{
			name:               "Multiple items without special pricing",
			items:              []string{"A", "B", "C", "D"},
			expectedTotal:      115,
			expectedPromotions: map[string]int{},
			errString:          "",
		},
		{
			name:               "Special pricing for A",
			items:              []string{"A", "A", "A"},
			expectedTotal:      130,
			expectedPromotions: map[string]int{"A": 1},
			errString:          "",
		},
		{
			name:               "Special pricing for B",
			items:              []string{"B", "B"},
			expectedTotal:      45,
			expectedPromotions: map[string]int{"B": 1},
			errString:          "",
		},
		{
			name:               "Mixed items with special pricing",
			items:              []string{"A", "A", "B", "B", "A", "C", "D"},
			expectedTotal:      130 + 45 + 20 + 15,
			expectedPromotions: map[string]int{"A": 1, "B": 1},
			errString:          "",
		},
		{
			name:               "Invalid item",
			items:              []string{"A", "E"},
			expectedTotal:      0,
			expectedPromotions: map[string]int{},
			errString:          "invalid SKU: E",
		},
		{
			name:               "Multiple special offers for A",
			items:              []string{"A", "A", "A", "A", "A", "A", "A"},
			expectedTotal:      130 + 130 + 50,
			expectedPromotions: map[string]int{"A": 2},
			errString:          "",
		},
		{
			name:               "No items",
			items:              []string{},
			expectedTotal:      0,
			expectedPromotions: map[string]int{},
			errString:          "",
		},
		{
			name:               "Just below special offer threshold",
			items:              []string{"A", "A", "B"},
			expectedTotal:      130,
			expectedPromotions: map[string]int{},
			errString:          "",
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
		var promotions map[string]int
		if err == nil {
			output, err = checkout.GetTotalPrice()
			promotions = checkout.GetAppliedPromotions()
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
		} else if output != test.expectedTotal {
			failCount++
			t.Errorf(`
---------------------------------
Test Failed: %s
 items: %v
 expected total: %d
 actual total: %d
`, test.name, test.items, test.expectedTotal, output)
		} else if !mapsEqual(promotions, test.expectedPromotions) {
			failCount++
			t.Errorf(`
---------------------------------
Test Failed: %s
 items: %v
 expected promotions: %v
 actual promotions: %v
`, test.name, test.items, test.expectedPromotions, promotions)
		} else {
			passCount++
			fmt.Printf(`
---------------------------------
Test Passed: %s
 items: %v
 expected total: %d
 actual total: %d
 expected promotions: %v
 actual promotions: %v
`, test.name, test.items, test.expectedTotal, output, test.expectedPromotions, promotions)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

func mapsEqual(map1, map2 map[string]int) bool {
	if len(map1) != len(map2) {
		return false
	}
	for k, v := range map1 {
		if map2[k] != v {
			return false
		}
	}
	return true
}
