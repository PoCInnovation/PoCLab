package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Reply struct {
	Content string
	Author  string
	Id      int
}

func getReplyInfos(reply string, postId int) (Reply, error) {
	regAuthor := regexp.MustCompile(`\\- \[([a-zA-Z0-9_.-@]+)\]`)
	regContent := regexp.MustCompile(`> ([^>\[\]]+)\n>`)
	matchAuthor := regAuthor.FindStringSubmatch(reply)
	matchContent := regContent.FindStringSubmatch(reply)

	fmt.Println(reply)

	if len(matchContent) < 2 || len(matchAuthor) < 2 {
		return Reply{}, fmt.Errorf("bad format reply :%s", reply)
	}
	p := Reply{
		Content: matchContent[1],
		Author:  matchAuthor[1],
		Id:      postId,
	}
	fmt.Println(p)
	return p, nil
}

func parsePostsReplies(postReplies string) []string {
	postReplies = strings.Split(postReplies, "\n\\- [@")[1]
	a := strings.Split(postReplies, "\n\n")[1:]
	var replies []string
	for _, c := range a {
		replies = append(replies, strings.Split(c, "> \n")...)
	}
	return replies
}

func parseNewReplies(postReplies string, post string, postTitle string, board string) []Embed {
	var replies []Reply
	newMaxId := maxId[post]
	a := parsePostsReplies(postReplies)
	for _, c := range a {
		nb, _ := getID(c)
		if nb > maxId[post] {
			reply, err := getReplyInfos(c, nb)
			if err != nil {
				continue
			}
			replies = append(replies, reply)
			if nb > newMaxId {
				newMaxId = nb
			}
		}
	}
	maxId[post] = newMaxId
	return embedReplies(replies, post, postTitle, board)
}

func getNewReplies(post string, board string, postTitle string) ([]Embed, error) {
	// this return the reply from the post
	PostReplies, err := getBoardsContents(post)
	if err != nil {
		return nil, err
	}

	ids := getMessagesIds(PostReplies, fmt.Sprintf("\\br\\/boards:%s\\/([0-9]+)", post))

	for _, i := range ids {
		if i > maxId[post] {
			return parseNewReplies(PostReplies, post, postTitle, board), nil
		}
	}
	return nil, nil
}
