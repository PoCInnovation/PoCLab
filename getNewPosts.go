package main

import (
	"fmt"
	embed "github.com/Clinet/discordgo-embed"
	abci "github.com/gnolang/gno/pkgs/bft/abci/types"
	"github.com/gnolang/gno/pkgs/bft/rpc/client"
	"regexp"
	"strconv"
	"strings"
)

var maxId []int

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

func parseNewPosts(BoardPosts string, index []string, indexMax int) []*embed.Embed {
	//TODO: replace by real parsing
	post := make([]Post, 0)
	//for _, i := range index {
	//	if i > indexMax {
	//		post = append(post, Post{Title: "", Author: ""})
	//	}
	//}
	return EmbedNewPosts(post)
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
	// fmt.Println(BoardPosts)

	re := regexp.MustCompile("\\bpostid=[0-9]+")
	fr := regexp.MustCompile("(\\b[0-9]+)")
	fmt.Println("Parsing new posts")
	// fmt.Println("------------------")
	var getId = re.FindAllString(BoardPosts, -1)
	// fmt.Println(getId)
	var newIdString = fr.FindAllString(strings.Join(getId, " "), -1)
	var newId []int

	for _, i := range newIdString {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		newId = append(newId, j)
	}
	if len(maxId) != 0 {
		if maxId[len(maxId)-1] < newId[len(newId)-1] {
			// fmt.Printf("(%d)(%d)\n", maxId[len(maxId)-1], newId[len(newId)-1])
			var prevMaxId = maxId[len(maxId)-1]

			for i := range newId {
				if newId[i] > prevMaxId {
					fmt.Printf("NewId %d\n", newId[i])
					maxId = newId
					return parseNewPosts(BoardPosts, newIdString, newId[i])
				}
				// fmt.Println(i, s)
			}
			// fmt.Printf("Il a y eu %d nouveaux msg\n", maxId[i])
			return nil
		}
	} else {
		maxId = newId
	}
	fmt.Println()

	// TODO: parse the posts && keep only the new ones && keep the highest id
	// TODO: return an array of posts (fill the function GetNewPosts())
	// brutPosts := parseNewPosts(BoardPosts, getId)
	// embedPosts := EmbedNewPosts(brutPosts)
	return nil
}

// 1 2 3
// 1 2 3 4 5
// 1 2 3 -> 4 5
