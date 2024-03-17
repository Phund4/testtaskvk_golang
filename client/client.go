package client

import ()

type Client struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float32 `json:"balance"`
}