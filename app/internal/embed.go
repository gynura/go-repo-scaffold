package internal

import "embed"

//go:embed config/dbschemas/*
var DatabaseMigrations embed.FS

//go:embed config/mappers/*
var DatabaseMappers embed.FS
