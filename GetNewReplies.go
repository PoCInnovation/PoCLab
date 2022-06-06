package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func GetNewReplies(board string, postId int) ([]Embed, error) {
	// this return the posts from the watched board
	BoardPosts, err := getBoardsContents(fmt.Sprintf("%s/%d", board, postId))
	if err != nil {
		return nil, err
	}
	fmt.Println(BoardPosts)
	re := regexp.MustCompile(fmt.Sprintf("\\br\\/boards:%s\\/[0-9]+\\/([0-9]+)", board))
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
