package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/http/middleware"
)

func (h *Handler) InvoicePage(c *gin.Context) {
  bookingID, err := parseID(c.Param("booking_id"))
  if err != nil {
    c.Redirect(http.StatusFound, "/bookings/my?err=invalid_id")
    return
  }

  inv, err := h.Invoice.GetByBookingID(c.Request.Context(), bookingID)
  if err != nil {
    c.Redirect(http.StatusFound, "/bookings/my?err=server_error")
    return
  }
  if inv == nil {
    c.Redirect(http.StatusFound, "/bookings/my?err=not_found")
    return
  }

  role, _ := c.Get(middleware.CtxRole)
  if role != "ADMIN" {
    userID, _ := getUserID(c)
    if inv.UserID != userID {
      c.Redirect(http.StatusFound, "/bookings/my?err=forbidden")
      return
    }
  }

  h.render(c, "Invoice", "invoice", gin.H{"Invoice": inv})
}

