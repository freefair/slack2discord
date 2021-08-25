package bidi_translation

import (
	"context"
	"slack2discord/discord"
)
import "slack2discord/slack"
import "github.com/go-redis/redis/v8"

type Slack2DiscordTranslation struct {
	SlackClient   *slack.Client
	DiscordClient *discord.Client
	redisClient   *redis.Client
	context		  context.Context
}

func NewTranslationClient(slackToken string, slackChannel string, discordToken string, discordChannel string, redis *redis.Client) *Slack2DiscordTranslation {
	translation := new(Slack2DiscordTranslation)

	translation.DiscordClient = discord.NewClient(discordToken, discordChannel)
	translation.SlackClient = slack.NewClient(slackToken, slackChannel)
	translation.redisClient = redis
	translation.context = context.Background()

	return translation
}

func saveMessagePointer(context context.Context, redisClient *redis.Client, mtype string, id string) {
	err := redisClient.Set(context, "message-id:"+mtype, id, 0).Err()
	if err != nil {
		panic(err)
	}
}

func getMessagePointer(context context.Context, client *redis.Client, mtype string) string {
	result, err := client.Get(context, "message-id:"+mtype).Result()
	if err == redis.Nil {
		return ""
	}
	if err != nil {
		panic(err)
	}
	return result
}

func (t Slack2DiscordTranslation) BiDiSendMessages() {
	slackMessages := t.SlackClient.GetNewMessages(getMessagePointer(t.context, t.redisClient, "slack"))
	for i := range slackMessages {
		message := slackMessages[i]
		t.DiscordClient.SendMessage(message)
		saveMessagePointer(t.context, t.redisClient, "slack", message.Id)
	}

	discordMessages := t.DiscordClient.GetNewMessages(getMessagePointer(t.context, t.redisClient, "discord"))
	for i := range discordMessages {
		message := discordMessages[i]
		t.SlackClient.SendMessage(message)
		saveMessagePointer(t.context, t.redisClient, "discord", message.Id)
	}
}