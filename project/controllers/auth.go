package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
	"project/models"
	"regexp"
	"strings"
)

type AuthController struct {
	web.Controller
}

func (c *AuthController) ShowRegisterForm() {
	flash := web.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		// Display settings successful
		c.TplName = "register.html"
	} else if _, ok = flash.Data["error"]; ok {
		// Display error messages
		c.TplName = "register.html"
	} else {
		c.TplName = "register.html"
	}
}

func (c *AuthController) ShowLoginForm() {
	flash := web.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		// Display settings successful
		c.TplName = "index.html"
	} else if _, ok = flash.Data["error"]; ok {
		// Display error messages
		c.TplName = "login.html"
	} else {
		c.TplName = "login.html"
	}

}

// Register Метод для обработки запроса на регистрацию нового пользователя
func (c *AuthController) Register() {
	flash := web.NewFlash()
	// Получение данных из запроса
	email := strings.TrimSpace(c.GetString("email"))
	password := strings.TrimSpace(c.GetString("password"))

	// Валидация email
	if !isValidEmail(email) {
		flash.Error("Invalid email address")
		flash.Store(&c.Controller)
		c.Redirect("/register", 302)
	}

	// Валидация пароля
	if len(password) < 8 || !isStrongPassword(password) {
		flash.Error("Password must be at least 8 characters long and contain uppercase letters, digits, and special characters")
		flash.Store(&c.Controller)
		c.Redirect("/register", 302)
	}

	// Проверка уникальности email
	IsEmailTaken, err := models.IsEmailTaken(email)
	if IsEmailTaken {
		flash.Error("Email address is already taken")
		flash.Store(&c.Controller)
		c.Redirect("/register", 302)
	}

	// Создание нового пользователя
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.Users{Email: email, Password: string(passwordHash)}
	_, err = models.CreateUser(&user)
	if err != nil {
		flash.Error("Failed to register user")
		flash.Store(&c.Controller)
		c.Redirect("/register", 302)
		return
	}

	// В случае успешной регистрации
	flash.Notice("Registered successfully")
	flash.Store(&c.Controller)
	c.Redirect("/", 302)
}

// Login Метод для обработки запроса на аутентификацию пользователя
func (c *AuthController) Login() {
	flash := web.NewFlash()
	// Получение данных из запроса
	email := strings.TrimSpace(c.GetString("email"))
	password := strings.TrimSpace(c.GetString("password"))

	// Проверка наличия пользователя в базе данных
	user, err := models.GetUserByEmail(email)
	if err != nil {
		flash.Error("User not found")
		flash.Store(&c.Controller)
		c.Redirect("/login", 302)
	}

	// Проверка правильности введенного пароля
	if !user.VerifyPassword(password) {
		flash.Error("Incorrect password" + password + " " + user.Password)
		flash.Store(&c.Controller)
		c.Redirect("/login", 302)
		return
	}

	// В случае успешной аутентификации
	flash.Notice("User authenticated successfully")
	err = c.SetSession("current_user", user)
	if err != nil {
		flash.Error("Error saving user to the session")
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}
	flash.Store(&c.Controller)
	c.Redirect("/recipes", 302)
}

// Валидация формата email
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Проверка сложности пароля
func isStrongPassword(password string) bool {
	// Проверяем, содержит ли пароль хотя бы одну заглавную букву, одну цифру и один специальный символ
	hasUpperCase := false
	hasDigit := false
	hasSpecialChar := false

	for _, char := range password {
		if 'A' <= char && char <= 'Z' {
			hasUpperCase = true
		}
		if '0' <= char && char <= '9' {
			hasDigit = true
		}
		if strings.Contains("!@#$%^&*()_+-=[]{}|;:,.<>?/", string(char)) {
			hasSpecialChar = true
		}
	}

	return hasUpperCase && hasDigit && hasSpecialChar
}
