package store

import (
	"errors"
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
	storeg map[uint32]*entity.User
	nextId uint32
}

type CommonStore struct {
	U UserStore
	P PostStore
	C CommentStore
}

func New() *CommonStore {
	return &CommonStore{
		U: UserStore{storeg: make(map[uint32]*entity.User)},
		P: PostStore{storeg: make(map[uint32]*entity.Post)},
		C: CommentStore{storeg: make(map[uint32]*entity.Comment)},
	}
}

//methods of user store
func (u *UserStore) Create(name string) uint32 {
	u.storeg[u.nextId] = &entity.User{Id: u.nextId, Name: name}
	u.nextId += 1
	return u.nextId - 1
}

func (u *UserStore) Get(id uint32) (*entity.User, error) {
	res, ok := u.storeg[id]
	if ok {
		return res, nil
	}
	return &entity.User{}, errors.New("Not found")
}

func (u *UserStore) Delete(id uint32) {
	delete(u.storeg, id)
}

//methods of post sore
func (p *PostStore) Add(post *entity.Post) uint32 {
	post.Id = p.nextId
	post.Time = time.Now().Format("2 Jan 2006 в 15:04")
	post.Anons = post.Data[0:150] + "..."
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
