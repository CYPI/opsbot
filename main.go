package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
)

// https://api.slack.com/slack-apps
// https://api.slack.com/internal-integrations
type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"3000"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID string `envconfig:"CHANNEL_ID" required:"true"`
}

func main() {
	var env envConfig
	err := envconfig.Process("cortex", &env)
	if err != nil {
		log.Fatal(err.Error())
		return 1
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)
	slackListener := &SlackListener{
		client:    client,
		botID:     env.BotID,
		channelID: env.ChannelID,
	}
	go slackListener.ListenAndResponse()

	// Register handler to receive interactive message
	// responses from slack (kicked by user action)
	http.Handle("/interaction", interactionHandler{
		verificationToken: env.VerificationToken,
	})

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0

}
