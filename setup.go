package main

import (
	"flag"
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

func setup() {
	if Seconde == 0 {
		Seconde = 5
	}
	//TODO: check if boards are valid
}
