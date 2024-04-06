package quest

import ()

type Quest struct {
	ID   int     `json:"quest_id"`
	Name string  `json:"name"`
	Cost float32 `json:"cost"`
}
