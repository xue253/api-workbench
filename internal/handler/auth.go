package handler

import (
	"net/http"
	"strconv"
	"time"

	"api-workbench/internal/config"
	"api-workbench/internal/model"
	"api-workbench/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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
	Token string     `json:"token"`
	User  model.User `json:"user"`
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
	expireHour := config.AppConfig.JWT.ExpireHour
	if expireHour <= 0 {
		expireHour = 168
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expireHour) * time.Hour).Unix(),
	})
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func getUintParam(c *gin.Context, key string) (uint, bool) {
	id, err := strconv.Atoi(c.Param(key))
	if err != nil || id <= 0 {
		errorResp(c, 400, "无效的ID")
		return 0, false
	}
	return uint(id), true
}

func getCurrentUserID(c *gin.Context) (uint, bool) {
	uid, exists := c.Get("user_id")
	if !exists {
		errorResp(c, 401, "未登录")
		return 0, false
	}
	userID, ok := uid.(uint)
	if !ok {
		errorResp(c, 500, "用户信息异常")
		return 0, false
	}
	return userID, true
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

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}
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

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": TokenResponse{Token: token, User: *user}})
}

// ---- Get Profile ----
func GetProfile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	user := &model.User{}
	err := repository.GetUserByID(userID, user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// ---- Update Profile ----
func UpdateProfile(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	var updates struct {
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &model.User{}
	err := repository.GetUserByID(userID, user)
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
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &model.User{}
	err := repository.GetUserByID(userID, user)
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
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	if err := repository.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "账号已删除"})
}

// ---- Project (需要 user_id) ----
func CreateProject(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	p.UserID = userID
	if err := repository.CreateProject(&p); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, p)
}

func ListProjects(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	var list []model.Project
	if err := repository.GetProjectsByUser(userID, &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func UpdateProject(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}
	var existing model.Project
	if err := repository.GetProjectByID(id, &existing); err != nil {
		errorResp(c, 404, "项目不存在")
		return
	}
	if existing.UserID != userID {
		errorResp(c, 403, "无权操作此项目")
		return
	}
	var p model.Project
	if err := c.ShouldBindJSON(&p); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	p.ID = id
	if err := repository.UpdateProject(&p); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, p)
}

func DeleteProject(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}
	var existing model.Project
	if err := repository.GetProjectByID(id, &existing); err != nil {
		errorResp(c, 404, "项目不存在")
		return
	}
	if existing.UserID != userID {
		errorResp(c, 403, "无权操作此项目")
		return
	}
	if err := repository.DeleteProject(id); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}
