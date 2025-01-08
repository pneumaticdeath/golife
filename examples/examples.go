package examples

import (
	"embed"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/pneumaticdeath/golife"
)

type Example struct {
	Title string
	Category string
	path  string
}

//go:embed files
var ExamplesFS embed.FS

func ListExamples() []Example {
	examples := make([]Example, 0, 64)

	fs.WalkDir(ExamplesFS, "files", func(dir string, file fs.DirEntry, err error) error {
		if !strings.HasSuffix(file.Name(), ".rle") {
			return nil
		}

		fields := strings.Split(dir,"/")
		category := fields[1]
		path := dir
		loader := golife.FindReader(path)
		filecontents, fileErr := ExamplesFS.ReadFile(path)
		if fileErr != nil {
			log.Print("Error reading embedded file ", fileErr)
			return nil
		}
		game, lifeErr := loader(strings.NewReader(string(filecontents)))
		if lifeErr != nil {
			log.Print("Error parsing embedded file ", path, " ", lifeErr)
			return nil
		}
		var title string = filepath.Base(path)
		//FixMe: The name comment doesn't have to be the first one.
		if len(game.Comments) > 0 && strings.HasPrefix(game.Comments[0], "N ") {
			title = game.Comments[0][2:]
		}

		examples = append(examples, Example{Title: title, Category: category, path: path})

		return nil
	})

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
