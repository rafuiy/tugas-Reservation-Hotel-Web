package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func (h *Handler) RoomsAPI(c *gin.Context) {
  rooms, err := h.Rooms.ListActive(c.Request.Context())
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  respondJSON(c, http.StatusOK, rooms, "")
}

func (h *Handler) RoomByIDAPI(c *gin.Context) {
  id, err := parseID(c.Param("id"))
  if err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_id")
    return
  }

  room, err := h.Rooms.GetByID(c.Request.Context(), id)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  if room == nil {
    respondJSON(c, http.StatusNotFound, nil, "not_found")
    return
  }
  respondJSON(c, http.StatusOK, room, "")
}

