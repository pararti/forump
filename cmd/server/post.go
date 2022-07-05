package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/internals/entity"
)

func (s *ServerForum) GetPostAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	post, err := s.store.GetPostByID(uint32(id))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, post)
}

func (s *ServerForum) CreatePostAPI(ctx *gin.Context) {
	post := &entity.Post{}
	if err := ctx.ShouldBindJSON(post); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	id, err := s.store.AddPost(post)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Id": id})
}

func (s *ServerForum) DeletePostAPI(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err = s.store.DeletePost(uint32(id))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Id": id})
}

func (s *ServerForum) ViewAll(ctx *gin.Context) {
	posts, err := s.store.GetAllPost()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "index", gin.H{
		"Title": "Posts",
		"posts": posts,
	})

}

func (s *ServerForum) CreatePost(ctx *gin.Context) {
	ok, err := s.SwitcherCookieStatus(s.CheckCookie(ctx), ctx)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		return
	}
	ctx.HTML(http.StatusOK, "creater", gin.H{
		"Title": "Create post",
	})
}

func (s *ServerForum) SavePost(ctx *gin.Context) {
	title := ctx.PostForm("title")
	data := ctx.PostForm("data")
	accessToken, err := ctx.Request.Cookie("access_token")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	id, err := ParseToken(accessToken.Value)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("id:", id)

	post := &entity.Post{
		Owner: id,
		Title: title,
		Data:  data,
	}
	_, err = s.store.AddPost(post)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}

func (s *ServerForum) ViewPost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	post, err := s.store.GetPostByID(uint32(id))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "page", post)
}
