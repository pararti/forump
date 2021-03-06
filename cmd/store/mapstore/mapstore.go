package store

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/pararti/forump/internals/entity"
)

type PostStore struct {
	m      sync.Mutex
	storeg map[uint32]*entity.Post
	nextId uint32
}

type CommentStore struct {
	m      sync.Mutex
	storeg map[uint32]*entity.Comment
	nextId uint32
}

type UserStore struct {
	m      sync.Mutex
	storeg map[string]*entity.User
	nextId uint32
}

type TokenStore struct {
	m      sync.Mutex
	storeg map[string]*entity.Token
}

type CommonStore struct {
	U UserStore
	P PostStore
	C CommentStore
	T TokenStore
}

func New() *CommonStore {
	return &CommonStore{
		U: UserStore{storeg: make(map[string]*entity.User)},
		P: PostStore{storeg: make(map[uint32]*entity.Post)},
		C: CommentStore{storeg: make(map[uint32]*entity.Comment)},
		T: TokenStore{storeg: make(map[string]*entity.Token)},
	}
}

//methods for user store
func (u *UserStore) Add(user *entity.User) uint32 {
	u.storeg[user.Email] = user
	user.Id = u.nextId
	u.nextId += 1
	return u.nextId - 1
}

/*
func (u *UserStore) Get(id uint32) (*entity.User, error) {
	res, ok := u.storeg[id]
	if ok {
		return res, nil
	}
	return &entity.User{}, errors.New("Not found")
}
*/

func (u *UserStore) GetByEmail(email string) (*entity.User, error) {
	user, ok := u.storeg[email]
	if ok {
		return user, nil
	}
	return &entity.User{}, errors.New("Not found")
}

func (u *UserStore) Delete(email string) {
	delete(u.storeg, email)
}

//methods of post sore
func (p *PostStore) Add(post *entity.Post) uint32 {
	post.Id = p.nextId
	post.URL = "/post/" + strconv.FormatUint(uint64(post.Id), 10)
	post.Time = time.Now().Format("2 Jan 2006 в 15:04")
	if len(post.Data) > 140 {
		post.Anons = post.Data[0:139] + "..."
	} else {
		post.Anons = post.Data
	}
	p.storeg[post.Id] = post
	p.nextId += 1
	return p.nextId - 1
}

func (p *PostStore) Get(id uint32) (*entity.Post, error) {
	post, ok := p.storeg[id]
	if ok {
		return post, nil
	}
	return &entity.Post{}, errors.New("Not found")
}

func (p *PostStore) GetAll() []*entity.Post {
	posts := make([]*entity.Post, 0, len(p.storeg))
	for i := len(p.storeg) - 1; i >= 0; i-- {
		posts = append(posts, p.storeg[uint32(i)])

	}
	return posts
}

func (p *PostStore) Delete(id uint32) {
	delete(p.storeg, id)
}

//metods of comment
func (c *CommentStore) Add(com *entity.Comment) uint32 {
	com.Id = c.nextId
	c.storeg[c.nextId] = com
	c.nextId += 1
	return c.nextId - 1
}

func (c *CommentStore) Get(id uint32) (*entity.Comment, error) {
	com, ok := c.storeg[id]
	if ok {
		return com, nil
	}
	return &entity.Comment{}, errors.New("Not found")
}

func (c *CommentStore) Delete(id uint32) {
	delete(c.storeg, id)
}

func (t *TokenStore) Add(token string, uid uint32) {
	newToken := &entity.Token{
		Token:  token,
		UserId: uid,
		Time:   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	t.storeg[token] = newToken
}

func (t *TokenStore) Delete(token string) {
	delete(t.storeg, token)
}

func (t *TokenStore) Get(token string) (*entity.Token, error) {
	tk, ok := t.storeg[token]
	if ok {
		return tk, nil
	}
	return &entity.Token{}, errors.New("Not found")
}

func (t *TokenStore) Check(token string) bool {
	if tk, ok := t.storeg[token]; ok {
		if tk.Time > time.Now().Unix() {
			return true
		}
	}
	return false
}
