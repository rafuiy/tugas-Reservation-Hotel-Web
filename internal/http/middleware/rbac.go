package middleware

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func RequireRole(role string, api bool) gin.HandlerFunc {
  return func(c *gin.Context) {
    ctxRole, ok := c.Get(CtxRole)
    if !ok || ctxRole == nil {
      forbidden(c, api)
      return
    }
    if ctxRole.(string) != role {
      forbidden(c, api)
      return
    }
    c.Next()
  }
}

func forbidden(c *gin.Context, api bool) {
  if api {
    c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": "forbidden"})
    return
  }
  c.Redirect(http.StatusFound, "/login?err=forbidden")
  c.Abort()
}
