# cmdflag : command support to the standard library flag package

## Overview [![GoDoc](https://godoc.org/github.com/pierrec/cmdflag?status.svg)](https://godoc.org/github.com/pierrec/cmdflag) [![Go Report Card](https://goreportcard.com/badge/github.com/pierrec/cmdflag)](https://goreportcard.com/report/github.com/pierrec/cmdflag)

Building on top of the excellent `flag`  package from the standard library, `cmdflag` adds a simple way of specifying nested commands. Its intent is to blend with the usage of `flag` by keeping its idioms and simply augment it with commands.

## Install

```
go get github.com/pierrec/cmdflag
```

## Usage

A `Command` is a set of flags (flag.FlagSet) with a (potentially empty) set of other Commands 
(sometimes referred to as subcommands). To define a subcommand, `Command.Add` an `Application` 
defining its properties:
  - Name - the subcommand name
  - Descr - a short desciption of the subcommand
  - Args - the list of arguments expected by the subcommand
  - Help - a long description of the subcommand
  - Err - what to do in case of error (same as in the flag package)
  - Init - the function to be run once the subcommand is encountered (lazily initialized)
  
Nested commands are supported, so a subcommand can also have its own subcommands.

## Example

```
./chef cook -toppings cheese,mushrooms pizza

c := cmdflag.New(nil)
c.Add(cmdflag.Application{
    Name: "cook",
    Init: func(fs *flag.FlagSet) cmdflag.Handler {
        var toppings string
        fs.StringVar(&toppings, "toppings", "", "coma separated list of toppings")
        
        return func(args ...string) (int, error) {
            fmt.Printf("courses: %v\n", args)
            fmt.Printf("toppings: %v\n", strings.Split(toppings, ","))
            return len(args), nil
        }
    }
}
```

## Extra features

To make life easier, a few common uses of a command line library are provided and can be easily activated: 
  - flags:
    defining them as a flag activates them
    - -version - see `VersionBoolFlag`
    - -fullversion - see `FullVersionBoolFlag`
    - the standard -h and -help flags are supported to display the usage of the command they apply to
  - commands:
    - help - provides a way to display `Application.Help` for a given command
      (activated by `Command.AddHelp`)

## Contributing

Contributions welcome via pull requests. Please provide tests.
