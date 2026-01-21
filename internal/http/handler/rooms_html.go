package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func (h *Handler) RoomsPage(c *gin.Context) {
  rooms, err := h.Rooms.ListActive(c.Request.Context())
  if err != nil {
    c.Redirect(http.StatusFound, "/rooms?err=server_error")
    return
  }
  h.render(c, "Rooms", "rooms", gin.H{"Rooms": rooms})
}

