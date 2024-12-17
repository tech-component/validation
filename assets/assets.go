package assets

import "embed"

//go:embed files/migrations
var migrationFs embed.FS

func MigrationFiles() embed.FS { return migrationFs }
