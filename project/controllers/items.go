package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	_ "google.golang.org/api/option"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"project/models"
	"strconv"
)

type ItemsController struct {
	web.Controller
}

// Item представляет модель данных элемента
type Item struct {
	ID                int
	Image             string
	Title             string
	Description       string
	Recipe            string
	Ingredients       map[string]int
	CookTimeInMinutes int
	AuthorID          int
}

func (c *ItemsController) ShowRecipeCreationForm() {
	flash := web.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		// Display settings successful
		c.TplName = "recipes.html"
	} else if _, ok = flash.Data["error"]; ok {
		// Display error messages
		c.TplName = "create_recipe.html"
	} else {
		c.TplName = "create_recipe.html"
	}
}

func (c *ItemsController) GetItemById() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	item, err := models.GetItemById(id)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = item
	c.ServeJSON()
}

func (c *ItemsController) GetItemsByAuthorId() {
	authorId, err := strconv.Atoi(c.Ctx.Input.Param(":authorId"))
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	items, err := models.GetItemsByAuthorId(authorId)
	if err != nil {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	c.Data["json"] = items
	c.ServeJSON()
}

// / GetAllItems получает все элементы.
func (c *ItemsController) ListItems() {
	items, err := models.GetAllItems()
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = items
	}
	c.ServeJSON()
}

// CreateItem создает новый элемент.
func (c *ItemsController) CreateItem() {
	var newItem models.Items
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &newItem)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		err = models.CreateItem(&newItem)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = "Items created successfully"
		}
	}
	c.ServeJSON()
}

// UpdateItem обновляет элемент по ID.
func (c *ItemsController) UpdateItem() {
	id := c.Ctx.Input.Param(":id")
	var updatedItem models.Items
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &updatedItem)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		err = models.UpdateItem(id, &updatedItem)
		if err != nil {
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = "Items updated successfully"
		}
	}
	c.ServeJSON()
}

// DeleteItem удаляет элемент по ID.
func (c *ItemsController) DeleteItem() {
	id := c.Ctx.Input.Param(":id")
	err := models.DeleteItem(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = "Items deleted successfully"
	}
	c.ServeJSON()
}

func (c *ItemsController) ShowRecipesWithCreatingItem() {
	go c.CreateRecipe()
	flash := web.NewFlash()
	flash.Notice("Рецепт создается. Пожалуйста, подождите.")
	flash.Store(&c.Controller)
	c.Redirect("/recipes", 302)
}

// CreateRecipe Метод для обработки POST запроса создания рецепта
func (c *ItemsController) CreateRecipe() {
	flash := web.NewFlash()
	var newRecipeID int
	// Получение последнего рецепта из базы данных
	lastRecipe, err := models.GetLastItem()
	if err != nil {
		log.Println("Error getting last recipe:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}
	// Увеличение ID последнего рецепта на 1
	if lastRecipe == nil {
		newRecipeID = 1
	} else {
		newRecipeID = lastRecipe.Id + 1
	}

	// Получение файла из запроса
	file, _, err := c.Ctx.Request.FormFile("image")
	if err != nil {
		log.Println("Error getting image from request:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}

	// Создание имени файла на основе увеличенного ID
	filename := fmt.Sprintf("%d.webp", newRecipeID)

	maxWidth := 300  // Максимальная ширина
	maxHeight := 300 // Максимальная высота

	compressedImage, err := compressImage(file, maxWidth, maxHeight)
	if err != nil {
		log.Println("Error compressing image:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}

	// Сохранение изображения в формате WebP на сервере
	imagePath := filepath.Join("static", "images", filename)
	err = saveWebPImage(compressedImage, imagePath)
	if err != nil {
		log.Println("Error saving image:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}
	defer file.Close()

	/*// Инициализация Firebase
	opt := option.WithCredentialsFile("google-services.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println("Error initializing Firebase app:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 500)
		return
	}

	// Получение ссылки на Bucket
	client, err := app.Storage(ctx)
	if err != nil {
		log.Println("Error initializing Storage client:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 500)
		return
	}
	bucket, err := client.Bucket("koocbook-82b63.appspot.com")
	if err != nil {
		log.Println("Error getting default bucket:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 500)
		return
	}

	// Загрузка изображения в Firebase Storage
	filename := filepath.Base(header.Filename)*/

	/*obj := bucket.Object("images/" + filename)
	wc := obj.NewWriter(ctx)

	if _, err := io.Copy(wc, compressedImage); err != nil {
		log.Println("Error uploading image to Firebase Storage:", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 500)
		return
	}

	defer wc.Close()

	time.Sleep(300 * time.Second)

	attrs, err := getAttrsWithRetry(ctx, obj)

	// Получение имени бакета
	bucketName := attrs.Bucket

	// Получение пути к объекту
	storagePath := attrs.Name

	token := "7e004a6d-9b79-4d9f-9b92-b6ec46f67253"

	imageURL := "https://firebasestorage.googleapis.com/v0/b/" + bucketName + "/o/" + url.PathEscape(storagePath) + "?alt=media&token=" + token

	log.Println(imageURL)*/
	title := c.GetString("title")
	description := c.GetString("description")
	recipe := c.GetString("recipe")
	cookTime, _ := c.GetInt("cookTime")
	ingredientNames := c.GetStrings("ingredient_names[]")
	ingredientQuantities := c.GetStrings("ingredient_quantities[]")
	ingredientUnits := c.GetStrings("ingredient_units[]")

	/*	ctx := context.Background()

		// Загрузка изображения в Firebase Storage и получение ссылки
		imageURL, err := uploadImageToFirebaseStorage(ctx, fileName)
		if err != nil {
			log.Printf("Error uploading image to Firebase Storage: %v", err)
			c.Redirect("/recipes", 302)
			return
		}*/

	// Получение текущего пользователя из сессии
	currentUser := c.GetSession("current_user").(*models.Users)

	// Создание нового рецепта
	item := &models.Items{
		Title:             title,
		Description:       description,
		Recipe:            recipe,
		CookTimeInMinutes: cookTime,
		Image:             imagePath, // Ссылка на изображение в Firebase Storage
		Author:            currentUser,
	}

	// Сохранение рецепта в базе данных
	if err := models.CreateItem(item); err != nil {
		log.Printf("Error creating recipe: %v", err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}

	// Создание ингредиентов и привязка их к рецепту
	for i, name := range ingredientNames {
		quantity, _ := strconv.Atoi(ingredientQuantities[i])
		unit := ingredientUnits[i]
		ingredient := &models.Ingredients{
			Name:     name,
			Quantity: quantity,
			Units:    unit,
			Recipe:   item,
		}
		if _, err := models.CreateIngredient(ingredient); err != nil {
			log.Printf("Error creating ingredient: %v", err)
			flash.Error(err.Error())
			flash.Store(&c.Controller)
			c.Redirect("/recipes", 302)
			return
		}
	}

	// Редирект на страницу с созданным рецептом или другую нужную вам страницу
	flash.Notice("Рецепт создан")
	flash.Store(&c.Controller)
	c.Redirect("/recipes", 302)
}

// Функция для сохранения изображения в формате WebP
func saveWebPImage(file multipart.File, imagePath string) error {
	// Декодируем изображение в формат image.Image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Создаем новый файл для сохранения изображения в формате WebP
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Кодируем и сохраняем изображение в формате WebP
	if err := webp.Encode(f, img, nil); err != nil {
		return err
	}

	return nil
}

/*func getAttrsWithRetry(ctx context.Context, obj *storage.ObjectHandle) (*storage.ObjectAttrs, error) {
	var err error
	for i := 0; i < 5; i++ {
		attrs, err := obj.Attrs(ctx)
		if err == nil {
			return attrs, nil
		}
		log.Printf("Error getting object attributes (attempt %d): %v", i+1, err)
		time.Sleep(time.Minute)
	}
	return nil, err
}*/

// Функция для сжатия изображения
func compressImage(file multipart.File, maxWidth, maxHeight int) (multipart.File, error) {
	// Открываем изображение
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Определяем размеры изображения
	size := img.Bounds().Size()

	// Вычисляем новые размеры с учетом максимальной ширины и высоты
	newWidth := size.X
	newHeight := size.Y
	if size.X > maxWidth {
		newWidth = maxWidth
		newHeight = (newWidth * size.Y) / size.X
	}
	if newHeight > maxHeight {
		newHeight = maxHeight
		newWidth = (newHeight * size.X) / size.Y
	}

	// Сжимаем изображение
	resizedImg := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	// Создаем буфер памяти для сохранения сжатого изображения
	buf := new(bytes.Buffer)

	// Кодируем сжатое изображение в буфер
	if err := png.Encode(buf, resizedImg); err != nil {
		return nil, err
	}

	// Создаем временный файл для сохранения сжатого изображения
	tempFile, err := os.CreateTemp("", "compressed_image_*.png")
	if err != nil {
		return nil, err
	}

	// Записываем буфер сжатого изображения во временный файл
	if _, err := io.Copy(tempFile, buf); err != nil {
		return nil, err
	}

	// Перемещаем указатель файла в начало
	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// Возвращаем временный файл как multipart.File
	return tempFile, nil
}

/*
	func uploadImageToFirebaseStorage(ctx context.Context, filePath string) (string, error) {
		config := &firebase.Config{
			StorageBucket: "your-storage-bucket-url",
		}

		app, err := firebase.NewApp(ctx, config, option.WithCredentialsFile("google-services.json"))
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
			return "", err
		}

		client, err := app.Storage(ctx)
		if err != nil {
			log.Fatalf("error initializing storage client: %v\n", err)
			return "", err
		}

		bucket, err := client.DefaultBucket()
		if err != nil {
			log.Fatalf("error getting bucket: %v\n", err)
			return "", err
		}

		// Создаем новый объект в Firebase Storage
		storagePath := "images/" + filepath.Base(filePath)
		wc := bucket.Object(storagePath).NewWriter(ctx)
		defer wc.Close()

		// Загружаем файл в Firebase Storage
		if _, err := wc.Write([]byte("")); err != nil {
			log.Fatalf("error writing to object: %v\n", err)
			return "", err
		}

		bucketName := "koocbook-82b63.appspot.com"

		// Получаем публичную ссылку на загруженное изображение
		imageURL := "https://firebasestorage.googleapis.com/" + bucketName + "/" + storagePath
		return imageURL, nil
	}
*/

func (c *ItemsController) ItemInfo() {
	flash := web.NewFlash()
	itemId, err := c.GetInt(":id")
	if err != nil {
		c.Abort("404")
		return
	}

	item, err := models.GetItemById(itemId)
	if err != nil {
		log.Println(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}

	ingredients, err := models.GetAllIngredientsByRecipeId(itemId)
	if err != nil {
		log.Println(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 302)
		return
	}

	author, err := models.GetUserById(item.Author.Id)
	if err != nil {
		log.Println(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		c.Redirect("/recipes", 500)
		return
	}

	c.Data["Recipe"] = item
	c.Data["Ingredients"] = ingredients
	c.Data["Author"] = author

	c.TplName = "recipe_by_id.html"
}

// Search Метод для обработки запроса поиска рецептов
func (c *ItemsController) Search() {
	query := c.GetString("query") // Получаем запрос из формы поиска
	// Выполняем поиск рецептов в базе данных по запросу
	recipes, err := models.SearchRecipes(query)
	if err != nil {
		// Обработка ошибки
		return
	}
	// Отображаем найденные рецепты в соответствующем шаблоне
	c.Data["Recipes"] = recipes
	c.TplName = "search_results.html"
}
