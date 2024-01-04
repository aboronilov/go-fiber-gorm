package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aboronilov/go-fiber-gorm/models"
	"github.com/aboronilov/go-fiber-gorm/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(ctx *fiber.Ctx) error {
	book := Book{}

	err := ctx.BodyParser(&book)
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"Message": "Couldn't parse request body"},
		)
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"Message": "Error creating item"},
		)
		return err
	}

	ctx.Status(http.StatusCreated).JSON(
		&fiber.Map{
			"Message": "Item created",
		},
	)
	return nil
}

func (r *Repository) ListBooks(ctx *fiber.Ctx) error {
	books := &[]models.Books{}

	err := r.DB.Find(books).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"Message": "Error finding items"},
		)
		return err
	}

	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"Message": "Items fetched",
			"Data":    books,
		},
	)
	return nil
}

func (r *Repository) DeleteBook(ctx *fiber.Ctx) error {
	bookModel := &models.Books{}

	id := ctx.Params("id")
	if id == "" {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "No id is provided"},
		)
		return nil
	}

	err := r.DB.Delete(bookModel, id).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Unable to delete item"},
		)
		return err
	}

	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Item successfully deleted"},
	)
	return nil
}

// func (r *Repository) GetBookByID(ctx *fiber.Ctx) error {
// 	id := ctx.Params("id")
// 	if id == "" {
// 		ctx.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "No id is provided"},
// 		)
// 		return nil
// 	}
// 	fmt.Println("Get book by id, the id is - ", id)

// 	bookModel := models.Books{}
// 	err := r.DB.First(&bookModel, id).Error
// 	if err != nil {
// 		ctx.Status(http.StatusNotFound).JSON(
// 			&fiber.Map{"message": "item not found"},
// 		)
// 		return err
// 	}

// 	ctx.Status(http.StatusOK).JSON(
// 		&fiber.Map{
// 			"Message": "Item succesfully fetched",
// 			"Data":    bookModel,
// 		},
// 	)
// 	return nil
// }

func (r *Repository) GetBookByID(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/books")
	api.Post("/", r.CreateBook)
	api.Delete("/{id}", r.DeleteBook)
	api.Get("/", r.ListBooks)
	api.Get("/{id}", r.GetBookByID)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error parsing env", err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalln("Error connecting to DB...", err)
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatalln("Error while migrationg...", err)
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8000")
}
