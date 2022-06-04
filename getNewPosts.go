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

func getBoardsPosts(board string) (string, error) {
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
	Title  string
	Author string
	Id     int
}

func GetPostInfos(post string, id int) Post {
	regAuthor := regexp.MustCompile(`\\- \[(@[a-z]+)\]`)
	regTitle := regexp.MustCompile(`## \[([^\[\]]+)\]`)
	matchTitle := regTitle.FindStringSubmatch(post)
	matchAuthor := regAuthor.FindStringSubmatch(post)

	fmt.Println(post)

	p := Post{
		Title:  matchTitle[1],
		Author: matchAuthor[1],
		Id:     id,
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

func parseNewPosts(BoardPosts string, indexMax int, board string) []*embed.Embed {
	var post []Post
	a := strings.Split(BoardPosts, "----------------------------------------")
	for _, c := range a {
		nb, _ := GetPostID(c)
		if nb > indexMax {
			post = append(post, GetPostInfos(c, nb))
			maxId[board] = nb
		}
	}
	return EmbedNewPosts(post, board)
}

func EmbedNewPosts(posts []Post, board string) []*embed.Embed {
	embeds := make([]*embed.Embed, 0)
	for _, post := range posts {
		embeds = append(embeds, embed.NewEmbed().
			SetTitle(fmt.Sprintf("New post on: %s", board)).
			SetDescription(fmt.Sprintf("**%s**\nhttps://gno.land/r/boards:%s/%d", post.Title, board, post.Id)).
			SetAuthor(post.Author).
			SetColor(0x00FF00))
	}
	fmt.Printf("THERE IS %d NEW POSTS\n", len(embeds))
	return embeds
}

func getNewPosts(board string) ([]*embed.Embed, error) {
	// this return the posts from the watched board
	BoardPosts, err := getBoardsPosts(board)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("\\bpostid=([0-9]+)")
	var newIdString = re.FindAllStringSubmatch(BoardPosts, -1)
	var newId []int

	for _, i := range newIdString {
		j, err := strconv.Atoi(i[1])
		if err != nil {
			panic(err)
		}
		newId = append(newId, j)
	}
	if maxId[board] != 0 {
		if maxId[board] < newId[len(newId)-1] {
			return parseNewPosts(BoardPosts, maxId[board], board), nil
		}
	} else {
		if len(newId) > 0 {
			fmt.Println("first setup for this board:", board)
			maxId[board] = newId[len(newId)-1]
		} else {
			fmt.Println("Empty board:", board)
		}
	}
	return nil, nil
}
