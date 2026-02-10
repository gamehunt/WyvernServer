package handler

import "wyvern/server/internal/service"

type AuthHandler struct {
    service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
    return &AuthHandler{service: s}
}
