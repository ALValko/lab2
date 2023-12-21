package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func search(str string, filename string) (bool, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == str {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// Функция для параллельного поиска строки в файлах
func parallelSearch(str string, files []string) {
	var wg sync.WaitGroup
	foundCh := make(chan string)
	errCh := make(chan error)

	for _, filename := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()

			found, err := search(str, filename)
			if err != nil {
				errCh <- err
			}
			if found {
				foundCh <- filename
			}
		}(filename)
	}

	go func() {
		wg.Wait()
		close(foundCh)
		close(errCh)
	}()

	for filename := range foundCh {
		fmt.Printf("Строка найдена в файле: %s\n", filename)
	}

	for err := range errCh {
		fmt.Printf("Ошибка поиска строки: %s\n", err)
	}
}

// Функция для последовательного поиска строки в файлах
func sequentialSearch(str string, files []string) {
	for _, filename := range files {
		found, err := search(str, filename)
		if err != nil {
			fmt.Printf("Ошибка поиска строки: %s\n", err)
		}
		if found {
			fmt.Printf("Строка найдена в файле: %s\n", filename)
		}
	}
}

func main() {
	files, err := filepath.Glob("files/*.txt")
	if err != nil {
		fmt.Println("Ошибка получения списка файлов:", err)
		return
	}

	str := "diagram"
	fmt.Println("Ищем в файлах строку:", str)

	fmt.Println("\nПоследовательный поиск:")
	startSequential := time.Now()
	sequentialSearch(str, files)
	fmt.Println("Время выполнения последовательного поиска:", time.Since(startSequential))

	fmt.Println("\nПараллельный поиск:")
	startParallel := time.Now()
	parallelSearch(str, files)
	fmt.Println("Время выполнения параллельного поиска:", time.Since(startParallel))

}
