package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Token     string
	ChannelID string
	Second    int
	Board     string

	// Message temporary
	Message string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ChannelID, "c", "", "Channel ID")
	flag.IntVar(&Second, "s", 5, "second  between refresh")
	flag.StringVar(&Board, "b", "announcement", "board to notify") // TODO: modify to add multiple boards

	// temporary
	flag.StringVar(&Message, "m", "", "message to be print")
	//
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	go func() {
		for true {
			time.Sleep(time.Duration(Second) * time.Second)

			msg := getNewPost()
			_, err := dg.ChannelMessageSendEmbed(ChannelID, msg)
			if err != nil {
				fmt.Println("error sending message,", err)
				return
			}
		}
	}()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = dg.Close()
	if err != nil {
		return
	}
}
