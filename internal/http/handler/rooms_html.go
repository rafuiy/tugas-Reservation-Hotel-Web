package handler

import (
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"
)

func (h *Handler) RoomsPage(c *gin.Context) {
  rooms, err := h.Rooms.ListAll(c.Request.Context())
  if err != nil {
    c.Redirect(http.StatusFound, "/rooms?err=server_error")
    return
  }
  h.render(c, "Rooms", "rooms", gin.H{"Rooms": rooms})
}

func (h *Handler) RoomDetailPage(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/rooms?err=invalid_room")
    return
  }

  room, err := h.Rooms.GetByID(c.Request.Context(), id)
  if err != nil {
    c.Redirect(http.StatusFound, "/rooms?err=server_error")
    return
  }
  if room == nil || room.Status != "ACTIVE" {
    c.Redirect(http.StatusFound, "/rooms?err=not_found")
    return
  }

  startDate := c.Query("start_date")
  endDate := c.Query("end_date")
  roomIDPtr := &id
  booked, err := h.Bookings.ListSchedule(c.Request.Context(), roomIDPtr)
  if err != nil {
    c.Redirect(http.StatusFound, "/rooms?err=server_error")
    return
  }

  h.render(c, "Room Detail", "room_detail", gin.H{
    "Room":       room,
    "Booked":     booked,
    "StartDate":  startDate,
    "EndDate":    endDate,
    "Facilities": splitFacilities(room.Facilities),
  })
}

func splitFacilities(raw string) []string {
  raw = strings.TrimSpace(raw)
  if raw == "" {
    return nil
  }
  parts := strings.FieldsFunc(raw, func(r rune) bool {
    return r == ',' || r == '\n' || r == ';'
  })
  out := make([]string, 0, len(parts))
  for _, part := range parts {
    item := strings.TrimSpace(part)
    if item == "" {
      continue
    }
    out = append(out, item)
  }
  return out
}

