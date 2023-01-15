package main

import (
	"github.com/DazFather/parrbot/message"
	"github.com/DazFather/parrbot/robot"
)

func main() {
	const (
		MSG     message.UpdateType = message.MESSAGE
		CLB                        = message.CALLBACK_QUERY
		MSG_CLB                    = MSG + CLB
	)

	// Start bot
	robot.Start([]robot.Command{
		{Trigger: "/start", CallFunc: startHandler, ReplyAt: MSG, Description: "Strat the bot"},
		{Trigger: "/list", CallFunc: listHandler, ReplyAt: MSG_CLB, Description: "Your shopping list"},
		{Trigger: "/open", CallFunc: openHandler, ReplyAt: CLB},
		{Trigger: "/drop", CallFunc: dropHandler, ReplyAt: CLB},
		{Trigger: "/add", CallFunc: addHandler, ReplyAt: CLB},
		{Trigger: "/sub", CallFunc: subHandler, ReplyAt: CLB},
		{Trigger: "/del", CallFunc: delHandler, ReplyAt: CLB},
		{CallFunc: messageHandler, ReplyAt: MSG},
	}...)
}
