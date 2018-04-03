package migrate

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestCockroachOptions_ConnectionString(t *testing.T) {
	type fields struct {
		ConnectionOptions *ConnectionOptions
		MigrationTable    string
		Secure            bool
		SSL               *CockroachSSL
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"Passing - insecure options",
			fields{
				&ConnectionOptions{
					User:     "test-user",
					Port:     9001,
					Host:     "test-host",
					Database: "test-database",
				},
				"migration_table",
				false,
				nil,
			},
			"postgresql://test-user@test-host:9001/test-database?sslmode=disable&x-migrations-table=migration_table",
			false,
		},
		{
			"Passing - defaults",
			fields{},
			"postgresql://root@localhost:26257/service?sslmode" +
				"=disable&x-migrations-table=schema_migrations",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := &CockroachConnectionOptions{
				ConnectionOptions: tt.fields.ConnectionOptions,
				MigrationTable:    tt.fields.MigrationTable,
				Secure:            tt.fields.Secure,
				SSL:               tt.fields.SSL,
			}
			got, err := co.ConnectionString()
			if (err != nil) != tt.wantErr {
				t.Errorf("CockroachConnectionOptions.ConnectionString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CockroachConnectionOptions.ConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
