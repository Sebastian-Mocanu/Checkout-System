package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	Scan(SKU string) (err error)
	GetTotalPrice() (totalPrice int, err error)
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

func (c *Checkout) GetAppliedPromotions() map[string]int {
	appliedPromotions := make(map[string]int)
	for sku, quantity := range c.ScannedItems {
		rule := c.PricingRules[sku]
		if rule.SpecialPrice.Quantity > 0 && quantity >= rule.SpecialPrice.Quantity {
			appliedPromotions[sku] = quantity / rule.SpecialPrice.Quantity
		}
	}
	return appliedPromotions
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pricingRules := make(map[string]PricingRule)

	fmt.Println("Welcome to the Supermarket Checkout System!")

	fmt.Println("\nLet's set up the product catalogue.")
	for {
		fmt.Print("Enter product SKU (or press Enter to finish): ")
		scanner.Scan()
		sku := scanner.Text()
		if sku == "" {
			break
		}

		fmt.Printf("Enter unit price for %s: ", sku)
		scanner.Scan()
		unitPrice, _ := strconv.Atoi(scanner.Text())

		rule := PricingRule{UnitPrice: unitPrice}

		fmt.Printf("Is there a special offer for %s? (y/n): ", sku)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) == "y" {
			fmt.Print("Enter special offer quantity: ")
			scanner.Scan()
			quantity, _ := strconv.Atoi(scanner.Text())

			fmt.Print("Enter special offer price: ")
			scanner.Scan()
			price, _ := strconv.Atoi(scanner.Text())

			rule.SpecialPrice.Quantity = quantity
			rule.SpecialPrice.Price = price
		}

		pricingRules[sku] = rule
	}

	checkout := NewCheckout(pricingRules)

	fmt.Println("\nNow, let's scan items.")
	for {
		fmt.Print("Scan an item (enter SKU or press Enter to finish): ")
		scanner.Scan()
		sku := scanner.Text()
		if sku == "" {
			break
		}

		err := checkout.Scan(sku)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Item scanned successfully.")
		}
	}

	total, err := checkout.GetTotalPrice()
	if err != nil {
		fmt.Println("Error calculating total price:", err)
		return
	}

	fmt.Println("\nCheckout Summary:")
	fmt.Println("------------------")
	for sku, quantity := range checkout.ScannedItems {
		fmt.Printf("%s: %d x %d = %d\n", sku, quantity, checkout.PricingRules[sku].UnitPrice, quantity*checkout.PricingRules[sku].UnitPrice)
	}

	fmt.Println("\nApplied Promotions:")
	fmt.Println("------------------")
	for sku, count := range checkout.GetAppliedPromotions() {
		rule := checkout.PricingRules[sku]
		saving := (rule.SpecialPrice.Quantity*rule.UnitPrice - rule.SpecialPrice.Price) * count
		fmt.Printf("%s: %d for %d applied %d times. You saved %d\n",
			sku, rule.SpecialPrice.Quantity, rule.SpecialPrice.Price, count, saving)
	}

	fmt.Println("\nTotal Price:", total)
}
