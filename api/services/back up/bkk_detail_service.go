package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// BKKDetailService service layer
type BKKDetailService struct {
	logger              lib.Logger
	casbinService       CasbinService
	bkkdetailRepository repository.BKKDetailRepository
	redis               lib.Redis
}

// NewBKKDetailService creates a new bkkdetailservice
func NewBKKDetailService(
	redis lib.Redis,
	logger lib.Logger,
	casbinService CasbinService,
	bkkdetailRepository repository.BKKDetailRepository,

) BKKDetailService {
	return BKKDetailService{
		logger:              logger,
		casbinService:       casbinService,
		bkkdetailRepository: bkkdetailRepository,
		redis:               redis,
	}
}

func wrapperBKKDetailTemp(key string) string {
	return fmt.Sprintf("bkkdetail%s:", key)
}

// WithTrx delegates transaction to repository database
func (a BKKDetailService) WithTrx(trxHandle *gorm.DB) BKKDetailService {
	a.bkkdetailRepository = a.bkkdetailRepository.WithTrx(trxHandle)

	return a
}

func (a BKKDetailService) Query(param *models.BKKDetailQueryParam) (bkkdetailQR *models.BKKDetailQueryResult, err error) {
	return a.bkkdetailRepository.Query(param)
}

func (a BKKDetailService) Get(id string) (*models.BKKDetail, error) {
	bkkdetail, err := a.bkkdetailRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return bkkdetail, nil
}

func (a BKKDetailService) Check(item *models.BKKDetail) error {
	qr, err := a.bkkdetailRepository.Query(&models.BKKDetailQueryParam{LinesDesc: item.LinesDesc})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.BKKDetailAlreadyExists
	}

	return nil
}

func (a BKKDetailService) Create(bkkdetail *models.BKKDetail) (id uint, err error) {
	if err = a.Check(bkkdetail); err != nil {
		return
	}

	if err = a.bkkdetailRepository.Create(bkkdetail); err != nil {
		return
	}
	return bkkdetail.RecordID, nil
}

func (a BKKDetailService) Update(id string, bkkdetail *models.BKKDetail) error {
	oBKKDetail, err := a.Get(id)
	if err != nil {
		return err
	} else if bkkdetail.LinesDesc != oBKKDetail.LinesDesc {
		if err = a.Check(bkkdetail); err != nil {
			return err
		}
	}

	if err := a.bkkdetailRepository.Update(id, bkkdetail); err != nil {
		return err
	}

	return nil
}

func (a BKKDetailService) Delete(id string) error {
	_, err := a.bkkdetailRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.bkkdetailRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a BKKDetailService) UpdateStatus(id string, status int) error {
	_, err := a.bkkdetailRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.bkkdetailRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}

func (a BKKDetailService) GetTamp(username string) (*models.BKKDetail, error) {
	var (
		key = wrapperBKKDetailTemp(username)
		val string
	)

	err := a.redis.Get(key, &val)
	if err != nil {
		return nil, err
	}
	bkkdetail := new(models.BKKDetail)
	//Deserialization
	errJson := json.Unmarshal([]byte(val), &bkkdetail)
	if errJson != nil {
		return nil, errJson
	}
	return bkkdetail, nil
}

var currentPostId int

func (a BKKDetailService) AddTamp(bkkdetail *models.BKKDetail, username string) (err error) {
	currentPostId += 1
	b, err := json.Marshal(bkkdetail)
	if err != nil {
		return err
	}
	key := wrapperBKKDetailTemp(username) + strconv.Itoa(currentPostId)

	errRedis := a.redis.Set(key, b, 0)
	if errRedis != nil {
		return err
	}

	return nil
}
