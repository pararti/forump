// user
package entity

type User struct {
	Id           uint32 `json:"id"`
	Name         string `json:"name"`
	RefreshToken string `json:"token,omitempty"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
}

type Post struct {
	Id    uint32 `json:"id"`
	Owner uint32 `json:"owner"`
	URL   string `json:"url"`
	Title string `json:"title"`
	Time  string `json:"time"`
	Anons string `json:"anons"`
	Data  string `json:"data"`
}

type Comment struct {
	Id     uint32 `json:"id"`
	PostId uint32 `json:"postid"`
	Owner  uint32 `json:"owner"`
	Name   string `json:"name"`
	Time   string `json:"time"`
	Data   string `json:"data"`
}

type Token struct {
	Token  string `json:"token"`
	UserId uint32 `json:"userid"`
	Time   int64  `json:"time"`
}

//methods of user
