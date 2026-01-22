package handler

import (
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/model"
  "project_ap3/internal/util"
)

func (h *Handler) AdminDashboard(c *gin.Context) {
  h.render(c, "Admin Dashboard", "admin_dashboard", nil)
}

func (h *Handler) AdminRoomsPage(c *gin.Context) {
  rooms, err := h.Rooms.ListAll(c.Request.Context())
  if err != nil {
    c.Redirect(http.StatusFound, "/admin?err=server_error")
    return
  }
  h.render(c, "Admin Rooms", "admin_rooms", gin.H{"Rooms": rooms})
}

func (h *Handler) AdminCreateRoom(c *gin.Context) {
  capacity, _ := strconv.Atoi(c.PostForm("capacity"))
  basePrice, _ := strconv.ParseInt(c.PostForm("base_price"), 10, 64)
  roomType := util.NormalizeRoomType(c.PostForm("type"))
  price, ok := util.ApplyRoomTypeMultiplier(basePrice, roomType)
  if !ok {
    c.Redirect(http.StatusFound, "/admin/rooms?err=invalid_input")
    return
  }

  room := &model.Room{
    RoomNo:       c.PostForm("room_no"),
    Name:         c.PostForm("name"),
    Type:         roomType,
    Capacity:     capacity,
    BasePrice:    basePrice,
    PricePerSlot: price,
    ImageURL:     c.PostForm("image_url"),
    Facilities:   c.PostForm("facilities"),
    Status:       c.PostForm("status"),
  }

  if room.RoomNo == "" || room.Name == "" || room.Type == "" || room.Capacity <= 0 || room.BasePrice < 0 || !util.ValidateRoomStatus(room.Status) || !util.ValidateRoomType(room.Type) {
    c.Redirect(http.StatusFound, "/admin/rooms?err=invalid_input")
    return
  }

  if _, err := h.Rooms.Create(c.Request.Context(), room); err != nil {
    c.Redirect(http.StatusFound, "/admin/rooms?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/rooms?msg=room_created")
}

func (h *Handler) AdminUpdateRoom(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/rooms?err=invalid_id")
    return
  }
  capacity, _ := strconv.Atoi(c.PostForm("capacity"))
  basePrice, _ := strconv.ParseInt(c.PostForm("base_price"), 10, 64)
  roomType := util.NormalizeRoomType(c.PostForm("type"))
  price, ok := util.ApplyRoomTypeMultiplier(basePrice, roomType)
  if !ok {
    c.Redirect(http.StatusFound, "/admin/rooms?err=invalid_input")
    return
  }

  room := &model.Room{
    ID:           id,
    RoomNo:       c.PostForm("room_no"),
    Name:         c.PostForm("name"),
    Type:         roomType,
    Capacity:     capacity,
    BasePrice:    basePrice,
    PricePerSlot: price,
    ImageURL:     c.PostForm("image_url"),
    Facilities:   c.PostForm("facilities"),
    Status:       c.PostForm("status"),
  }

  if room.RoomNo == "" || room.Name == "" || room.Type == "" || room.Capacity <= 0 || room.BasePrice < 0 || !util.ValidateRoomStatus(room.Status) || !util.ValidateRoomType(room.Type) {
    c.Redirect(http.StatusFound, "/admin/rooms?err=invalid_input")
    return
  }

  if err := h.Rooms.Update(c.Request.Context(), room); err != nil {
    c.Redirect(http.StatusFound, "/admin/rooms?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/rooms?msg=room_updated")
}

func (h *Handler) AdminDeleteRoom(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/rooms?err=invalid_id")
    return
  }
  if err := h.Rooms.Delete(c.Request.Context(), id); err != nil {
    c.Redirect(http.StatusFound, "/admin/rooms?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/rooms?msg=room_deleted")
}

func (h *Handler) AdminAvailabilityPage(c *gin.Context) {
  var roomIDPtr *int64
  if roomParam := c.Query("room_id"); roomParam != "" {
    if id, err := parseID(roomParam); err == nil {
      roomIDPtr = &id
    }
  }

  rows, err := h.Bookings.ListSchedule(c.Request.Context(), roomIDPtr)
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=server_error")
    return
  }
  rooms, _ := h.Rooms.ListAll(c.Request.Context())

  h.render(c, "Admin Availability", "admin_availability", gin.H{
    "Booked":       rows,
    "Rooms":        rooms,
    "FilterRoomID": roomIDPtr,
    "FilterRoomIDValue": func() int64 {
      if roomIDPtr == nil {
        return 0
      }
      return *roomIDPtr
    }(),
  })
}

func (h *Handler) AdminCreateAvailability(c *gin.Context) {
  roomID, err := parseID(c.PostForm("room_id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=invalid_room")
    return
  }
  date := c.PostForm("date")
  timeStart := c.PostForm("time_start")
  timeEnd := c.PostForm("time_end")
  isOpen := c.PostForm("is_open") == "on"

  if date == "" || !util.ValidateTimeRange(timeStart, timeEnd) {
    c.Redirect(http.StatusFound, "/admin/availability?err=invalid_input")
    return
  }

  a := &model.Availability{
    RoomID:    roomID,
    Date:      date,
    TimeStart: timeStart,
    TimeEnd:   timeEnd,
    IsOpen:    isOpen,
  }

  if _, err := h.Availability.Create(c.Request.Context(), a); err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/availability?msg=availability_created")
}

func (h *Handler) AdminUpdateAvailability(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=invalid_id")
    return
  }
  roomID, err := parseID(c.PostForm("room_id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=invalid_room")
    return
  }
  date := c.PostForm("date")
  timeStart := c.PostForm("time_start")
  timeEnd := c.PostForm("time_end")
  isOpen := c.PostForm("is_open") == "on"

  if date == "" || !util.ValidateTimeRange(timeStart, timeEnd) {
    c.Redirect(http.StatusFound, "/admin/availability?err=invalid_input")
    return
  }

  a := &model.Availability{
    ID:        id,
    RoomID:    roomID,
    Date:      date,
    TimeStart: timeStart,
    TimeEnd:   timeEnd,
    IsOpen:    isOpen,
  }

  if err := h.Availability.Update(c.Request.Context(), a); err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/availability?msg=availability_updated")
}

func (h *Handler) AdminDeleteAvailability(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=invalid_id")
    return
  }
  if err := h.Availability.Delete(c.Request.Context(), id); err != nil {
    c.Redirect(http.StatusFound, "/admin/availability?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/availability?msg=availability_deleted")
}

func (h *Handler) AdminBookingsPage(c *gin.Context) {
  status := c.Query("status")
  rows, err := h.Bookings.ListAdmin(c.Request.Context(), status)
  if err != nil {
    h.setDebugErrorHeader(c, err)
    h.render(c, "Admin Bookings", "admin_bookings", gin.H{
      "Bookings": []interface{}{},
      "Status":   status,
      "Error":    "server_error",
    })
    return
  }
  h.render(c, "Admin Bookings", "admin_bookings", gin.H{
    "Bookings": rows,
    "Status":   status,
  })
}

func (h *Handler) AdminApproveBooking(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/bookings?err=invalid_id")
    return
  }
  if err := h.Bookings.UpdateStatus(c.Request.Context(), id, "APPROVED"); err != nil {
    c.Redirect(http.StatusFound, "/admin/bookings?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/bookings?msg=approved")
}

func (h *Handler) AdminRejectBooking(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/bookings?err=invalid_id")
    return
  }
  if err := h.Bookings.UpdateStatus(c.Request.Context(), id, "REJECTED"); err != nil {
    c.Redirect(http.StatusFound, "/admin/bookings?err=server_error")
    return
  }
  c.Redirect(http.StatusFound, "/admin/bookings?msg=rejected")
}

func (h *Handler) AdminPaymentsPage(c *gin.Context) {
  status := c.Query("status")
  rows, err := h.Payments.ListAdmin(c.Request.Context(), status)
  if err != nil {
    h.setDebugErrorHeader(c, err)
    h.render(c, "Admin Payments", "admin_payments", gin.H{
      "Payments": []interface{}{},
      "Status":   status,
      "Error":    "server_error",
    })
    return
  }
  h.render(c, "Admin Payments", "admin_payments", gin.H{
    "Payments": rows,
    "Status":   status,
  })
}

func (h *Handler) AdminUsersPage(c *gin.Context) {
  users, err := h.Users.List(c.Request.Context())
  if err != nil {
    c.Redirect(http.StatusFound, "/admin/users?err=server_error")
    return
  }
  h.render(c, "Admin Users", "admin_users", gin.H{"Users": users})
}

