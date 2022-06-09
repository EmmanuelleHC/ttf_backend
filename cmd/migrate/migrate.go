package migrate

import (
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/spf13/cobra"
)

var configFile string

func init() {
	pf := StartCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c",
		"config/config.yaml", "this parameter is used to start the service application")

	cobra.MarkFlagRequired(pf, "config")
}

var StartCmd = &cobra.Command{
	Use:          "migrate",
	Short:        "Migrate database",
	Example:      "{execfile} migrate -c config/config.yaml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		lib.SetConfigPath(configFile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		config := lib.NewConfig()
		logger := lib.NewLogger(config)
		db := lib.NewDatabase(config, logger)

		if err := db.ORM.AutoMigrate(
			&models.User{},
			&models.UserRole{},
			&models.Role{},
			&models.RoleMenu{},
			&models.Menu{},
			&models.MenuAction{},
			&models.MenuActionResource{},
			&models.Company{},
			&models.Branch{},
			&models.Department{},
			&models.Account{},
			&models.CostCentre{},
			&models.Saldo{},
			&models.Trx{},
			&models.Kasbon{},
			&models.Counter{},
			&models.BKKHeader{},
			&models.BKKDetail{},
			&models.Employee{},
			&models.InvoiceHeader{},
			&models.InvoiceDetail{},
			&models.SaldoHistory{},
			&models.SaldoMonth{},
			&models.TarikDana{},
		); err != nil {
			logger.Zap.Fatalf("Error to migrate database: %v", err)
		}
	},
}
