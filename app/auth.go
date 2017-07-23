package app

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type authSettings struct {
	Realm           string
	SigningAlorithm string
	Secret          []byte
	Timeout         time.Duration
	RefreshTimeout  time.Duration
}

var auth authSettings

func initAuth() {
	key := Env("SECRET", "")
	if key == "" {
		Log.Critical("secret missing from env")
		return
	}

	//TODO configurable refresh and timeout
	auth = authSettings{
		Realm:          "imgblog",
		Secret:         []byte(key),
		Timeout:        24 * time.Hour,
		RefreshTimeout: 24 * 30 * time.Hour,
	}
}

func CreateToken(user string, role string) string {
	claims := jwt.MapClaims{
		"user": user,
		"role": role,
		"exp":  time.Now().Add(auth.Timeout).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	result, _ := token.SignedString(auth.Secret)
	return result
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("WWW-Authenticate", "JWT realm="+auth.Realm)

		if authenticateUser(c) == false {
			c.Abort()
			return
		}
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	//TODO better error messages
	return func(c *gin.Context) {
		c.Header("WWW-Authenticate", "JWT realm="+auth.Realm)

		if authenticateUser(c) == false {
			Unauthorized(c)
			return
		}

		role, exists := c.Get("role")
		if !exists {
			Unauthorized(c)
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return
			}
		}

		Unauthorized(c)
	}
}

func authenticateUser(c *gin.Context) bool {
	token, err := parseHeader(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	user := claims["user"].(string)
	role := claims["role"].(string)

	//TODO handle refresh
	// iat := claims["iat"].(float64)
	// Log.Debug("found iat: " + strconv.FormatFloat(iat, 'f', 6, 64))

	if user == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user or role"})
		return false
	}

	c.Set("user", user)
	c.Set("role", role)
	return true
}

func parseHeader(header string) (*jwt.Token, error) {
	if header == "" {
		return nil, errors.New("empty authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("invalid authorization header")
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid singing algorithm")
		}

		return auth.Secret, nil
	})
}
