package main

import (
	"errors"
	"fmt"
	"io"
	config "is/internal/domain"
	"is/internal/handlers"
	conf "is/pkg/config"
	"is/pkg/middleware"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type justFilesFilesystem struct {
	fs               http.FileSystem
	readDirBatchSize int
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredStatFile{File: f, readDirBatchSize: fs.readDirBatchSize}, nil
}

type neuteredStatFile struct {
	http.File
	readDirBatchSize int
}

func (e neuteredStatFile) Stat() (os.FileInfo, error) {
	s, err := e.File.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
	LOOP:
		for {
			fl, err := e.File.Readdir(e.readDirBatchSize)
			switch err {
			case io.EOF:
				break LOOP
			case nil:
				for _, f := range fl {
					if f.Name() == "index.html" {
						return s, err
					}
				}
			default:
				return nil, err
			}
		}
		return nil, os.ErrNotExist
	}
	return s, err
}

func main() {
	envConfig, err := conf.Config[config.EnvConfig](".env")
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	r := gin.New()

	r.MaxMultipartMemory = 8 << 20

	r.StaticFS("/image", http.Dir("public"))

	api := r.Group("/api/v3")
	{
		api.POST("/image",
			middleware.CORSMiddleware(),
			middleware.BasicAuth(),
			handlers.AddImage(envConfig.ImagePath))
	}

	r.SetTrustedProxies(nil)

	fmt.Println("Server started on http://localhost:3333")

	if err := r.Run(":3333"); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed gracefully")
		} else {
			fmt.Printf("error starting server: %v\n", err)
			os.Exit(1)
		}
	}
}
