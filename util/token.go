package util

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateToken(claims jwt.MapClaims, key string, tmout time.Duration) (string, error) {
	claims["times"] = time.Now()
	if tmout > 0 {
		claims["timeout"] = time.Now().Add(tmout)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokens, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokens, nil
}
func SetToken(c *gin.Context, p jwt.MapClaims, key string, rem bool, doman ...string) (string, error) {
	tmout := time.Hour * 5
	if rem {
		tmout = time.Hour * 24 * 5
	}
	tokens, err := CreateToken(p, key, tmout)
	if err != nil {
		return "", err
	}
	cke := http.Cookie{
		Name:     "gokinstk",
		Value:    tokens,
		HttpOnly: true,
	}
	if len(doman) > 0 {
		cke.Domain = doman[0]
	}

	cke.MaxAge = 60 * 60 * 5
	if rem {
		cke.MaxAge = 60 * 60 * 24 * 5
	}
	c.Writer.Header().Add("Set-Cookie", cke.String())
	return tokens, nil
}

func ClearToken(c *gin.Context, doman ...string) error {
	cke := http.Cookie{
		Name:     "gokinstk",
		HttpOnly: true,
	}
	if len(doman) > 0 {
		cke.Domain = doman[0]
	}
	cke.MaxAge = -1
	c.Writer.Header().Set("Set-Cookie", cke.String())
	return nil
}

func getToken(c *gin.Context) string {
	tkc, err := c.Request.Cookie("gokinstk")
	if err != nil {
		return ""
	}
	return tkc.Value
}
func getTokenAuth(c *gin.Context) string {
	ats := c.GetHeader("Authorization")
	if ats == "" {
		return ""
	}
	aths, err := url.PathUnescape(ats)
	if err != nil {
		return ""
	}
	aths = strings.Replace(aths, "TOKEN ", "", 1)
	return aths
}
func GetTokens(s string, key string) jwt.MapClaims {
	if s == "" {
		return nil
	}
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err == nil {
		claim, ok := token.Claims.(jwt.MapClaims)
		if ok {
			return claim
		}
	}
	return nil
}
func GetToken(c *gin.Context, key string) jwt.MapClaims {
	tk := getTokenAuth(c)
	if tk == "" {
		tk = getToken(c)
	}
	if tk == "" {
		tk = c.Query("authToken")
	}
	return GetTokens(tk, key)
}
