package server

import (
	"crypto/sha256"
	_ "errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/internals/entity"
)

const (
	COOKIE_NOT_FOUND = iota
	COOKIE_NEED_REFRESH
	COOKIE_NEED_AUTH
	COOKIE_IS_OK
)

const (
	salt      = "hjkoiopoemohparartisalt32g" //change this values if u will use that code
	secretKey = "ajfnjnivmdo4jmfsmcsm2miemfsmvkbotlsoewjfjvmxmnmi8mamiemfmaadfvmiw"
	ttl       = 12 * 60 * 60
)

type tokenClaims struct {
	jwt.StandardClaims
	id uint32 `json:"id"`
}

func hashPasswd(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return string(hash.Sum([]byte(salt)))
}

func generateRefrashToken() (string, error) {
	b := make([]byte, 256)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	if _, err := r.Read(b); err != nil {
		return "", err
	}
	return string(b), nil
}

func generateAccessToken(id uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl * time.Second).Unix(),
			Id:        id,
			IssuedAt:  time.Now().Unix(),
		},
	})
	return token.SignedString([]byte(secretKey))
}

func Refreshing(id uint32) (string, string, error) {
	refresh, err := generateRefrashToken()
	if err != nil {
		return "", "", err
	}
	access, err := generateAccessToken(id)
	if err != nil {
		return "", "", nil
	}
	return access, refresh, nil
}

func SetCookieAuth(access, refresh string, ctx *gin.Context) {
	ctx.SetCookie("access_token", access, ttl, "/", false, true)
	ctx.SetCookie("refresh_token", refresh, ttl, "/", false, true)

}

//func SwitcherCookieStatus(status int)

func (s *serverForum) CheckCookie(ctx *gin.Context) int {
	cookie, err := ctx.Request.Cookie("access_token")
	if err != nil {
		return COOKIE_NOT_FOUND
	}
	if cookie.MaxAge < int(time.Now().Unix()-time.Now().Add(-1*cookie.MaxAge).Unix()) {
		return COOKIE_NEED_REFRESH
	}
	cookie2, err := ctx.Request.Cookie("refresh_token")
	if err != nil {
		return COOKIE_NOT_FOUND
	}
	ok := s.store.T.Check(cookie2.Value)
	if !ok {
		return COOKIE_NEED_AUTH
	}
	return COOKIE_IS_OK
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
		refreshToken, err := generateRefrashToken()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
		}
		user := &entity.User{
			Name:              name,
			RefreshTokenToken: refreshToken,
			Email:             email,
			Password:          passwd,
		}
		id := s.store.U.Add(user)
		s.store.T.Add(refreshToken, id)
		accessToken, err := generateAccessToken(id)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
		}
		SetCookieAuth(accessToken, refreshToken, ctx)
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
