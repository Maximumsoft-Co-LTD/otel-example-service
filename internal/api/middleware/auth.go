package middleware

import (
	"context"
	"example-otelgo-miniservice/internal/application/service"
	"example-otelgo-miniservice/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	CtxClaimsKey = "claims"
	CtxUserIDKey = "user_id"
)

func AuthMiddleware(jwe *service.JWE) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := utils.BearerTokenFrom(c.GetHeader("Authorization"))
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "0",
				"message": "missing bearer token",
			})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		claims, err := jwe.DecryptToken(ctx, raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "0",
				"message": "invalid token",
			})
			c.Abort()
			return
		}

		if claims.ID != "" {
			if revoked, err := jwe.IsBlacklisted(ctx, claims.ID); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    "0",
					"message": "blacklist check error",
				})
				c.Abort()
				return
			} else if revoked {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    "0",
					"message": "token is revoked",
				})
				c.Abort()
				return
			}
		}

		var userIDHex string
		if claims.Subject != "" {
			if oid, err := bson.ObjectIDFromHex(claims.Subject); err == nil {
				userIDHex = oid.Hex()
			} else {
				userIDHex = claims.Subject
			}
		}

		c.Set(CtxClaimsKey, claims)    // *domain.TokenPayload
		c.Set(CtxUserIDKey, userIDHex) // string
		c.Next()
	}
}
