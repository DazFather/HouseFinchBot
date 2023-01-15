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

var shoppingList = map[string]*ShoppingItem{}

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

func buildName(rawName string) string {
	rawName = strings.ToLower(strings.TrimSpace(rawName))
	var name = strings.ToUpper(string(rawName[0])) + rawName[1:]

	for _, piece := range strings.Split(rawName, " ") {
		if emoji, found := EMOJI[strings.TrimSuffix(piece, "s")]; found {
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

func GetItem(itemID string) *ShoppingItem {
	return shoppingList[itemID]
}

func GetAllItems() []ShoppingItem {
	return GenItemList(func(itemID string, item ShoppingItem) ShoppingItem {
		return item
	})
}

func Save(name string, quantity uint8, important bool) (itemID string, item *ShoppingItem) {
	// Generate itemID
	itemID = name
	item = &ShoppingItem{
		quantity:  quantity,
		name:      buildName(name),
		important: important,
	}
	shoppingList[itemID] = item
	return
}

func DropList() {
	shoppingList = make(map[string]*ShoppingItem)
}

func DropItem(itemID string) bool {
	_, ok := shoppingList[itemID]
	if ok {
		delete(shoppingList, itemID)
	}

	return ok
}

func AddQuantity(itemID string, modifier uint8) int {
	item := shoppingList[itemID]
	item.quantity += modifier
	return int(item.quantity)
}

func RemoveQuantity(itemID string, modifier uint8) (left int) {
	var item = shoppingList[itemID]
	if item == nil {
		return 0
	}

	left = int(item.quantity) - int(modifier)

	if modifier >= item.quantity {
		delete(shoppingList, itemID)
		return
	}

	item.quantity -= modifier
	return
}

func GenItemList[T any](mapper func(string, ShoppingItem) T) []T {
	var list = make([]T, len(shoppingList))

	i := 0
	for itemID, item := range shoppingList {
		list[i] = mapper(itemID, *item)
		i++
	}

	return list
}
