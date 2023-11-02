package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var files []string
	var resA int
	var resB int
	var resC int
	var resFF int

	wg := sync.WaitGroup{}

	err := filepath.Walk("./testing", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			name := info.Name()
			files = append(files, name)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed with open directory")
	}

	chA := make(chan int, len(files))
	chB := make(chan int, len(files))
	chC := make(chan int, len(files))
	chFF := make(chan int, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			readFile, err := os.Open("./testing/" + file)
			if err != nil {
				fmt.Println(err)
			}
			fileScanner := bufio.NewScanner(readFile)

			fileScanner.Split(bufio.ScanLines)

			for fileScanner.Scan() {
				if strings.Contains(fileScanner.Text(), "a") {
					str := regexp.MustCompile(`\D+`).ReplaceAllString(fileScanner.Text(), "")
					val, _ := strconv.Atoi(str)
					chA <- val
				}
				if strings.Contains(fileScanner.Text(), "b") {
					str := regexp.MustCompile(`\D+`).ReplaceAllString(fileScanner.Text(), "")
					val, _ := strconv.Atoi(str)
					chB <- val
				}
				if strings.Contains(fileScanner.Text(), "c") {
					str := regexp.MustCompile(`\D+`).ReplaceAllString(fileScanner.Text(), "")
					val, _ := strconv.Atoi(str)
					chC <- val
				}
				if strings.Contains(fileScanner.Text(), "ff") {
					str := regexp.MustCompile(`\D+`).ReplaceAllString(fileScanner.Text(), "")
					val, _ := strconv.Atoi(str)
					chFF <- val
				}
			}
			readFile.Close()
		}(file)
	}
	wg.Wait()
	close(chA)
	close(chB)
	close(chC)
	close(chFF)

	for a := range chA {
		resA += a
	}
	for b := range chB {
		resB += b
	}
	for c := range chC {
		resC += c
	}
	for ff := range chFF {
		resFF += ff
	}

	file, _ := os.Create("./testing/result.txt")
	w := bufio.NewWriter(file)
	lines := []string{"a:", "b:", "c:", "ff:"}
	for _, s := range lines {
		if s == "a:" {
			w.WriteString(s + strconv.Itoa(resA) + "\n")
		}
		if s == "b:" {
			w.WriteString(s + strconv.Itoa(resB) + "\n")
		}
		if s == "c:" {
			w.WriteString(s + strconv.Itoa(resC) + "\n")
		}
		if s == "ff:" {
			w.WriteString(s + strconv.Itoa(resFF) + "\n")
		}
	}
	w.Flush()

	fmt.Printf("RESULT A = %d\n", resA)
	fmt.Printf("RESULT B = %d\n", resB)
	fmt.Printf("RESULT C = %d\n", resC)
	fmt.Printf("RESULT FF = %d\n", resFF)
}
