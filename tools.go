package main

import (
	"regexp"
	"strconv"
)

func formatAuthor(author string) string {
	if author[0] != '@' {
		return "@randomUser"
	}
	return author
}

func getID(s string) (int, error) {
	re := regexp.MustCompile("\\bpostid=([0-9]+)")
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return 0, nil
	}
	return strconv.Atoi(match[1])
}
