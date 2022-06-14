// user
package entity

import (
	"time"
)

type User struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	Id    uint32 `json:"id"`
	Owner uint32 `json:"owner"`
	Title string `json:"title"`
	Time  string `json:"time"`
	Anons string `json:"anons"`
	Data  string `json:"data"`
}

type Comment struct {
	Id     uint32    `json:"id"`
	PostId uint32    `json:"postid"`
	Owner  uint32    `json:"owner"`
	Name   string    `json:"name"`
	Time   time.Time `json:"time"`
	Data   string    `json:"data"`
}

//methods of user
