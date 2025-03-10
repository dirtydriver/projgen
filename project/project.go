package project

import (
	"errors"
	"path/filepath"

	"github.com/dirtydriver/projgen/filescheck"
	"github.com/dirtydriver/projgen/templater"
)

func Generate(templateDir, outputDir string, paramsMap map[string]interface{}) error {
	projectType, ok := paramsMap["type"].(string)

	if !ok {
		return errors.New("type is not a string")
	}

	templatePath := filepath.Join(templateDir, projectType)
	filelist, err := filescheck.FilesInDirectories(templatePath)

	if err != nil {
		return err
	}

	for _, file := range filelist {

		if templater.IsTemplate(file) {
			rendered, err := templater.RenderTemplate(file, paramsMap)

			if err != nil {
				return err
			}

		}

	}

}
