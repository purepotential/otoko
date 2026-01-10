package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

type valueCmd struct {
	Currency string `kong:"help='Currency rate for total conversion (e.g., USD, EUR, GBP)',default=USD"`
}

func (cmd *valueCmd) Run(c *Client) error {
	items, err := c.GetCollection(c.Fan.ID)
	if err != nil {
		return err
	}

	cost, err := c.Value(c.Fan, items, cmd.Currency)
	if err != nil {
		return err
	}
	fmt.Printf("%d items = ~%.0f%s\n",
		len(items), math.Ceil(cost), cmd.Currency)
	return nil
}

type listCmd struct {
	JSON bool `kong:"short=j,help='Output all known data using a JSON object array'"`
}

func (cmd *listCmd) Run(c *Client) error {
	items, err := c.GetCollection(c.Fan.ID)
	if err != nil {
		return err
	}

	if cmd.JSON {
		enc := json.NewEncoder(os.Stderr)
		enc.SetIndent("", "  ")
		return enc.Encode(items)
	}

	for _, item := range items {
		fmt.Printf("%-10d %s %s - %s\n",
			item.ID,
			item.Purchased.Format(time.DateOnly),
			item.BandName, item.Title)
	}
	return nil
}
