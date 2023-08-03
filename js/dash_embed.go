package dash_embed

import "embed"

//go:embed packages/dash/dist/*
var DashFolder embed.FS

//go:generate pnpm --filter "dash" build
