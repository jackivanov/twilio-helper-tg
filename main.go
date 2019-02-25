package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// twilio

type StudioFlows struct {
	Meta  FlowsMeta
	Flows []StudioFlow
}

type FlowsMeta struct {
	Page              int
	Page_size         int
	First_page_url    string
	Previous_page_url string
	Url               string
	Next_page_url     string
	Key               string
}

type StudioFlow struct {
	Status        string
	Date_updated  string
	Friendly_name string
	Account_sid   string
	Url           string
	Version       int
	Sid           string
	Date_created  string
	Links         StudioFlowLinks
}

type StudioFlowLinks struct {
	Engagements string
	Executions  string
}

// getToken returns the value of given environment variable
func getToken(v string) string {
	if value, resp := os.LookupEnv(v); resp {
		return value
	}
	log.Panic("Telegram bot token should be specified in ", v)
	return ""
}

// ArgsToSlice returns a lice of given strings separated by spaces
func ArgsToSlice(args string) []string {
	return strings.Fields(args)
}

// StudioCallTo calls to the Twilio API and configured given Studio flow forwarding calls to given number
func StudioCallTo(flow, to string) bool {
	return true
}

// getFlows gets all existing Studio flows
func getFlows() []StudioFlow {
	Flows := new(StudioFlows)
	accountSid := getToken("TWILIO_ACCOUNT_SID")
	accountToken := getToken("TWILIO_AUTH_TOKEN")

	TwilioGetJson("https://studio.twilio.com/v1/Flows?PageSize=1000&Page=0", accountSid, accountToken, Flows)

	return Flows.Flows
}

// getFlows gets all existing Studio flows
func getFlowInfo() []StudioFlow {
	Flows := new(StudioFlows)
	accountSid := getToken("TWILIO_ACCOUNT_SID")
	accountToken := getToken("TWILIO_AUTH_TOKEN")

	TwilioGetJson("https://studio.twilio.com/v1/Flows?PageSize=1000&Page=0", accountSid, accountToken, Flows)

	return Flows.Flows
}

// Twilio http
func TwilioGetJson(url, accountSid, accountToken string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(accountSid, accountToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
		log.Printf("Can't connect to the Twilio server")
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

//
// func main() {
// 	Flows := new(StudioFlows)
// 	accountSid := getToken("TWILIO_ACCOUNT_SID")
// 	accountToken := getToken("TWILIO_AUTH_TOKEN")
//
// 	TwilioGetJson("https://studio.twilio.com/v1/Flows?PageSize=1000&Page=0", accountSid, accountToken, Flows)
//
// 	fmt.Println(Flows.Flows)
// }

func main() {
	BotAPIToken := getToken("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(BotAPIToken)
	if err != nil {
		log.Printf(". Bot not found or can not be authenticated")
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		who := update.Message.From.ID
		if who == 139541756 || who == 150631890 {
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "help", "start":
					msg.Text = "/call (flow sid) (number) - redirect calls in a given flow to a given numner. eg: /call FWf.. +79..\n/flows - get existing flows (name and sid)\n/status - check the status of the bot\n/help - show help message"
				case "call":
					args := ArgsToSlice(update.Message.CommandArguments())
					if len(args) == 2 {
						flow := args[0]
						to := args[1]

						if resp := StudioCallTo(flow, to); resp {
							m := fmt.Sprint("Redirectings call to ", to, " in flow ", flow)
							msg.Text = m
						}
					} else {
						msg.Text = "Error. Exact 2 arguments should be specified"
					}
				case "flows":
					var m string
					for _, v := range getFlows() {
						m = fmt.Sprint(v.Friendly_name, " ", v.Sid, "\n")
					}
					msg.Text = m
				case "status":
					msg.Text = "I'm ok."
				default:
					msg.Text = "I don't know that command"
				}
				bot.Send(msg)
			}
		}
	}
}
