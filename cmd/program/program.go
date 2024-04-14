package program

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/DenisBytes/go-create-templ/cmd/flags"
	myTemplate "github.com/DenisBytes/go-create-templ/cmd/template"
	"github.com/DenisBytes/go-create-templ/cmd/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

//this is where all the template placeholder values are
type Project struct {
	ProjectName  string
	Exit         bool
	AbsolutePath string
	ProjectType  flags.Framework
	FrameworkMap map[flags.Framework]Framework
}

type Framework struct {
	packageName []string
	templater   Templater
}

// TODO: update the fields
type Templater interface {
	Main() []byte
}

func (p *Project) ExitCLI(tprogram *tea.Program) {
	if p.Exit {
		// logo render here
		if err := tprogram.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

var (
	chiPackage   = []string{"github.com/go-chi/chi/v5"}
	ginPackage   = []string{"github.com/gin-gonic/gin"}
	fiberPackage = []string{"github.com/gofiber/fiber/v2"}
	echoPackage  = []string{"github.com/labstack/echo/v4", "github.com/labstack/echo/v4/middleware"}

	godotenvPackage = []string{"github.com/joho/godotenv"}
	templPackage    = []string{"github.com/a-h/templ"}
)

const (
	root                 = "/"
	cmdApiPath           = "cmd/api"
	cmdWebPath           = "cmd/web"
	internalServerPath   = "internal/server"
	internalDatabasePath = "internal/database"
	gitHubActionPath     = ".github/workflows"
	testHandlerPath      = "tests"
)

func (p *Project) createFrameworkMap() {

	if p.FrameworkMap == nil {
		p.FrameworkMap = make(map[flags.Framework]Framework)
	}

	p.FrameworkMap[flags.Chi] = Framework{
		packageName: chiPackage,
		templater:   myTemplate.ChiTemplates{},
	}
}

func (p *Project) CreateProject() error {
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			log.Printf("Could not create directory: %v", err)
			return err
		}
	}
	nameSet, err := utils.CheckGitConfig("user.name")
	if err != nil {
		cobra.CheckErr(err)
	}
	if !nameSet {
		fmt.Println("user.name is not set in git config.")
		fmt.Println("Please set up git config before trying again.")
		panic("\nGIT CONFIG ISSUE: user.name is not set in git config.\n")
	}

	emailSet, err := utils.CheckGitConfig("user.email")
	if err != nil {
		cobra.CheckErr(err)
	}
	if !emailSet {
		fmt.Println("user.email is not set in git config.")
		fmt.Println("Please set up git config before trying again.")
		panic("\nGIT CONFIG ISSUE: user.email is not set in git config.\n")
	}

	p.ProjectName = strings.TrimSpace(p.ProjectName)

	projectPath := filepath.Join(p.AbsolutePath, p.ProjectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		err := os.MkdirAll(projectPath, 0751)
		if err != nil {
			log.Printf("Error creating root project directory %v\n", err)
			return err
		}
	}

	p.createFrameworkMap()

	err = utils.InitGoMod(p.ProjectName, projectPath)
	if err != nil {
		log.Printf("Could not initialize go.mod in new project %v\n", err)
		cobra.CheckErr(err)
	}

	if p.ProjectType != flags.StandardLibrary {
		err = utils.GoGetPackage(projectPath, p.FrameworkMap[p.ProjectType].packageName)
		if err != nil {
			log.Printf("Could not install go dependency for the chosen framework %v\n", err)
			cobra.CheckErr(err)
		}
	}

	err = utils.GoGetPackage(projectPath, godotenvPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = p.CreatePath(cmdApiPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(cmdApiPath, projectPath, "main.go", "main")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = utils.ExecuteCmd("git", []string{"init"}, projectPath)
	if err != nil {
		log.Printf("Error initializing git repo: %v", err)
		cobra.CheckErr(err)
		return err
	}

	gitignoreFile, err := os.Create(filepath.Join(projectPath, ".gitignore"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}
	defer gitignoreFile.Close()

	gitignoreTemplate := template.Must(template.New(".gitignore").Parse(string(myTemplate.GitIgnoreTemplate())))
	err = gitignoreTemplate.Execute(gitignoreFile, p)
	if err != nil {
		return err
	}

	airTomlFile, err := os.Create(filepath.Join(projectPath, ".air.toml"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	defer airTomlFile.Close()

	airTomlTemplate := template.Must(template.New("airtoml").Parse(string(myTemplate.AirTomlTemplate())))
	err = airTomlTemplate.Execute(airTomlFile, p)
	if err != nil {
		return err
	}

	err = utils.GoTidy(projectPath)
	if err != nil {
		log.Printf("Could not go tidy in new project %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoFmt(projectPath)
	if err != nil {
		log.Printf("Could not gofmt in new project %v\n", err)
		cobra.CheckErr(err)
		return err
	}

	err = utils.ExecuteCmd("git", []string{"add", "."}, projectPath)
	if err != nil {
		log.Printf("Error adding files to git repo: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = utils.ExecuteCmd("git", []string{"commit", "-m", "Initial commit"}, projectPath)
	if err != nil {
		log.Printf("Error committing files to git repo: %v", err)
		cobra.CheckErr(err)
		return err
	}

	return nil
}

func (p *Project) CreatePath(pathToCreate string, projectPath string) error {
	path := filepath.Join(projectPath, pathToCreate)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0751)
		if err != nil {
			log.Printf("Error creating directory %v\n", err)
			return err
		}
	}

	return nil
}

func (p *Project) CreateFileWithInjection(pathToCreate string, projectPath string, fileName string, methodName string) error {
	createdFile, err := os.Create(filepath.Join(projectPath, pathToCreate, fileName))
	if err != nil {
		return err
	}

	defer createdFile.Close()

	switch methodName {
	case "main":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.Main())))
		err = createdTemplate.Execute(createdFile, p)
		// case "server":
		// 	createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.Server())))
		// 	err = createdTemplate.Execute(createdFile, p)
		// case "routes":
		// 	routeFileBytes := p.FrameworkMap[p.ProjectType].templater.Routes()
		// 	createdTemplate := template.Must(template.New(fileName).Parse(string(routeFileBytes)))
		// 	err = createdTemplate.Execute(createdFile, p)
		// case "releaser":
		// 	createdTemplate := template.Must(template.New(fileName).Parse(string(advanced.Releaser())))
		// 	err = createdTemplate.Execute(createdFile, p)
		// case "releaser-config":
		// 	createdTemplate := template.Must(template.New(fileName).Parse(string(advanced.ReleaserConfig())))
		// 	err = createdTemplate.Execute(createdFile, p)
		// case "env":
		// 	createdTemplate := template.Must(template.New(fileName).Parse(string(tpl.GlobalEnvTemplate())))
		// 	err = createdTemplate.Execute(createdFile, p)
	}

	if err != nil {
		return err
	}

	return nil
}
