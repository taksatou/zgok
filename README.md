# zgok

Zgok is a simple GO command line toolKit

features:

* very simple cli tool
* generating pretty help
* nested subcommand

# Getting Started

get the code:
```
go get github.com/taksatou/zgok
```

example:
```
package main

import (
	"flag"
	"fmt"
	"mojavy.com/zgok"
	"os"
)

type Echo struct{}

// implement zgok.Command interface
func (e Echo) Name() string          { return "echo" }
func (e Echo) Desc() string          { return "echo the input" }
func (e Echo) Do(s string) error     { fmt.Println(s); return nil }
func (e Echo) SubCommand() *zgok.Cli { return nil }
func (e Echo) Flag() *flag.FlagSet   { return nil }

var (
	cmd *zgok.Cli
)

func init() {
	// register echo
	cmd = zgok.NewCli("example", "This is an example command")
	cmd.Register(&Echo{})
}

func main() {
	// run with arguments
	e := cmd.Run(os.Args[1:])
	if e != nil {
		os.Exit(1)
	}
}
```

 
```
$ ./example
example:

  This is an example command
  
  Commands:
  
    echo        echo the input
```

