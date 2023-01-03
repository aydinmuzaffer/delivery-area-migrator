package vendorlist_test

import (
	"context"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	gormtransaction "github.com/aydinmuzaffer/migration-tool-service/src/db/gormtransaction"
	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
	vendorlist "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/storage/vendor_list"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type VendorDeliveryAreaGetterSuite struct {
	suite.Suite
	mockDB        *gorm.DB
	gormDBWrapper gormtransaction.DBWrapper
	sqlMock       sqlmock.Sqlmock
	dag           domain.VendorDeliveryAreaGetter
}

func (s *VendorDeliveryAreaGetterSuite) SetupSuite() {
	testDB, mock, err := sqlmock.New()
	if err != nil {
		s.Error(err)
	}

	s.sqlMock = mock
	dialector := postgres.New(postgres.Config{
		DSN:                  "migrationtool",
		DriverName:           "postgres",
		Conn:                 testDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		s.Error(err)
	}
	s.mockDB = db
}

func (s *VendorDeliveryAreaGetterSuite) generateMockVendorDeliveryAreaModelRows(
	count int,
) (*[]domain.VendorDeliveryAreaModel, *sqlmock.Rows) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randomer := rand.New(randSource)
	deliveryAreas := []domain.VendorDeliveryAreaModel{}
	for i := 0; i < count; i++ {
		deliveryAreas = append(deliveryAreas, domain.VendorDeliveryAreaModel{
			ID: randomer.Int63(),
		})
	}

	rows := sqlmock.NewRows(
		[]string{
			"last_updated_on",
			"tourist_tax_value",
			"tourist_tax_type",
			"municipality_tax_value",
			"municipality_tax_type",
			"minimum_order_value",
			"delivery_fee_value",
			"delivery_fee_fype",
			"delivery_area_status",
			"porygon_delivery_area_id",
			"vendor_id",
			"id",
		},
	)
	for i := range deliveryAreas {
		rows.AddRow(
			deliveryAreas[i].LastUpdatedOn,
			deliveryAreas[i].TouristTaxValue,
			deliveryAreas[i].TouristTaxType,
			deliveryAreas[i].MunicipalityTaxValue,
			deliveryAreas[i].MunicipalityTaxType,
			deliveryAreas[i].MinimumOrderValue,
			deliveryAreas[i].DeliveryFeeValue,
			deliveryAreas[i].DeliveryFeeType,
			deliveryAreas[i].DeliveryAreaStatus,
			deliveryAreas[i].PorygonDeliveryAreaID,
			deliveryAreas[i].VendorID,
			deliveryAreas[i].ID)
	}

	return &deliveryAreas, rows
}

func (s *VendorDeliveryAreaGetterSuite) generateMockVendorDeliveryAreaIdsModelRows(
	count int,
) (*[]int64, *sqlmock.Rows) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randomer := rand.New(randSource)
	deliveryAreaIds := []int64{}
	for i := 0; i < count; i++ {
		deliveryAreaIds = append(deliveryAreaIds, randomer.Int63())
	}

	rows := sqlmock.NewRows(
		[]string{
			"id",
		},
	)
	for i := range deliveryAreaIds {
		rows.AddRow(deliveryAreaIds[i])
	}

	return &deliveryAreaIds, rows
}

func (s *VendorDeliveryAreaGetterSuite) SetupTest() {
	s.gormDBWrapper = gormtransaction.NewGormDBWrapper(s.mockDB)
	s.dag = vendorlist.NewVendorDeliveryAreaGetter(s.gormDBWrapper)
}

func (s *VendorDeliveryAreaGetterSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.sqlMock.ExpectationsWereMet())
}

func TestVendorDeliveryAreaGetterSuite(t *testing.T) {
	suite.Run(t, new(VendorDeliveryAreaGetterSuite))
}

func (s *VendorDeliveryAreaGetterSuite) TestGet_Success() {
	models, rows := s.generateMockVendorDeliveryAreaModelRows(2)
	deliveryAreaIDs := []int64{1, 2}
	queryString := `
	SELECT
		id AS id,
		vendorid AS vendor_id,
		porygondeliveryareaid AS porygon_delivery_area_id,
		deliveryareastatus AS delivery_area_status,
		deliveryfeetype AS delivery_fee_fype,
		deliveryfeevalue AS delivery_fee_value,
		minimumordervalue::numeric(9,2) AS minimum_order_value,
		municipalitytaxtype AS municipality_tax_type,
		municipalitytaxvalue AS municipality_tax_value,
		touristtaxtype AS tourist_tax_type,
		touristtaxvalue AS tourist_tax_value,
		lastupdatedon::timestamptz AS last_updated_on
	FROM public.tlb_vendor_deliveryarea
	WHERE id IN `
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(queryString)).WillReturnRows(rows)

	results, err := s.dag.Get(context.Background(), &deliveryAreaIDs)

	assert.Equal(s.T(), len(*models), len(*results))
	assert.ElementsMatch(s.T(), *models, *results)
	assert.Nil(s.T(), err)
}

func (s *VendorDeliveryAreaGetterSuite) TestGetIds_Success() {
	models, rows := s.generateMockVendorDeliveryAreaIdsModelRows(2)
	queryString := `SELECT id FROM public.tlb_vendor_deliveryarea WHERE isdeleted is false`
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(queryString)).WillReturnRows(rows)

	results, err := s.dag.GetIds(context.Background())

	assert.Equal(s.T(), len(*models), len(*results))
	assert.ElementsMatch(s.T(), *models, *results)
	assert.Nil(s.T(), err)
}
