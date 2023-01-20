package main

import (
	"fmt"

	"github.com/DazFather/parrbot/message"
	"github.com/DazFather/parrbot/robot"
	"github.com/DazFather/parrbot/tgui"
)

var closeCommand robot.Command = tgui.Closer("/close", false)

var infoHandler = tgui.Sender(message.Text{fmt.Sprintln(
	"🐦", bold("HouseFinchBot"), "is a free and", repo("open source", "HouseFinchBot"),
	"Telegram's bot that allows to better manage your house.",
	"\nIt's still in work in progress and is being actively developed with ❤️ by @DazFather",
	"\n💪Powerade by", repo("Parr(B)ot", "parrbot"), "framework and written in pure ", link("Go", "https://go.dev/"),
	emph("Feel free to contact me on Telegram or contribute to the projects"),
), defaultOpt()})

var startHandler = tgui.Sender(message.Text{fmt.Sprintln(
	"👋 Welcome, I'm", bold("House Finch Bot"), "🐦",
	"\nYour personal robo-passerine that will help you take care of the house",
	"\nFor now I only know how to keep track of your shopping list but I'm learning 😅",
	"\nJust type the name of the item and it will be auto"+emph("magically✨")+" added to the /list",
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
	return warn(callback, "🫤 This item has been removed")
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
	return warn(callback, "🫤 This item has been already removed")
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

func shareHandler(b *robot.Bot, u *message.Update) message.Any {
	var (
		callback  = u.CallbackQuery // Probably nil
		msg       = u.FromMessage()
		house, _  = SelectHouse(b.ChatID)
		strUserID = trimCommand(msg.Text)
		name      = extractName(msg.From)
	)

	if callback != nil || strUserID == "" {
		show(u, fmt.Sprintln(
			"📨", bold("Invite new roomers"),
			"\nTo invite someone in your house use the command like this:\n",
			mono("/share TELEGRAMID"), "replacing", emph("\"TELEGRAMID\""), "with the Telegram's unique chat ID of the", bold("PERSON YOU WANT"), "to invite",
			"\n", emph("es:"), mono("/share ", b.ChatID),
			"\nThe other person must have started the bot already. If instead you want to know your ID use the command /id",
		), inlineCallerRow("🔙 Home menu", "/home"))
		return nil
	}

	// Foword request to owner if user is not
	if !house.IsOwner(b.ChatID) {
		invite := send(house.ownerID,
			"🔔 "+name+" your house member would like to invite user: "+user(strUserID, strUserID),
			inlineCallerRow("📨 Send invitation", "/share", strUserID),
		)
		if invite != nil {
			return warn(callback, "📨 Request sent to the house owner")
		}
	}

	// Send invitation
	if userID, err := toUserID(strUserID); err == nil {
		invite := send(userID,
			fmt.Sprintln(
				"💌 You have been invited to join the house of", name,
				emph("\n⚠️ If you join you loose all the datas (shopping list) of your current house"),
			),
			[]tgui.InlineButton{
				tgui.InlineCaller("✅ Accept", "/join", toString(b.ChatID)),
				tgui.InlineCaller("❌ Reject", "/close", "❌ Invitation rejected"),
			},
		)
		if invite != nil {
			return warn(callback, "📨 Request sent to "+extractChatName(invite.Chat))
		}
	}

	return warn(callback, "🫤 Something went wrong")
}

func joinHandler(b *robot.Bot, u *message.Update) message.Any {
	var (
		welcome  = `🐦 - "` + emph("Welcome to the house!") + `"`
		callback = u.CallbackQuery
		name     = extractName(callback.From)
	)

	if userID, err := toUserID(trimCommand(callback.Data)); err == nil {
		if house, registered := SelectHouse(userID); registered && house.Join(b.ChatID, name) {
			go send(userID, "🔔 "+name+" accepted the invitation!\n"+welcome)

			show(u, fmt.Sprintln("✅", bold("Invitation accepted!"), "\n"+welcome,
				"\nNow you are an house member too, use /list to manage the shopping list",
			))
			return nil
		}
	}

	return warn(callback, "🫤 Something went wrong")
}

func kickHandler(b *robot.Bot, u *message.Update) message.Any {
	var (
		callback = u.CallbackQuery
		name     = extractName(callback.From)
	)

	if userID, err := toUserID(trimCommand(callback.Data)); err == nil {
		if house, registered := SelectHouse(b.ChatID); registered && house.IsOwner(b.ChatID) && house.Kick(userID) {
			go send(userID, "🔔 "+name+" removed you from the house")
			return warn(callback, "User kicked successfully!")
		}
	}

	return warn(callback, "🫤 Something went wrong")
}

func idHandler(b *robot.Bot, u *message.Update) message.Any {
	return message.Text{"🆔 Your Telegram ID: " + mono(b.ChatID), defaultOpt()}
}

func homeHandler(b *robot.Bot, u *message.Update) message.Any {
	var (
		text string = bold("🏠 Home")
		kbd         = [][]tgui.InlineButton{
			inlineCallerRow("🛒 Shopping list", "/list"),
			inlineCallerRow("📨 Invite someone", "/share"),
			inlineCallerRow("✖️ Close menu", "/close"),
		}
	)

	if house, registered := SelectHouse(b.ChatID); registered {
		if size := house.Members(); size > 0 {
			text += " (👥" + mono(size+1) + ")"
			if house.IsOwner(b.ChatID) {
				kbd = append(kbd, inlineCallerRow("👥 Manage roomers", "/roomers"))
			}
		}
	}

	show(u,
		text+"\nUse the button below to help you navigate the bot"+"\n🐦 - \""+emph("Home sweet home")+"\"",
		kbd...,
	)
	return nil
}

func roomerHandler(b *robot.Bot, u *message.Update) message.Any {
	var house, registered = SelectHouse(b.ChatID)

	if !registered || house.Members() == 0 {
		return warn(u.CallbackQuery, "😐 You are the only roomer")
	}

	if !house.IsOwner(b.ChatID) {
		return warn(u.CallbackQuery, "🚫 You are NOT the owner")
	}

	size := house.Members()
	kbd := make([][]tgui.InlineButton, size)
	house.EachMembers(func(userID int64, name string) {
		size--
		kbd[size] = inlineCallerRow("🚷 "+name, "/kick", toString(userID))
	})

	show(u, bold("👥 Manage your roomers")+"\nTap on the ID of the roomer to kick him out of the house", kbd...)
	return nil
}
