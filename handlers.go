package main

import (
	"fmt"

	"github.com/DazFather/parrbot/message"
	"github.com/DazFather/parrbot/robot"
	"github.com/DazFather/parrbot/tgui"
)

var startHandler = tgui.Sender(message.Text{fmt.Sprint(
	"👋 Welcome, </b>I'm House Finch Bot</b>\n",
	"Your personal robo-passerine that will help you take care of the house\n",
	"\nFor now I only know how to keep track of your shopping list but I'm learning 😅\n",
	"Just type the name of the item and it will be auto<i>magically</i>✨ added to the /list",
), defaultOpt()})

func listHandler(b *robot.Bot, u *message.Update) message.Any {
	btns := GenItemList(func(itemID string, item ShoppingItem) tgui.InlineButton {
		return tgui.InlineCaller(item.Caption(), "/open", itemID)
	})

	if len(btns) == 0 {
		btns = tgui.Wrap(tgui.InlineCaller("🔄 Refresh", "/list"))
		show(u, "🕸 Your list is empty at the moment\nText me anything and I'll add it", btns)
		return nil
	}

	var kbd = make([][]tgui.InlineButton, len(btns)+1)
	for i := range btns {
		kbd[i] = tgui.Wrap(btns[i])
	}
	kbd[len(btns)] = []tgui.InlineButton{
		tgui.InlineCaller("🗑 Delete all", "/drop"),
		tgui.InlineCaller("🔄 Refresh list", "/list"),
	}

	show(u, "🛒 Your current shopping list:", kbd...)
	return nil
}

func openHandler(b *robot.Bot, u *message.Update) message.Any {
	var itemID, item = extractItem("/open", u.CallbackQuery)
	if item == nil {
		return warn(u.CallbackQuery, "Invalid item")
	}

	showItemMenu(item.Caption(), "Collapse", itemID, *u)
	return nil
}

func dropHandler(b *robot.Bot, u *message.Update) message.Any {
	DropList()
	return warn(u.CallbackQuery, "Shopping list has been deleted")
}

func addHandler(b *robot.Bot, u *message.Update) message.Any {
	var itemID, item = extractItem("/add", u.CallbackQuery)
	if item == nil {
		return warn(u.CallbackQuery, "Invalid item")
	}

	return warn(u.CallbackQuery, fmt.Sprintln("x1", item.name, "added successfully,", AddQuantity(itemID, 1), "left"))
}

func subHandler(b *robot.Bot, u *message.Update) message.Any {
	var itemID, item = extractItem("/sub", u.CallbackQuery)
	if item == nil {
		return warn(u.CallbackQuery, "Invalid item")
	}

	return warn(u.CallbackQuery, fmt.Sprintln("x1", item.name, "removed successfully,", RemoveQuantity(itemID, 1), "left"))
}

func delHandler(b *robot.Bot, u *message.Update) message.Any {
	var itemID, item = extractItem("/del", u.CallbackQuery)
	if item == nil {
		return warn(u.CallbackQuery, "Invalid item")
	}

	DropItem(itemID)
	return warn(u.CallbackQuery, "Item deleted successfully")
}

func messageHandler(b *robot.Bot, u *message.Update) message.Any {
	var itemID, item = Save(u.Message.Text, 1, false)
	if item == nil {
		return warn(u.CallbackQuery, "Something went wrong :/")
	}

	showItemMenu(item.name, "📄 Show list", itemID, *u)
	return nil
}
