package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var BlackListTokens = []string{}

func TokenAuthMiddleware(tokenRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if tokenRequired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Missing Authorization header",
				})
				c.Abort()
				return
			} else {
				c.Set("isAuthenticated", false)
				c.Next()
				return
			}

		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			if tokenRequired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Invalid or missing Bearer token",
				})
				c.Abort()
				return
			} else {
				c.Set("isAuthenticated", false)
				c.Next()
				return
			}
		}

		tokenString := authHeaderParts[1]

		for _, blackToken := range BlackListTokens {
			if blackToken == tokenString {
				if tokenRequired {
					c.JSON(http.StatusUnauthorized, gin.H{
						"status":  http.StatusUnauthorized,
						"message": "Token blacklisted",
					})
					c.Abort()
					return
				} else {
					c.Set("isAuthenticated", false)
					c.Next()
					return
				}
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			if tokenRequired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": err.Error(),
				})
				c.Abort()
				return
			} else {
				c.Set("isAuthenticated", false)
				c.Next()
				return
			}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
			if time.Now().UTC().After(expirationTime) {
				if tokenRequired {
					c.JSON(http.StatusUnauthorized, gin.H{
						"status":  http.StatusUnauthorized,
						"message": "Token has expired",
					})
					c.Abort()
					return
				} else {
					c.Set("isAuthenticated", false)
					c.Next()
					return
				}
			}

			c.Set("token", tokenString)
			c.Set("isAuthenticated", true)
			c.Set("userId", claims["id"])
			c.Next()
		} else {
			if tokenRequired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Invalid token",
				})
				c.Abort()
				return
			} else {
				c.Set("isAuthenticated", false)
				c.Next()
				return
			}
		}
	}
}
