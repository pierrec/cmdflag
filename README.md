# cmdflag : command support to the standard library flag package

## Overview [![GoDoc](https://godoc.org/github.com/pierrec/cmdflag?status.svg)](https://godoc.org/github.com/pierrec/cmdflag) [![Go Report Card](https://goreportcard.com/badge/github.com/pierrec/cmdflag)](https://goreportcard.com/report/github.com/pierrec/cmdflag)

Building on top of the excellent `flag`  package from the standard library, `cmdflag` adds a simple way of specifying nested commands. Its intent is to blend with the usage of `flag` by keeping its idioms and simply augment it with commands.

## Install

```
go get github.com/pierrec/cmdflag
```

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
            return nil
        }
    }
}
```

## Contributing

Contributions welcome via pull requests. Please provide tests.

## License

MIT.

