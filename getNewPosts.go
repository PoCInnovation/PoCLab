package main

import (
	"fmt"
	"github.com/enescakir/emoji"
	abci "github.com/gnolang/gno/pkgs/bft/abci/types"
	"github.com/gnolang/gno/pkgs/bft/rpc/client"
	"regexp"
	"strconv"
	"strings"
)

var maxId = make(map[string]int)

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

func getBoardsContents(board string) (string, error) {
	qpath := "vm/qrender"
	data := []byte(fmt.Sprintf("%s\n%s", "gno.land/r/boards", board))
	res, err := makeRequest(qpath, data)

	if err != nil {
		fmt.Println("Error: ", res.Log)
		return "", err
	}
	return string(res.Data), nil
}

type Post struct {
	Title       string
	Author      string
	Description string
	Id          int
}

type Author struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      Author `json:"author"`
	Color       int    `json:"color"`
}

func GetPostInfos(post string, id int) Post {
	regAuthor := regexp.MustCompile(`\\- \[(@[a-z]+)\]`)
	regTitle := regexp.MustCompile(`## \[([^\[\]]+)\]`)
	regDescription := regexp.MustCompile(`(?s)\)\n\n.*\n\\`)
	matchTitle := regTitle.FindStringSubmatch(post)
	matchAuthor := regAuthor.FindStringSubmatch(post)
	matchDescription := regDescription.FindStringSubmatch(post)[0][3:]
	matchDescription = matchDescription[:len(matchDescription)-2]

	fmt.Println(post)

	p := Post{
		Title:       matchTitle[1],
		Author:      matchAuthor[1],
		Description: matchDescription,
		Id:          id,
	}
	fmt.Println(p)
	return p
}

func GetPostID(s string) (int, error) {
	re := regexp.MustCompile("\\bpostid=([0-9]+)")
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return 0, nil
	}
	return strconv.Atoi(match[1])
}

func parseNewPosts(BoardPosts string, board string) []Embed {
	var post []Post
	newMaxId := maxId[board]
	a := strings.Split(BoardPosts, "----------------------------------------")
	for _, c := range a {
		nb, _ := GetPostID(c)
		//GetNewReplies(board, nb)
		if nb > maxId[board] {
			post = append(post, GetPostInfos(c, nb))
			if nb > newMaxId {
				newMaxId = nb
			}
		}
	}
	maxId[board] = newMaxId
	return EmbedNewPosts(post, board)
}

func EmbedNewPosts(posts []Post, board string) []Embed {
	embeds := make([]Embed, 0)
	for _, post := range posts {
		embeds = append(embeds, Embed{
			Title:       fmt.Sprintf("New post on: %s %v ", board, emoji.OpenMailboxWithRaisedFlag),
			Description: fmt.Sprintf("**%s**\n%s\n\nhttps://gno.land/r/boards:%s/%d", post.Title, post.Description, board, post.Id),
			Author: Author{
				Name:    post.Author,
				IconUrl: "https://cdn.discordapp.com/attachments/981266192390049846/983052513932607488/unknown.png",
			},
			Color: 7212552,
		})
	}
	fmt.Printf("THERE IS %d NEW POSTS\n", len(embeds))
	return embeds
}

func getNewPosts(board string) ([]Embed, error) {
	// this return the posts from the watched board
	BoardPosts, err := getBoardsContents(board)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("\\bpostid=[0-9]+/([0-9]+)")
	var newIdString = re.FindAllStringSubmatch(BoardPosts, -1)
	// var newId []int

	for _, i := range newIdString {
		j, err := strconv.Atoi(i[1])
		if err != nil {
			panic(err)
		}
		if j > maxId[board] {
			return parseNewPosts(BoardPosts, board), nil
		}
	}
	return nil, nil
}
