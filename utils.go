package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/DazFather/parrbot/message"
	"github.com/DazFather/parrbot/tgui"

	"github.com/NicoNex/echotron/v3"
)

func showItemMenu(text, listCaption, itemID string, u message.Update) {
	show(&u, text, []tgui.InlineButton{
		tgui.InlineCaller("➕ 1", "/add", itemID),
		tgui.InlineCaller("➖ 1", "/sub", itemID),
		tgui.InlineCaller("🗑", "/del", itemID),
	}, tgui.Wrap(tgui.InlineCaller(listCaption, "/list")))
}

func showEmptyListMenu(u *message.Update) {
	show(u,
		"🕸 Your list is empty at the moment\nText me anything and I'll add it",
		tgui.Wrap(tgui.InlineCaller("🔄 Refresh", "/list")),
	)
}

func showListMenu(chart *ShoppingList, u *message.Update) {
	if chart == nil {
		showEmptyListMenu(u)
		return
	}
	size := chart.Items()
	if size == 0 {
		showEmptyListMenu(u)
		return
	}

	kbd := make([][]tgui.InlineButton, size+1)
	chart.ForEach(func(itemID string, item *ShoppingItem) {
		size--
		kbd[size] = tgui.Wrap(tgui.InlineCaller(item.Caption(), "/open", itemID))
	})
	kbd[len(kbd)-1] = []tgui.InlineButton{
		tgui.InlineCaller("🗑 Delete all", "/drop"),
		tgui.InlineCaller("🔄 Refresh list", "/list"),
	}
	show(u, "🛒 Your current shopping list:", kbd...)
	return
}

func extractKbd(msg message.UpdateMessage) [][]tgui.InlineButton {
	if markup := msg.InlineKeyboard; markup != nil {
		return markup.InlineKeyboard
	}
	return nil
}

func extractItemID(callback *message.CallbackQuery) string {
	if callback == nil {
		return ""
	}
	command := callback.Data
	return strings.TrimSpace(command[strings.IndexRune(command, ' '):])
}

func extractItem(callback *message.CallbackQuery) (itemID string, item *ShoppingItem) {
	if callback == nil {
		return
	}

	itemID = extractItemID(callback)
	if house, registered := SelectHouse(callback.From.ID); registered {
		item = house.Chart.GetItem(itemID)
	}
	return
}

func warn(callback *message.CallbackQuery, text string) message.Any {
	if callback == nil {
		return message.Text{text, nil}
	}

	// TODO: fix weird error "<nil>" on Parrbot
	callback.AnswerToast(text, 0)
	return nil
}

func sameKbd(kbd1, kbd2 [][]tgui.InlineButton) bool {
	if len(kbd1) != len(kbd2) {
		return false
	}

	for i := range kbd1 {
		if len(kbd1[i]) != len(kbd2[i]) {
			return false
		}
		for j := range kbd1[i] {
			if kbd1[i][j] != kbd2[i][j] {
				return false
			}
		}
	}

	return true
}

func show(u *message.Update, text string, buttons ...[]tgui.InlineButton) (sent *message.UpdateMessage) {
	var err error

	// Check if there is nothing to edit
	if u.CallbackQuery != nil {
		msg := u.CallbackQuery.Message
		if msg != nil && text == msg.Text && sameKbd(extractKbd(*msg), buttons) {
			u.CallbackQuery.AnswerToast("🫤 Nothing changed", 0)
			return nil
		}
	}

	sent, err = tgui.ShowMessage(*u, text, defaultEditOpt(buttons...))
	if err != nil && fmt.Sprint(err) != "<nil>" {
		log.Println("Error (show):", err)
	}

	return sent
}

func defaultEditOpt(buttons ...[]tgui.InlineButton) (opt *tgui.EditOptions) {
	opt = tgui.ParseModeOpt(nil, "HTML")
	if len(buttons) > 0 {
		tgui.InlineKbdOpt(opt, buttons)
	}
	return
}

func defaultOpt(buttons ...[]tgui.InlineButton) (opt *echotron.MessageOptions) {
	opt = &echotron.MessageOptions{ParseMode: "HTML"}
	if len(buttons) > 0 {
		opt.ReplyMarkup = echotron.InlineKeyboardMarkup{InlineKeyboard: buttons}
	}
	return
}
