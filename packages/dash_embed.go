package dash_embed

import "embed"

//go:embed dash/dist/*
var DashFolder embed.FS

//go:generate pnpm --filter "dash" build
