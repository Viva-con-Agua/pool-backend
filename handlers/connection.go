package handlers

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcago/vmdb"
)

var Database = vmdb.NewDatabase(
	"pool-user",
	vcago.Settings.String("DB_HOST", "w", "localhost"),
	vcago.Settings.String("DB_PORT", "w", "27017"),
)
