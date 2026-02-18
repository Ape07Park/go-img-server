package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth - API 키 인증 미들웨어
//
// Spring 비유:
//   OncePerRequestFilter 또는 HandlerInterceptor 와 동일
//   헤더에 X-API-Key 값이 없거나 틀리면 401 반환 후 요청 차단
//
// 사용법 (main.go):
//   api := r.Group("/api/v1")
//   api.Use(middleware.APIKeyAuth(cfg.APIKey))
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "X-API-Key 헤더가 필요합니다",
			})
			return
		}
		if key != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "유효하지 않은 API 키입니다",
			})
			return
		}
		// 인증 통과 → 다음 핸들러 실행 (Spring의 chain.doFilter()와 동일)
		c.Next()
	}
}
