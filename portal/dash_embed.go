package portal

import "embed"

//go:embed packages/dash/dist/*
var DashFolder embed.FS
