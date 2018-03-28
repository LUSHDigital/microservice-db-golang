package migrate

import (
	"testing"
	"os"
	"fmt"
	"io/ioutil"
)

func Test_migrationsInPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			"passing",
			"{CALCULATED}",
			true,
		},
		{
			"failing - no path",
			"",
			false,
		},
		{
			"failing - invalid path",
			"/etc/ihope/i/dont/exist	",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Path setup for passing tests
			if tt.path == "{CALCULATED}" {
				var (
					dir  = os.TempDir()
					file = fmt.Sprintf("%s/temp.sql", dir)
				)
				defer os.RemoveAll(dir)

				if err := ioutil.WriteFile(file, nil, 0644); err != nil {
					t.Fatal(err)
				}

				// Update path
				tt.path = dir
			}

			exists := migrationsInPath(tt.path)
			if tt.expected != exists {
				t.Errorf("Expected %s but got %d", tt.expected, exists)
			}
		})
	}
}
