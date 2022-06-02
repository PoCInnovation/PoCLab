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
	Seconde   int
	Board     string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ChannelID, "c", "", "Channel ID")
	flag.IntVar(&Seconde, "s", 5, "second  between refresh")
	flag.StringVar(&Board, "b", "announcement", "board to notify") // TODO: modify to add multiple boards

	flag.Parse()
}

func main() {
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
			time.Sleep(time.Duration(Seconde) * time.Second)

			newPosts := getNewPosts()
			for _, v := range newPosts {
				_, err := dg.ChannelMessageSendEmbed(ChannelID, v.MessageEmbed)
				if err != nil {
					return
				}
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
