package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/leekchan/accounting"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// ReportService service layer
type ReportService struct {
	logger                 lib.Logger
	casbinService          CasbinService
	saldoService           SaldoService
	userRepository         repository.UserRepository
	bkkheaderRepository    repository.BKKHeaderRepository
	saldohistoryRepository repository.SaldoHistoryRepository
	saldoRepository        repository.SaldoRepository
	branchRepository       repository.BranchRepository
}

// NewReportService creates a new reportservice
func NewReportService(
	logger lib.Logger,
	casbinService CasbinService,
	saldoService SaldoService,
	userRepository repository.UserRepository,
	bkkheaderRepository repository.BKKHeaderRepository,
	saldohistoryRepository repository.SaldoHistoryRepository,
	saldoRepository repository.SaldoRepository,
	branchRepository repository.BranchRepository,
) ReportService {
	return ReportService{
		logger:                 logger,
		casbinService:          casbinService,
		saldoService:           saldoService,
		userRepository:         userRepository,
		bkkheaderRepository:    bkkheaderRepository,
		saldohistoryRepository: saldohistoryRepository,
		saldoRepository:        saldoRepository,
		branchRepository:       branchRepository,
	}
}

const (
	layoutID     = "02 January 2006"
	layoutDetail = "02-01-2006"
)

func (a ReportService) GenerateLmdp(param *dto.ReportMonitoring) error {
	cbn, err := a.branchRepository.Get(param.BranchID)
	if err != nil {
		return err
	}

	begin := time.Now()

	period1, _ := time.Parse(time.RFC3339, param.DateParams[0])
	period2, _ := time.Parse(time.RFC3339, param.DateParams[1])
	var tglCetak string = "Tgl Cetak: " + begin.Format(layoutID)
	var pukulCetak string = "Pkl Cetak: " + begin.Format("15:04:05")
	var userIdCetak string = "Pkl Cetak: " + param.UserId
	var kodeNama string = "Kode - Nama: " + cbn.Code + " - " + cbn.Name
	var periodeParam string = "Periode: " + period1.Format(layoutID) + " - " + period2.Format(layoutID)

	darkGrayColor := getDarkGrayColor()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeaderLmdp()
	contents, totalA, totalI, totalO, totalAh := getDataSaldo(a, param.CompanyID, param.BranchID, param.DateParams)

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	m.SetAliasNbPages("{nb}")
	m.SetFirstPageNb(1)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				m.Text("Leader Barista Point Coffe. SA", props.Text{
					Size:        8,
					Align:       consts.Left,
					Extrapolate: false,
					Color:       darkGrayColor,
				})
			})

			m.ColSpace(6)

			m.Col(3, func() {
				m.Text(tglCetak, props.Text{
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: darkGrayColor,
				})
				m.Text(pukulCetak, props.Text{
					Top:   3,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: darkGrayColor,
				})
				m.Text(userIdCetak, props.Text{
					Top:   6,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: darkGrayColor,
				})
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Laporan Monitoring Dana Petty Cash", props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(kodeNama, props.Text{
				Top:   1,
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
			})
			m.Text(periodeParam, props.Text{
				Top:   6,
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.Row(2, func() {
		m.ColSpace(7)
		m.Col(5, func() {
			m.Text("Halaman: "+strconv.Itoa(m.GetCurrentPage())+"/{nb}", props.Text{
				Align: consts.Right,
				Size:  8,
			})
		})
	})

	m.Line(10)
	m.SetBackgroundColor(whiteColor)

	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{1, 1, 2, 2, 2, 2, 2},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{1, 1, 2, 2, 2, 2, 2},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   2,
		Line:                 false,
	})

	m.Line(5)

	m.Row(4, func() {
		m.ColSpace(8)
		m.Col(2, func() {
			m.Text("Total Saldo Awal:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(2, func() {
			m.Text(totalA, props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	m.Row(4, func() {
		m.ColSpace(8)
		m.Col(2, func() {
			m.Text("Total Saldo Masuk:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(2, func() {
			m.Text(totalI, props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	m.Row(4, func() {
		m.ColSpace(8)
		m.Col(2, func() {
			m.Text("Total Saldo Keluar:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(2, func() {
			m.Text(totalO, props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	m.Row(4, func() {
		m.ColSpace(8)
		m.Col(2, func() {
			m.Text("Total Saldo Akhir:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(2, func() {
			m.Text(totalAh, props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	m.Line(10)

	m.Row(35, func() {
		m.ColSpace(5)
		m.Col(14, func() {
			_ = m.FileImage("assets/images/ttd.png", props.Rect{
				Percent: 100,
			})
		})
	})

	errFile := m.OutputFileAndClose("pdfs/lmdp.pdf")
	if errFile != nil {
		fmt.Println("Could not save PDF:", errFile)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	return nil
}

func getHeaderLmdp() []string {
	return []string{"No", "Tanggal", "Deskripsi", "Saldo Awal (Rp.)", "Uang Masuk (Rp.)", "Uang Keluar (Rp.)", "Saldo Akhir (Rp.)"}
}

func getDataSaldo(a ReportService, companyId string, branchId string, dateQuery []string) (data [][]string, totalA string, totalI string, totalO string, totalAh string) {
	var result [][]string
	var totalSaldoAwal int64 = 0
	var totalInAmount int64 = 0
	var totalOutAmount int64 = 0
	var totalSaldoAkhir int64 = 0
	t := time.Now()
	firstday := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	dateStart, _ := time.Parse(time.RFC3339, dateQuery[0])

	acc := accounting.Accounting{Precision: 2, Thousand: ".", Decimal: ","}
	qr, err := a.saldohistoryRepository.Query(&models.SaldoHistoryQueryParam{
		CompanyID:  companyId,
		BranchID:   branchId,
		DateQuery:  dateQuery,
		OrderParam: dto.OrderParam{Key: "created_at", Direction: dto.OrderByASC},
	})

	saldoAwal, _ := a.saldoService.GetSaldoAwal(companyId, branchId, dateQuery[0], firstday.String())
	if saldoAwal != 0 {
		awl := []string{strconv.Itoa(1),
			dateStart.Format(layoutDetail),
			"Saldo Awal",
			acc.FormatMoney(saldoAwal),
			acc.FormatMoney(0),
			acc.FormatMoney(0),
			acc.FormatMoney(saldoAwal),
		}
		result = append(result, awl)
		totalSaldoAwal = saldoAwal
		totalSaldoAkhir = saldoAwal
	}

	if err == nil {

		for i, e := range qr.List {
			// i is the index, e the element
			tgl, _ := e.CreatedAt.ValueDate()

			totalSaldoAkhir = totalSaldoAkhir + e.InAmount - e.OutAmount

			x := []string{strconv.Itoa(i + 1),
				tgl.(string),
				e.Desc,
				acc.FormatMoney(e.SaldoAwal),
				acc.FormatMoney(e.InAmount),
				acc.FormatMoney(e.OutAmount),
				acc.FormatMoney(totalSaldoAkhir),
			}
			result = append(result, x)
			totalSaldoAwal = totalSaldoAwal + e.SaldoAwal
			totalInAmount = totalInAmount + e.InAmount
			totalOutAmount = totalOutAmount + e.OutAmount

		}

	}

	return result, acc.FormatMoney(totalSaldoAwal), acc.FormatMoney(totalInAmount), acc.FormatMoney(totalOutAmount), acc.FormatMoney(totalSaldoAkhir)
}

func (a ReportService) GenerateLrdp(param *dto.ReportMonitoring) error {
	cbn, err := a.branchRepository.Get(param.BranchID)
	if err != nil {
		return err
	}

	begin := time.Now()

	period1, _ := time.Parse(time.RFC3339, param.DateParams[0])
	period2, _ := time.Parse(time.RFC3339, param.DateParams[1])
	var tglCetak string = "Tgl Cetak: " + begin.Format(layoutID)
	var pukulCetak string = "Pkl Cetak: " + begin.Format("15:04:05")
	var userIdCetak string = "Pkl Cetak: " + param.UserId
	var kodeNama string = "Kode - Nama: " + cbn.Code + " - " + cbn.Name
	var periodeParam string = "Periode: " + period1.Format(layoutID) + " - " + period2.Format(layoutID)

	darkGrayColor := getDarkGrayColor()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeaderLrdp()
	contents, totalAll := getDataBkk(a, param.CompanyID, param.BranchID, param.DateParams)

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	m.SetAliasNbPages("{nb}")
	m.SetFirstPageNb(1)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				m.Text("Leader Barista Point Coffe. SA", props.Text{
					Size:        8,
					Align:       consts.Left,
					Extrapolate: false,
					Color:       darkGrayColor,
				})
			})

			m.ColSpace(6)

			m.Col(3, func() {
				m.Text(tglCetak, props.Text{
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: darkGrayColor,
				})
				m.Text(pukulCetak, props.Text{
					Top:   3,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: darkGrayColor,
				})
				m.Text(userIdCetak, props.Text{
					Top:   6,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Right,
					Color: darkGrayColor,
				})
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Laporan Reimbursement Dana Petty Cash", props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(kodeNama, props.Text{
				Top:   1,
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
			})
			m.Text(periodeParam, props.Text{
				Top:   6,
				Size:  10,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.Row(2, func() {
		m.ColSpace(7)
		m.Col(5, func() {
			m.Text("Halaman: "+strconv.Itoa(m.GetCurrentPage())+"/{nb}", props.Text{
				Align: consts.Right,
				Size:  8,
			})
		})
	})

	m.Line(10)
	m.SetBackgroundColor(whiteColor)

	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      10,
			GridSizes: []uint{1, 3, 5, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{1, 3, 5, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   2,
		Line:                 false,
	})

	m.Line(5)

	m.Row(4, func() {
		m.ColSpace(8)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(2, func() {
			m.Text(totalAll, props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	m.Line(10)

	m.Row(35, func() {
		m.ColSpace(5)
		m.Col(14, func() {
			_ = m.FileImage("assets/images/ttd.png", props.Rect{
				Percent: 100,
			})
		})
	})

	errFile := m.OutputFileAndClose("pdfs/lrdp.pdf")
	if errFile != nil {
		fmt.Println("Could not save PDF:", errFile)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	return nil
}

func getHeaderLrdp() []string {
	return []string{"No", "Tanggal", "Deskripsi", "Nilai (Rp.)"}
}

func getDataBkk(a ReportService, companyId string, branchId string, dateQuery []string) (data [][]string, totalAll string) {
	var result [][]string
	var total int64 = 0

	acc := accounting.Accounting{Precision: 2, Thousand: ".", Decimal: ","}
	qr, err := a.bkkheaderRepository.Query(&models.BKKHeaderQueryParam{
		CompanyID:  companyId,
		BranchID:   branchId,
		DateQuery:  dateQuery,
		OrderParam: dto.OrderParam{Key: "created_at", Direction: dto.OrderByASC},
	})

	if err == nil {
		count := 0
		for _, e := range qr.List {
			tgl, _ := e.CreatedAt.ValueDate()
			for _, d := range e.BKKDetails {
				count++
				x := []string{strconv.Itoa(count),
					tgl.(string),
					d.LinesDesc,
					acc.FormatMoney(d.LinesAmount),
				}
				result = append(result, x)
				total = total + d.LinesAmount
			}

		}
	}

	return result, acc.FormatMoney(total)
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}
