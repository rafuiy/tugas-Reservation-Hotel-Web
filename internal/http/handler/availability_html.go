package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/model"
)

func (h *Handler) AvailabilityPage(c *gin.Context) {
  startDate := c.Query("start_date")
  endDate := c.Query("end_date")
  var roomIDPtr *int64
  if roomParam := c.Query("room_id"); roomParam != "" {
    if id, err := parseID(roomParam); err == nil {
      roomIDPtr = &id
    }
  }

  rows := []model.BookingView{}
  if roomIDPtr != nil {
    booked, err := h.Bookings.ListSchedule(c.Request.Context(), roomIDPtr)
    if err != nil {
      c.Redirect(http.StatusFound, "/availability?err=server_error")
      return
    }
    rows = booked
  }
  rooms, _ := h.Rooms.ListActive(c.Request.Context())

  h.render(c, "Availability", "availability", gin.H{
    "Booked":       rows,
    "Rooms":        rooms,
    "StartDate":    startDate,
    "EndDate":      endDate,
    "FilterRoomID": roomIDPtr,
    "FilterRoomIDValue": func() int64 {
      if roomIDPtr == nil {
        return 0
      }
      return *roomIDPtr
    }(),
  })
}

