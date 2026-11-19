package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func downloadImage(id int, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	url = strings.TrimSpace(url)
	if url == "" {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Download failed:", resp.Status)
		return
	}

	file, err := os.Create(fmt.Sprintf("images/pic_%d.jpg", id+1))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Downloaded image:", id+1)
}

func main() {
	os.MkdirAll("images", 0755)

	data, _ := os.ReadFile("urls.txt")
	urls := strings.Split(string(data), "\n")

	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go downloadImage(i, url, &wg)
	}

	wg.Wait()

	fmt.Println("All images downloaded!")
}
