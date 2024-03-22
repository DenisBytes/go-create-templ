package program

import (
	"log"
	"os"

	"github.com/DenisBytes/go-create-app/cmd/flags"
	tea "github.com/charmbracelet/bubbletea"
)

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
	Server() []byte
	Routes() []byte
	TestHandler() []byte
	HtmxTemplRoutes() []byte
	HtmxTemplImports() []byte
	WebsocketImports() []byte
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
