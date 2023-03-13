package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	urls := []string{
		"https://filesamples.com/samples/document/txt/sample1.txt",
		"https://filesamples.com/samples/document/txt/sample2.txt",
		"https://filesamples.com/samples/document/txt/sample3.txt",
		"https://www.folgerdigitaltexts.org/download/txt/AYL.txt",
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	fileChan := make(chan string)

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error downloading file:", err)
				return
			}
			defer resp.Body.Close()

			fileName := getFileName(url)
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				fmt.Println("Error writing file to disk:", err)
				return
			}

			fileChan <- fileName
		}(url)
	}

	go func() {
		for fileName := range fileChan {
			fmt.Println("Downloaded file:", fileName)
		}
	}()

	wg.Wait()

	close(fileChan)
}

func getFileName(url string) string {
	components := strings.Split(url, "/")

	fileName := components[len(components)-1]

	return fileName
}
