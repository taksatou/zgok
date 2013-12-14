package main

import (
	"bufio"
	"flag"
	"fmt"
	"mojavy.com/zgok"
	"os"
	"regexp"
)

type Base struct{}

func (b Base) SubCommand() *zgok.Cli { return nil }
func (b Base) Flag() *flag.FlagSet   { return nil }
func (b Base) Do(s string) error     { return nil }

// simple example
type Echo struct {
	Base
}

func (e Echo) Name() string      { return "echo" }
func (e Echo) Desc() string      { return "echo the input" }
func (e Echo) Do(s string) error { fmt.Println(s); return nil }

// example of flag
type Grep struct {
	Base
	f *flag.FlagSet

	lineNumber  bool
	invertMatch bool
}

func (g *Grep) Name() string { return "grep" }
func (g *Grep) Desc() string { return "simple grep" }
func (g *Grep) Flag() *flag.FlagSet {
	if g.f == nil {
		g.f = flag.NewFlagSet("grep", flag.ExitOnError)
		g.f.BoolVar(&g.lineNumber, "n", false, "display line number")
		g.f.BoolVar(&g.invertMatch, "v", false, "invert the sense of matching")
	}
	return g.f
}

func (g *Grep) fmt(i int, s string) {
	if g.lineNumber {
		fmt.Printf("%d:%s\n", i, s)
	} else {
		fmt.Println(s)
	}
}

func (g *Grep) Do(s string) error {
	r := regexp.MustCompile(s)

	for i, s := 1, bufio.NewScanner(os.Stdin); s.Scan(); i++ {
		if r.Match(s.Bytes()) {
			if !g.invertMatch {
				g.fmt(i, s.Text())
			}
		} else {
			if g.invertMatch {
				g.fmt(i, s.Text())
			}
		}
	}
	return nil
}

// example of nested subcommand
// nested subcommand need to implement SubCommand method
type Nested struct {
	Base
}

func (n Nested) Name() string { return "nested" }
func (n Nested) Desc() string { return "example for nested commands" }
func (n Nested) SubCommand() *zgok.Cli {
	cmd := zgok.NewCli("nexted example", "This is nested sub commands example")
	cmd.Register(&Echo{})
	cmd.Register(&Grep{})
	cmd.Register(&Nested{})
	return cmd
}

var (
	cmd *zgok.Cli
)

func init() {
	cmd = zgok.NewCli("example", "This is an example command")
	cmd.Register(&Echo{})
	cmd.Register(&Grep{})
	cmd.Register(&Nested{})
}

func main() {
	e := cmd.Run(os.Args[1:])
	if e != nil {
		os.Exit(1)
	}
}
