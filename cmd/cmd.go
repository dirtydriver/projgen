package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/dirtydriver/projgen/filescheck"
	"github.com/dirtydriver/projgen/project"
	"github.com/dirtydriver/projgen/templater"
	"github.com/spf13/cobra"
)

var (
	projectType    string
	projectName    string
	outputDir      string
	parameters     []string
	getTemplate    bool
	templateFile   string
	templateDir    string
	parametersFile string
)

func getRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "projgen",
		Short: "Project generator that renders project skeletons from templates",
		Run: func(cmd *cobra.Command, args []string) {
			// If user wants to see the expected template parameters.
			if getTemplate {
				if projectType == "" && templateDir == "" {
					log.Fatal("Please provide a project type and template directory")
				}
				templatePath := path.Join(templateDir, projectType)
				files, err := filescheck.FindTemplateFiles(templatePath, "tmpl")
				if err != nil {
					log.Fatalf("Error collection template files: %v", err)
				}

				params, err := templater.CollectParameters(files)
				if err != nil {
					log.Fatalf(err.Error())
				}

				fmt.Println("Template requires the following parameters:")
				for _, p := range params {
					fmt.Println(" -", p)
				}
				return
			}

			// Collect additional parameters passed via --parameter flags.
			paramsMap := make(map[string]interface{})
			for _, p := range parameters {
				parts := strings.SplitN(p, "=", 2)
				if len(parts) != 2 {
					log.Fatalf("Invalid parameter format: %s. Expected key=value", p)
				}
				key, value := parts[0], parts[1]
				paramsMap[key] = value
			}

			if parametersFile != "" {
				if err := filescheck.ReadParamsFromFile(parametersFile, &paramsMap); err != nil {
					log.Fatal(err.Error())
				}
			}

			_, NameExists := paramsMap["name"]
			_, TypeExists := paramsMap["type"]

			// Optionally, you can also merge explicit flags (like projectType or projectName)
			// into the parameters map if you want them to be available in the template.
			if projectName != "" && !NameExists {
				paramsMap["name"] = projectName
			}
			if projectType != "" && !TypeExists {
				paramsMap["type"] = projectType
			}

			// Proceed with generating the project using the merged parameters.
			if projectType == "" && !NameExists || projectName == "" && !TypeExists {
				log.Fatal("Project type and project name are required")
			}
			err := project.Generate(templateDir, outputDir, paramsMap)
			if err != nil {
				log.Fatalf("Error generating project: %v", err)
			}
			fmt.Println("Project generated successfully!")
		},
	}

	// Basic project flags.
	rootCmd.Flags().StringVarP(&projectType, "type", "t", "", "Type of project (e.g. maven, gradle, angular)")
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	rootCmd.Flags().StringVarP(&outputDir, "out", "o", ".", "Output directory")
	rootCmd.Flags().StringVarP(&parametersFile, "file", "f", "", "Path to the parameters file to inspect")

	// Parameter flags: allow multiple values.
	rootCmd.Flags().StringArrayVarP(&parameters, "parameter", "p", []string{}, "Additional parameters in key=value format")

	// Flags to list expected template parameters.
	rootCmd.Flags().BoolVar(&getTemplate, "get-template-params", false, "List the parameters required by the template")
	rootCmd.Flags().StringVar(&templateDir, "template-dir", "", "Path to the template directory")
	rootCmd.Flags().StringVar(&templateFile, "template-file", "", "Path to the template file to inspect")

	return rootCmd
}

func RunRootCmd() {
	rootCmd := getRootCmd()
	if len(os.Args) < 2 {
		// Print help and exit if no flags or arguments are provided.
		_ = rootCmd.Help()
		os.Exit(0)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
