package portal

import "embed"

//go:embed packages/dash/out/*
var DashFolder embed.FS
