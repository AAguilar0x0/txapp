package models

import "time"

type User interface {
	GetID() string
	GetPassword() string
	GetRole() string
}

type Token interface {
	GetID() string
	GetSub() string
	GetIss() string
	GetExp() time.Time
}
