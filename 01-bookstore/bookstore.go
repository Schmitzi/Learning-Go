package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const libraryFile = "library.json"

type Book struct {
	ID       string
	Title    string
	Author   string
	Price    string
	ImageURL string
}

type Library struct {
	Books []Book
}

func loadLibrary(filename string) Library {
	data, err := os.ReadFile(filename)
	if err != nil {
		// File doesn't exist yet, return empty library
		return Library{}
	}
	var library Library
	json.Unmarshal(data, &library)
	return library
}

func saveLibrary(filename string, library *Library) {
	data, err := json.MarshalIndent(library, "", "  ")
	if err != nil {
		log.Fatal("could not encode library: ", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Fatal("could not write library file: ", err)
	}
}

func help() {
	fmt.Println("--- Bookstore Help ---")
	fmt.Println("get --all: Return all books in library")
	fmt.Println("get --id X: Return book X from library")
	fmt.Println("add X: Add book to library")
	fmt.Println("update: Update book in library")
	fmt.Println("delete: Delete book from library")

}

func format_check() bool {
	// Expected: add --id X --title X --author X --price X --imageurl X
	expected := []struct {
		index int
		flag  string
	}{
		{2, "--id"},
		{4, "--title"},
		{6, "--author"},
		{8, "--price"},
		{10, "--imageurl"},
	}

	if len(os.Args) < 12 {
		log.Println("missing parameters, try again")
		fmt.Println("Usage: add --id X --title X --author X --price X --imageurl X")
		return false
	}

	for _, e := range expected {
		if os.Args[e.index] != e.flag {
			log.Printf("expected %s at position %d, got %s\n", e.flag, e.index, os.Args[e.index])
			return false
		}
	}

	return true
}

func remove_book(slice []Book, index int) []Book {
	return append(slice[:index], slice[index+1:]...)
}

func add_cmd(library *Library) {
	if !format_check() {
		return
	}
	for i := 0; i < len(library.Books); i++ {
		if os.Args[3] == library.Books[i].ID {
			log.Println("ID already in use")
			return
		}
	}
	newBook := Book{
		ID:       os.Args[3],
		Title:    os.Args[5],
		Author:   os.Args[7],
		Price:    os.Args[9],
		ImageURL: os.Args[11],
	}

	library.Books = append(library.Books, newBook)
	saveLibrary(libraryFile, library)
}

func get_cmd(library *Library) {
	switch os.Args[2] {
	case "--all":
		fmt.Println("ID\tTitle\t\tAuthor\tPrice\tImageURL")
		for i := 0; i < len(library.Books); i++ {
			b := library.Books[i]
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", b.ID, b.Title, b.Author, b.Price, b.ImageURL)
		}
	case "--id":
		fmt.Println("ID\tTitle\t\tAuthor\tPrice\tImageURL")
		for i := 0; i < len(library.Books); i++ {
			if os.Args[3] == library.Books[i].ID {
				b := library.Books[i]
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n", b.ID, b.Title, b.Author, b.Price, b.ImageURL)
				return
			}
		}
		log.Println("Book not found")
	default:
		log.Println("malformed parameters")
	}
}

func update_cmd(library *Library) {
	if !format_check() {
		return
	}
	for i := 0; i < len(library.Books); i++ {
		if os.Args[3] == library.Books[i].ID {
			library.Books[i] = Book{
				ID:       os.Args[3],
				Title:    os.Args[5],
				Author:   os.Args[7],
				Price:    os.Args[9],
				ImageURL: os.Args[11],
			}
			saveLibrary(libraryFile, library)
			return
		}
	}
	log.Println("Book not found")
}

func delete_cmd(library *Library) {
	if os.Args[2] != "--id" {
		log.Println("malformed parameters")
		return
	}
	for i := 0; i < len(library.Books); i++ {
		if os.Args[3] == library.Books[i].ID {
			library.Books = remove_book(library.Books, i)
			saveLibrary(libraryFile, library)
			return
		}
	}
	log.Println("book not found")
}

func main() {
	log.SetPrefix("Error: ")
	// Initial error handling
	if len(os.Args) < 3 {
		if len(os.Args) == 2 && os.Args[1] == "help" {
			help()
		} else {
			log.Println("too few args. Try ./bookstore help")
		}
		os.Exit(0)
	}

	library := loadLibrary(libraryFile)

	switch os.Args[1] {
	case "add":
		add_cmd(&library)
	case "get":
		get_cmd(&library)
	case "update":
		update_cmd(&library)
	case "delete":
		delete_cmd(&library)
	default:
		log.Println("invalid command")
	}
}
