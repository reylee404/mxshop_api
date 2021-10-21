package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/response"
	"mxshop_api/user-web/models"
	"net/http"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "请登录"))
			c.Abort()
			return
		}

		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			switch err {
			case TokenExpiredErr:
				c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "授权过期"))
			default:
				c.JSON(http.StatusOK, response.NewFailedBaseResponse(400, "请登录"))
			}
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.Id)
		c.Next()
	}
}

type JWT struct {
	signKey []byte
}

func NewTestJWT()*JWT  {
	return &JWT{
		signKey: []byte("P8g9PjUx3wfIMRcyAU0HCfbvIbe4sSsb"),
	}
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTConfig.SignKey),
	}
}

func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.signKey)
}

var (
	TokenExpiredErr     = errors.New("Token is expired")
	TokenNotValidYetErr = errors.New("Token not active yet")
	TokenMalformedErr   = errors.New("That's not even a token")
	TokenInvalidErr     = errors.New("Couldn't handle this token:")
)

func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformedErr
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, TokenExpiredErr
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, TokenNotValidYetErr
			default:
				return nil, TokenInvalidErr
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalidErr
	} else {
		return nil, TokenInvalidErr
	}
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	//jwt.TimeFunc = func() time.Time {
	//	return time.Unix(0, 0)
	//}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalidErr
}
