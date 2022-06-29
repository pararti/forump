package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/internals/entity"
	"github.com/pararti/forump/pkg/random"
	"golang.org/x/crypto/bcrypt"
)

const (
	COOKIE_NOT_FOUND = iota
	COOKIE_NOT_FOUND_ACCESS
	COOKIE_NOT_FOUND_REFRESH
	COOKIE_IS_OK
)

const (
	//change this values if u will use that code
	secretKey = "ajfnjnivmdo4jmfsmcsm2miemfsmvkbotlsoewjfjvmxmnmi8mamiemfmaadfvmiw"
	ttl       = 12 * 60 * 60
)

func hashPasswd(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func generateRefrashToken() (string, error) {
	s := random.RandomString(64)
	return s, nil
}

func generateAccessToken(id uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl * time.Second).Unix(),
		Id:        strconv.FormatUint(uint64(id), 10),
		IssuedAt:  time.Now().Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func GetRefreshToken(id uint32) (string, string, error) {
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

func (s *ServerForum) Refreshing(ctx *gin.Context) error {
	cookie, err := ctx.Request.Cookie("refresh_token")
	if err != nil {
		return err
	}
	id, err := s.store.GetTokenID(cookie.Value)
	if err != nil {
		return err
	}
	a, r, err := GetRefreshToken(id)
	if err != nil {
		return err
	}
	SetCookieAuth(a, r, ctx)
	err = s.store.UpdateToken(a, id)
	if err != nil {
		return err
	}
	return nil
}

func SetCookieAuth(access, refresh string, ctx *gin.Context) {
	ctx.SetCookie("access_token", access, ttl, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh, 60*60*24*30, "/", "localhost", false, true)

}

func (s *ServerForum) SwitcherCookieStatus(status int, ctx *gin.Context) (bool, error) {
	switch status {
	case COOKIE_NOT_FOUND:
		f := ErrorHandler(http.StatusUnauthorized, "Пожалуйста зарегистрируйтесь", "/auth/signup", "Регистрация")
		f(ctx)
		return false, nil
	case COOKIE_NOT_FOUND_ACCESS:
		err := s.Refreshing(ctx)
		if err != nil {
			return false, err
		}
	case COOKIE_NOT_FOUND_REFRESH:
		f := ErrorHandler(http.StatusUnauthorized, "Пожалуйста войдите", "/auth/signin", "Войти")
		f(ctx)
		return false, nil
	case COOKIE_IS_OK:
		return true, nil
	}
	return true, nil
}

func ParseToken(access string) (uint32, error) {
	token, err := jwt.ParseWithClaims(access, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		id, err := strconv.ParseUint(claims.Id, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(id), nil
	}
	return 0, err
}

func (s *ServerForum) CheckCookie(ctx *gin.Context) int {
	_, err := ctx.Request.Cookie("access_token")
	_, err2 := ctx.Request.Cookie("refresh_token")
	if err != nil && err2 != nil {
		return COOKIE_NOT_FOUND
	}
	if err != nil && err2 == nil {
		return COOKIE_NOT_FOUND_ACCESS
	}
	if err2 != nil && err == nil {
		return COOKIE_NOT_FOUND_REFRESH
	}
	return COOKIE_IS_OK
}

func (s *ServerForum) GetAuthSignUpPage(ctx *gin.Context) {
	stat := s.CheckCookie(ctx)
	if stat == COOKIE_NOT_FOUND_ACCESS {
		err := s.Refreshing(ctx)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		stat = COOKIE_IS_OK
	}
	if stat == COOKIE_IS_OK {
		f := SuccessHandler("Вы уже авторизированы")
		f(ctx)
		return
	}
	if stat == COOKIE_NOT_FOUND_REFRESH {
		f := ErrorHandler(http.StatusUnauthorized, "Пожалуйста войдите", "/auth/signin", "Войти")
		f(ctx)
		return
	}
	ctx.HTML(http.StatusOK, "signup", gin.H{
		"Title": "Registration",
	})

}

func (s *ServerForum) GetAuthSignInPage(ctx *gin.Context) {
	stat := s.CheckCookie(ctx)
	if stat == COOKIE_NOT_FOUND_ACCESS {
		err := s.Refreshing(ctx)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		stat = COOKIE_IS_OK
	}
	if stat == COOKIE_IS_OK {
		f := SuccessHandler("Вы уже авторизированы")
		f(ctx)
		return
	}
	if stat == COOKIE_NOT_FOUND {
		f := ErrorHandler(http.StatusUnauthorized, "Пожалуйста авторизируйтесь", "/auth/signup", "Регистрация")
		f(ctx)
		return
	}
	ctx.HTML(http.StatusOK, "signin", gin.H{
		"Title": "sign in",
	})

}

func (s *ServerForum) SignUp(ctx *gin.Context) {
	/*stat := s.CheckCookie(ctx)
	if stat != COOKIE_NOT_FOUND {
		f := SuccessHandler("Вы уже авторизированы")
		f(ctx)
		return
	}
	*/
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	passwd, err := hashPasswd(ctx.PostForm("password"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("email, name, passwd", email, name, passwd)
	b, err := s.store.CheckUserByEmail(email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if !b {
		refreshToken, err := generateRefrashToken()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		s.store.AddToken(refreshToken, 0)
		user := &entity.User{
			Name:         name,
			RefreshToken: refreshToken,
			Email:        email,
			Password:     passwd,
		}
		id, err := s.store.AddUser(user)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		err = s.store.SetTokenUserID(refreshToken, id)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		accessToken, err := generateAccessToken(id)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		SetCookieAuth(accessToken, refreshToken, ctx)
		f := SuccessHandler("Регистрация прошла успешно!")
		f(ctx)
	} else {
		f := ErrorHandler(http.StatusConflict, "Пользователь уже существует!", "/auth/signin/", "Войти")
		f(ctx)
	}
}

func (s *ServerForum) SignIn(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password, err := s.store.GetUserPasswordByEmail(email)
	if err != nil {
		f := ErrorHandler(http.StatusNotFound, "Пользователь не найден", "/auth/signup/", "Регистрация")
		f(ctx)
		return
	} else {
		passwd := ctx.PostForm("password")
		err := bcrypt.CompareHashAndPassword([]byte(password), []byte(passwd))
		if err == nil {
			f := SuccessHandler("Добро пожаловать!")
			f(ctx)
		} else {
			f := ErrorHandler(http.StatusNotFound, "Данные введены неверно!", "/auth/signin/", "Попробовать снова")
			f(ctx)
		}
	}
}
