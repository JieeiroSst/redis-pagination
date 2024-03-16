package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/go-redis/redis/v8"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	key := "books"
	ctx := context.Background()

	books := []Book{
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
		{ID: "32", Title: "Pride and Prejudice", Author: "Jane Austen"},
	}
	jsonBooks := make([]string, len(books))
	for i, book := range books {
		jsonBytes, err := json.Marshal(book)
		if err != nil {
			fmt.Println("Error marshaling book:", err)
			continue
		}
		jsonBooks[i] = string(jsonBytes)
	}
	err := rdb.LPush(ctx, key, jsonBooks).Err()
	if err != nil {
		fmt.Println("Error storing books in Redis:", err)
		return
	}

	limit := 50
	page := 2
	totalElements, err := rdb.LLen(ctx, key).Result()
	if err != nil {
		fmt.Println("Error getting total elements:", err)
		return
	}
	fmt.Println("======= totalElements: ", totalElements)
	fmt.Println("======= totalPage: ", math.Ceil(float64(totalElements)/float64(limit)))
	if int64(page) > (totalElements/int64(limit))+1 {
		fmt.Println("Requested page exceeds total number of pages")
		return
	}
	startIndex := (page - 1) * limit
	endIndex := min(int64(startIndex)+int64(limit)-1, totalElements-1)
	results, err := rdb.LRange(ctx, key, int64(startIndex), int64(endIndex)).Result()
	if err != nil {
		fmt.Println("Error fetching books from Redis:", err)
		return
	}
	var booksOnPage []Book
	for _, result := range results {
		var book Book
		err := json.Unmarshal([]byte(result), &book)
		if err != nil {
			fmt.Println("Error unmarshaling book data:", err)
			continue
		}
		booksOnPage = append(booksOnPage, book)
	}
	fmt.Println("Books on page", page, ":")
	for _, book := range booksOnPage {
		fmt.Printf("  - ID: %s, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
