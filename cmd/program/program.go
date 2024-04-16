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

// this is where all the template placeholder values are
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
	HandlerAuth() []byte
	HandlerHome() []byte
	HandlerMiddleware() []byte
	HandlerSettings() []byte
	HandlerUtil() []byte
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

	godotenvPackage       = []string{"github.com/joho/godotenv"}
	golangMigratePackage  = []string{"github.com/golang-migrate/migrate/v4"}
	templPackage          = []string{"github.com/a-h/templ"}
	uuidPackage           = []string{"github.com/google/uuid"}
	sessionPackage        = []string{"github.com/gorilla/sessions"}
	postgresDriverPackage = []string{"github.com/lib/pq"}
	supabasePackage       = []string{"github.com/nedpals/supabase-go"}
	bunPackage1           = []string{"github.com/uptrace/bun"}
	bunPackage2           = []string{"github.com/uptrace/bun/dialect/pgdialect"}
	bunPackage3           = []string{"github.com/uptrace/bun/extra/bundebug"}
	airPackage            = []string{"github.com/cosmtrek/air@latest"}
)

const (
	root                 = "/"
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

	// checking if absolute path exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			log.Printf("Could not create directory: %v", err)
			return err
		}
	}

	// checking git config --get username
	nameSet, err := utils.CheckGitConfig("user.name")
	if err != nil {
		cobra.CheckErr(err)
	}
	if !nameSet {
		fmt.Println("user.name is not set in git config.")
		fmt.Println("Please set up git config before trying again.")
		panic("\nGIT CONFIG ISSUE: user.name is not set in git config.\n")
	}

	// checking git config --get email
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

	// initializing go project (go mod init <name of project>)
	err = utils.InitGoMod(p.ProjectName, projectPath)
	if err != nil {
		log.Printf("Could not initialize go.mod in new project %v\n", err)
		cobra.CheckErr(err)
	}

	// importing framework if not standard library (go get <framework>)
	if p.ProjectType != flags.StandardLibrary {
		err = utils.GoGetPackage(projectPath, p.FrameworkMap[p.ProjectType].packageName)
		if err != nil {
			log.Printf("Could not install go dependency for the chosen framework %v\n", err)
			cobra.CheckErr(err)
		}
	}

	// importing godotenv package (go get <godotenvPackage>)
	err = utils.GoGetPackage(projectPath, godotenvPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoInstallPackage(projectPath, templPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, templPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoInstallPackage(projectPath, airPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, airPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, golangMigratePackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, postgresDriverPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, sessionPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, supabasePackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, uuidPackage)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, bunPackage1)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, bunPackage2)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoGetPackage(projectPath, bunPackage3)
	if err != nil {
		log.Printf("Could not install go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = p.CreateFileWithInjection("/", projectPath, "main.go", "main")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/handler", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/handler", projectPath, "auth.go", "handler/auth")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/handler", projectPath, "home.go", "handler/home")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/handler", projectPath, "middleware.go", "handler/middleware")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/handler", projectPath, "settings.go", "handler/settings")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/handler", projectPath, "util.go", "handler/util")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/db", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/db", projectPath, "db.go", "db/db")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/db", projectPath, "query.go", "db/query")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/cmd/migrate", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/cmd/migrate", projectPath, "main.go", "migrate/migrate")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/cmd/reset", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/cmd/reset", projectPath, "main.go", "migrate/reset")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/types", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/types", projectPath, "user.go", "types/user")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/types", projectPath, "account.go", "types/account")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/pkg/validate", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/pkg/validate", projectPath, "validate.go", "pkg/validate")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/pkg/sb", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/pkg/sb", projectPath, "supabase.go", "pkg/sb")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/pkg/util", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/pkg/util", projectPath, "util.go", "pkg/util")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/view", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view", projectPath, "util.go", "view/util")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/view/auth", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view/auth", projectPath, "auth.templ", "view/auth")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/view/css", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view/css", projectPath, "app.css", "view/css")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/view/home", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view/home", projectPath, "index.templ", "view/home")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/view/layout", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view/layout", projectPath, "app.templ", "view/layout")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/view/settings", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view/settings", projectPath, "account.templ", "view/settings")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath("/view/ui", projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view/ui", projectPath, "navigation.templ", "view/ui/navigation")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection("/view/ui", projectPath, "toast.templ", "view/ui/toast")
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

	err = utils.ExecuteCmd("npm", []string{"install", "-D", "tailwindcss"}, projectPath)
	if err != nil {
		log.Printf("Error initializing git repo: %v", err)
		cobra.CheckErr(err)
		return err
	}

	err = utils.ExecuteCmd("npm", []string{"install", "-D", "daisyui@latest"}, projectPath)
	if err != nil {
		log.Printf("Error initializing git repo: %v", err)
		cobra.CheckErr(err)
		return err
	}

	tailwindConfigFile, err := os.Create(filepath.Join(projectPath, "tailwind.config.js"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}
	defer tailwindConfigFile.Close()

	err = utils.ExecuteCmd("npx", []string{"tailwind", "init", "-p"}, projectPath)
	if err != nil {
		log.Printf("Error initializing git repo: %v", err)
		cobra.CheckErr(err)
		return err
	}

	tailwindConfigTemplate := template.Must(template.New("tailwind.config.js").Parse(string(myTemplate.TailwindConfigTemplate())))
	err = tailwindConfigTemplate.Execute(tailwindConfigFile, p)
	if err != nil {
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

	makefileFile, err := os.Create(filepath.Join(projectPath, "Makefile"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}
	defer makefileFile.Close()

	makefileTemplate := template.Must(template.New("Makefile").Parse(string(myTemplate.MakefileTemplate())))
	err = makefileTemplate.Execute(makefileFile, p)
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
	case "handler/auth":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.HandlerAuth())))
		err = createdTemplate.Execute(createdFile, p)
	case "handler/home":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.HandlerHome())))
		err = createdTemplate.Execute(createdFile, p)
	case "handler/middleware":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.HandlerMiddleware())))
		err = createdTemplate.Execute(createdFile, p)
	case "handler/settings":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.HandlerSettings())))
		err = createdTemplate.Execute(createdFile, p)
	case "handler/util":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.HandlerUtil())))
		err = createdTemplate.Execute(createdFile, p)
	case "db/db":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.DbTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "db/query":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.DbQueryTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "migrate/migrate":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.CmdMigrateTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "migrate/reset":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.CmdMigrateReset())))
		err = createdTemplate.Execute(createdFile, p)
	case "types/user":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.TypesUserTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "types/account":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.TypesAccountTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "pkg/validate":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ValidateTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "pkg/sb":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.SupabaseTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/util":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewUtilTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/auth":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewAuthTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/css":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewCssTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/home":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewHomeTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/layout":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewLayoutTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/settings":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewAccountTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/ui/navigation":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewNavigationTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	case "view/ui/toast":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewToastTemplate())))
		err = createdTemplate.Execute(createdFile, p)
	}

	if err != nil {
		return err
	}

	return nil
}
