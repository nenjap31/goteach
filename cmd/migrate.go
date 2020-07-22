package cmd

import (
	"database/sql"
	"goteach/config"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.SetHelpFunc(func(*cobra.Command, []string) {
		log.Print(usageCommands)
	})
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
	Run:   migrateHandler,
}

var migrateHandler = func(cmd *cobra.Command, args []string) {
	mysql := config.LoadDBConfig("mysql")
	goose.SetDialect("mysql")
	db, err := sql.Open("mysql", mysql.User+":"+mysql.Password+"@tcp("+mysql.Host+":"+strconv.Itoa(mysql.Port)+")/"+mysql.DBName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	dir := config.AppPath + "/database/migration"
	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}
	command := args[0]
	cmdArgs := args[1:]
	if err := goose.Run(command, db, dir, cmdArgs...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

var usageCommands = `
Run database migrations

Usage:
    goteach migrate [command]

Available Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
`
