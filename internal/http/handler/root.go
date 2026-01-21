package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/http/middleware"
)

func (h *Handler) Root(c *gin.Context) {
  roleVal, ok := c.Get(middleware.CtxRole)
  if !ok || roleVal == nil {
    c.Redirect(http.StatusFound, "/login")
    return
  }
  if roleVal.(string) == "ADMIN" {
    c.Redirect(http.StatusFound, "/admin")
    return
  }
  c.Redirect(http.StatusFound, "/rooms")
}

