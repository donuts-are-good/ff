package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	namePtr := flag.String("name", "", "keyword to search in file names")
	contentsPtr := flag.String("contents", "", "keyword to search in file contents")
	hiddenPtr := flag.Bool("hidden", false, "include hidden files in the search")

	flag.Parse()

	path := "/"
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	if *namePtr != "" {
		findByName(*namePtr, path, *hiddenPtr)
	} else if *contentsPtr != "" {
		findByContents(*contentsPtr, path, *hiddenPtr)
	} else {
		fmt.Println("Invalid search type. Use --name or --contents.")
	}
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func findByName(keyword string, path string, hidden bool) {
	var count int
	start := time.Now()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if !strings.Contains(err.Error(), "operation not permitted") {
				log.Printf("Error accessing path %s: %v\n", path, err)
			}
			return nil
		}
		if !hidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.Contains(info.Name(), keyword) {
			fmt.Printf("%d) %s\n", count, path)
			count++
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %s: %v\n", path, err)
		return
	}

	if count == 0 {
		fmt.Printf("No files found with the name containing '%s'.\n", keyword)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n%d files found in %s\n", count, formatDuration(elapsed))

}

func findByContents(keyword string, path string, hidden bool) {
	var count int
	start := time.Now()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if !strings.Contains(err.Error(), "operation not permitted") {
				log.Printf("Error accessing path %s: %v\n", path, err)
			}
			return nil
		}
		if !hidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				if !strings.Contains(err.Error(), "operation not permitted") {
					log.Printf("Error opening file %s: %v\n", path, err)
				}
				return nil
			}

			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), keyword) {
					fmt.Printf("%d) %s\n", count, path)
					count++
					break
				}
			}
			if err := scanner.Err(); err != nil {
				log.Printf("Error reading file %s: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking the path %s: %v\n", path, err)
		return
	}

	if count == 0 {
		fmt.Printf("No files found with contents containing '%s'.\n", keyword)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n%d files found in %s\n", count, formatDuration(elapsed))

}
