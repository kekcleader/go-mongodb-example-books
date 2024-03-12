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

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Замените "your_username:your_password" на ваши фактические учетные данные
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://your_username:your_password@localhost:27017/library"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Connected to MongoDB!")

	booksCollection := client.Database("library").Collection("books")

	AddBook(ctx, booksCollection, Book{"The Hobbit", "J.R.R. Tolkien", 1937})
	AddBook(ctx, booksCollection, Book{"1984", "George Orwell", 1949})

	GetAllBooks(ctx, booksCollection)

	UpdateBook(ctx, booksCollection, "The Hobbit", Book{"The Hobbit", "J.R.R. Tolkien", 1951})

	DeleteBook(ctx, booksCollection, "1984")

	GetAllBooks(ctx, booksCollection)
}

func AddBook(ctx context.Context, collection *mongo.Collection, book Book) {
	_, err := collection.InsertOne(ctx, book)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added a new book: ", book.Title)
}

func GetAllBooks(ctx context.Context, collection *mongo.Collection) {
	cursor, err := collection.Find(ctx, bson.D{{}})
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

func UpdateBook(ctx context.Context, collection *mongo.Collection, title string, updatedBook Book) {
	filter := bson.D{{"title", title}}
	update := bson.D{{"$set", bson.D{{"author", updatedBook.Author}, {"year", updatedBook.Year}}}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Book updated: ", title)
}

func DeleteBook(ctx context.Context, collection *mongo.Collection, title string) {
	_, err := collection.DeleteOne(ctx, bson.D{{"title", title}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Book deleted: ", title)
}
