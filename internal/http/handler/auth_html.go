package handler

import (
  "net/http"

  "github.com/gin-gonic/gin"

  "project_ap3/internal/model"
  "project_ap3/internal/util"
)

func (h *Handler) ShowLogin(c *gin.Context) {
  h.render(c, "Login", "login", nil)
}

func (h *Handler) PostLogin(c *gin.Context) {
  email := c.PostForm("email")
  password := c.PostForm("password")

  if !util.ValidateEmail(email) || password == "" {
    c.Redirect(http.StatusFound, "/login?err=invalid_input")
    return
  }

  user, err := h.Users.GetByEmail(c.Request.Context(), email)
  if err != nil || user == nil || !util.CheckPassword(user.PasswordHash, password) {
    c.Redirect(http.StatusFound, "/login?err=invalid_credentials")
    return
  }

  token, err := util.GenerateToken(h.Cfg.JWTSecret, user.ID, user.Role, tokenTTL())
  if err != nil {
    c.Redirect(http.StatusFound, "/login?err=token_error")
    return
  }

  c.SetCookie("auth_token", token, int(tokenTTL().Seconds()), "/", "", false, true)
  if user.Role == "ADMIN" {
    c.Redirect(http.StatusFound, "/admin")
    return
  }
  c.Redirect(http.StatusFound, "/rooms")
}

func (h *Handler) ShowRegister(c *gin.Context) {
  h.render(c, "Register", "register", nil)
}

func (h *Handler) PostRegister(c *gin.Context) {
  name := c.PostForm("name")
  email := c.PostForm("email")
  password := c.PostForm("password")
  phone := c.PostForm("phone")

  if name == "" || !util.ValidateEmail(email) || !util.ValidatePassword(password) {
    c.Redirect(http.StatusFound, "/register?err=invalid_input")
    return
  }

  existing, err := h.Users.GetByEmail(c.Request.Context(), email)
  if err != nil {
    c.Redirect(http.StatusFound, "/register?err=server_error")
    return
  }
  if existing != nil {
    c.Redirect(http.StatusFound, "/register?err=email_taken")
    return
  }

  hash, err := util.HashPassword(password)
  if err != nil {
    c.Redirect(http.StatusFound, "/register?err=server_error")
    return
  }

  user := &model.User{
    Name:         name,
    Email:        email,
    PasswordHash: hash,
    Phone:        phone,
    Role:         "USER",
  }
  if _, err := h.Users.Create(c.Request.Context(), user); err != nil {
    c.Redirect(http.StatusFound, "/register?err=server_error")
    return
  }

  c.Redirect(http.StatusFound, "/login?msg=register_success")
}

func (h *Handler) Logout(c *gin.Context) {
  c.SetCookie("auth_token", "", -1, "/", "", false, true)
  c.Redirect(http.StatusFound, "/login?msg=logout_success")
}

