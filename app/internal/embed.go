package app

import "embed"

//go:embed include/dbschemas/*
var DatabaseMigrations embed.FS

//go:embed include/config/*
var DatabaseMappers embed.FS
