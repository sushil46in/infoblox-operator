package controllers

import (
	"fmt"

	"github.com/slack-go/slack"
)

func sendmsg(hn string, ip string) {
	api := slack.New("xoxp-1745509759573-1748801275043-1792816798822-82a66ebdbf2b807f65fd3d6bf4e77b6f")
	attachment := slack.Attachment{
		Pretext: "New DNS Record",
		Text:    hn + " " + ip,
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "a",
					Value: "no",
				},
			},
		*/
	}

	channelID, timestamp, err := api.PostMessage(
		"demo",
		slack.MsgOptionText("Some text", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

}
