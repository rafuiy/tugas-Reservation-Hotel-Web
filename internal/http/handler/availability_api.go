package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func (h *Handler) AvailabilityAPI(c *gin.Context) {
  date := c.Query("date")
  var roomIDPtr *int64
  if roomParam := c.Query("room_id"); roomParam != "" {
    if id, err := parseID(roomParam); err == nil {
      roomIDPtr = &id
    }
  }

  rows, err := h.Availability.List(c.Request.Context(), date, roomIDPtr, true)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rows, "")
}

