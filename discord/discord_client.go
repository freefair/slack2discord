package discord

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"slack2discord/helper"
	"slack2discord/model"
	"strings"
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

func (c Client) mapAttachments(attachments []*discordgo.MessageAttachment) []model.Attachment {
	var result []model.Attachment
	for i := range attachments {
		attachment := attachments[i]
		download, err := helper.HTTPDownload(attachment.URL)
		if err != nil {
			panic(err)
		}
		result = append(result, model.Attachment{
			Data:        bytes.NewReader(download),
			ContentType: "",
			Name:        attachment.Filename,
		})
	}
	return result
}

func (c Client) cleanupMessage(text string, message *discordgo.Message) string {
	regex := regexp.MustCompile(`<@!([a-zA-Z0-9]*)>`)
	matches := regex.FindAllStringSubmatch(text, -1)
	for i := range matches {
		match := matches[i]
		for mi := range message.Mentions {
			mention := message.Mentions[mi]
			if mention.ID == match[1] {
				text = strings.Replace(text, match[0], "@" + mention.Username, -1)
			}
		}
	}
	return text
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
			Id:          message.ID,
			Text:        c.cleanupMessage(message.Content, message),
			User:        message.Author.Username,
			Time:        formatTimestamp(string(message.Timestamp)),
			Attachments: c.mapAttachments(message.Attachments),
		}}, result...)
	}

	return result
}

func (c Client) SendMessage(message model.Message) {
	if len(message.Attachments) > 0 {
		var files []*discordgo.File
		for i := range message.Attachments {
			attachment := message.Attachments[i]
			files = append(files, &discordgo.File{
				Name:        strings.ToLower(attachment.Name),
				ContentType: attachment.ContentType,
				Reader:      attachment.Data,
			})
		}
		_, err := c.DiscordClient.ChannelMessageSendComplex(c.ChannelId, &discordgo.MessageSend{
			Content: message.User + ": " + message.Text,
			Files: files,
		})
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	} else {
		_, err := c.DiscordClient.ChannelMessageSend(c.ChannelId, message.User+": "+message.Text)
		if err != nil {
			panic(err)
		}
	}
}
