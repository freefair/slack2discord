package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/heptiolabs/healthcheck"
	"net/http"
	"runtime"
	"slack2discord/bidi_translation"
	"time"
)

func redisCheck(redisClient *redis.Client) healthcheck.Check {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		defer cancel()
		if redisClient == nil {
			return fmt.Errorf("redisClient is nil")
		}
		return redisClient.Ping(ctx).Err()
	}
}

func slackCheck(translation *bidi_translation.Slack2DiscordTranslation) healthcheck.Check {
	return func() error {
		_, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		defer cancel()
		if translation == nil {
			return fmt.Errorf("translation is nil")
		}
		return translation.SlackClient.Ping()
	}
}

func discordCheck(translation *bidi_translation.Slack2DiscordTranslation) healthcheck.Check {
	return func() error {
		_, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		defer cancel()
		if translation == nil {
			return fmt.Errorf("translation is nil")
		}
		return translation.DiscordClient.Ping()
	}
}

func main() {
	redisHost := flag.String("redis-host", "127.0.0.1:6379", "Redis hostname")
	redisPassword := flag.String("redis-pw", "", "Redis password")
	redisDatabase := flag.Int("redis-db", 10, "Redis database")

	slackToken := flag.String("slack-token", "", "Slack auth token")
	slackChannel := flag.String("slack-channel", "", "Slack channel id")

	discordToken := flag.String("discord-token", "", "Discord bot auth token")
	discordChannel := flag.String("discord-channel", "", "Discord channel id")

	flag.Parse()

	health := healthcheck.NewHandler()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     *redisHost,
		Password: *redisPassword, // no password set
		DB:       *redisDatabase,  // use default DB
	})
	translation := bidi_translation.NewTranslationClient(
		*slackToken,
		*slackChannel,
		*discordToken,
		*discordChannel,
		redisClient)

	health.AddReadinessCheck("redis", redisCheck(redisClient))
	health.AddReadinessCheck("slack", slackCheck(translation))
	health.AddReadinessCheck("discord", discordCheck(translation))

	health.AddLivenessCheck("redis", redisCheck(redisClient))
	health.AddLivenessCheck("slack", slackCheck(translation))
	health.AddLivenessCheck("discord", discordCheck(translation))

	go http.ListenAndServe("0.0.0.0:8080", health)
	go func() {
		for true {
			translation.BiDiSendMessages()
			time.Sleep(time.Second * 15)
		}
	}()

	runtime.Goexit()
}
