package main

import (
	"fmt"
	"time"
)

type Body struct {
	Embeds []Embed `json:"embeds"`
}

func main() {
	err := setup(Boards)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for true {
		time.Sleep(time.Duration(Seconde))

		for _, board := range Boards {
			newPosts, warn := getNewPosts(board)
			if warn != nil {
				fmt.Println("error in getNewPosts.")
				return
			}
			//request 10 by 10
			QueueRequest(newPosts)
		}
	}
}
