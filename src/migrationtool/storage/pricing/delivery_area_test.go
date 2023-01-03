package pricing_test

// import (
// 	"context"
// 	"fmt"
// 	"math/rand"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	gormtransaction "github.com/aydinmuzaffer/migration-tool-service/src/db/gormtransaction"
// 	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
// 	pricing "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/storage/pricing"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type VenderDeliveryAreaSetterSuite struct {
// 	suite.Suite
// 	mockDB        *gorm.DB
// 	gormDBWrapper gormtransaction.DBWrapper
// 	sqlMock       sqlmock.Sqlmock
// 	das           domain.VenderDeliveryAreaSetter
// }

// func (s *VenderDeliveryAreaSetterSuite) SetupSuite() {
// 	testDB, mock, err := sqlmock.New()
// 	if err != nil {
// 		s.Error(err)
// 	}

// 	s.sqlMock = mock
// 	dialector := postgres.New(postgres.Config{
// 		DSN:                  "migrationtool",
// 		DriverName:           "postgres",
// 		Conn:                 testDB,
// 		PreferSimpleProtocol: true,
// 	})

// 	db, err := gorm.Open(dialector, &gorm.Config{})
// 	if err != nil {
// 		s.Error(err)
// 	}
// 	s.mockDB = db
// }

// func (s *VenderDeliveryAreaSetterSuite) generateMockVendorDeliveryAreaModelRows(
// 	count int,
// ) (*[]domain.VendorDeliveryAreaModel, []interface{}, *sqlmock.Rows) {
// 	deliveryAreas := []domain.VendorDeliveryAreaModel{}
// 	randSource := rand.NewSource(time.Now().UnixNano())
// 	randomer := rand.New(randSource)
// 	for i := 0; i < count; i++ {
// 		deliveryAreaa := domain.VendorDeliveryAreaModel{
// 			ID:                    4444,
// 			VendorID:              23233,
// 			PorygonDeliveryAreaID: uuid.New().String(),
// 			DeliveryAreaStatus:    0,
// 			DeliveryFeeType:       0,
// 			DeliveryFeeValue:      12.0,
// 			MinimumOrderValue:     33.22,
// 			MunicipalityTaxType:   0,
// 			MunicipalityTaxValue:  10.4,
// 			TouristTaxType:        0,
// 			TouristTaxValue:       0,
// 		}
// 		deliveryAreas = append(deliveryAreas, deliveryAreaa)
// 	}

// 	columnNumber := 11
// 	valueArgs := make([]interface{}, len(deliveryAreas)*columnNumber)
// 	rows := sqlmock.NewRows([]string{"id", "created_at", "last_update_on"})
// 	for i, deliveryArea := range deliveryAreas {

// 		valueArgs[i*columnNumber+0] = deliveryArea.ID
// 		valueArgs[i*columnNumber+1] = deliveryArea.VendorID
// 		valueArgs[i*columnNumber+2] = deliveryArea.PorygonDeliveryAreaID
// 		valueArgs[i*columnNumber+3] = deliveryArea.DeliveryAreaStatus
// 		valueArgs[i*columnNumber+4] = fmt.Sprintf("%v", deliveryArea.DeliveryFeeType)
// 		valueArgs[i*columnNumber+5] = deliveryArea.DeliveryFeeValue
// 		valueArgs[i*columnNumber+6] = deliveryArea.MinimumOrderValue
// 		valueArgs[i*columnNumber+7] = fmt.Sprintf("%v", deliveryArea.MunicipalityTaxType)
// 		valueArgs[i*columnNumber+8] = deliveryArea.MunicipalityTaxValue
// 		valueArgs[i*columnNumber+9] = deliveryArea.TouristTaxType
// 		valueArgs[i*columnNumber+10] = fmt.Sprintf("%v", deliveryArea.TouristTaxValue)

// 		rows.AddRow(randomer.Int63(), time.Time{}, nil)
// 	}
// 	return &deliveryAreas, valueArgs, rows
// }

// func (s *VenderDeliveryAreaSetterSuite) SetupTest() {
// 	s.gormDBWrapper = gormtransaction.NewGormDBWrapper(s.mockDB)
// 	s.das = pricing.NewVenderDeliveryAreaSetter(s.gormDBWrapper)
// }

// func (s *VenderDeliveryAreaSetterSuite) AfterTest(_, _ string) {
// 	require.NoError(s.T(), s.sqlMock.ExpectationsWereMet())
// }

// func TestVendorDeliveryAreaGetterSuite(t *testing.T) {
// 	suite.Run(t, new(VenderDeliveryAreaSetterSuite))
// }

// func (s *VenderDeliveryAreaSetterSuite) TestCreate_DeliveryArea() {
// 	model, values, _ := s.generateMockVendorDeliveryAreaModelRows(1)

// 	s.sqlMock.MatchExpectationsInOrder(false)
// 	//s.sqlMock.ExpectBegin()
// 	s.sqlMock.ExpectExec(regexp.QuoteMeta(`INSERT INTO dbo."vendor_delivery_area" ("id", "vendor_id", "porygon_delivery_area_id", "delivery_area_status","delivery_fee_type", "delivery_fee_value", "minimum_order_value", "municipality_tax_type", "municipality_tax_value", "tourist_tax_type", "tourist_tax_value")    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`)).
// 		WithArgs(values)
// 		//WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 	//s.sqlMock.ExpectCommit()
// 	err := s.das.Insert(context.Background(), model)

// 	s.Nil(err)
// }
