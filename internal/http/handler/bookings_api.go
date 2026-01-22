package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/service"
)

type bookingCreateRequest struct {
  RoomID     int64  `json:"room_id"`
  StartDate  string `json:"start_date"`
  EndDate    string `json:"end_date"`
  GuestName  string `json:"guest_name"`
  GuestPhone string `json:"guest_phone"`
  Notes      string `json:"notes"`
}

func (h *Handler) CreateBookingAPI(c *gin.Context) {
  userID, _ := getUserID(c)

  var req bookingCreateRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if req.RoomID == 0 || req.StartDate == "" || req.EndDate == "" || req.GuestName == "" || req.GuestPhone == "" {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  id, err := h.BookingService.Create(c.Request.Context(), service.BookingCreateInput{
    UserID:     userID,
    RoomID:     req.RoomID,
    StartDate:  req.StartDate,
    EndDate:    req.EndDate,
    GuestName:  req.GuestName,
    GuestPhone: req.GuestPhone,
    Notes:      req.Notes,
  })
  if err != nil {
    if handled := handleBookingServiceErrorAPI(c, err); handled {
      return
    }
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }

  respondJSON(c, http.StatusCreated, gin.H{"id": id}, "")
}

func (h *Handler) MyBookingsAPI(c *gin.Context) {
  userID, _ := getUserID(c)

  rows, err := h.Bookings.ListByUser(c.Request.Context(), userID)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rows, "")
}

