package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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

//replace with the diff function
func getNewPost() string {
	return "MESSAGE2"
}

func pinner(s *discordgo.Session, m *discordgo.MessageCreate) {
	//TODO: channel use to act async
	time.Sleep(10 * time.Second)

	_, err := s.ChannelMessageSend(ChannelID, getNewPost())
	if err != nil {
		return
	}
}
