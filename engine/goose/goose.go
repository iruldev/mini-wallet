package main

import (
	_ "github.com/iruldev/mini-wallet/src/config"
	_ "github.com/iruldev/mini-wallet/src/schema/migrations"

	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"text/template"

	"github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

var (
	flags        = flag.NewFlagSet("goose", flag.ExitOnError)
	dir          = flags.String("dir", defaultMigrationDir, "migrations directory")
	table        = flags.String("table", "goose_db_version", "migrations table name")
	verbose      = flags.Bool("v", false, "enable verbose mode")
	help         = flags.Bool("h", false, "print help")
	version      = flags.Bool("version", false, "print version")
	certfile     = flags.String("certfile", "", "file path to root CA's certificates in pem format (only support on mysql)")
	sequential   = flags.Bool("s", false, "use sequential numbering for new migrations")
	allowMissing = flags.Bool("allow-missing", false, "applies missing (out-of-order) migrations")
	sslcert      = flags.String("ssl-cert", "", "file path to SSL certificates in pem format (only support on mysql)")
	sslkey       = flags.String("ssl-key", "", "file path to SSL key in pem format (only support on mysql)")
	noVersioning = flags.Bool("no-versioning", false, "apply migration commands with no versioning, in file order, from directory pointed to")
	schema       = flags.String("schema", "migration", "for types (migration) or (seeder)")
)
var (
	gooseVersion = "Goose 1.0.0"
	migration    = "migration"
	seed         = "seeder"
)

type stdLogger struct{}

func (*stdLogger) Fatal(v ...interface{})                 { fmt.Println(v...) }
func (*stdLogger) Fatalf(format string, v ...interface{}) { fmt.Printf(format, v...) }
func (*stdLogger) Print(v ...interface{})                 { fmt.Print(v...) }
func (*stdLogger) Println(v ...interface{})               { fmt.Println(v...) }
func (*stdLogger) Printf(format string, v ...interface{}) { fmt.Printf(format, v...) }

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	if *version {
		if buildInfo, ok := debug.ReadBuildInfo(); ok && buildInfo != nil && gooseVersion == "" {
			gooseVersion = buildInfo.Main.Version
		}
		fmt.Printf("goose version:%s\n", gooseVersion)
		return
	}
	if *verbose {
		goose.SetVerbose(true)
	}
	if *sequential {
		goose.SetSequential(true)
	}
	if *schema != "" {
		switch *schema {
		case migration:
			fmt.Println(viper.GetString("DB_MIGRATION"))
			*dir = viper.GetString("DB_MIGRATION")
		case seed:
			fmt.Println(viper.GetString("DB_SEED"))
			*noVersioning = true
			*dir = viper.GetString("DB_SEED")
		}
	}
	goose.SetTableName(*table)

	args := flags.Args()
	if len(args) == 0 || *help {
		flags.Usage()
		return
	}
	// The -dir option has not been set, check whether the env variable is set
	// before defaulting to ".".
	if *dir == defaultMigrationDir && os.Getenv(envGooseMigrationDir) != "" {
		*dir = viper.GetString("DB_MIGRATION")
	}

	switch args[0] {
	case "init":
		if err := gooseInit(*dir); err != nil {
			fmt.Printf("goose run: %v", err)
		}
		return
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			fmt.Printf("goose run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			fmt.Printf("goose run: %v", err)
		}
		return
	}

	args = mergeArgs(args)
	if len(args) < 1 {
		flags.Usage()
		return
	}

	goose.SetLogger(&stdLogger{})

	command := args[0]
	dbstring := viper.GetString("DB_USER") + ":" + viper.GetString("DB_PASS") + "@tcp(" + viper.GetString("DB_HOST") + ":" + viper.GetString("DB_PORT") + ")/" + viper.GetString("DB_NAME")

	if *certfile == "" {
		cf := viper.GetString("CA_CERT")
		certfile = &cf
	}
	if *sslcert == "" {
		clCert := viper.GetString("CLIENT_CERT")
		sslcert = &clCert
	}
	if *sslkey == "" {
		clKey := viper.GetString("CLIENT_KEY")
		sslkey = &clKey
	}

	db, err := goose.OpenDBWithDriver("mysql", normalizeDBString("mysql", dbstring, *certfile, *sslcert, *sslkey))
	if err != nil {
		fmt.Printf("-dbstring %v, %v\n", dbstring, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	options := []goose.OptionsFunc{}
	if *allowMissing {
		options = append(options, goose.WithAllowMissing())
	}
	if *noVersioning {
		options = append(options, goose.WithNoVersioning())
	}
	log.Println(*dir)
	if err := goose.RunWithOptions(
		command,
		db,
		*dir,
		arguments,
		options...,
	); err != nil {
		fmt.Printf("goose run: %v", err)
	}
}

func normalizeDBString(driver string, str string, certfile string, sslcert string, sslkey string) string {
	if driver == "mysql" {
		var isTLS = certfile != ""
		if isTLS {
			if err := registerTLSConfig(certfile, sslcert, sslkey); err != nil {
				fmt.Printf("goose run: %v", err)
			}
		}
		var err error
		str, err = normalizeMySQLDSN(str, isTLS)
		if err != nil {
			fmt.Printf("failed to normalize MySQL connection string: %v", err)
		}
	}
	return str
}

const tlsConfigKey = ""

var tlsReg = regexp.MustCompile(`(\?|&)tls=[^&]*(?:&|$)`)

func normalizeMySQLDSN(dsn string, tls bool) (string, error) {
	// If we are sharing a DSN in a different environment, it may contain a TLS
	// setting key with a value name that is not "custom," so clear it.
	dsn = tlsReg.ReplaceAllString(dsn, `$1`)
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", err
	}
	config.ParseTime = true
	if tls {
		config.TLSConfig = tlsConfigKey
	}
	return config.FormatDSN(), nil
}

func registerTLSConfig(pemfile string, sslcert string, sslkey string) error {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(pemfile)
	if err != nil {
		return err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to append PEM: %q", pemfile)
	}

	tlsConfig := &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}
	if sslcert != "" && sslkey != "" {
		cert, err := tls.LoadX509KeyPair(sslcert, sslkey)
		if err != nil {
			return fmt.Errorf("failed to load x509 keypair: %w", err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	return mysql.RegisterTLSConfig(tlsConfigKey, tlsConfig)
}

const (
	envGooseDriver       = "GOOSE_DRIVER"
	envGooseDBString     = "GOOSE_DBSTRING"
	envGooseMigrationDir = "GOOSE_MIGRATION_DIR"
)

const (
	defaultMigrationDir = "src/schema/migrations"
)

func mergeArgs(args []string) []string {
	if len(args) < 1 {
		return args
	}
	if d := os.Getenv(envGooseDriver); d != "" {
		args = append([]string{d}, args...)
	}
	if d := os.Getenv(envGooseDBString); d != "" {
		args = append([]string{args[0], d}, args[1:]...)
	}
	return args
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND
or
Set environment key
GOOSE_DRIVER=DRIVER
GOOSE_DBSTRING=DBSTRING
Usage: goose [OPTIONS] COMMAND
Examples:
    goose mysql "user:password@/dbname?parseTime=true" status
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)

var sqlMigrationTemplate = template.Must(template.New("goose.sql-migration").Parse(`-- Thank you for giving goose a try!
--
-- This file was automatically created running goose init. If you're familiar with goose
-- feel free to remove/rename this file, write some SQL and goose up. Briefly,
--
-- Documentation can be found here: https://pressly.github.io/goose
--
-- A single goose .sql file holds both Up and Down migrations.
--
-- All goose .sql files are expected to have a -- +goose Up directive.
-- The -- +goose Down directive is optional, but recommended, and must come after the Up directive.
--
-- The -- +goose NO TRANSACTION directive may be added to the top of the file to run statements
-- outside a transaction. Both Up and Down migrations within this file will be run without a transaction.
--
-- More complex statements that have semicolons within them must be annotated with
-- the -- +goose StatementBegin and -- +goose StatementEnd directives to be properly recognized.
--
-- Use GitHub issues for reporting bugs and requesting features, enjoy!
-- +goose Up
SELECT 'up SQL query';
-- +goose Down
SELECT 'down SQL query';
`))

// initDir will create a directory with an empty SQL migration file.
func gooseInit(dir string) error {
	if dir == "" || dir == defaultMigrationDir {
		dir = "migrations"
	}
	_, err := os.Stat(dir)
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err == nil, errors.Is(err, fs.ErrExist):
		return fmt.Errorf("directory already exists: %s", dir)
	default:
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return goose.CreateWithTemplate(nil, dir, sqlMigrationTemplate, "initial", "sql")
}
