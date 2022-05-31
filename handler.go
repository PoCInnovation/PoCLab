package main

import (
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
func getNewPost(c chan string) {
	time.Sleep(5 * time.Second)
	c <- Message
}

func clock(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := make(chan string)
	go getNewPost(c)

	msg := <-c
	_, err := s.ChannelMessageSend(ChannelID, msg)
	if err != nil {
		return
	}
}
