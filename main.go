package main

import (
	"fmt"
)

type PricingRule struct {
	UnitPrice    int
	SpecialPrice struct {
		Quantity int
		Price    int
	}
}

type Checkout struct {
	PricingRules map[string]PricingRule
	ScannedItems map[string]int
}

type ICheckout interface {
	Scan(SKU string) error
	GetTotalPrice() (int, error)
}

func NewCheckout(pricingRules map[string]PricingRule) *Checkout {
	return &Checkout{
		PricingRules: pricingRules,
		ScannedItems: make(map[string]int),
	}
}

func (c *Checkout) Scan(SKU string) error {
	if _, exists := c.PricingRules[SKU]; !exists {
		return fmt.Errorf("invalid SKU: %s", SKU)
	}
	c.ScannedItems[SKU]++
	return nil
}

func (c *Checkout) GetTotalPrice() (int, error) {
	total := 0
	for sku, quantity := range c.ScannedItems {
		rule, exists := c.PricingRules[sku]
		if !exists {
			return 0, fmt.Errorf("no pricing rule for SKU: %s", sku)
		}
		if rule.SpecialPrice.Quantity > 0 && quantity >= rule.SpecialPrice.Quantity {
			specialPriceCount := quantity / rule.SpecialPrice.Quantity
			remainingCount := quantity % rule.SpecialPrice.Quantity
			total += specialPriceCount * rule.SpecialPrice.Price
			total += remainingCount * rule.UnitPrice
		} else {
			total += quantity * rule.UnitPrice
		}
	}
	return total, nil
}

func main() {
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

	checkout := NewCheckout(pricingRules)

	items := []string{"A", "B", "A", "A", "B", "C", "D"}
	for _, item := range items {
		err := checkout.Scan(item)
		if err != nil {
			fmt.Printf("Error scanning item %s: %v\n", item, err)
		}
	}

	total, err := checkout.GetTotalPrice()
	if err != nil {
		fmt.Printf("Error calculating total price: %v\n", err)
	} else {
		fmt.Printf("Total price: %d\n", total)
	}
}
