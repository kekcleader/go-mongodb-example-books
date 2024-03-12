package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Title  string `bson:"title"`
	Author string `bson:"author"`
	Year   int    `bson:"year"`
}

var (
	ctx             context.Context
	booksCollection *mongo.Collection
)

func main() {
	// Установка контекста для подключения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Получение коллекции
	booksCollection = client.Database("library").Collection("books")

	fmt.Println("Connected to MongoDB!")

	// Добавление книг
	AddBook(Book{"The Hobbit", "J.R.R. Tolkien", 1937})
	AddBook(Book{"1984", "George Orwell", 1949})

	// Получение и вывод всех книг
	GetAllBooks()

	// Обновление книги
	UpdateBook("The Hobbit", Book{"The Hobbit", "J.R.R. Tolkien", 1951})

	// Удаление книги
	DeleteBook("1984")

	// Повторное получение и вывод всех книг
	GetAllBooks()
}

// Добавление книги
func AddBook(book Book) {
	_, err := booksCollection.InsertOne(ctx, book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added a new book: ", book.Title)
}

// Получение списка всех книг
func GetAllBooks() {
	cursor, err := booksCollection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var books []Book
	if err = cursor.All(ctx, &books); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Books found:")
	for _, book := range books {
		fmt.Printf("%s by %s, %d\n", book.Title, book.Author, book.Year)
	}
}

// Обновление книги
func UpdateBook(title string, updatedBook Book) {
	filter := bson.D{{"title", title}}
	update := bson.D{{"$set", bson.D{{"author", updatedBook.Author}, {"year", updatedBook.Year}}}}
	_, err := booksCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Book updated: ", title)
}

// Удаление книги
func DeleteBook(title string) {
	_, err := booksCollection.DeleteOne(ctx, bson.D{{"title", title}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Book deleted: ", title)
}
