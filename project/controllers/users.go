package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
	"project/models"
	_ "project/models"
	"strconv"
)

// UsersController определяет контроллер для пользователей
type UsersController struct {
	web.Controller
}

// User представляет модель данных пользователя
type User struct {
	ID       int
	Email    string
	Password string
}

func (c *UsersController) GetUserById() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	user, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = user
	c.ServeJSON()
}

func (c *UsersController) GetUserByEmail() {
	email := c.Ctx.Input.Param(":email")

	user, err := models.GetUserByEmail(email)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = user
	c.ServeJSON()
}

func (c *UsersController) Register() {
	var user models.Users
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// Валидация данных пользователя (например, проверка наличия email и пароля)
	if user.Email == "" || user.Password == "" {
		c.Data["json"] = "Email and Password are required fields"
		c.ServeJSON()
		return
	}

	// Хеширование пароля перед сохранением в базе данных
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	user.Password = string(hashedPassword)

	// Сохранение пользователя в базе данных
	userId, err := models.CreateUser(&user)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]int{"user_id": int(userId)}
	c.ServeJSON()
}

func (c *UsersController) Login() {
	var user models.Users
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// Поиск пользователя по email
	existingUser, err := models.GetUserByEmail(user.Email)
	if err != nil {
		c.Data["json"] = "Invalid email or password"
		c.ServeJSON()
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		c.Data["json"] = "Invalid email or password"
		c.ServeJSON()
		return
	}

	// Вход успешен
	// Создайте здесь токен аутентификации и верните его в ответе

	c.Data["json"] = "Login successful"
	c.ServeJSON()
}

// Метод для получения всех пользователей
func (c *UsersController) ListUsers() {
	var users []*models.Users
	users, err := models.GetAllUsers()
	if err != nil {
		c.Data["json"] = err
	} else {
		c.Data["json"] = users
	}
	c.ServeJSON()
}

// Метод для получения пользователя по ID
func (c *UsersController) GetUser() {
	id, _ := strconv.Atoi(c.GetString(":id"))
	user, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = err
	} else {
		c.Data["json"] = user
	}
	c.ServeJSON()
}

// Метод для создания нового пользователя
func (c *UsersController) CreateUser() {
	var user models.Users
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		_, err = models.CreateUser(&user)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = "Users created successfully"
		}
	}
	c.ServeJSON()
}

// Метод для обновления пользователя
func (c *UsersController) UpdateUser() {
	id := c.GetString(":id")
	var user models.Users
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		err := models.UpdateUser(id, &user)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = "Users updated successfully"
		}
	}
	c.ServeJSON()
}

// Метод для удаления пользователя
func (c *UsersController) DeleteUser() {
	id := c.GetString(":id")
	err := models.DeleteUser(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "Users deleted successfully"
	}
	c.ServeJSON()
}
