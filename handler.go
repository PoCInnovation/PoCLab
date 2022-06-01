package main

import (
	"fmt"
	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"

	"time"
)

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			return
		}
	}

	if m.Content == "pong" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			return
		}
	}
}

//TODO: replace with the diff function
func getNewPost(c chan *discordgo.MessageEmbed) {
	time.Sleep(time.Duration(Second) * time.Second)
	c <- embed.NewGenericEmbed("New post on board: "+Board, "voila le message")
	return
}

func clock(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := make(chan *discordgo.MessageEmbed, 1)
	go getNewPost(c)

	msg := <-c
	_, err := s.ChannelMessageSendEmbed(ChannelID, msg)
	fmt.Println("check")
	if err != nil {
		return
	}
}
