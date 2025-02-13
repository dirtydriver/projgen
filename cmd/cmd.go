package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	projectType  string
	projectName  string
	outputDir    string
	parameters   []string
	getTemplate  bool
	templateFile string
)

func getRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "projgen",
		Short: "Project generator that renders project skeletons from templates",
		Run: func(cmd *cobra.Command, args []string) {
			// If user wants to see the expected template parameters.
			/* 	if getTemplate {
				if templateFile == "" {
					log.Fatal("Please provide a template file using --template-file")
				}
				params, err := project.ExtractTemplateParams(templateFile)
				if err != nil {
					log.Fatalf("Error extracting template parameters: %v", err)
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

			// Optionally, you can also merge explicit flags (like projectType or projectName)
			// into the parameters map if you want them to be available in the template.
			if projectName != "" {
				paramsMap["ProjectName"] = projectName
			}
			if projectType != "" {
				paramsMap["ProjectType"] = projectType
			}

			// Proceed with generating the project using the merged parameters.
			if projectType == "" || projectName == "" {
				log.Fatal("Project type and project name are required")
			}
			err := project.Generate(projectType, projectName, outputDir, paramsMap)
			if err != nil {
				log.Fatalf("Error generating project: %v", err)
			}
			fmt.Println("Project generated successfully!") */
		},
	}

	// Basic project flags.
	rootCmd.Flags().StringVarP(&projectType, "type", "t", "", "Type of project (e.g. maven, gradle, angular)")
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	rootCmd.Flags().StringVarP(&outputDir, "out", "o", ".", "Output directory")

	// Parameter flags: allow multiple values.
	rootCmd.Flags().StringArrayVarP(&parameters, "parameter", "p", []string{}, "Additional parameters in key=value format")

	// Flags to list expected template parameters.
	rootCmd.Flags().BoolVar(&getTemplate, "get-template-params", false, "List the parameters required by the template")
	rootCmd.Flags().StringVar(&templateFile, "template-file", "", "Path to the template file to inspect")

	return rootCmd
}

func RunRootCmd() {

	if err := getRootCmd().Execute(); err != nil {
		os.Exit(1)
	}

}
