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

func GetPostInfos(post string) Post {
	// cath regex "\\- \[(@[a-z]+)\]"
	// matchTitle[1] is the author
	// matchTitle[2] is the title
	regAuthor := regexp.MustCompile(`\\- \[(@[a-z]+)\]`)
	regTitle := regexp.MustCompile(`## \[([^\[\]]+)\]`)
	// get group 1 of regex in post
	matchTitle := regTitle.FindStringSubmatch(post)
	matchAuthor := regAuthor.FindStringSubmatch(post)

	fmt.Println(post)
	fmt.Println("ICI")
	fmt.Println(matchTitle)
	fmt.Println(matchAuthor)

	p := Post{
		Title:  matchTitle[1],
		Author: matchAuthor[1],
	}
	fmt.Println(p)
	return p
}

func GetPostID(s string) (int, error) {
	re := regexp.MustCompile("\\bpostid=([0-9]+)")
	match := re.FindStringSubmatch(s)
	//fmt.Println(s)
	if len(match) == 0 {
		return 0, nil
	}
	//fmt.Printf("THIS IS %s\n", match[1])
	return strconv.Atoi(match[1])
}

func parseNewPosts(BoardPosts string, index []string, indexMax int) []*embed.Embed {
	//TODO: replace by real parsing
	var post []Post
	a := strings.Split(BoardPosts, "----------------------------------------")
	for _, c := range a {
		nb, _ := GetPostID(c)
		fmt.Printf("ID MAX IS %d\n", indexMax)
		if nb > indexMax {
			fmt.Println("New post HERE")
			post = append(post, GetPostInfos(c))
			maxId[len(maxId)-1] = nb
		}
	}
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
	fmt.Printf("THERE IS %d NEW POSTS\n", len(embeds))
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
	//fmt.Println(BoardPosts)

	re := regexp.MustCompile("\\bpostid=([0-9]+)")
	fr := regexp.MustCompile("(\\b[0-9]+)")
	if len(maxId) > 0 {
		fmt.Printf("Parsing new posts %d\n", maxId[len(maxId)-1])
	}
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
					return parseNewPosts(BoardPosts, newIdString, maxId[len(maxId)-1])
				}
				// fmt.Println(i, s)
			}
			// fmt.Printf("Il a y eu %d nouveaux msg\n", maxId[i])
			fmt.Println("no new id")
			return nil
		}
	} else {
		fmt.Println("first setup")
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
