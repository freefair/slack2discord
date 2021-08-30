package slack

import (
	"bytes"
	"fmt"
	"github.com/slack-go/slack"
	"regexp"
	"slack2discord/helper"
	"slack2discord/model"
	"strconv"
	"strings"
	"time"
)

var fileMap = map[string]string {
	"jpg": "image/jpeg",
	"jpeg": "image/jpeg",
	"png": "image/png",
}

type Client struct {
	Token       string
	ChannelId   string
	SlackClient *slack.Client
	users       map[string]string
	ownUserId   string
}

func NewClient(token string, channelId string) *Client {
	c := new(Client)
	c.Token = token
	c.ChannelId = channelId
	c.SlackClient = slack.New(token)
	c.users = make(map[string]string)
	identity, err := c.SlackClient.AuthTest()
	if err != nil {
		panic(err)
	}
	c.ownUserId = identity.UserID
	return c
}

func (c Client) Ping() error {
	_, err := c.SlackClient.AuthTest()
	return err
}

func (c Client) GetUserName(userId string) string {
	if _, ok := c.users[userId]; !ok {
		user, err := c.SlackClient.GetUserInfo(userId)
		if err != nil {
			return "None"
		}
		c.users[userId] = user.Name
	}
	return c.users[userId]
}

func (c Client) cleanupMessage(text string) string {
	regex := regexp.MustCompile(`<@([a-zA-Z0-9]*)>`)
	matches := regex.FindAllStringSubmatch(text, -1)
	for i := range matches {
		match := matches[i]
		text = strings.Replace(text, match[0], "@" + c.GetUserName(match[1]), -1)
	}
	return text
}

func formatTimestamp(val string) time.Time {
	i, err := strconv.ParseFloat(val, 10)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(int64(i), 0)
	return tm
}

func (c Client) mapAttachments(attachments []slack.File) []model.Attachment {
	result := []model.Attachment{}
	for i := range attachments {
		attachment := attachments[i]
		url := attachment.URLDownload
		if len(url) <= 0 {
			url = attachment.URLPrivateDownload
		}
		download, err := helper.HTTPDownload(url)
		if err != nil {
			panic(err)
		}
		filetype := "application/octet-stream"
		if t, ok := fileMap[attachment.Filetype]; ok {
			filetype = t
		}
		result = append(result, model.Attachment{
			Data:        bytes.NewReader(download),
			ContentType: filetype,
			Name:        attachment.Name,
		})
	}
	return result
}

func (c Client) GetNewMessages(last string) []model.Message {
	var result []model.Message

	messageParams := new(slack.GetConversationHistoryParameters)
	messageParams.Oldest = last
	messageParams.Limit = 100
	messageParams.ChannelID = c.ChannelId
	messageParams.Inclusive = false
	history, err := c.SlackClient.GetConversationHistory(messageParams)
	if err != nil {
		fmt.Printf("%s\n", err)
		return result
	}

	for i := range history.Messages {
		message := history.Messages[i]
		if message.User == c.ownUserId || len(message.BotID) > 0 {
			continue
		}
		result = append([]model.Message{{
			Id:          message.Timestamp,
			Text:        c.cleanupMessage(message.Text),
			User:        c.GetUserName(message.User),
			Time:        formatTimestamp(message.Timestamp),
			Attachments: c.mapAttachments(message.Files),
		}}, result...)
	}

	return result
}

func (c Client) SendMessage(message model.Message) {
	if len(message.Attachments) > 0 {
		comment := message.Text
		for i := range message.Attachments {
			attachment := message.Attachments[i]
			_, err := c.SlackClient.UploadFile(slack.FileUploadParameters{
				Reader:   attachment.Data,
				Filename: "file",
				Channels: []string{c.ChannelId},
				InitialComment: comment,
				Filetype: attachment.ContentType,
				Title: message.User,
			})

			comment = ""

			if err != nil {
				panic(err)
			}
		}
	} else {
		_, _, err := c.SlackClient.PostMessage(c.ChannelId, slack.MsgOptionText(message.Text, true), slack.MsgOptionAsUser(false), slack.MsgOptionUsername(message.User))
		if err != nil {
			panic(err)
		}
	}
}
