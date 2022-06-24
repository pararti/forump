package server

import (
	"crypto/sha256"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pararti/forump/internals/entity"
	"github.com/pararti/forump/pkg/random"
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

func hashPasswd(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return string(hash.Sum([]byte(salt)))
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
	accessToken, _ := ctx.Request.Cookie("access_token")
	refreshToken, _ := ctx.Request.Cookie("refresh_token")
	id, err := ParseToken(accessToken.Value)
	if err != nil {
		return err
	}
	a, r, err := GetRefreshToken(id)
	if err != nil {
		return err
	}
	SetCookieAuth(a, r, ctx)
	s.store.AddToken(r, id)
	s.store.T.Delete(refreshToken.Value)
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
	case COOKIE_NEED_REFRESH:
		err := s.Refreshing(ctx)
		if err != nil {
			return false, err
		}
	case COOKIE_NEED_AUTH:
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
	cookie, err := ctx.Request.Cookie("access_token")
	if err != nil {
		return COOKIE_NOT_FOUND
	}
	if cookie.MaxAge < int(time.Now().Unix()-time.Now().Add(-1*time.Duration(cookie.MaxAge)*time.Second).Unix()) {
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

func (s *ServerForum) GetAuthSignUpPage(ctx *gin.Context) {
	stat := s.CheckCookie(ctx)
	if stat == COOKIE_NEED_REFRESH {
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
	if stat == COOKIE_NEED_AUTH {
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
	if stat == COOKIE_NEED_REFRESH {
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
	passwd := hashPasswd(ctx.PostForm("password"))
	_, err := s.store.U.GetByEmail(email)
	if err != nil {
		refreshToken, err := generateRefrashToken()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
		}
		user := &entity.User{
			Name:         name,
			RefreshToken: refreshToken,
			Email:        email,
			Password:     passwd,
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

func (s *ServerForum) SignIn(ctx *gin.Context) {
	/*stat := s.CheckCookie(ctx)
	if stat != COOKIE_NOT_FOUND || stat != COOKIE_NEED_AUTH {
		f := SuccessHandler("Вы уже авторизированы")
		f(ctx)
		return
	}
	*/
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
