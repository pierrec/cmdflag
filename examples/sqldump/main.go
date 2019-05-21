// This example shows how to create a simple cli tool to export data from an SQL database.
// It connects to the database and runs a SELECT query and displays the fetched rows.
package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pierrec/cmdflag"
	"github.com/xo/dburl"
)

var cli = cmdflag.New(nil)
var db *sql.DB

func init() {
	cli.AddHelp()
	connect := cli.MustAdd(cmdflag.Application{
		Name:  "connect",
		Descr: "connect to an SQL database",
		Help: `The database url format is:
protocol+transport://user:pass@host/dbname?option1=a&option2=b

protocol  - driver name or alias (see below)
transport - "tcp", "udp", "unix" or driver name (odbc/oleodbc)                                  |
user      - username
pass      - password
host      - host
dbname*   - database, instance, or service name/id to connect to
?opt1=... - additional database driver options
              (see respective SQL driver for available options)

(for more information, see https://godoc.org/github.com/xo/dburl)`,
		Args: "<database url>",
		Init: func(fset *flag.FlagSet) cmdflag.Initializer {
			return func(args ...string) (int, error) {
				if len(args) == 0 {
					return 0, fmt.Errorf("missing database url")
				}
				var err error
				db, err = dburl.Open(args[0])
				if err != nil {
					return 1, err
				}
				return 1, db.Ping()
			}
		},
	})
	connect.MustAdd(cmdflag.Application{
		Name:  "export",
		Descr: "export table rows",
		Args:  "table",
		Init: func(fset *flag.FlagSet) cmdflag.Initializer {
			var output string
			fset.StringVar(&output, "o", "", "output file name (default=standard output)")
			var columns string
			fset.StringVar(&columns, "select", "*", "columns to be selected")
			var filter string
			fset.StringVar(&filter, "filter", "true", "rows to be selected")

			return func(args ...string) (int, error) {
				if len(args) == 0 {
					return 0, fmt.Errorf("no table specified")
				}
				if filter == "" {
					filter = "true"
				}
				out := os.Stdout
				if output != "" {
					f, err := os.Create(output)
					if err != nil {
						return 0, err
					}
					defer f.Close()
					out = f
				}
				buf := bufio.NewWriter(out)
				defer buf.Flush()

				q := fmt.Sprintf("select %s from %s where %s", columns, args[0], filter)
				rows, err := db.Query(q)
				if err != nil {
					return 1, err
				}
				defer rows.Close()

				cols, err := rows.Columns()
				if err != nil {
					return 1, err
				}
				_, err = fmt.Fprintln(buf, cols)
				if err != nil {
					return 1, err
				}
				res := make([]interface{}, len(cols))
				res2print := make([]string, len(cols))
				for i := range cols {
					res[i] = &res2print[i]
				}
				for rows.Next() {
					err := rows.Scan(res...)
					if err != nil {
						return 1, err
					}
					_, err = fmt.Fprintln(buf, res2print)
					if err != nil {
						return 1, err
					}
				}
				return 1, rows.Err()
			}
		},
	})
}

func main() {
	if err := cli.Parse(); err != nil {
		log.Fatal(err)
	}
}
