package handler

import (
	"net/http"
	"time"

	"github.com/bsyrlhabibi/airdrop/internal/model"
	"github.com/bsyrlhabibi/airdrop/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo  *repository.UserRepo
	JWTSecret string
}

func NewAuthHandler(ur *repository.UserRepo, secret string) *AuthHandler {
	return &AuthHandler{UserRepo: ur, JWTSecret: secret}
}

type RegisterRequest struct {
	Email    string `json:"email" example:"user@email.com"`
	Password string `json:"password" example:"secret123"`
	Name     string `json:"name" example:"Bita"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@email.com"`
	Password string `json:"password" example:"secret123"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  model.User  `json:"user"`
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body RegisterRequest true "Register data"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &model.User{Email: req.Email, Password: string(hash), Name: req.Name}

	if err := h.UserRepo.Create(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registered",
		"user":    gin.H{"id": user.ID, "email": user.Email, "name": user.Name},
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and get JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body LoginRequest true "Login credentials"
// @Success      200  {object}  AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserRepo.FindByEmail(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(h.JWTSecret))

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
		"user":  gin.H{"id": user.ID, "email": user.Email, "name": user.Name},
	})
}
