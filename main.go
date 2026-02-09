package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func readFile(fileName string) ([]string, error){
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
        url := strings.TrimSpace(scanner.Text())
        if url != "" {
            urls = append(urls, url)
        }
    }

	if err := scanner.Err(); err != nil {
        return nil, err
    }

    return urls, nil
}

func checkUrl(url string, timeout time.Duration, wg *sync.WaitGroup, results chan<- string){
	defer wg.Done()

	start := time.Now()

    client := &http.Client{
        Timeout: timeout,
    }

    resp, err := client.Get(url)
    if err != nil {
        elapsed := time.Since(start)
        results <- fmt.Sprintf("[ERROR] `%s` (%v) - %v", url, elapsed, err)
        return
    }
    defer resp.Body.Close()
    elapsed := time.Since(start)

	syze_bites, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	} 

    result := fmt.Sprintf("[%d] `%s` (%v) (size: %v)", resp.StatusCode, url, elapsed, len(syze_bites))
    results <- result

}

func main() {
	fileName := "urls.txt"
	fmt.Println(readFile(fileName))

	timeout := time.Duration(10 * time.Second)

    urls, err := readFile("urls.txt")
    if err != nil {
        fmt.Printf("Error reading URLs: %v\n", err)
        os.Exit(1)
    }
    
    if len(urls) == 0 {
        fmt.Println("No URLs found in urls.txt")
        return
    }
    
    fmt.Printf("Checking %d URLs with timeout %v...\n\n", len(urls), timeout)
    
    var wg sync.WaitGroup
    results := make(chan string, len(urls))

    for _, url := range urls {
        wg.Add(1)
        go checkUrl(url, timeout, &wg, results)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        fmt.Println(result)
    }
}
