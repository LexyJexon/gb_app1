package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"project/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/users",
			beego.NSRouter("/", &controllers.UsersController{}, "get:ListUsers"),
			beego.NSRouter("/:id", &controllers.UsersController{}, "get:GetUser"),
			beego.NSRouter("/", &controllers.UsersController{}, "post:CreateUser"),
			beego.NSRouter("/:id", &controllers.UsersController{}, "put:UpdateUser"),
			beego.NSRouter("/:id", &controllers.UsersController{}, "delete:DeleteUser"),
		),
		beego.NSNamespace("/items",
			beego.NSRouter("/", &controllers.ItemsController{}, "get:ListItems"),
			beego.NSRouter("/:id", &controllers.ItemsController{}, "get:GetItemById"),
			beego.NSRouter("/", &controllers.ItemsController{}, "post:CreateItem"),
			beego.NSRouter("/:id", &controllers.ItemsController{}, "put:UpdateItem"),
			beego.NSRouter("/:id", &controllers.ItemsController{}, "delete:DeleteItem"),
		),
	)
	beego.AddNamespace(ns)
	// Роуты для пользователей
	beego.Router("/users", &controllers.UsersController{}, "get:ListUsers")
	beego.Router("/users", &controllers.UsersController{}, "post:CreateUser")
	beego.Router("/users/:id([0-9]+)", &controllers.UsersController{}, "get:GetUserById")
	beego.Router("/users/email/:email", &controllers.UsersController{}, "get:GetUserByEmail")
	beego.Router("/users/:id", &controllers.UsersController{}, "get:GetUser;put:UpdateUser;delete:DeleteUser")

	// Роуты для элементов
	beego.Router("/items", &controllers.ItemsController{}, "get:ListItems")
	beego.Router("/items", &controllers.ItemsController{}, "post:CreateItem")
	beego.Router("/items/:id([0-9]+)", &controllers.ItemsController{}, "get:GetItemById")
	beego.Router("/items/author/:authorId([0-9]+)", &controllers.ItemsController{}, "get:GetItemsByAuthorId")
	beego.Router("/items/:id", &controllers.ItemsController{}, "get:GetItemById;put:UpdateItem;delete:DeleteItem")

	// Отображение начальной страницы
	beego.Router("/", &controllers.MainController{})
	// Отображение главной страницы
	beego.Router("/index", &controllers.MainController{}, "get:Home")
	beego.Router("/recipes", &controllers.MainController{}, "get:Recipes")
	beego.Router("/create-recipe", &controllers.ItemsController{}, "get:ShowRecipeCreationForm")
	beego.Router("/create-recipe", &controllers.ItemsController{}, "post:CreateRecipe")
	beego.Router("/recipe/:id", &controllers.ItemsController{}, "get:ItemInfo")
	beego.Router("/search", &controllers.ItemsController{}, "get:Search")

	// Отображение страницы регистрации
	beego.Router("/register", &controllers.AuthController{}, "get:ShowRegisterForm")
	beego.Router("/register", &controllers.AuthController{}, "post:Register")
	// Отображение страницы входа
	beego.Router("/login", &controllers.AuthController{}, "get:ShowLoginForm")
	beego.Router("/login", &controllers.AuthController{}, "post:Login")
}
