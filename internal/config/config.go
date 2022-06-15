package config

import "flag"

// DataDirectory is the path used for loading templates/database migrations
var DataDirectory = flag.String("data-directory", "", "Path for loading templates and migrations scripts.")