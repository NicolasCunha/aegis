package main

import (
	migrations "nfcunha/aegis/database"
	api "nfcunha/aegis/api"
)

func main() {
	migrations.Migrate()
	api.RegisterApis()
}