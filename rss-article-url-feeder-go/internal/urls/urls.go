package urls

import (
	"bufio"
	"fmt"
	"github.com/heussd/nats-news-keyword-matcher.go/internal/config"
	"math/rand"
	"os"
	"strings"
)

var Urls []string

func init() {
	ReloadUrls()
}

func ReloadUrls() {
	Urls = nil

	readFile, err := os.Open(config.UrlsFile)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		var text = fileScanner.Text()

		if text == "" ||
			strings.HasPrefix(text, "#") {
			continue
		}

		fmt.Printf("Adding URL \"%s\"\n", text)

		Urls = append(Urls, text)
	}

	// Randomize URLs
	for i := range Urls {
		j := rand.Intn(i + 1)
		Urls[i], Urls[j] = Urls[j], Urls[i]
	}

	if len(Urls) == 0 {
		fmt.Println("Error: No URLS found")
		os.Exit(1)
	} else {
		fmt.Printf("URL init complete. There are %d urls.\n", len(Urls))
	}
}
