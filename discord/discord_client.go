package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"slack2discord/model"
	"time"
)

type Client struct {
	Token         string
	ChannelId     string
	DiscordClient *discordgo.Session
	users         map[string]string
}

func NewClient(token string, channelId string) *Client {
	c := new(Client)
	c.Token = token
	c.ChannelId = channelId
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
	}
	c.DiscordClient = session
	c.users = make(map[string]string)
	return c
}

func formatTimestamp(val string) time.Time {
	parse, err := time.Parse(time.RFC3339, val)
	if err != nil {
		panic(err)
	}
	return parse
}

func (c Client) Ping() error {
	_, err := c.DiscordClient.GatewayBot()
	return err
}

func (c Client) GetNewMessages(last string) []model.Message {
	var result []model.Message

	messages, err := c.DiscordClient.ChannelMessages(c.ChannelId, 100, "", last, "")
	if err != nil {
		fmt.Println("error while getting messages", err)
	}

	for i := range messages {
		message := messages[i]
		if message.Author.Bot {
			continue
		}
		result = append([]model.Message{{
			Id:   message.ID,
			Text: message.Content,
			User: message.Author.Username,
			Time: formatTimestamp(string(message.Timestamp)),
		}}, result...)
	}

	return result
}

func (c Client) SendMessage(message model.Message) {
	_, err := c.DiscordClient.ChannelMessageSend(c.ChannelId, message.User+": "+message.Text)
	if err != nil {
		panic(err)
	}
}
