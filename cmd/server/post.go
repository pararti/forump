package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/internals/entity"
)

func (s *serverForum) GetPostAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	post, err := s.store.P.Get(uint32(id))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, post)
}

func (s *serverForum) CreatePostAPI(ctx *gin.Context) {
	post := &entity.Post{}
	if err := ctx.ShouldBindJSON(post); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	id := s.store.P.Add(post)
	ctx.JSON(http.StatusOK, gin.H{"Id": id})
}

func (s *serverForum) DeletePostAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	s.store.P.Delete(uint32(id))
	ctx.JSON(http.StatusOK, gin.H{"Id": id})
}

func (s *serverForum) ViewAll(ctx *gin.Context) {
	posts := s.store.P.GetAll()
	ctx.HTML(http.StatusOK, "index", gin.H{
		"Title": "Posts",
		"posts": posts,
	})

}

func (s *serverForum) CreatePost(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "creater", gin.H{
		"Title": "Create post",
	})
}

func (s *serverForum) SavePost(ctx *gin.Context) {
	title := ctx.PostForm("title")
	data := ctx.PostForm("data")
	post := &entity.Post{
		Title: title,
		Data:  data,
	}
	id := s.store.P.Add(post)
	ctx.JSON(http.StatusOK, gin.H{"id": id})
	ctx.Redirect(http.StatusSeeOther, "/")

}

func (s *serverForum) ViewPost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	post, err := s.store.P.Get(uint32(id))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "page", post)
}
