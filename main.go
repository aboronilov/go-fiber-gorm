package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Piblisher string `json:"publisher"`
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

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalln("Error connecting to DB...", err)
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8000")
}
