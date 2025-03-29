package gateway

import (
	"discord-go/api"
	"fmt"
)

func MessageCreateHandler(event interface{}) {
	message, ok := event.(*api.Message)
	if !ok {
		fmt.Println("Error: Incorrect event type for MESSAGE_CREATE")
		return
	}

	fmt.Printf("New message from %s: %s\n", message.Author.Username, message.Content)
}

func ChannelCreateHandler(event interface{}) {
	channel, ok := event.(*api.Channel)
	if !ok {
		fmt.Println("Error: Incorrect event type for CHANNEL_CREATE")
		return
	}

	fmt.Printf("New channel created: %s (ID: %s)\n", channel.Name, channel.ID)
}

func ChannelUpdateHandler(event interface{}) {
	channel, ok := event.(*api.Channel)
	if !ok {
		fmt.Println("Error: Incorrect event type for CHANNEL_UPDATE")
		return
	}

	fmt.Printf("Channel updated: %s (ID: %s)\n", channel.Name, channel.ID)
}

func ChannelDeleteHandler(event interface{}) {
	channel, ok := event.(*api.Channel)
	if !ok {
		fmt.Println("Error: Incorrect event type for CHANNEL_DELETE")
		return
	}

	fmt.Printf("Channel deleted: %s (ID: %s)\n", channel.Name, channel.ID)
}

func PresenceUpdateHandler(event interface{}) {
	presence, ok := event.(*api.PresenceUpdate)
	if !ok {
		fmt.Println("Error: Incorrect event type for PRESENCE_UPDATE")
		return
	}

	fmt.Printf("User %s is now %s\n", presence.User.Username, presence.Status)
}
