package delete

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
	Use:          "deletetrans",
	Short:        "Delete Transactions database",
	Example:      "{execfile} deletetrans -c config/config.yaml",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		lib.SetConfigPath(configFile)
	},
	Run: func(cmd *cobra.Command, args []string) {
		config := lib.NewConfig()
		logger := lib.NewLogger(config)
		db := lib.NewDatabase(config, logger)

		db.ORM.Find(&models.BKKDetail{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.BKKDetail{})

		db.ORM.Find(&models.BKKHeader{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.BKKHeader{})

		db.ORM.Find(&models.SaldoHistory{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.SaldoHistory{})

		db.ORM.Find(&models.SaldoMonth{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.SaldoMonth{})

		db.ORM.Find(&models.Saldo{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.Saldo{})

		db.ORM.Find(&models.InvoiceDetail{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.InvoiceDetail{})

		db.ORM.Find(&models.InvoiceHeader{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.InvoiceHeader{})

		db.ORM.Find(&models.TarikDana{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.TarikDana{})

		db.ORM.Find(&models.Kasbon{})
		db.ORM.Unscoped().Where("1 = 1").Delete(&models.Kasbon{})
	},
}
