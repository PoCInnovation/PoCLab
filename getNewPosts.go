package main

import (
	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

//TODO: replace with the diff function
func getNewPost() *discordgo.MessageEmbed {
	return embed.NewGenericEmbed("New post on board: "+Board, "voila le message")
}
