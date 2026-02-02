package main

import (
	"fmt"
	"io"
	"net/http"
)

type RealInspector struct{}

type Inspector interface {
	Inspect(url string)
}

func (r RealInspector) Inspect(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	defer resp.Body.Close()
	size, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("URL: %s | Размер: %d", url, len(size))
}

func main() {
	url := "https://example.com"
	r := RealInspector{}
	r.Inspect(url)
}
