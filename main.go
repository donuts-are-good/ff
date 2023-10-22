package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	namePtr := flag.String("name", "", "keyword to search in file names")
	contentsPtr := flag.String("contents", "", "keyword to search in file contents")

	flag.Parse()

	if *namePtr != "" {
		findByName(*namePtr)
	} else if *contentsPtr != "" {
		findByContents(*contentsPtr)
	} else {
		fmt.Println("Invalid search type. Use --name or --contents.")
	}
}

func findByName(keyword string) {
	var count int
	filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.Contains(info.Name(), keyword) {
			fmt.Printf("%d) %s\n", count, path)
			count++
		}
		return nil
	})

	if count == 0 {
		fmt.Printf("No files found with the name containing '%s'.\n", keyword)
	} else {
		fmt.Printf("%d files found\n", count)
	}
}

func findByContents(keyword string) {
	var count int
	filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			buf := make([]byte, 4096)
			for {
				n, err := file.Read(buf)
				if err != nil && err != io.EOF {
					return nil
				}
				if n == 0 {
					break
				}

				if strings.Contains(string(buf[:n]), keyword) {
					fmt.Printf("%d) %s\n", count, path)
					count++
					break
				}
			}
		}
		return nil
	})

	if count == 0 {
		fmt.Printf("No files found with contents containing '%s'.\n", keyword)
	} else {
		fmt.Printf("%d files found\n", count)
	}
}
