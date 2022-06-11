package main

import (
	"fmt"
	"github.com/enescakir/emoji"
)

func embedPosts(posts []Post, board string) []Embed {
	embeds := make([]Embed, 0)
	for _, post := range posts {
		embeds = append(embeds, Embed{
			Title:       fmt.Sprintf("New post on: %s %v ", board, emoji.OpenMailboxWithRaisedFlag),
			Description: fmt.Sprintf("**%s**\n%s\n\nhttps://gno.land/r/boards:%s/%d", post.Title, post.Description, board, post.Id),
			Author: Author{
				Name:    formatAuthor(post.Author),
				IconUrl: "https://cdn.discordapp.com/attachments/981266192390049846/983052513932607488/unknown.png",
			},
			Color: 7212552,
		})
	}
	fmt.Printf("THERE IS %d NEW POSTS\n", len(embeds))
	return embeds
}

func embedReplies(replies []Reply, post string, postTitle string, board string) []Embed {
	var embeds []Embed
	for _, r := range replies {
		embed := Embed{
			Title:       fmt.Sprintf("New reply on board: %s", board),
			Description: fmt.Sprintf("post: ***%s***\nreply: %s\n\nhttps://gno.land/r/boards:%s", postTitle, r.Content, post),
			Author: Author{
				Name:    formatAuthor(r.Author),
				IconUrl: "https://cdn.discordapp.com/attachments/981266192390049846/983052513932607488/unknown.png",
			},
			Color: 7212552,
		}
		embeds = append(embeds, embed)
	}
	return embeds
}
