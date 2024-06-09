package models

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type Users struct {
	Id       int    `orm:"auto"`
	Email    string `orm:"unique"`
	Password string
}

type Items struct {
	Id          int `orm:"auto"`
	Image       string
	Title       string
	Description string
	Recipe      string
	//Ingredients       map[string]int `orm:"-"`
	CookTimeInMinutes int
	Author            *Users `orm:"rel(fk);on_delete(cascade)"`
}

type Ingredients struct {
	Id       int `orm:"auto"`
	Name     string
	Quantity int
	Units    string // Единицы измерения: шт., г., мл.
	Recipe   *Items `orm:"rel(fk);on_delete(cascade)"`
}

func init() {
	// Need to register model in init
	orm.RegisterModel(new(Users))
	orm.RegisterModel(new(Items))
	orm.RegisterModel(new(Ingredients))
}

func CreateUser(user *Users) (int64, error) {
	o := orm.NewOrm()
	userId, err := o.Insert(user)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func GetUserById(id int) (*Users, error) {
	o := orm.NewOrm()
	user := Users{Id: id}
	err := o.Read(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*Users, error) {
	o := orm.NewOrm()
	user := Users{Email: email}
	err := o.Read(&user, "Email")
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetItemById(id int) (*Items, error) {
	o := orm.NewOrm()
	item := Items{Id: id}
	err := o.Read(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func GetItemsByAuthorId(authorId int) ([]*Items, error) {
	o := orm.NewOrm()
	var items []*Items
	_, err := o.QueryTable("items").Filter("AuthorId", authorId).All(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetAllItems возвращает все элементы.
func GetAllItems() ([]*Items, error) {
	o := orm.NewOrm()
	var items []*Items
	_, err := o.QueryTable("items").All(&items)
	return items, err
}

// GetAllUsers возвращает все users.
func GetAllUsers() ([]*Users, error) {
	o := orm.NewOrm()
	var users []*Users
	_, err := o.QueryTable("users").All(&users)
	return users, err
}

// CreateItem создает новый элемент.
func CreateItem(item *Items) error {
	o := orm.NewOrm()
	_, err := o.Insert(item)
	return err
}

// UpdateItem обновляет существующий элемент.
func UpdateItem(id string, item *Items) error {
	o := orm.NewOrm()
	itemId, _ := strconv.Atoi(id)
	existingItem := &Items{Id: itemId}
	if o.Read(existingItem) == nil {
		item.Id = existingItem.Id
		_, err := o.Update(item)
		return err
	}
	return nil
}

// UpdateUser обновляет существующий элемент.
func UpdateUser(id string, user *Users) error {
	o := orm.NewOrm()
	userId, _ := strconv.Atoi(id)
	existingUser := &Users{Id: userId}
	if o.Read(existingUser) == nil {
		user.Id = existingUser.Id
		_, err := o.Update(user)
		return err
	}
	return nil
}

// DeleteItem удаляет элемент по его ID.
func DeleteItem(id string) error {
	o := orm.NewOrm()
	itemId, _ := strconv.Atoi(id)
	item := &Items{Id: itemId}
	_, err := o.Delete(item)
	return err
}

// DeleteUser удаляет элемент по его ID.
func DeleteUser(id string) error {
	o := orm.NewOrm()
	userId, _ := strconv.Atoi(id)
	user := &Users{Id: userId}
	_, err := o.Delete(user)
	return err
}

// VerifyPassword проверяет, соответствует ли указанный пароль хэшу пароля пользователя
func (u *Users) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsEmailTaken проверяет, занят ли указанный адрес электронной почты в базе данных
func IsEmailTaken(email string) (bool, error) {
	o := orm.NewOrm()
	user := Users{Email: email}
	err := o.Read(&user, "Email")
	if errors.Is(err, orm.ErrNoRows) {
		// Адрес электронной почты не найден, поэтому он не занят
		return false, nil
	} else if err == nil {
		// Адрес электронной почты найден, он занят
		return true, nil
	} else {
		// Произошла ошибка при выполнении запроса к базе данных
		return false, err
	}
}

func CreateIngredient(ingredient *Ingredients) (int64, error) {
	o := orm.NewOrm()
	ingredientId, err := o.Insert(ingredient)
	if err != nil {
		return 0, err
	}
	return ingredientId, nil
}

// UpdateIngredient Метод для обновления информации об ингредиенте
func UpdateIngredient(ingredient *Ingredients) error {
	o := orm.NewOrm()
	_, err := o.Update(ingredient)
	return err
}

// DeleteIngredientById Метод для удаления ингредиента по его ID
func DeleteIngredientById(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Ingredients{Id: id})
	return err
}

// GetIngredientById Метод для получения информации об ингредиенте по его ID
func GetIngredientById(id int) (*Ingredients, error) {
	o := orm.NewOrm()
	ingredient := &Ingredients{Id: id}
	err := o.Read(ingredient)
	return ingredient, err
}

// GetAllIngredientsByRecipeId Метод для получения всех ингредиентов по ID рецепта
func GetAllIngredientsByRecipeId(recipeId int) ([]*Ingredients, error) {
	o := orm.NewOrm()
	var ingredients []*Ingredients
	_, err := o.QueryTable("ingredients").Filter("recipe_id", recipeId).All(&ingredients)
	return ingredients, err
}

// SearchRecipes Метод для поиска рецептов по запросу
func SearchRecipes(query string) ([]*Items, error) {
	var recipes []*Items

	// Используем ORM Beego для выполнения запроса к базе данных
	o := orm.NewOrm()
	qs := o.QueryTable("items").SetCond(orm.NewCondition().Or("Title__icontains", query).Or("Description__icontains", query).Or("Recipe__icontains", query))
	_, err := qs.All(&recipes)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func GetLastItem() (*Items, error) {
	var item Items
	err := orm.NewOrm().QueryTable("items").OrderBy("-Id").Limit(1).One(&item)
	if err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			// В таблице нет строк, возвращаем nil
			return &Items{}, nil
		}
		// В случае другой ошибки возвращаем её
		return nil, err
	}

	return &item, nil
}
