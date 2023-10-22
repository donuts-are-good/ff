package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	// define the flags
	namePtr := flag.String("name", "", "keyword to search in file names")
	contentsPtr := flag.String("contents", "", "keyword to search in file contents")
	hiddenPtr := flag.Bool("hidden", false, "include hidden files in the search")

	// parse the flags
	flag.Parse()

	// unless a path is specified, assume it is /
	path := "/"
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	// check if we're looking in names or contents
	if *namePtr != "" {
		findByName(*namePtr, path, *hiddenPtr)
	} else if *contentsPtr != "" {
		findByContents(*contentsPtr, path, *hiddenPtr)
	} else {
		fmt.Println("Invalid search type. Use --name or --contents.")
	}
}

// formatDuration measures how long the operation took and returns it as a human readable amount of time
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

// findByName will search for a substring in the filename
func findByName(keyword string, path string, hidden bool) {

	// we use count for the result numbers and final count
	var count int

	// this is when we start counting the duration of the operation
	start := time.Now()

	// walk the path recursively through the files
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// if we're not doing hidden files, skip it
		if !hidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// if we have a positive result...
		if strings.Contains(info.Name(), keyword) {
			fmt.Printf("%d) %s\n", count, path)
			count++
		}
		return nil
	})

	if err != nil {
		return
	}

	if count == 0 {
		fmt.Printf("No files found with the name containing '%s'.\n", keyword)
	}

	// this is how long the operation took
	elapsed := time.Since(start)
	fmt.Printf("\n%d files found in %s\n", count, formatDuration(elapsed))

}

func findByContents(keyword string, path string, hidden bool) {
	var count int
	start := time.Now()

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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
				return nil
			}
			defer file.Close()

			// Read the entire file
			content, err := io.ReadAll(file)
			if err != nil {
				return nil
			}

			// If the file contains the keyword, print the path
			if strings.Contains(string(content), keyword) {
				fmt.Printf("%d) %s\n", count, path)
				count++
			}
		}
		return nil
	})

	if err != nil {
		return
	}

	if count == 0 {
		fmt.Printf("No files found with contents containing '%s'.\n", keyword)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n%d files found in %s\n", count, formatDuration(elapsed))
}
