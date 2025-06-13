package storage

import "github.com/google/uuid"

type Entry struct {
	Id   uuid.UUID
	Text string
}

func (e Entry) String() string {
	return e.Text
}
