package program

import (
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
	HandlerApp() []byte
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
	// ginPackage   = []string{"github.com/gin-gonic/gin"}
	// fiberPackage = []string{"github.com/gofiber/fiber/v2"}
	echoPackage  = []string{"github.com/labstack/echo/v4", "github.com/labstack/echo/v4/middleware"}

	godotenvPackage             = []string{"github.com/joho/godotenv"}
	golangMigratePackage        = []string{"github.com/golang-migrate/migrate/v4"}
	golangMigrateInstallPackage = []string{"github.com/golang-migrate/migrate/v4/cmd/migrate@latest"}
	templInstallPackage         = []string{"github.com/a-h/templ/cmd/templ@latest"}
	uuidPackage                 = []string{"github.com/google/uuid"}
	sessionPackage              = []string{"github.com/gorilla/sessions"}
	postgresDriverPackage       = []string{"github.com/lib/pq"}
	supabasePackage             = []string{"github.com/nedpals/supabase-go"}
	bunPackage1                 = []string{"github.com/uptrace/bun"}
	bunPackage2                 = []string{"github.com/uptrace/bun/dialect/pgdialect"}
	bunPackage3                 = []string{"github.com/uptrace/bun/extra/bundebug"}
	airPackage                  = []string{"github.com/cosmtrek/air@latest"}
)

const (
	root             = "/"
	handlerPath      = "/handler"
	dbPath           = "/db"
	migratePath      = "/cmd/migrate"
	migrateResetPath = "/cmd/reset"
	typesPath        = "/types"
	validatePath     = "/pkg/validate"
	sbPath           = "/pkg/sb"
	pkgUtilPath      = "/pkg/util"
	viewPath         = "/view"
	viewAppPath      = "/view/app"
	viewAuthPath     = "/view/auth"
	viewCssPath      = "/view/css"
	viewHomePath     = "/view/home"
	viewLayoutPath   = "/view/layout"
	viewSettingsPath = "/view/settings"
	viewUIPath       = "/view/ui"
)

func (p *Project) createFrameworkMap() {

	if p.FrameworkMap == nil {
		p.FrameworkMap = make(map[flags.Framework]Framework)
	}

	p.FrameworkMap[flags.Chi] = Framework{
		packageName: chiPackage,
		templater:   myTemplate.ChiTemplates{},
	}

	p.FrameworkMap[flags.Echo] = Framework{
		packageName: echoPackage,
		templater:   myTemplate.EchoTemplates{},
	}

	p.FrameworkMap[flags.StandardLibrary] = Framework{
		packageName: []string{},
		templater:   myTemplate.HttpstdTemplates{},
	}

	// p.FrameworkMap[flags.Gin] = Framework{
	// 	packageName: ginPackage,
	// 	templater:   myTemplate.GinTemplates{},
	// }

	// p.FrameworkMap[flags.Fiber] = Framework{
	// 	packageName: fiberPackage,
	// 	templater:   myTemplate.FiberTemplates{},
	// }
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
		panic("\nGIT CONFIG ISSUE: user.name is not set in git config.\n")
	}

	// checking git config --get email
	emailSet, err := utils.CheckGitConfig("user.email")
	if err != nil {
		cobra.CheckErr(err)
	}
	if !emailSet {
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
			log.Printf("Could not get go dependency for the chosen framework %v\n", err)
			cobra.CheckErr(err)
		}
	}

	err = p.CreateFileWithInjection(root, projectPath, "main.go", "main")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(handlerPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(handlerPath, projectPath, "auth.go", "handler/auth")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(handlerPath, projectPath, "home.go", "handler/home")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(handlerPath, projectPath, "app.go", "handler/app")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(handlerPath, projectPath, "middleware.go", "handler/middleware")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(handlerPath, projectPath, "settings.go", "handler/settings")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(handlerPath, projectPath, "util.go", "handler/util")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(dbPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(dbPath, projectPath, "db.go", "db/db")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(dbPath, projectPath, "query.go", "db/query")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(migratePath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(migratePath, projectPath, "main.go", "migrate/migrate")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(migrateResetPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(migrateResetPath, projectPath, "main.go", "migrate/reset")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(typesPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(typesPath, projectPath, "user.go", "types/user")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(typesPath, projectPath, "account.go", "types/account")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(validatePath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(validatePath, projectPath, "validate.go", "pkg/validate")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(sbPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(sbPath, projectPath, "supabase.go", "pkg/sb")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(pkgUtilPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(pkgUtilPath, projectPath, "util.go", "pkg/util")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewPath, projectPath, "util.go", "view/util")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewAuthPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewAuthPath, projectPath, "auth.templ", "view/auth")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewCssPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewCssPath, projectPath, "app.css", "view/css")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewHomePath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewHomePath, projectPath, "index.templ", "view/home")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewAppPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewAppPath, projectPath, "index.templ", "view/app")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewLayoutPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewLayoutPath, projectPath, "app.templ", "view/layout")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewSettingsPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewSettingsPath, projectPath, "account.templ", "view/settings")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreatePath(viewUIPath, projectPath)
	if err != nil {
		log.Printf("Error creating path: %s", projectPath)
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewUIPath, projectPath, "navigation.templ", "view/ui/navigation")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	err = p.CreateFileWithInjection(viewUIPath, projectPath, "toast.templ", "view/ui/toast")
	if err != nil {
		cobra.CheckErr(err)
		return err
	}

	// importing godotenv package (go get <godotenvPackage>)
	err = utils.GoGetPackage(projectPath, godotenvPackage)
	if err != nil {
		log.Printf("Could not get go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoInstallPackage(projectPath, templInstallPackage)
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
		log.Printf("Could not get go dependency %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.GoInstallPackage(projectPath, golangMigrateInstallPackage)
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
		log.Printf("Could not get go dependency %v\n", err)
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

	err = utils.ExecuteCmd("templ", []string{"generate", "view"}, projectPath)
	if err != nil {
		log.Printf("Error building templ view: %v", err)
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

	envFile, err := os.Create(filepath.Join(projectPath, ".env"))
	if err != nil {
		cobra.CheckErr(err)
		return err
	}
	defer envFile.Close()

	envFileTemplate := template.Must(template.New(".env").Parse(string(myTemplate.EnvTemplate())))
	err = envFileTemplate.Execute(envFile, p)
	if err != nil {
		return err
	}

	err = utils.GoTidy(projectPath)
	if err != nil {
		log.Printf("Could not go tidy in new project %v\n", err)
		cobra.CheckErr(err)
	}

	err = utils.ExecuteCmd("npx", []string{"tailwindcss", "-i", "view/css/app.css", "-o", "public/styles.css"}, projectPath)
	if err != nil {
		log.Printf("Error adding files to git repo: %v", err)
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

	err = utils.ExecuteCmd("go", []string{"get", "-u", "github.com/a-h/templ"}, projectPath)
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
	case "handler/app":
		createdTemplate := template.Must(template.New(fileName).Parse(string(p.FrameworkMap[p.ProjectType].templater.HandlerApp())))
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
	case "pkg/util":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.PkgUtilTemplate())))
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
	case "view/app":
		createdTemplate := template.Must(template.New(fileName).Parse(string(myTemplate.ViewAppTemplate())))
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
