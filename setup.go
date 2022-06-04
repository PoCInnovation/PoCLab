package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "isma"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type timer time.Duration

func (t *timer) String() string {
	return time.Duration(*t).String()
}

func (t *timer) Set(value string) error {
	d, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*t = timer(time.Duration(d) * time.Second)
	return nil
}

var (
	Token     string
	ChannelID string
	Seconde   timer
	Boards    arrayFlags
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ChannelID, "c", "", "Channel ID")
	flag.Var(&Seconde, "s", "second  between refresh")
	flag.Var(&Boards, "b", "board to notify")
	flag.Parse()
}

func setup(Boards []string) error {
	if Seconde == 0 {
		Seconde = 5
	}
	for _, board := range Boards {
		qpath := "vm/qrender"
		data := []byte(fmt.Sprintf("%s\n%s", "gno.land/r/boards", board))
		res, err := makeRequest(qpath, data)

		if err != nil {
			return err
		}

		re := regexp.MustCompile("\\b(board does not exist:)")
		match := re.FindStringSubmatch(string(res.Data))

		if match != nil {
			return fmt.Errorf("%s", string(res.Data))
		}
	}
	return nil
}
