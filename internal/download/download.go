package download

import (
	"fmt"
	"is/pkg/transport"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Entity struct {
	Brand   string
	Product string
}

type Request struct {
	Url        string
	ImagePath  string
	EntityType string
	Shop       string
}

func (d Request) Download() string {
	transportClient, _ := transport.NewClient(time.Minute * 10)
	fileName := filepath.Base(strings.Split(d.Url, "?")[0])

	i := d.EntityType + "/" + strings.ToLower(d.Shop) + "/" + fileName
	imgPath := "/image/" + i

	path := fmt.Sprintf("%s/%s/%s/%s", d.ImagePath, d.EntityType, d.Shop, fileName)

	_, err := os.Stat(path)
	if err == nil {
		fmt.Println("Файл существует")
		return i
	}

	err = transportClient.DownloadFile(path, d.Url)
	if err != nil {
		return ""
	}
	return imgPath
}
