package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/internal/slashcommand"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./configs") // path to look for the config file in
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// FIXME what should I do with changed file?
		fmt.Println("Config file changed:", e.Name)
	})

	var cfg config.Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config file: %s \n", err))
	}

	appToken := cfg.SlackToken

	api := slack.New("",
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken))

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	// FIXME remove it and use client.Run(context.Context) instead
	ctx := context.Background()

	clusterConfigs := cfg.ClusterConfigs
	go func() {
		for evt := range client.Events {
			select {
			case <-ctx.Done():
				os.Exit(0)
			default:
				switch evt.Type {
				case socketmode.EventTypeConnecting:
					fmt.Println("Connecting to Slack with Socket Mode...")
				case socketmode.EventTypeConnectionError:
					fmt.Println("Connection failed. Retrying later...")
				case socketmode.EventTypeConnected:
					fmt.Println("Connected to Slack with Socket Mode.")
				case socketmode.EventTypeEventsAPI:
					eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
					if !ok {
						fmt.Printf("Ignored %+v\n", evt)

						continue
					}

					fmt.Printf("Event received: %+v\n", eventsAPIEvent)

					client.Ack(*evt.Request)

					switch eventsAPIEvent.Type {
					case slackevents.CallbackEvent:
						innerEvent := eventsAPIEvent.InnerEvent
						switch ev := innerEvent.Data.(type) {
						case *slackevents.AppMentionEvent:
							_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
							if err != nil {
								fmt.Printf("failed posting message: %v", err)
							}
						case *slackevents.MemberJoinedChannelEvent:
							fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
						}
					default:
						client.Debugf("unsupported Events API event received")
					}
				case socketmode.EventTypeInteractive:
					callback, ok := evt.Data.(slack.InteractionCallback)
					if !ok {
						fmt.Printf("Ignored %+v\n", evt)

						continue
					}

					fmt.Printf("Interaction received: %+v\n", callback)

					var payload interface{}

					switch callback.Type {
					case slack.InteractionTypeBlockActions:
						// See https://api.slack.com/apis/connections/socket-implement#button

						client.Debugf("button clicked!")
					case slack.InteractionTypeShortcut:
					case slack.InteractionTypeViewSubmission:
						// See https://api.slack.com/apis/connections/socket-implement#modal
					case slack.InteractionTypeDialogSubmission:
					default:

					}

					client.Ack(*evt.Request, payload)
				case socketmode.EventTypeSlashCommand:
					cmd, ok := evt.Data.(slack.SlashCommand)
					if !ok {
						fmt.Printf("Ignored %+v\n", evt)

						continue
					}

					client.Debugf("Slash command received: %+v", cmd)

					payload := map[string]interface{}{
						"response_type": "in_channel",
					}

					client.Ack(*evt.Request, payload)

					messages, err := slashcommand.Run(ctx, clusterConfigs, cmd.Command, cmd.Text)
					if messages == nil {
						messages = []string{}
					}
					if err != nil {
						messages = append([]string{err.Error()}, messages...)
					}

					message := strings.Join(messages, "\n")

					api.PostMessage(cmd.ChannelID, slack.MsgOptionResponseURL(cmd.ResponseURL, "in_channel"), slack.MsgOptionText(message, false), slack.MsgOptionText(err.Error(), false))

				default:
					fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
				}
			}
		}
	}()

	client.Run()

}
