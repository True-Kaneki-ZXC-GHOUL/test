package main

import (
    "bufio"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"
)

// readURLs читает URL из файла и возвращает их список
func readURLs(filename string) ([]string, error) {
    file, err := os.Open(filename) // Открываем файл для чтения
    if err != nil {
        return nil, err // Если ошибка - возвращаем её
    }
    defer file.Close() // Гарантируем закрытие файла при выходе из функции

    var urls []string
    scanner := bufio.NewScanner(file) // Создаём сканер для построчного чтения
    
    // Читаем файл построчно
    for scanner.Scan() {
        url := strings.TrimSpace(scanner.Text()) // Убираем лишние пробелы
        if url != "" { // Пропускаем пустые строки
            urls = append(urls, url)
        }
    }
    
    // Проверяем, не было ли ошибок при сканировании
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    
    return urls, nil
}

// checkURL проверяет доступность одного URL
func checkURL(url string, timeout time.Duration, wg *sync.WaitGroup, results chan<- string) {
    defer wg.Done() // Уменьшаем счётчик WaitGroup при завершении функции

    start := time.Now() // Запоминаем время начала запроса
    
    // Создаём HTTP клиент с таймаутом
    client := &http.Client{
        Timeout: timeout,
    }
    
    // Отправляем GET запрос
    resp, err := client.Get(url)
    if err != nil {
        // Если ошибка - выводим сообщение об ошибке
        elapsed := time.Since(start)
        results <- fmt.Sprintf("[ERROR] `%s` (%v) - %v", url, elapsed, err)
        return
    }
    defer resp.Body.Close() // Гарантируем закрытие тела ответа
    
    // Вычисляем время выполнения запроса
    elapsed := time.Since(start)
    
    // Формируем результат в требуемом формате
    result := fmt.Sprintf("[%d] `%s` (%v)", resp.StatusCode, url, elapsed)
    results <- result
}

func main() {
    // Определяем флаг для таймаута (extra задание)
    timeout := flag.Duration("t", 10*time.Second, "Timeout for HTTP requests")
    flag.Parse()
    
    // Читаем URL из файла
    urls, err := readURLs("urls.txt")
    if err != nil {
        fmt.Printf("Error reading URLs: %v\n", err)
        os.Exit(1)
    }
    
    // Проверяем, что файл не пустой
    if len(urls) == 0 {
        fmt.Println("No URLs found in urls.txt")
        return
    }
    
    fmt.Printf("Checking %d URLs with timeout %v...\n\n", len(urls), *timeout)
    
    // Создаём WaitGroup для ожидания завершения всех горутин
    var wg sync.WaitGroup
    // Создаём канал для сбора результатов
    results := make(chan string, len(urls))
    
    // Запускаем проверку каждого URL в отдельной горутине
    for _, url := range urls {
        wg.Add(1) // Увеличиваем счётчик WaitGroup
        go checkURL(url, *timeout, &wg, results) // Запускаем горутину
    }
    
    // Запускаем горутину, которая закроет канал после завершения всех проверок
    go func() {
        wg.Wait() // Ждём завершения всех горутин
        close(results) // Закрываем канал результатов
    }()
    
    // Выводим результаты по мере их поступления
    for result := range results {
        fmt.Println(result)
    }
}