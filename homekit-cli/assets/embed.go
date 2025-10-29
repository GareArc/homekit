package assets

import "embed"

// FS bundles scripts and templates that ship with the CLI.
//
//go:embed scripts/* templates/*
var FS embed.FS
