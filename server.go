package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/cmd/store"
	"github.com/pararti/forump/internals/entity"
)

type serverForum struct {
	store   *store.CommonStore
	posts10 []entity.Post
}

func NewServer() *serverForum {
	return &serverForum{
		store: store.New(),
	}
}

func (s *serverForum) getPostAPI(ctx *gin.Context) {
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

func (s *serverForum) createPostAPI(ctx *gin.Context) {
	post := &entity.Post{}
	if err := ctx.ShouldBindJSON(post); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	id := s.store.P.Add(post)
	ctx.JSON(http.StatusOK, gin.H{"Id": id})
}

func (s *serverForum) deletePostAPI(ctx *gin.Context) {
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
		"title": "Posts",
		"posts": posts,
	})

}

func main() {
	server := NewServer()
	router := gin.Default()
	router.LoadHTMLGlob("ui/html/*")
	router.Static("/css", "ui/css")
	router.GET("/", server.ViewAll)
	routerAPI := router.Group("/api")
	{
		routerAPI.GET("/posts/:id", server.getPostAPI)
		routerAPI.POST("/posts/", server.createPostAPI)
		routerAPI.DELETE("/posts/:id", server.deletePostAPI)
	}

	router.Run(":8080")
}
