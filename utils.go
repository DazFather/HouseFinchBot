package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/DazFather/parrbot/message"
	"github.com/DazFather/parrbot/tgui"

	"github.com/NicoNex/echotron/v3"
)

/*--- regex ---*/

var (
	trigger = regexp.MustCompile(`^\s*/\w+\s*`)
	plurals = regexp.MustCompile("e?s$")
)

func rgxTrim(rgx *regexp.Regexp, s string) string {
	return strings.TrimSpace(rgx.ReplaceAllLiteralString(s, ""))
}

/*--- menues & messages ---*/

func showItemMenu(text, listCaption, itemID string, u message.Update) {
	show(&u, text, []tgui.InlineButton{
		tgui.InlineCaller("➕ 1", "/add", itemID),
		tgui.InlineCaller("➖ 1", "/sub", itemID),
		tgui.InlineCaller("🗑", "/del", itemID),
	}, inlineCallerRow(listCaption, "/list"))
}

func showEmptyListMenu(u *message.Update) {
	show(u,
		"🕸 Your list is empty at the moment\nText me anything and I'll add it",
		inlineCallerRow("🔄 Refresh", "/list"),
		inlineCallerRow("🔙 Home menu", "/home"),
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

	kbd := make([][]tgui.InlineButton, size+2)
	chart.ForEach(func(itemID string, item *ShoppingItem) {
		size--
		kbd[size] = inlineCallerRow(item.Caption(), "/open", itemID)
	})
	size = len(kbd) - 1
	kbd[size-1] = []tgui.InlineButton{
		tgui.InlineCaller("🗑 Delete all", "/drop"),
		tgui.InlineCaller("🔄 Refresh list", "/list"),
	}
	kbd[size] = inlineCallerRow("🔙 Home menu", "/home")

	show(u, "🛒 Your current shopping list:", kbd...)
	return
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

func send(to int64, text string, buttons ...[]tgui.InlineButton) *message.UpdateMessage {
	var (
		msg       = message.Text{text, defaultOpt(buttons...)}
		sent, err = msg.Send(to)
	)

	if err != nil {
		log.Println("Unable to send message:", err)
	}
	return sent
}

func warn(callback *message.CallbackQuery, text string) message.Any {
	if callback == nil {
		return message.Text{text, nil}
	}

	// TODO: fix weird error "<nil>" on Parrbot
	callback.AnswerToast(text, 0)
	return nil
}

/*--- extractors ---*/

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

	return trimCommand(callback.Data)
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

// Extract @username or name + surname of user
func extractName(user *echotron.User) string {
	if user == nil {
		return ""
	}

	var name = user.Username
	if name != "" {
		return "@" + name
	}

	name = user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}
	return name
}

// Extract @username or name + surname of user
func extractChatName(chat *echotron.Chat) string {
	if chat == nil {
		return ""
	}

	var name = chat.Username
	if name != "" {
		return "@" + name
	}

	name = chat.FirstName
	if chat.LastName != "" {
		name += " " + chat.LastName
	}
	return name
}

func trimCommand(text string) string {
	return rgxTrim(trigger, text)
}

func toUserID(rawUserId string) (userID int64, err error) {
	var intID int

	intID, err = strconv.Atoi(rawUserId)
	if err == nil {
		userID = int64(intID)
	}

	return
}

func toString(userID int64) string {
	return fmt.Sprint(userID)
}

/*--- markup, keyboards and message options ---*/

const defParseMode echotron.ParseMode = "HTML"

func link(caption, url string) string {
	return `<a href="` + url + `">` + caption + `</a>`
}

func user[T int | int64 | string](caption string, userID T) string {
	return link(caption, "tg://user?id="+fmt.Sprint(userID))
}

func repo(caption, name string) string {
	return link(caption, "https://github.com/DazFather/"+name)
}

func bold(caption ...any) string {
	return "<b>" + fmt.Sprint(caption...) + "</b>"
}

func mono(caption ...any) string {
	return "<code>" + fmt.Sprint(caption...) + "</code>"
}

func emph(caption ...any) string {
	return "<i>" + fmt.Sprint(caption...) + "</i>"
}

func defaultEditOpt(buttons ...[]tgui.InlineButton) (opt *tgui.EditOptions) {
	opt = tgui.ParseModeOpt(nil, defParseMode)
	if len(buttons) > 0 {
		tgui.InlineKbdOpt(opt, buttons)
	}
	return
}

func defaultOpt(buttons ...[]tgui.InlineButton) (opt *echotron.MessageOptions) {
	opt = &echotron.MessageOptions{ParseMode: defParseMode}
	if len(buttons) > 0 {
		opt.ReplyMarkup = echotron.InlineKeyboardMarkup{InlineKeyboard: buttons}
	}
	return
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

func inlineCallerRow(caption, trigger string, payload ...string) []tgui.InlineButton {
	return tgui.Wrap(tgui.InlineCaller(caption, trigger, payload...))
}

/*--- emoji ---*/

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
	"chip":       "🍟",
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

func findEmoji(s string) string {
	for _, piece := range strings.Split(strings.ToLower(strings.TrimSpace(s)), " ") {
		if emoji, found := EMOJI[rgxTrim(plurals, piece)]; found {
			return emoji
		}
	}
	return ""
}
