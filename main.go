package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/cmd/server"
)

func main() {
	s := server.NewServer()
	router := gin.Default()
	router.LoadHTMLGlob("ui/html/*")
	router.Static("./css", "ui/css")
	router.Static("/auth/css", "ui/css")
	router.Static("/img", "ui/img")
	router.GET("/", s.ViewAll)
	router.GET("/create/", s.CreatePost)
	router.POST("/save_post/", s.SavePost)
	router.GET("/post/:id", s.ViewPost)
	router.GET("/auth/", s.GetAuthSignUpPage)
	router.POST("/auth/signup/", s.SignUp)
	routerAPI := router.Group("/api")
	{
		routerAPI.GET("/posts/:id", s.GetPostAPI)
		routerAPI.POST("/posts/", s.CreatePostAPI)
		routerAPI.DELETE("/posts/:id", s.DeletePostAPI)
	}

	router.Run(":8080")
}
