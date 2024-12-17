package assets

import "embed"

//go:embed files
var filesFs embed.FS

//go:embed files/migrations
var migrationFs embed.FS

// Files returns an instance of embedded files.
func Files() embed.FS {
	return filesFs
}

func MigrationFiles() embed.FS { return migrationFs }
