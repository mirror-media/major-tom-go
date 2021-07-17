package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	gootkitconfig "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/slashcommand"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(
		&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})

	var cfg config.Config

	// parse flag
	botFlags := gootkitconfig.New("bot config")
	err := botFlags.LoadFlags([]string{"c:string"})
	if err != nil {
		panic(errors.Wrap(err, "loading flags for config file has error"))
	}

	c := botFlags.String("c")
	logrus.Infof("config file is %s", c)

	botConfig := gootkitconfig.New("bot config")
	botConfig.AddDriver(yaml.Driver)

	// load config file
	err = botConfig.LoadFiles(c)
	if err != nil {
		panic(errors.Wrap(err, "loading config file has error"))
	}

	err = botConfig.BindStruct("", &cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error binding config file to struct: %s", err))
	}

	var k8sRepoCFG config.KubernetesConfigsRepo

	k := botFlags.String("k")
	logrus.Infof("k8s repo config file is %s", k)

	k8sConfig := gootkitconfig.New("bot config")
	k8sConfig.AddDriver(yaml.Driver)

	// load config file
	err = k8sConfig.LoadFiles(k)
	if err != nil {
		panic(errors.Wrap(err, "loading k8s repo config file has error"))
	}

	err = k8sConfig.BindStruct("", &k8sRepoCFG)
	if err != nil {
		panic(fmt.Errorf("fatal error binding config file to struct: %s", err))
	}

	appToken := cfg.SlackAppToken

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

	// TODO
	// clusterConfigs := cfg.ClusterConfigs
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

					messages, err := slashcommand.Run(ctx, k8sRepoCFG, cmd.Command, cmd.Text, cmd.UserName)
					if messages == nil {
						messages = []string{}
					}
					if err != nil {
						messages = append([]string{err.Error()}, messages...)
					}

					message := fmt.Sprintf("<@%s> on ground control\n```%s```", cmd.UserID, strings.Join(messages, "\n"))

					api.PostMessage(cmd.ChannelID, slack.MsgOptionResponseURL(cmd.ResponseURL, "in_channel"), slack.MsgOptionText(message, false))

				default:
					fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
				}
			}
		}
	}()

	client.Run()

}
