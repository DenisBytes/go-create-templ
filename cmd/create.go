package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/DenisBytes/go-create-templ/cmd/flags"
	"github.com/DenisBytes/go-create-templ/cmd/program"
	"github.com/DenisBytes/go-create-templ/cmd/steps"
	"github.com/DenisBytes/go-create-templ/cmd/ui/multiinput"
	"github.com/DenisBytes/go-create-templ/cmd/ui/spinner"
	"github.com/DenisBytes/go-create-templ/cmd/ui/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var logo = `
 _____                              _        _
|_   _|   ___   _ __ ___    _ __   | |      / \     _ __    _ __
  | |    / _ \ | '_ ' _ \  | '_ \  | |     / _ \   | '_ \  | '_ \
  | |   |  __/ | | | | | | | |_) | | |    / ___ \  | |_) | | |_) |
  |_|    \___| |_| |_| |_| | .__/  |_|   /_/   \_\ | .__/  | .__/
						   |_|                     |_|     |_|
`

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)

type Options struct {
	ProjectName *textinput.Output
	ProjectType *multiinput.Selection
}

func init() {
	var flagFramework flags.Framework
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("name", "n", "", "Name of project to create")
	createCmd.Flags().VarP(&flagFramework, "framework", "f", fmt.Sprintf("Framework to use. Allowed values: %s", strings.Join(flags.AllowedProjectTypes, ", ")))
}

var createCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a Go project and don't worry about the structure",
	Long:  "Go Blueprint is a CLI tool that allows you to focus on the actual Go code, and not the project structure. Perfect for someone new to the Go language",

	Run: func(cmd *cobra.Command, args []string) {
		var tprogram *tea.Program
		var err error

		flagName := cmd.Flag("name").Value.String()
		if flagName != "" && doesDirectoryExistAndIsNotEmpty(flagName) {
			err = fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different name", flagName)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		flagFramework := flags.Framework(cmd.Flag("framework").Value.String())

		options := Options{
			ProjectName: &textinput.Output{},
			ProjectType: &multiinput.Selection{},
		}

		project := &program.Project{
			ProjectName: flagName,
			ProjectType: flagFramework,
		}

		steps := steps.InitSteps(flagFramework)
		fmt.Printf("%s\n", logoStyle.Render(logo))

		if project.ProjectName == "" {
			tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.ProjectName, "What is the name of your project?", project))
			if _, err := tprogram.Run(); err != nil {
				log.Printf("Name of project contains an error: %v", err)
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			if doesDirectoryExistAndIsNotEmpty(options.ProjectName.Output) {
				err = fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different name", options.ProjectName.Output)
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			project.ExitCLI(tprogram)

			project.ProjectName = options.ProjectName.Output
			err := cmd.Flag("name").Value.Set(project.ProjectName)
			if err != nil {
				log.Fatal("failed to set the name flag value", err)
			}
		}

		if project.ProjectType == "" {
			step := steps.Steps["framework"]
			tprogram = tea.NewProgram(multiinput.InitialModelMulti(step.Options, options.ProjectType, step.Headers, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			project.ExitCLI(tprogram)

			step.Field = options.ProjectType.Choice

			project.ProjectType = flags.Framework(strings.ToLower(options.ProjectType.Choice))
			err := cmd.Flag("framework").Value.Set(project.ProjectType.String())
			if err != nil {
				log.Fatal("failed to set the framework flag value", err)
			}
		}
		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Printf("could not get current working directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		project.AbsolutePath = currentWorkingDir

		spinner := tea.NewProgram(spinner.InitialModelNew())
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			if _, err := spinner.Run(); err != nil {
				cobra.CheckErr(err)
			}
		}()

		// defer func() {
		// 	if r := recover(); r != nil {
		// 		fmt.Println("The program encountered an unexpected issue and had to exit. The error was:", r)
		// 		fmt.Println("If you continue to experience this issue, please post a message on our GitHub page or join our Discord server for support.")
		// 		if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
		// 			log.Printf("Problem releasing terminal: %v", releaseErr)
		// 		}
		// 	}
		// }()
		fmt.Println(project.ProjectName, project.ProjectType)

		err = project.CreateProject()
		if err != nil {
			if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
				log.Printf("Problem releasing terminal: %v", releaseErr)
			}
			log.Printf("Problem creating files for project. %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		fmt.Println(endingMsgStyle.Render("\nNext steps:"))
		fmt.Println(endingMsgStyle.Render(fmt.Sprintf("• cd into the newly created project with: `cd %s`\n", project.ProjectName)))

		err = spinner.ReleaseTerminal()
		if err != nil {
			log.Printf("Could not release terminal: %v", err)
			cobra.CheckErr(err)
		}
	},
}

func doesDirectoryExistAndIsNotEmpty(name string) bool {
	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			log.Printf("could not read directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err))
		}
		if len(dirEntries) > 0 {
			return true
		}
	}
	return false
}
