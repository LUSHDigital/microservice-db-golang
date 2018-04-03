package migrate

import (
	"strings"
	"strconv"
	"errors"
	"io/ioutil"
	"path/filepath"
)

// Migrator defines database-specific migrators ready to run migrations
// against a single database.
type Migrator interface {
	Migrate() error
}

// ConnectionOptions define a collection of information that's required to
// connect to any database in order to be able to run migrations.
// This will be used by all `Migrator`s, either embedded in database-specific
// options, or wholly.
type ConnectionOptions struct {
	Host     string
	Port     int
	User     string
	Pass     string
	Database string

	protocol string
}

// ConnectionString concatenates the CockroachConnectionOptions down to a string,
// applying defaults to options that were not set ready to be used in a
// connection to the database
func (co *ConnectionOptions) String() (string, error) {
	// Start a string builder for the connection string.
	var conn strings.Builder

	// Ensure a protocol exists. If it's gibberish it is not of our concern,
	// as an error will be returned later when trying to connect via a
	// gibberish protocol.
	if len(co.protocol) == 0 {
		return "", errors.New("microservice-db: unrecognised protocol for database connection")
	}

	// Write the protocol
	conn.WriteString(co.protocol)
	conn.WriteString("://")

	// User
	if len(co.User) != 0 {
		conn.WriteString(co.User)
	} else {
		conn.WriteString(DefaultUser)
	}

	conn.WriteString("@")

	// Host
	if len(co.Host) != 0 {
		conn.WriteString(co.Host)
	} else {
		conn.WriteString(DefaultHost)
	}

	conn.WriteString(":")

	// Port
	if co.Port != 0 {
		conn.WriteString(strconv.Itoa(co.Port))
	} else {
		conn.WriteString(DefaultPort)
	}

	conn.WriteString("/")

	if len(co.Database) != 0 {
		conn.WriteString(co.Database)
	} else {
		conn.WriteString(DefaultDatabase)
	}

	return conn.String(), nil
}

// migrationsInPath will traverse a path and look for migration files ending
// in '.sql', if none are found we assume there are no migrations in the path
func migrationsInPath(path string) bool {
	// Ensure we have a migration path
	if len(path) == 0 {
		return false
	}

	// Pull all files up from the migrations path
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}

	// If one sql file exists we can consider this a valid migration directory.
	for _, file := range files {
		if file.Mode().IsRegular() &&
			filepath.Ext(file.Name()) == ".sql" {
			return true
		}
	}

	// If no SQL files were found whilst traversing a path - then none exists!
	return false
}
