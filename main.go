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
			var b []Embed
			for i, p := range newPosts {
				b = append(b, p)
				if (i+1)%10 == 0 {
					err := sendRequest(b)
					if err != nil {
						return
					}
					b = nil
				}
			}
			if len(b) > 0 {
				err := sendRequest(b)
				if err != nil {
					return
				}
			}
		}
	}
}
