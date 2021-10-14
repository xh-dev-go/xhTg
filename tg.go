package main

import (
	"errors"
	"flag"
	"github.com/xh-dev/xhTg/getMe"
	"github.com/xh-dev/xhTg/notification"
	"os"
)

func main() {
	const cmd_token = "token"
	const cmd_notification = "notify"
	const cmd_get_me = "get-me"
	var isNotification, isGetMe bool
	var token string

	flag.StringVar(&token, cmd_token, "", "token")
	flag.BoolVar(&isNotification, cmd_notification, false, "notification to someone in telegram")
	flag.BoolVar(&isGetMe, cmd_get_me, false, "get me function")
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	flag.Parse()

	if token == "" {
		panic(errors.New("Token is not allow empty!!"))
	}

	args := flag.Args()
	if len(args) == 0 {
		panic(errors.New("no operation found"))
	} else if args[0]==cmd_notification {
		notification.Notify(token)
	} else if args[0] == cmd_get_me {
		getMe.GetMe(token)
	} else {
		flag.Usage()
	}
	os.Exit(0)
}
