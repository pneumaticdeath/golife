package examples

import (
	"embed"
	"log"
	"path/filepath"
	"strings"

	"github.com/pneumaticdeath/golife"
)

type Example struct {
	Title string
	path  string
}

//go:embed files
var ExamplesFS embed.FS

func ListExamples() []Example {
	examples := make([]Example, 0, 64)

	files, err := ExamplesFS.ReadDir("files")
	if err != nil {
		log.Print("Can't ready embedded files:", err)
		return examples
	}

	for _, ent := range files {
		path := "files/" + ent.Name()
		loader := golife.FindReader(path)
		filecontents, fileErr := ExamplesFS.ReadFile(path)
		if fileErr != nil {
			log.Print("Error reading embedded file", path, fileErr)
			continue
		}
		game, lifeErr := loader(strings.NewReader(string(filecontents)))
		if lifeErr != nil {
			log.Print("Error parsing embedded file", path, lifeErr)
			continue
		}
		var title string = filepath.Base(path)
		if len(game.Comments) > 0 && strings.HasPrefix(game.Comments[0], "N") {
			title = game.Comments[0][1:]
		}

		examples = append(examples, Example{Title: title, path: path})
	}

	return examples
}

func LoadExample(e Example) *golife.Game {
	loader := golife.FindReader(e.path)
	contents, fileErr := ExamplesFS.ReadFile(e.path)
	if fileErr != nil {
		log.Print("Unable to read golife example file", e.path, fileErr)
		return nil
	}
	g, err := loader(strings.NewReader(string(contents)))
	if err != nil {
		log.Print("Unable to load golife example", e.path, err)
	} else {
		g.Filename = e.path
	}

	return g
}
