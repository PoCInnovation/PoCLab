package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Body struct {
	Embeds []Embed `json:"embeds"`
}

func main() {
	err := setup(Boards)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

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

	for true {
		time.Sleep(time.Duration(Seconde))

		for _, board := range Boards {
			newPosts, warn := getNewPosts(board)
			if warn != nil {
				fmt.Println("error in getNewPosts.")
				err = dg.Close()
				return
			}
			// post request on webhook
			for _, p := range newPosts {
				e := Body{Embeds: []Embed{p}}
				jsonData, _ := json.Marshal(e)
				_, err = http.Post(`https://discord.com/api/webhooks/982746889868935198/CVXU-yDUWej-1WbgIYETpIaehMU1GY37u8hDrYqBZJl4UItt313C-J_t-WQ9L8Ey0wG8`, "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					return
				}
			}
		}
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = dg.Close()
	if err != nil {
		return
	}
}
