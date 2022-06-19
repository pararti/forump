package server

import (
	"crypto/sha256"
	_ "errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/internals/entity"
)

const (
	salt      = "hjkoiopoemohparartisalt32g" //change this values if u will use that code
	secretKey = "ajfnjnivmdo4jmfsmcsm2miemfsmvkbotlsoewjfjvmxmnmi8mamiemfmaadfvmiw"
)

type tokenClaims struct {
	jwt.StandardClaims
	email string `json:"email"`
}

func hashPasswd(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return string(hash.Sum([]byte(salt)))
}

func generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		email,
	})

	return token.SignedString([]byte(secretKey))
}

func (s *serverForum) GetAuthSignUpPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup", gin.H{
		"Title": "Registration",
	})

}

func (s *serverForum) GetAuthSignInPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signin", gin.H{
		"Title": "sign in",
	})

}

func (s *serverForum) SignUp(ctx *gin.Context) {
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	passwd := hashPasswd(ctx.PostForm("password"))
	_, err := s.store.U.GetByEmail(email)
	if err != nil {
		token, err := generateToken(email)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
		}
		user := &entity.User{
			Name:     name,
			Token:    token,
			Email:    email,
			Password: passwd,
		}
		s.store.U.Add(user)
		f := SuccessHandler("Регистрация прошла успешно!")
		f(ctx)
	} else {
		f := ErrorHandler(http.StatusConflict, "Пользователь уже существует!", "/auth/signin/", "Войти")
		f(ctx)
	}
}

func (s *serverForum) SignIn(ctx *gin.Context) {
	email := ctx.PostForm("email")
	user, err := s.store.U.GetByEmail(email)
	if err != nil {
		f := ErrorHandler(http.StatusNotFound, "Пользователь не найден", "/auth/signup/", "Регистрация")
		f(ctx)
	} else {
		passwd := hashPasswd(ctx.PostForm("password"))
		if passwd == user.Password {
			f := SuccessHandler("Добро пожаловать!")
			f(ctx)
		} else {
			f := ErrorHandler(http.StatusNotFound, "Данные введены неверно!", "/auth/signin/", "Попробовать снова")
			f(ctx)
		}
	}
}
