package main

import (
	"fmt"

	"github.com/DazFather/parrbot/message"
	"github.com/DazFather/parrbot/robot"
	"github.com/DazFather/parrbot/tgui"
)

var startHandler = tgui.Sender(message.Text{fmt.Sprint(
	"👋 Welcome, I'm <b>House Finch Bot</b>🐦\n",
	"Your personal robo-passerine that will help you take care of the house\n",
	"\nFor now I only know how to keep track of your shopping list but I'm learning 😅\n",
	"Just type the name of the item and it will be auto<i>magically</i>✨ added to the /list",
), defaultOpt()})

func listHandler(b *robot.Bot, u *message.Update) message.Any {
	var chart *ShoppingList = nil
	if house, registered := SelectHouse(b.ChatID); registered {
		chart = &house.Chart
	}

	showListMenu(chart, u)
	return nil
}

func openHandler(b *robot.Bot, u *message.Update) message.Any {
	var itemID, item = extractItem(u.CallbackQuery)
	if item == nil {
		return warn(u.CallbackQuery, "🫤 This item has been removed")
	}

	showItemMenu(item.Caption(), "🔙 Back to list", itemID, *u)
	return nil
}

func dropHandler(b *robot.Bot, u *message.Update) message.Any {
	var callback = u.CallbackQuery
	if house, registered := SelectHouse(b.ChatID); registered {
		house.Chart.DropAll()
		callback.Delete()
		return warn(callback, "🗑 Shopping list emptied successfully!")
	}
	return warn(callback, "🫤 Shopping list is already empty")
}

func addHandler(b *robot.Bot, u *message.Update) message.Any {
	if house, registered := SelectHouse(b.ChatID); registered {
		itemID := extractItemID(u.CallbackQuery)
		if item := house.Chart.GetItem(itemID); item != nil {
			n := house.Chart.AddQuantity(itemID, 1)
			return warn(u.CallbackQuery, fmt.Sprintln("➕1", item.name, "added successfully,", n, "left"))
		}
	}
	return warn(u.CallbackQuery, "🫤 This item has been removed")
}

func subHandler(b *robot.Bot, u *message.Update) message.Any {
	var callback = u.CallbackQuery
	if house, registered := SelectHouse(b.ChatID); registered {
		itemID := extractItemID(callback)
		if item := house.Chart.GetItem(itemID); item != nil {
			itemName := item.name
			if n := house.Chart.RemoveQuantity(itemID, 1); n > 0 {
				return warn(callback, fmt.Sprintln("➖1", itemName, "removed successfully,", n, "left"))
			}
			callback.Delete()
			return warn(callback, fmt.Sprintln("🗑", itemName, "deleted from list successfully!"))
		}
	}
	return warn(u.CallbackQuery, "🫤 This item has been removed")
}

func delHandler(b *robot.Bot, u *message.Update) message.Any {
	var callback = u.CallbackQuery
	if house, registered := SelectHouse(b.ChatID); registered {
		itemID := extractItemID(callback)
		chart := house.Chart
		if item := chart.GetItem(itemID); item != nil && chart.DropItem(itemID) {
			showListMenu(&chart, u)
			return warn(callback, fmt.Sprintln("🗑", item.name, "deleted from list successfully!"))
		}
	}
	return warn(u.CallbackQuery, "🫤 This item has been already removed")
}

func messageHandler(b *robot.Bot, u *message.Update) message.Any {
	var (
		house, _     = SelectHouse(b.ChatID)
		itemID, item = house.Chart.Save(u.Message.Text, 1, false)
	)
	if item == nil {
		return warn(u.CallbackQuery, "🫤 Something went wrong")
	}

	showItemMenu(item.Caption(), "📄 Show list", itemID, *u)
	return nil
}
