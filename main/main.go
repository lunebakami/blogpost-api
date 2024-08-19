package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/lunebakami/futtodos-api/config"
	"github.com/lunebakami/futtodos-api/handlers"
	"github.com/lunebakami/futtodos-api/storage"
)


func main() {
  cfg := config.LoadConfig()

  postStorage, err := storage.NewPostStorage(cfg)
  if err != nil {
    log.Fatalf("Failed to initialize storage: %v", err)
  }

  defer postStorage.Close()

  postHandler := handlers.NewPostHandler(postStorage)

  e := echo.New()

  e.GET("/", func (c echo.Context) error {
    return c.String(200, "Hello, World!")
  })
  e.POST("/posts", postHandler.CreatePost)
  e.GET("/posts", postHandler.GetAllPost)
  e.GET("/posts/:id", postHandler.GetPost)
  e.PUT("/posts/:id", postHandler.UpdatePost)
  e.DELETE("/posts/:id", postHandler.DeletePost)

  e.Logger.Fatal(e.Start(":8080"))
}
