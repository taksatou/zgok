package zgok

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type Command interface {
	// Sub-command name. Name must return one word, otherwise it goes undefined.
	Name() string

	// Short description for this command
	Desc() string

	// If you want to use flags for this command, return initialized FlagSet object,
	// and then remaining arguments are passed to it.
	// Otherwise return nil
	Flag() *flag.FlagSet

	// If you want to use nested sub command, return initialized Cli object.
	// Otherwise return nil
	SubCommand() *Cli

	// Execute the body with remaining comman line argument
	Do(string) error
}

type Cli struct {
	Name     string
	Doc      string
	Commands map[string]Command
}

func NewCli(name, doc string) *Cli {
	return &Cli{name, doc, make(map[string]Command)}
}

func (cli *Cli) Register(c Command) {
	cli.Commands[c.Name()] = c
}

var tpl = `{{.Name}}:

{{.Doc | indent}}

Commands:
{{range .Commands}}
  {{.Name | printf "%-11s"}} {{.Desc}}{{end}}

`

func (cli *Cli) PrintHelp() {
	t := template.New("help")
	r := regexp.MustCompile("^")
	t.Funcs(template.FuncMap{
		"indent": func(s string) string { return string(r.ReplaceAll([]byte(s), []byte("  "))) },
	})

	template.Must(t.Parse(tpl))
	t.Execute(os.Stdout, cli)
}

func (cli *Cli) Run(args []string) error {
	if len(args) < 1 {
		cli.PrintHelp()
		return errors.New("no sub command specified")
	}

	if c, ok := cli.Commands[args[0]]; ok {
		c2 := c.SubCommand()
		if c2 != nil {
			return c2.Run(args[1:])
		} else {
			f := c.Flag()
			if f != nil {
				e := f.Parse(args[1:])
				if e != nil {
					c.Flag().Usage()
					return e
				}
				return c.Do(f.Arg(0))
			} else {
				return c.Do(strings.Join(args[1:], " "))
			}
		}
	} else {
		cli.PrintHelp()
		return errors.New(fmt.Sprintf("unknown command: %s", args[0]))
	}

	return nil
}
