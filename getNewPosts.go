package main

import (
	"fmt"
	"regexp"
	"strings"
)

var maxId = make(map[string]int)

type Post struct {
	Title       string
	Author      string
	Description string
	Id          int
}

func getPostInfos(post string, id int) (Post, error) {
	regTitle := regexp.MustCompile(`## \[([^\[\]]+)\]`)
	regAuthor := regexp.MustCompile(`\\- \[([a-zA-Z0-9_.-@]+)\]`)
	regDescription := regexp.MustCompile(`(?s)\)\n\n.*\n\\`)
	matchTitle := regTitle.FindStringSubmatch(post)
	matchAuthor := regAuthor.FindStringSubmatch(post)
	matchDescription := regDescription.FindStringSubmatch(post)[0][3:]
	matchDescription = matchDescription[:len(matchDescription)-2]

	fmt.Println(post)

	if len(matchTitle) < 2 || len(matchAuthor) < 2 {
		return Post{}, fmt.Errorf("bad format reply :%s", post)
	}
	p := Post{
		Title:       matchTitle[1],
		Author:      matchAuthor[1],
		Description: matchDescription,
		Id:          id,
	}
	fmt.Println(p)
	return p, nil
}

func getPostTitle(s string) string {
	regTitle := regexp.MustCompile(`## \[([^\[\]]+)\]`)
	match := regTitle.FindStringSubmatch(s)
	if len(match) == 0 {
		return ""
	}
	return match[1]
}

func parseNewPosts(BoardPosts string, board string) []Embed {
	var posts []Post
	newMaxId := maxId[board]
	a := strings.Split(BoardPosts, "----------------------------------------")
	for _, c := range a {
		nb, _ := getID(c)
		if DoesReply {
			r, err := getNewReplies(fmt.Sprintf("%s/%d", board, nb), board, getPostTitle(c))
			if err != nil {
				return nil
			}
			if len(r) > 0 {
				queueRequest(r)
			}
		}
		if nb > maxId[board] {
			post, err := getPostInfos(c, nb)
			if err != nil {
				continue
			}
			posts = append(posts, post)
			if nb > newMaxId {
				newMaxId = nb
			}
		}
	}
	maxId[board] = newMaxId
	return embedPosts(posts, board)
}

func getNewPosts(board string) ([]Embed, error) {
	// this return the posts from the watched board
	BoardPosts, err := getBoardsContents(board)
	if err != nil {
		return nil, err
	}

	ids := getMessagesIds(BoardPosts, fmt.Sprintf("\\bpostid=([0-9]+)"))

	for _, i := range ids {
		if i > maxId[board] || DoesReply {
			return parseNewPosts(BoardPosts, board), nil
		}
	}
	return nil, nil
}
