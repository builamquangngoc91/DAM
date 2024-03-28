package middlewares

import (
	"context"
	"dam/apis"
	"dam/enums"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Authentication(rdClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, apis.ErrorResponse{
				Message: "Authorization header is required",
				Code:    enums.AuthorizationHeaderRequiredError,
			})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, apis.ErrorResponse{
				Message: "Invalid authorization header",
				Code:    enums.InvalidAuthorizationHeaderError,
			})
			c.Abort()
			return
		}

		tokenStr := authParts[1]
		if _, err := rdClient.Get(ctx, tokenStr).Result(); err != nil {
			c.JSON(http.StatusUnauthorized, apis.ErrorResponse{
				Message: "Invalid JWT",
				Code:    enums.InvalidTokenError,
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, apis.ErrorResponse{
				Message: "Invalid JWT",
				Code:    enums.InvalidTokenError,
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, apis.ErrorResponse{
				Message: "Invalid JWT claims",
				Code:    enums.InvalidTokenError,
			})
			c.Abort()
			return
		}

		ctx = context.WithValue(ctx, enums.UserIDCtxKey, claims[string(enums.UserIDCtxKey)])

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
