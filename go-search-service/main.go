package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Year        int    `json:"year"`
	IsAvailable bool   `json:"isAvailable"`
}

func main() {
	r := gin.Default()

	r.GET("/books/search", func(c *gin.Context) {
		title := c.Query("title")
		author := c.Query("author")

		resp, err := http.Get("http://dotnetcorelibrary:5096/books")
		if err != nil {
			println("Ошибка при запросе к dotnetcorelibrary:", err.Error())
			c.JSON(500, gin.H{"error": "Не удалось получить список книг из основного сервиса", "details": err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			println("dotnetcorelibrary не вернул 200, а:", resp.Status)
			c.JSON(500, gin.H{"error": "dotnetcorelibrary не вернул 200", "status": resp.Status})
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			println("Ошибка чтения тела:", err.Error())
			c.JSON(500, gin.H{"error": "Ошибка чтения тела ответа", "details": err.Error()})
			return
		}

		println("Ответ dotnetcorelibrary:", string(body)) // Вот тут мы увидим, что реально приходит!

		var books []Book
		if err := json.Unmarshal(body, &books); err != nil {
			println("Ошибка разбора JSON:", err.Error(), "Тело:", string(body))
			c.JSON(500, gin.H{"error": "Ошибка разбора JSON", "details": err.Error(), "raw": string(body)})
			return
		}

		var result []Book
		for _, b := range books {
			matchTitle := title == "" || containsIgnoreCase(b.Title, title)
			matchAuthor := author == "" || containsIgnoreCase(b.Author, author)
			if matchTitle && matchAuthor {
				result = append(result, b)
			}
		}
		c.JSON(200, result)
	})

	r.Run(":8080")
}

func containsIgnoreCase(str, substr string) bool {
	return len(substr) == 0 || (len(str) >= len(substr) &&
		strings.Contains(strings.ToLower(str), strings.ToLower(substr)))
}
