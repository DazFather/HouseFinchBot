package main

import (
	"github.com/DazFather/parrbot/message"
	"github.com/DazFather/parrbot/robot"
)

func main() {
	// Start bot using following commands...
	robot.Start(
		// ... launched by inline buttons
		btnReply("/open", openHandler),
		btnReply("/open", openHandler),
		btnReply("/drop", dropHandler),
		btnReply("/join", joinHandler),
		btnReply("/add", addHandler),
		btnReply("/sub", subHandler),
		btnReply("/del", delHandler),
		closeCommand, // delete message and show toast alert
		// ... that can be launched also directly from the user
		userMenu("/start", startHandler, "▶️ Start the bot"),
		userMenu("/home", homeHandler, "🏠 House info"),
		userMenu("/list", listHandler, "🛒 Your shopping list"),
		userMenu("/share", shareHandler, "👥 Share with someone"),
		userMenu("/id", idHandler, ""),
		userMenu("/info", infoHandler, "ℹ️ Bot infos"),
		// ... reply without any explicit /trigger and only by user
		robot.Command{CallFunc: messageHandler, ReplyAt: message.MESSAGE},
	)
}

func btnReply(trigger string, handler robot.CommandFunc) robot.Command {
	return robot.Command{Trigger: trigger, CallFunc: handler, ReplyAt: message.CALLBACK_QUERY}
}

func userMenu(trigger string, handler robot.CommandFunc, description string) robot.Command {
	return robot.Command{
		Trigger:     trigger,
		CallFunc:    handler,
		ReplyAt:     message.MESSAGE + message.CALLBACK_QUERY,
		Description: description,
	}
}
