package main

import (
	"archive/zip"
	"bytes"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	// API trả về file ZIP
	app.Get("/download", func(ctx iris.Context) {
		// Tạo nội dung file trong zip
		var buffer bytes.Buffer
		zipWriter := zip.NewWriter(&buffer)

		files := []struct {
			Name    string
			Content string
		}{
			{"file1.txt", "Hello from file1!"},
			{"file2.txt", "Hello from file2!"},
		}

		for _, file := range files {
			writer, err := zipWriter.Create(file.Name)
			if err != nil {
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString("Error creating ZIP file")
				return
			}
			writer.Write([]byte(file.Content))
		}

		zipWriter.Close()

		// Thiết lập header cho response
		ctx.Header("Content-Type", "application/zip")
		ctx.Header("Content-Disposition", "attachment; filename=example.zip")
		ctx.Write(buffer.Bytes())
	})

	app.Get("/api/status", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"status": "Server is running",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	app.Post("/api/echo", func(ctx iris.Context) {
		var input map[string]interface{}
		if err := ctx.ReadJSON(&input); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Invalid JSON"})
			return
		}
		ctx.JSON(iris.Map{"received": input})
	})

	app.Post("/api/upload", func(ctx iris.Context) {
		file, _, err := ctx.FormFile("file")
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Failed to upload file"})
			return
		}
		defer file.Close()
		ctx.JSON(iris.Map{"message": "File uploaded successfully"})
	})

	app.Get("/api/download/{filename}", func(ctx iris.Context) {
		filename := ctx.Params().Get("filename")
		filePath := filepath.Join("data", filename)

		// Kiểm tra nếu file không tồn tại
		if ctx.SendFile(filePath, filename) != nil {
			ctx.StatusCode(http.StatusNotFound)
			ctx.JSON(iris.Map{"error": "File not found"})
		}
	})

	app.Listen(":8080")
}
