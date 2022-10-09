package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shomali11/slacker"
)

// print event log info from channel of commandevents
func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events: ")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	// set slack tokens in environment
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-4208651927009-4196022262594-Mr2SnwNY3mPdvy2EG2TSxH6h")
	os.Setenv("SLACK_APP_TOKEN", "xapp-1-A0464K7D56V-4195932217059-1ebf4b97afd8eb47d7cd0486b5daf18b01c10fbecac8d55ebec42716a3f7317a")

	// create bot
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// concurrently print stream of command events
	go printCommandEvents(bot.CommandEvents())

	// define bot command for determining user age
	bot.Command("My year of birth is <year>", &slacker.CommandDefinition{
		Description: "yob calculator",
		Examples:    []string{"My year of birth is 2022"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {

			// get year from request, convert to int
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				println("error when parsing year")
			}

			// compute year
			age := 2022 - yob

			// format response as string, pass to response writer
			r := fmt.Sprintf("Age is %d", age)
			response.Reply(r)
		},
	})

	// return copy of background context, defer cancel call when program exits
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listen for events in context
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
