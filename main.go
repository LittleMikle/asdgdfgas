package main

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

type Config struct {
	PathDir  string
	ReadFile string
	ResDir   string
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	var files []string

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	ch := make(chan map[string]int)

	err := initConfig()
	if err != nil {
		log.Fatal().Msg("error with viper")
	} else {
		log.Info().Msg("Config initiation successful")
	}
	cfg := Config{
		PathDir:  viper.GetString("pathDir"),
		ReadFile: viper.GetString("readFile"),
		ResDir:   viper.GetString("resDir"),
	}

	err = filepath.Walk(cfg.PathDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			name := info.Name()
			files = append(files, name)
		}
		return nil
	})
	if err != nil {
		log.Fatal().Msg("failed with open directory")
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		resMap := map[string]int{}

		for _, file := range files {
			readFile, err := os.Open(cfg.ReadFile + file)
			if err != nil {
				log.Error().Err(err).Msg("failed with readfile")
			}
			fileScanner := bufio.NewScanner(readFile)

			fileScanner.Split(bufio.ScanLines)

			for fileScanner.Scan() {
				key := regexp.MustCompile(`[^a-zA-Z]+`).ReplaceAllString(fileScanner.Text(), "")
				valueStr := regexp.MustCompile(`\D+`).ReplaceAllString(fileScanner.Text(), "")
				value, _ := strconv.Atoi(valueStr)
				mu.Lock()
				resMap[key] += value
				mu.Unlock()
			}
			readFile.Close()
		}
		ch <- resMap
		close(ch)
	}()
	wg.Wait()

	resMap := <-ch

	file, _ := os.Create(cfg.ResDir)
	w := bufio.NewWriter(file)

	for k, v := range resMap {
		w.WriteString(k + ":" + strconv.Itoa(v) + "\n")
	}
	w.Flush()
	log.Info().Msg("File created successfully")
}
