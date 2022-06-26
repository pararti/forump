package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/cmd/server"
	"github.com/pararti/forump/internals/entity"
)

func main() {

	s, err := server.NewServer(entity.DefaultPSQLConfig())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	router := gin.Default()
	router.LoadHTMLGlob("ui/html/*")
	router.Static("./css", "ui/css")
	router.Static("/auth/css", "ui/css")
	router.Static("/img", "ui/img")
	router.GET("/", s.ViewAll)
	router.GET("/create/", s.CreatePost)
	router.POST("/save_post/", s.SavePost)
	router.GET("/post/:id", s.ViewPost)
	router.GET("/auth/signup", s.GetAuthSignUpPage)
	router.POST("/auth/signup", s.SignUp)
	router.GET("/auth/signin", s.GetAuthSignInPage)
	router.POST("/auth/signin", s.SignIn)
	routerAPI := router.Group("/api")
	{
		routerAPI.GET("/posts/:id", s.GetPostAPI)
		routerAPI.POST("/posts/", s.CreatePostAPI)
		routerAPI.DELETE("/posts/:id", s.DeletePostAPI)
	}

	router.Run(":8080")
}
