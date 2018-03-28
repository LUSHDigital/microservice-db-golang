package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/cockroachdb"
	// file is imported for its side-effect of loading migration files from
	// disk when specifying the migrations directory
	_ "github.com/mattes/migrate/source/file"

	// pq is imported for its  side-effect of using postgresql:// the sql
	// connection string
	_ "github.com/lib/pq"
)

// Defaults for cockroachdb connection string
const (
	DefaultHost           = "localhost"
	DefaultUser           = "root"
	DefaultPort           = "26257"
	DefaultDatabase       = "service"
	DefaultMigrationTable = "schema_migrations"
)

// Cockroach wraps a database connection along with the necessary fields
// required to run migrations against said database.
type Cockroach struct {
	db *sql.DB

	MigrationsPath string
}

// CockroachConnectionOptions expands on the standard ConnectionOptions and
// adds support for cockroachdb's secure connection along with a specified
// table to store migration information under
type CockroachConnectionOptions struct {
	*ConnectionOptions

	// MigrationTable optionally specifies the table name to store migration
	// history in. If not specified it'll default to `schema_migrations`.
	MigrationTable string
	// Secure informs the migrator that the connection is secured via SSL,
	// and that it should use the certificates from `SSL` struct paths when
	// verifying the identity and connecting to the database.
	Secure bool
	// SSL is used in conjunction with `Secure`.
	SSL *CockroachSSL
}

// CockroachSSL contains paths to the SSL keys on the disk for making a
// connection to a secured database.
type CockroachSSL struct {
	CertPath string
	KeyPath  string
	Mode     string
	RootCert string
}

// ConnectionString concatenates the CockroachConnectionOptions down to a string,
// applying defaults to options that were not set ready to be used in a
// connection to the database
func (co *CockroachConnectionOptions) ConnectionString() (string, error) {
	// Prevent panics and just return the exact default upon nil options
	if co.ConnectionOptions == nil {
		return fmt.Sprintf(
			"postgresql://%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
			DefaultUser, DefaultHost, DefaultPort, DefaultDatabase, DefaultMigrationTable,
		), nil
	}

	// Set the correct protocol used by cockroachdb
	co.ConnectionOptions.protocol = "postgresql"

	// Build a connection string from the underlying connection options ready
	// to apply cockroachdb specific options to
	prefab, err := co.ConnectionOptions.String()
	if err != nil {
		return "", err
	}

	// Create a string builder and write the prefabricated connection string
	// to it
	var conn strings.Builder
	conn.WriteString(prefab)

	// Set sslmode and the paths to the keys if secure or disable sslmode if
	// not secure
	if co.Secure {
		conn.WriteString("?sslmode=enable")
		conn.WriteString(fmt.Sprintf("&sslcert=%s", co.SSL.CertPath))
		conn.WriteString(fmt.Sprintf("&sslkey=%s", co.SSL.KeyPath))
		conn.WriteString(fmt.Sprintf("&sslmode=%s", co.SSL.Mode))
		conn.WriteString(fmt.Sprintf("&sslrootcert=%s", co.SSL.RootCert))
	} else {
		conn.WriteString("?sslmode=disable")
	}

	// Optional migrations table name used by github.com/mattes/migrate
	if len(co.MigrationTable) != 0 {
		conn.WriteString(fmt.Sprintf("&x-migrations-table=%s", co.MigrationTable))
	} else {
		conn.WriteString(fmt.Sprintf("&x-migrations-table=%s", DefaultMigrationTable))
	}

	return conn.String(), nil
}

// NewCockroach provides a Migrator to be ran against a CockroachDB databases
func NewCockroach(path string, opts *CockroachConnectionOptions) (*Cockroach, error) {
	// Validate the migrations path
	if !migrationsInPath(path) {
		return nil, errors.New("migrate: no migration files in given path")
	}

	// Construct connection string from options/defaults
	conn, err := opts.ConnectionString()
	if err != nil {
		return nil, err
	}

	// Open an SQL conneciton to cockroach via the postgresql protocol
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &Cockroach{
		db:             db,
		MigrationsPath: path,
	}, nil
}

// Migrate performs the migrations within the given migration directory for
// CockroachDB.
func (c *Cockroach) Migrate() error {
	// Use the existing database instance when setting up the driver for
	// migrations
	driver, err := cockroachdb.WithInstance(c.db, &cockroachdb.Config{})
	if err != nil {
		log.Fatalf("could not get migrations driverName: %s", err)
	}

	// Stage the migrations ready to bring up
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", c.MigrationsPath),
		"sql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migrate: could not initialise migrations: %s", err)
	}

	// Bring the migrations up against the database
	return m.Up()
}
