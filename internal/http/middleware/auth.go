package middleware

import (
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/util"
)

const (
  CtxUserID = "user_id"
  CtxRole   = "role"
)

func RequireAuth(secret string, api bool) gin.HandlerFunc {
  return func(c *gin.Context) {
    tokenString := extractToken(c)
    if tokenString == "" {
      unauthorized(c, api)
      return
    }
    claims, err := util.ParseToken(secret, tokenString)
    if err != nil {
      unauthorized(c, api)
      return
    }
    c.Set(CtxUserID, claims.UserID)
    c.Set(CtxRole, claims.Role)
    c.Next()
  }
}

func TryAuth(secret string) gin.HandlerFunc {
  return func(c *gin.Context) {
    tokenString := extractToken(c)
    if tokenString == "" {
      c.Next()
      return
    }
    claims, err := util.ParseToken(secret, tokenString)
    if err == nil {
      c.Set(CtxUserID, claims.UserID)
      c.Set(CtxRole, claims.Role)
    }
    c.Next()
  }
}

func extractToken(c *gin.Context) string {
  if cookie, err := c.Cookie("auth_token"); err == nil && cookie != "" {
    return cookie
  }
  authHeader := c.GetHeader("Authorization")
  if strings.HasPrefix(authHeader, "Bearer ") {
    return strings.TrimPrefix(authHeader, "Bearer ")
  }
  return ""
}

func unauthorized(c *gin.Context, api bool) {
  if api {
    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
    return
  }
  c.Redirect(http.StatusFound, "/login?err=unauthorized")
  c.Abort()
}
