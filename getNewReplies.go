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

func getRepliesInfos(replies string, postId int) Reply {
	regAuthor := regexp.MustCompile(`\\- \[([a-z0-9@]+)\]`)
	regContent := regexp.MustCompile(`> ([^>\[\]]+)\n>`)
	matchAuthor := regAuthor.FindStringSubmatch(replies)
	matchContent := regContent.FindStringSubmatch(replies)

	fmt.Println(replies)

	p := Reply{
		Content: matchContent[1],
		Author:  matchAuthor[1],
		Id:      postId,
	}
	fmt.Println(p)
	return p
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
			replies = append(replies, getRepliesInfos(c, nb))
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
	re := regexp.MustCompile(fmt.Sprintf("\\br\\/boards:%s\\/([0-9]+)", post))
	var newIdString = re.FindAllStringSubmatch(PostReplies, -1)

	for _, i := range newIdString {
		j, err := strconv.Atoi(i[1])
		if err != nil {
			panic(err)
		}
		if j > maxId[post] {
			return parseNewReplies(PostReplies, post, postTitle, board), nil
		}
	}
	return nil, nil
}
