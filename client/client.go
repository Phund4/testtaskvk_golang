package client

import ()

type Client struct {
	ID      int     `json:"client_id"`
	Name    string  `json:"name"`
	Balance float32 `json:"balance"`
}