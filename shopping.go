package main

import (
	"fmt"
	"strings"
)

type ShoppingItem struct {
	quantity  uint8
	name      string
	important bool
}

type ShoppingList struct {
	items map[string]*ShoppingItem
}

func newShoppingList() ShoppingList {
	return ShoppingList{items: make(map[string]*ShoppingItem)}
}

func buildItemName(rawName string) string {
	rawName = strings.ToLower(strings.TrimSpace(rawName))
	var name = strings.ToUpper(string(rawName[0])) + rawName[1:]

	if emoji := findEmoji(rawName); emoji != "" {
		name = emoji + " " + name
	}
	return name
}

func (item ShoppingItem) Caption() (caption string) {
	if item.important {
		caption = "<b>" + item.name + "</b>"
	} else {
		caption = item.name
	}

	return fmt.Sprint(" 🔻 ", caption, " x", item.quantity)
}

func (chart ShoppingList) GetItem(itemID string) *ShoppingItem {
	return chart.items[itemID]
}

func (chart *ShoppingList) Save(name string, quantity uint8, important bool) (itemID string, item *ShoppingItem) {
	// Generate itemID
	itemID = name
	item = &ShoppingItem{
		quantity:  quantity,
		name:      buildItemName(name),
		important: important,
	}

	chart.items[itemID] = item
	return
}

func (chart ShoppingList) ForEach(do func(itemID string, item *ShoppingItem)) {
	for itemID, item := range chart.items {
		do(itemID, item)
	}
}

func (chart ShoppingList) Items() int {
	return len(chart.items)
}

func (chart *ShoppingList) DropAll() {
	chart.items = make(map[string]*ShoppingItem)
}

func (chart *ShoppingList) DropItem(itemID string) bool {
	_, ok := chart.items[itemID]
	if ok {
		delete(chart.items, itemID)
	}

	return ok
}

func (chart *ShoppingList) AddQuantity(itemID string, modifier uint8) int {
	item := chart.items[itemID]
	item.quantity += modifier
	return int(item.quantity)
}

func (chart *ShoppingList) RemoveQuantity(itemID string, modifier uint8) (left int) {
	var item = chart.items[itemID]
	if item == nil {
		return 0
	}

	left = int(item.quantity) - int(modifier)
	if modifier >= item.quantity {
		chart.DropItem(itemID)
		return
	}

	item.quantity -= modifier
	return
}
