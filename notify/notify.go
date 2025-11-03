package notify

import "log"

func Send(channel, msg string) {
	switch channel {
	case "slack":
		// call Slack webhook
	case "email":
		// send SMTP mail
	default:
		log.Println("ALERT:", msg)
	}
}
