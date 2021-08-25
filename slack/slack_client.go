package slack

import (
	"fmt"
	"github.com/slack-go/slack"
	"slack2discord/model"
	"strconv"
	"time"
)

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

func formatTimestamp(val string) time.Time {
	i, err := strconv.ParseFloat(val, 10)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(int64(i), 0)
	return tm
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
		if message.User == c.ownUserId {
			continue
		}
		result = append([]model.Message{{
			Id:   message.Timestamp,
			Text: message.Text,
			User: c.GetUserName(message.User),
			Time: formatTimestamp(message.Timestamp),
		}}, result...)
	}

	return result
}

func (c Client) SendMessage(message model.Message) {
	_, _, err := c.SlackClient.PostMessage(c.ChannelId, slack.MsgOptionText(message.User+": "+message.Text, true), slack.MsgOptionAsUser(true))
	if err != nil {
		panic(err)
	}
}
