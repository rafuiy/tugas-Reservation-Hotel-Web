package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/http/middleware"
)

func (h *Handler) InvoiceAPI(c *gin.Context) {
  bookingID, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }

  inv, err := h.Invoice.GetByBookingID(c.Request.Context(), bookingID)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  if inv == nil {
    respondJSON(c, http.StatusNotFound, nil, "not_found")
    return
  }

  role, _ := c.Get(middleware.CtxRole)
  if role != "ADMIN" {
    userID, _ := getUserID(c)
    if inv.UserID != userID {
      respondJSON(c, http.StatusForbidden, nil, "forbidden")
      return
    }
  }

  respondJSON(c, http.StatusOK, inv, "")
}

