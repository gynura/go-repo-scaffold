package internal

import "embed"

//go:embed config/mappers/*
var DatabaseMappers embed.FS
