package chat

import (
	"lection05/user"
	"sort"
	"time"
)

type Message struct {
	Author user.Login `json:"author"`
	Text   string     `json:"text"`
	TS     time.Time  `json:"ts"`
}

type Chat map[time.Time]Message

func (c Chat) String() string {
	times := make([]time.Time, 0, len(c))

	for t := range c {
		times = append(times, t)
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i].Unix() < times[j].Unix()
	})

	output := ""

	for _, t := range times {
		output += "[" + t.Format(time.RFC822) + "] " + c[t].Author + " > " + c[t].Text + "\n"
	}

	return output
}

func (c Chat) SendMessage(author user.Login, text string) {
	c.sendMessage(author, text, time.Now())
}

func (c Chat) sendMessage(author user.Login, text string, ts time.Time) {
	c[ts] = Message{
		Author: author,
		Text:   text,
		TS:     ts,
	}
}

type PersonalChat map[user.Login]map[user.Login]Chat

func (pc PersonalChat) SendMessage(from user.Login, to user.Login, text string) {
	ts := time.Now()

	if _, ok := pc[from]; !ok {
		pc[from] = map[user.Login]Chat{}
		pc[from][to] = Chat{}

		pc[to] = map[user.Login]Chat{}
		pc[to][from] = Chat{}
	}

	pc[from][to].sendMessage(from, text, ts)
	pc[to][from].sendMessage(from, text, ts)
}
