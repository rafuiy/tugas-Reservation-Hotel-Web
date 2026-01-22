package handler

import (
  "database/sql"
  "net/http"
  "time"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/model"
  "project_ap3/internal/util"
)

type adminRoomRequest struct {
  RoomNo       string `json:"room_no"`
  Name         string `json:"name"`
  Type         string `json:"type"`
  Capacity     int    `json:"capacity"`
  BasePrice    int64  `json:"base_price"`
  PricePerSlot int64  `json:"price_per_slot"`
  ImageURL     string `json:"image_url"`
  Facilities   string `json:"facilities"`
  Status       string `json:"status"`
}

type adminAvailabilityRequest struct {
  RoomID    int64  `json:"room_id"`
  Date      string `json:"date"`
  TimeStart string `json:"time_start"`
  TimeEnd   string `json:"time_end"`
  IsOpen    bool   `json:"is_open"`
}

type adminBookingUpdateRequest struct {
  Status string `json:"status"`
}

type adminPaymentUpdateRequest struct {
  Status string `json:"status"`
}

func (h *Handler) AdminRoomsAPI(c *gin.Context) {
  rooms, err := h.Rooms.ListAll(c.Request.Context())
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rooms, "")
}

func (h *Handler) AdminCreateRoomAPI(c *gin.Context) {
  var req adminRoomRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  basePrice := req.BasePrice
  if basePrice == 0 && req.PricePerSlot > 0 {
    basePrice = req.PricePerSlot
  }
  roomType := util.NormalizeRoomType(req.Type)
  price, ok := util.ApplyRoomTypeMultiplier(basePrice, roomType)
  if req.RoomNo == "" || req.Name == "" || roomType == "" || req.Capacity <= 0 || basePrice < 0 || !util.ValidateRoomStatus(req.Status) || !util.ValidateRoomType(roomType) || !ok {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  id, err := h.Rooms.Create(c.Request.Context(), &model.Room{
    RoomNo:       req.RoomNo,
    Name:         req.Name,
    Type:         roomType,
    Capacity:     req.Capacity,
    BasePrice:    basePrice,
    PricePerSlot: price,
    ImageURL:     req.ImageURL,
    Facilities:   req.Facilities,
    Status:       req.Status,
  })
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusCreated, gin.H{"id": id}, "")
}

func (h *Handler) AdminUpdateRoomAPI(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }

  var req adminRoomRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  basePrice := req.BasePrice
  if basePrice == 0 && req.PricePerSlot > 0 {
    basePrice = req.PricePerSlot
  }
  roomType := util.NormalizeRoomType(req.Type)
  price, ok := util.ApplyRoomTypeMultiplier(basePrice, roomType)
  if req.RoomNo == "" || req.Name == "" || roomType == "" || req.Capacity <= 0 || basePrice < 0 || !util.ValidateRoomStatus(req.Status) || !util.ValidateRoomType(roomType) || !ok {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  if err := h.Rooms.Update(c.Request.Context(), &model.Room{
    ID:           id,
    RoomNo:       req.RoomNo,
    Name:         req.Name,
    Type:         roomType,
    Capacity:     req.Capacity,
    BasePrice:    basePrice,
    PricePerSlot: price,
    ImageURL:     req.ImageURL,
    Facilities:   req.Facilities,
    Status:       req.Status,
  }); err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, gin.H{"updated": true}, "")
}

func (h *Handler) AdminDeleteRoomAPI(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }
  if err := h.Rooms.Delete(c.Request.Context(), id); err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, gin.H{"deleted": true}, "")
}

func (h *Handler) AdminAvailabilityAPI(c *gin.Context) {
  var roomIDPtr *int64
  if roomParam := c.Query("room_id"); roomParam != "" {
    if id, err := parseID(roomParam); err == nil {
      roomIDPtr = &id
    }
  }

  rows, err := h.Bookings.ListSchedule(c.Request.Context(), roomIDPtr)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rows, "")
}

func (h *Handler) AdminCreateAvailabilityAPI(c *gin.Context) {
  var req adminAvailabilityRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if req.RoomID == 0 || req.Date == "" || !util.ValidateTimeRange(req.TimeStart, req.TimeEnd) {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  id, err := h.Availability.Create(c.Request.Context(), &model.Availability{
    RoomID:    req.RoomID,
    Date:      req.Date,
    TimeStart: req.TimeStart,
    TimeEnd:   req.TimeEnd,
    IsOpen:    req.IsOpen,
  })
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusCreated, gin.H{"id": id}, "")
}

func (h *Handler) AdminUpdateAvailabilityAPI(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }

  var req adminAvailabilityRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if req.RoomID == 0 || req.Date == "" || !util.ValidateTimeRange(req.TimeStart, req.TimeEnd) {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  if err := h.Availability.Update(c.Request.Context(), &model.Availability{
    ID:        id,
    RoomID:    req.RoomID,
    Date:      req.Date,
    TimeStart: req.TimeStart,
    TimeEnd:   req.TimeEnd,
    IsOpen:    req.IsOpen,
  }); err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, gin.H{"updated": true}, "")
}

func (h *Handler) AdminDeleteAvailabilityAPI(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }
  if err := h.Availability.Delete(c.Request.Context(), id); err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, gin.H{"deleted": true}, "")
}

func (h *Handler) AdminBookingsAPI(c *gin.Context) {
  status := c.Query("status")
  rows, err := h.Bookings.ListAdmin(c.Request.Context(), status)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rows, "")
}

func (h *Handler) AdminUpdateBookingAPI(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }

  var req adminBookingUpdateRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if req.Status != "APPROVED" && req.Status != "REJECTED" {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_status")
    return
  }

  if err := h.Bookings.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
    if err == sql.ErrNoRows {
      respondJSON(c, http.StatusNotFound, nil, "not_found")
      return
    }
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, gin.H{"status": req.Status}, "")
}

func (h *Handler) AdminPaymentsAPI(c *gin.Context) {
  status := c.Query("status")
  rows, err := h.Payments.ListAdmin(c.Request.Context(), status)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rows, "")
}

func (h *Handler) AdminUpdatePaymentAPI(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }

  var req adminPaymentUpdateRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if !util.ValidatePaymentStatus(req.Status) {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_status")
    return
  }

  var paidAt *time.Time
  if req.Status == "PAID" {
    now := time.Now()
    paidAt = &now
  }

  if err := h.Payments.UpdateStatus(c.Request.Context(), id, req.Status, paidAt); err != nil {
    if err == sql.ErrNoRows {
      respondJSON(c, http.StatusNotFound, nil, "not_found")
      return
    }
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, gin.H{"status": req.Status}, "")
}

func (h *Handler) AdminUsersAPI(c *gin.Context) {
  users, err := h.Users.List(c.Request.Context())
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, users, "")
}

