package embedded

import "embed"

//go:embed migrations
var Migrations embed.FS

//go:embed AppVersion.txt
var AppVersion string
