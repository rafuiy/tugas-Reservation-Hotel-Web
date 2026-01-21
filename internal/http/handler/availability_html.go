package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func (h *Handler) AvailabilityPage(c *gin.Context) {
  date := c.Query("date")
  var roomIDPtr *int64
  if roomParam := c.Query("room_id"); roomParam != "" {
    if id, err := parseID(roomParam); err == nil {
      roomIDPtr = &id
    }
  }

  rows, err := h.Availability.List(c.Request.Context(), date, roomIDPtr, true)
  if err != nil {
    c.Redirect(http.StatusFound, "/availability?err=server_error")
    return
  }
  rooms, _ := h.Rooms.ListActive(c.Request.Context())

  h.render(c, "Availability", "availability", gin.H{
    "Availability": rows,
    "Rooms":        rooms,
    "FilterDate":   date,
    "FilterRoomID": roomIDPtr,
  })
}

