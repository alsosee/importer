// importer scrapes web pages and creates YAML files
// Arguments:
// - info: path to "info" directory where YAML files are stored
// - scraper: name of the scraper to use
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alsosee/finder/structs"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

type Scraper interface {
	Scrape(page string) (map[string]structs.Content, error)
}

type config struct {
	Info    string `long:"info" short:"i" description:"path to 'info' directory where YAML files are stored" required:"true"`
	Scraper string `long:"scraper" short:"s" description:"name of the scraper to use" required:"true"`
	Page    string `long:"page" short:"p" description:"page to scrape"`
	// YearFrom int    `long:"year-from" short:"f" description:"year to start scraping from (including)" default:"2023"`
	// YearTo   int    `long:"year-to" short:"t" description:"year to stop scraping at (including)" default:"2023"`
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
	}

	log.Println("Done")
}

func run() error {
	var cfg config
	_, err := flags.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	scraperName := strings.ToLower(cfg.Scraper)
	var scraper Scraper

	switch scraperName {
	case "wikipedia":
		scraper = NewWikipedia()
	// case "bafta":
	// 	scraper, err = NewBAFTA(cfg.YearFrom, cfg.YearTo)
	// 	if err != nil {
	// 		return fmt.Errorf("creating BAFTA scraper: %w", err)
	// 	}
	default:
		return fmt.Errorf("unknown scraper: %s", cfg.Scraper)
	}

	contents, err := scraper.Scrape(cfg.Page)
	if err != nil {
		return fmt.Errorf("scraping: %w", err)
	}

	// write contents to YAML files
	for path, content := range contents {
		filePath := fmt.Sprintf("%s/%s.yml", cfg.Info, path)
		b, err := yaml.Marshal(content)
		if err != nil {
			return fmt.Errorf("marshaling content: %w", err)
		}

		err = os.WriteFile(filePath, b, 0644)
		if err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}
	}

	return nil
}
