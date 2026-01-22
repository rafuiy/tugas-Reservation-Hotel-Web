package handler

import (
  "database/sql"
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/service"
)

func (h *Handler) MyBookingsPage(c *gin.Context) {
  userID, _ := getUserID(c)
  rows, err := h.Bookings.ListByUser(c.Request.Context(), userID)
  if err != nil {
    c.Redirect(http.StatusFound, "/bookings/my?err=server_error")
    return
  }
  h.render(c, "My Bookings", "my_bookings", gin.H{"Bookings": rows})
}

func (h *Handler) CreateBooking(c *gin.Context) {
  userID, _ := getUserID(c)
  roomID, err := parseID(c.PostForm("room_id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/availability?err=invalid_room")
    return
  }
  startDate := c.PostForm("start_date")
  endDate := c.PostForm("end_date")
  guestName := c.PostForm("guest_name")
  guestPhone := c.PostForm("guest_phone")
  notes := c.PostForm("notes")

  if guestName == "" || guestPhone == "" || startDate == "" || endDate == "" {
    c.Redirect(http.StatusFound, "/availability?err=invalid_input")
    return
  }

  _, err = h.BookingService.Create(c.Request.Context(), service.BookingCreateInput{
    UserID:     userID,
    RoomID:     roomID,
    StartDate:  startDate,
    EndDate:    endDate,
    GuestName:  guestName,
    GuestPhone: guestPhone,
    Notes:      notes,
  })
  if err != nil {
    switch err {
    case service.ErrConflict:
      c.Redirect(http.StatusFound, "/availability?err=conflict")
      return
    case service.ErrNotFound:
      c.Redirect(http.StatusFound, "/availability?err=room_not_found")
      return
    case service.ErrInvalid:
      c.Redirect(http.StatusFound, "/availability?err=invalid_range")
      return
    default:
      c.Redirect(http.StatusFound, "/availability?err=server_error")
      return
    }
  }

  c.Redirect(http.StatusFound, "/bookings/my?msg=booking_created")
}

func handleBookingServiceErrorAPI(c *gin.Context, err error) bool {
  switch err {
  case service.ErrConflict:
    respondJSON(c, http.StatusConflict, nil, "conflict")
  case service.ErrNotFound:
    respondJSON(c, http.StatusNotFound, nil, "not_found")
  case service.ErrInvalid:
    respondJSON(c, http.StatusBadRequest, nil, "invalid")
  case sql.ErrNoRows:
    respondJSON(c, http.StatusNotFound, nil, "not_found")
  default:
    return false
  }
  return true
}

