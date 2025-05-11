package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/dirtydriver/projgen/filescheck"
	"github.com/dirtydriver/projgen/project"
	"github.com/dirtydriver/projgen/templater"
	"github.com/dirtydriver/projgen/utils"
	"github.com/dirtydriver/projgen/version"
	"github.com/spf13/cobra"
)

var (
	projectType    string
	projectName    string
	outputDir      string
	parameters     []string
	templateDir    string
	parametersFile string
)

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "projgen",
		Short: "Project generator that renders project skeletons from templates",
	}

	// Add shared flags that apply to multiple commands
	rootCmd.PersistentFlags().StringVar(&templateDir, "template-dir", "", "Path to the template directory")
	rootCmd.PersistentFlags().StringVarP(&projectType, "type", "t", "", "Type of project (e.g. maven, gradle, angular)")

	// Add all subcommands
	rootCmd.AddCommand(
		getVersionCmd(),
		getGenerateCmd(),
		getInspectCmd(),
	)

	return rootCmd
}

func getGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project from a template",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if templateDir == "" {
				return fmt.Errorf("required flag \"template-dir\" not set")
			}
			if projectType == "" {
				return fmt.Errorf("required flag \"type\" not set")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Collect additional parameters passed via --parameter flags
			paramsMap := make(map[string]interface{})

			if parametersFile != "" {
				if err := filescheck.ReadParamsFromYaml(parametersFile, &paramsMap); err != nil {
					log.Fatal(err.Error())
				}
			}

			utils.ApplyOverrides(paramsMap, parameters)
			templatePath := path.Join(templateDir, projectType)

			files, err := filescheck.FindTemplateFiles(templatePath, "tmpl")

			if err != nil {
				log.Fatalf("Error collecting template files: %v", err)
			}

			params, err := templater.CollectParameters(files)
			if err != nil {
				log.Fatalf(err.Error())
			}

			err = utils.CheckMissingKeys(paramsMap, params)
			if err != nil {
				log.Fatalf(err.Error())
			}

			err = project.Generate(templatePath, outputDir, paramsMap)
			if err != nil {
				log.Fatalf("Error generating project: %v", err)
			}
			fmt.Println("Project generated successfully!")
		},
	}

	// Add flags specific to generate command
	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project (can also be provided via --parameter name=value)")
	cmd.Flags().StringVarP(&outputDir, "out", "o", ".", "Output directory")
	cmd.Flags().StringVarP(&parametersFile, "file", "f", "", "Path to the parameters file")
	cmd.Flags().StringArrayVarP(&parameters, "parameter", "p", []string{}, "Additional parameters in key=value format")

	return cmd
}

func getInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect",
		Short: "Inspect template parameters and requirements",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if templateDir == "" {
				return fmt.Errorf("required flag \"template-dir\" not set")
			}
			if projectType == "" {
				return fmt.Errorf("required flag \"type\" not set")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			templatePath := path.Join(templateDir, projectType)
			files, err := filescheck.FindTemplateFiles(templatePath, "tmpl")
			if err != nil {
				log.Fatalf("Error collecting template files: %v", err)
			}

			params, err := templater.CollectParameters(files)
			if err != nil {
				log.Fatalf(err.Error())
			}

			fmt.Println("Template requires the following parameters:")
			for _, p := range params {
				fmt.Println(" -", p)
			}
		},
	}
}

func getVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the current version of projgen",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("projgen version %s\n", version.Version)
		},
	}
}

// RunRootCmd executes the root command of the projgen CLI tool.
// If no arguments are provided, it displays the help information.
// The function handles command execution and exits with an error if the command fails.
func RunRootCmd() {
	rootCmd := getRootCmd()
	if len(os.Args) < 2 {
		// Print help and exit if no flags or arguments are provided
		_ = rootCmd.Help()
		os.Exit(0)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
