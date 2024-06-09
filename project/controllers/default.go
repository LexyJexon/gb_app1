package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"project/models"
)

type MainController struct {
	web.Controller
}

func (c *MainController) Get() {
	//c.Data["Website"] = "beego.vip"
	//c.Data["Email"] = "astaxie@gmail.com"
	flash := web.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		// Display settings successful
		c.TplName = "main.html"
	}
	c.TplName = "main.html"
}

func (c *MainController) Home() {
	flash := web.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		// Display settings successful
		c.TplName = "index.html"
	} else if _, ok = flash.Data["error"]; ok {
		// Display error messages
		c.TplName = "index.html"
	}
	c.TplName = "index.html"
}

func (c *MainController) Recipes() {
	flashErr := web.NewFlash()
	flash := web.ReadFromRequest(&c.Controller)
	recipes, err := models.GetAllItems()
	if err != nil {
		flashErr.Error("Can not get recipes")
		flashErr.Store(&c.Controller)
		c.Redirect("/recipes", 500)
		return
	}
	currentUserInterface := c.GetSession("current_user") //.(*models.Users)
	if currentUserInterface == nil {
		c.Redirect("/", 302)
		return
	}
	//currentUser := currentUserInterface.(*models.Users)
	if _, ok := flash.Data["notice"]; ok {
		c.Data["Recipes"] = recipes
		// Display settings successful
		c.TplName = "recipes.html"
	} else if _, ok := flash.Data["error"]; ok {
		c.Data["Recipes"] = recipes
		// Display settings successful
		c.TplName = "recipes.html"
	} else {
		c.Data["Recipes"] = recipes
		c.TplName = "recipes.html"
	}

}
