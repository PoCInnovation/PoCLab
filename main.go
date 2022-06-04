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

type arrayFlags []string

var (
	Token     string
	ChannelID string
	Seconde   int
	Boards    arrayFlags
)

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ChannelID, "c", "", "Channel ID")
	flag.IntVar(&Seconde, "s", 5, "second  between refresh")
	flag.Var(&Boards, "b", "board to notify")
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

			for _, board := range Boards {
				newPosts := getNewPosts(board)
				for _, v := range newPosts {
					_, err := dg.ChannelMessageSendEmbed(ChannelID, v.MessageEmbed)
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
