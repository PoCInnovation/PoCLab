package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	setup()

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
			time.Sleep(time.Duration(Seconde))

			for _, board := range Boards {
				newPosts, warn := getNewPosts(board)
				if warn != nil {
					fmt.Println("error in getNewPosts", err)
					return
				}
				for _, p := range newPosts {
					_, err := dg.ChannelMessageSendEmbed(ChannelID, p.MessageEmbed)
					if err != nil {
						return
					}
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
