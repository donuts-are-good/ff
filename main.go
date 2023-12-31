package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	switch {
	case *namePtr != "":
		findByName(*namePtr, path, *hiddenPtr)
	case *contentsPtr != "":
		findByContents(*contentsPtr, path, *hiddenPtr)
	default:
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			findByName(flag.Arg(1), path, *hiddenPtr)
		}()
		go func() {
			defer wg.Done()
			findByContents(flag.Arg(1), path, *hiddenPtr)
		}()
		wg.Wait()
	}
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
	fmt.Printf("\n%d files found in %s\n", count, elapsed.Truncate(10*time.Millisecond))
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
	fmt.Printf("\n%d files found in %s\n", count, elapsed.Truncate(10*time.Millisecond))
}
