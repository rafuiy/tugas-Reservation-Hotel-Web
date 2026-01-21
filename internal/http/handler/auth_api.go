package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/model"
  "project_ap3/internal/util"
)

type authRegisterRequest struct {
  Name     string `json:"name"`
  Email    string `json:"email"`
  Password string `json:"password"`
  Phone    string `json:"phone"`
}

type authLoginRequest struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}

func (h *Handler) RegisterAPI(c *gin.Context) {
  var req authRegisterRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if req.Name == "" || !util.ValidateEmail(req.Email) || !util.ValidatePassword(req.Password) {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  existing, err := h.Users.GetByEmail(c.Request.Context(), req.Email)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  if existing != nil {
    respondJSON(c, http.StatusConflict, nil, "email_taken")
    return
  }

  hash, err := util.HashPassword(req.Password)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }

  user := &model.User{
    Name:         req.Name,
    Email:        req.Email,
    PasswordHash: hash,
    Phone:        req.Phone,
    Role:         "USER",
  }
  id, err := h.Users.Create(c.Request.Context(), user)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }

  respondJSON(c, http.StatusCreated, gin.H{"id": id}, "")
}

func (h *Handler) LoginAPI(c *gin.Context) {
  var req authLoginRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_request")
    return
  }
  if !util.ValidateEmail(req.Email) || req.Password == "" {
    respondJSON(c, http.StatusBadRequest, nil, "invalid_input")
    return
  }

  user, err := h.Users.GetByEmail(c.Request.Context(), req.Email)
  if err != nil || user == nil || !util.CheckPassword(user.PasswordHash, req.Password) {
    respondJSON(c, http.StatusUnauthorized, nil, "invalid_credentials")
    return
  }

  token, err := util.GenerateToken(h.Cfg.JWTSecret, user.ID, user.Role, tokenTTL())
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "token_error")
    return
  }

  respondJSON(c, http.StatusOK, gin.H{"token": token}, "")
}

func (h *Handler) MeAPI(c *gin.Context) {
  userID, ok := getUserID(c)
  if !ok {
    respondJSON(c, http.StatusUnauthorized, nil, "unauthorized")
    return
  }

  user, err := h.Users.GetByID(c.Request.Context(), userID)
  if err != nil {
    respondJSON(c, http.StatusInternalServerError, nil, "server_error")
    return
  }
  if user == nil {
    respondJSON(c, http.StatusNotFound, nil, "not_found")
    return
  }
  user.PasswordHash = ""

  respondJSON(c, http.StatusOK, user, "")
}

