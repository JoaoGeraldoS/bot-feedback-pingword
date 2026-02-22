package models

import "time"

type FeedBack struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type Update struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}
