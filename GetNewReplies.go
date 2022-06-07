package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Reply struct {
	Content string
	Author  string
	Id      int
}

func GetRepliesInfos(replies string, postId int) Reply {
	regAuthor := regexp.MustCompile(`\\- \[([a-z1-9@]+)\]`)
	regContent := regexp.MustCompile(`> ([^>\[\]]+)\n>`)
	matchAuthor := regAuthor.FindStringSubmatch(replies)
	matchContent := regContent.FindStringSubmatch(replies)

	//fmt.Println(replies)

	p := Reply{
		Content: matchContent[1],
		Author:  matchAuthor[1],
		Id:      postId,
	}
	fmt.Println(p)
	return p
}

func EmbedNewReplies(replies []Reply, post string) []Embed {
	var embeds []Embed
	for _, r := range replies {
		embed := Embed{
			Title:       fmt.Sprintf("New reply on %s", post),
			Description: fmt.Sprintf("%s said: %s", r.Author, r.Content),
			Author: Author{
				Name:    r.Author,
				IconUrl: "https://cdn.discordapp.com/attachments/981266192390049846/983052513932607488/unknown.png",
			},
			Color: 0x00ff00,
		}
		embeds = append(embeds, embed)
	}
	return embeds
}

func ParsePostsReplies(postReplies string) []string {
	postReplies = strings.Split(postReplies, "\n\\- [@")[1]
	a := strings.Split(postReplies, "\n\n")[1:]
	var replies []string
	for _, c := range a {
		replies = append(replies, strings.Split(c, "> \n")...)
	}
	return replies
}

func ParseNewReplies(postReplies string, post string) []Embed {
	var replies []Reply
	newMaxId := maxId[post]
	// parse the replies
	a := ParsePostsReplies(postReplies)
	for _, c := range a {
		nb, _ := GetPostID(c)
		if nb > maxId[post] {
			replies = append(replies, GetRepliesInfos(c, nb))
			if nb > newMaxId {
				newMaxId = nb
			}
		}
	}
	maxId[post] = newMaxId
	return EmbedNewReplies(replies, post)
}

func GetNewReplies(post string, board string) ([]Embed, error) {
	// this return the posts from the watched board
	PostReplies, err := getBoardsContents(post)
	if err != nil {
		return nil, err
	}
	//fmt.Println(PostReplies)
	re := regexp.MustCompile(fmt.Sprintf("\\br\\/boards:%s\\/([0-9]+)", post))
	var newIdString = re.FindAllStringSubmatch(PostReplies, -1)

	for _, i := range newIdString {
		j, err := strconv.Atoi(i[1])
		if err != nil {
			panic(err)
		}
		if j > maxId[post] {
			return ParseNewReplies(PostReplies, post), nil
		}
	}
	return nil, nil
}
