package main

import (
	"fmt"
	embed "github.com/Clinet/discordgo-embed"
	abci "github.com/gnolang/gno/pkgs/bft/abci/types"
	"github.com/gnolang/gno/pkgs/bft/rpc/client"
	"regexp"
	"strings"
)

func makeRequest(qpath string, data []byte) (res *abci.ResponseQuery, err error) {
	opts2 := client.ABCIQueryOptions{
		// Height: height, XXX
		// Prove: false, XXX
	}
	remote := "gno.land:36657"
	cli := client.NewHTTP(remote, "/websocket")
	qres, err := cli.ABCIQueryWithOptions(qpath, data, opts2)
	if err != nil {
		return nil, err
	}
	if qres.Response.Error != nil {
		fmt.Printf("Log: %s\n",
			qres.Response.Log)
		return nil, qres.Response.Error
	}
	return &qres.Response, nil
}

func getBoardsPosts() (string, error) {
	qpath := "vm/qrender"
	data := []byte(fmt.Sprintf("%s\n%s", "gno.land/r/boards", Board))
	res, err := makeRequest(qpath, data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(res.Data), nil
}

type Post struct {
	Title  string
	Author string
}

func parseNewPosts(BoardPosts string, index []string) []Post {
	//TODO: replace by real parsing
	post := make([]Post, 0)
	//for _, i := range index {
	//	if i > indexMax {
	//		post = append(post, Post{Title: "", Author: ""})
	//	}
	//}
	return post
}

func EmbedNewPosts(posts []Post) []*embed.Embed {
	embeds := make([]*embed.Embed, 0)
	for _, post := range posts {
		embeds = append(embeds, embed.NewEmbed().
			SetTitle(fmt.Sprintf("New post on: %s", Board)).
			SetDescription(post.Title).
			SetAuthor(post.Author).
			SetColor(0x00FF00))
	}
	return embeds
}

// TODO: replace with the diff function
func getNewPosts() []*embed.Embed {
	// this return the posts from the watched board
	BoardPosts, err := getBoardsPosts()
	if err != nil {
		return nil
	}
	// debug
	fmt.Println(BoardPosts)

	re := regexp.MustCompile("\\bpostid=[0-9]+")
	fr := regexp.MustCompile("(\\b[0-9]+)")
	fmt.Println("Parsing new posts")
	fmt.Println("------------------")
	var getId = re.FindAllString(BoardPosts, -1)
	fmt.Println(getId)
	fmt.Println(fr.FindAllString(strings.Join(getId, " "), -1))

	// TODO: parse the posts && keep only the new ones && keep the highest id
	// TODO: return an array of posts (fill the function GetNewPosts())
	brutPosts := parseNewPosts(BoardPosts, getId)
	embedPosts := EmbedNewPosts(brutPosts)
	return embedPosts
}

// 1 2 3
// 1 2 3 4 5
// 1 2 3 -> 4 5
