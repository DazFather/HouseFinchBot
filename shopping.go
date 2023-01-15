package main

import (
	"fmt"
	"regexp"
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

type House struct {
	Chart ShoppingList
}

func newHouse() *House {
	return &House{Chart: newShoppingList()}
}

func newShoppingList() ShoppingList {
	return ShoppingList{items: make(map[string]*ShoppingItem)}
}

var saved = make(map[int64]*House)

func SelectHouse(chatID int64) (house *House, registeredUser bool) {
	house, registeredUser = saved[chatID]
	if !registeredUser {
		house = newHouse()
		saved[chatID] = house
	}
	return
}

var EMOJI = map[string]string{
	"mushroom":   "🍄",
	"apple":      "🍎",
	"pear":       "🍐",
	"orange":     "🍊",
	"lemon":      "🍋",
	"banana":     "🍌",
	"watermelon": "🍉",
	"grape":      "🍇",
	"strawberry": "🍓",
	"blackberry": "🫐",
	"melon":      "🍈",
	"cherry":     "🍒",
	"peach":      "🍑",
	"mango":      "🥭",
	"pineapple":  "🍍",
	"coco":       "🥥",
	"kiwi":       "🥝",
	"tomato":     "🍅",
	"eggplant":   "🍆",
	"avocado":    "🥑",
	"broccoli":   "🥦",
	"lettuce":    "🥬",
	"pickle":     "🥒",
	"spicy":      "🌶",
	"pepper":     "🫑",
	"corn":       "🌽",
	"carrot":     "🥕",
	"olive":      "🫒",
	"garlic":     "🧄",
	"onion":      "🧅",
	"chestnut":   "🌰",
	"potato":     "🥔",
	"bread":      "🍞",
	"baguette":   "🥖",
	"pretzel":    "🥨",
	"cheese":     "🧀",
	"egg":        "🥚",
	"meat":       "🥩",
	"bacon":      "🥓",
	"waffle":     "🧇",
	"panckake":   "🥞",
	"chips":      "🍟",
	"pizza":      "🍕",
	"candy":      "🍬",
	"chocolate":  "🍫",
	"popcorn":    "🍿",
	"cookie":     "🍪",
	"penaut":     "🥜",
	"beer":       "🍺",
	"milk":       "🥛",
	"wine":       "🍷",
	"coffe":      "☕️",
	"tea":        "🫖",
	"ice":        "🧊",
}

var plurals = regexp.MustCompile("e?s$")

func buildName(rawName string) string {
	rawName = strings.ToLower(strings.TrimSpace(rawName))
	var name = strings.ToUpper(string(rawName[0])) + rawName[1:]

	for _, piece := range strings.Split(rawName, " ") {
		if emoji, found := EMOJI[plurals.ReplaceAllLiteralString(piece, "")]; found {
			name = emoji + " " + name
			break
		}
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
		name:      buildName(name),
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
