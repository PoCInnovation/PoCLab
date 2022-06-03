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

var maxId []int 

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
 
func parseNewPosts(BoardPosts string, index []string, indexMax int) []Post {
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
	fr := regexp.MustCompile("((\\b[0-9]+)$)")
	fmt.Println("Parsing new posts")
	fmt.Println("------------------")
	var test = re.FindAllString(BoardPosts, -1)
	fmt.Println(test)
	// fmt.Println(fr.FindAllString(strings.Join(test, " "), -1))
	newId, err := strconv.Atoi(fr.FindAllString(strings.Join(test, " "), -1))
	var newId = [1, 2, 3, 4]
	if len(maxId) != 0 {
		if maxId[len(maxId)-1] < newId[len(newId)-1] {
			// fmt.Printf("(%d)(%d)\n", maxId[len(maxId)-1], newId[len(newId)-1])
			var prevMaxId = maxId[len(maxId)-1]
			
			for i := range newId {
				if newId[i] > prevMaxId {
					fmt.Printf("NewId %d\n", newId[i])
					parseNewPosts(BoardPosts, test, newId[i])
					return 
				}
				// fmt.Println(i, s)
			}
			// fmt.Printf("Il a y eu %d nouveaux msg\n", maxId[i])
			maxId = newId
			return 
		}
	} else {
		maxId = newId
	}
	fmt.Println(maxId)
	// TODO: parse the posts && keep only the new ones && keep the highest id
	// TODO: return an array of posts (fill the function GetNewPosts())
	brutPosts := parseNewPosts(BoardPosts, test)
	embedPosts := EmbedNewPosts(brutPosts)
	return embedPosts
}