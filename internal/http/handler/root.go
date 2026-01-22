package handler

import (
  "github.com/gin-gonic/gin"
)

func (h *Handler) Root(c *gin.Context) {
  rooms, err := h.Rooms.ListActive(c.Request.Context())
  if err != nil {
    rooms = nil
  }
  h.render(c, "Home", "landing", gin.H{"Rooms": rooms})
}

