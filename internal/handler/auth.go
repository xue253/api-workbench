package handler

import (
	"net/http"
	"strconv"
	"time"

	"api-workbench/internal/model"
	"api-workbench/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("api-workbench-secret-key-2026")

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"`
}

type TokenResponse struct {
	Token string      `json:"token"`
	User  model.User  `json:"user"`
}

func hashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hashed, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd))
	return err == nil
}

func generateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// ---- Register ----
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: hashed,
		Email:    req.Email,
	}

	if err := repository.CreateUser(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	token, _ := generateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"data": TokenResponse{Token: token, User: user}})
}

// ---- Login ----
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{}
	err := repository.GetUserByUsername(req.Username, user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if !checkPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, _ := generateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"data": TokenResponse{Token: token, User: *user}})
}

// ---- Get Profile ----
func GetProfile(c *gin.Context) {
	uid, _ := c.Get("user_id")
	user := &model.User{}
	err := repository.GetUserByID(uid.(uint), user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ---- Update Profile ----
func UpdateProfile(c *gin.Context) {
	uid, _ := c.Get("user_id")
	var updates struct {
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &model.User{}
	err := repository.GetUserByID(uid.(uint), user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	user.Email = updates.Email
	user.Avatar = updates.Avatar
	repository.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ---- Change Password ----
func ChangePassword(c *gin.Context) {
	uid, _ := c.Get("user_id")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &model.User{}
	err := repository.GetUserByID(uid.(uint), user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if !checkPassword(user.Password, req.OldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "原密码错误"})
		return
	}
	hashed, _ := hashPassword(req.NewPassword)
	user.Password = hashed
	repository.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// ---- Delete Account ----
func DeleteAccount(c *gin.Context) {
	uid, _ := c.Get("user_id")
	if err := repository.DeleteUser(uid.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "账号已删除"})
}

// ---- List Users (admin) ----
func ListUsers(c *gin.Context) {
	var users []model.User
	if err := repository.GetUsers(&users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// ---- Project (需要 user_id) ----
func CreateProject(c *gin.Context) {
	uid, _ := c.Get("user_id")
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	p.UserID = uid.(uint)
	if err := repository.CreateProject(&p); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, p)
}

func ListProjects(c *gin.Context) {
	uid, _ := c.Get("user_id")
	var list []model.Project
	if err := repository.GetProjectsByUser(uid.(uint), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func UpdateProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	p.ID = uint(id)
	if err := repository.UpdateProject(&p); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, p)
}

func DeleteProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repository.DeleteProject(uint(id)); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}
