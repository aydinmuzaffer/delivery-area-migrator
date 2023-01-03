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
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type VenderDeliveryAreaPolygonGetterSuite struct {
	suite.Suite
	mockDB        *gorm.DB
	gormDBWrapper gormtransaction.DBWrapper
	sqlMock       sqlmock.Sqlmock
	dapg          domain.VenderDeliveryAreaPolygonGetter
}

func (s *VenderDeliveryAreaPolygonGetterSuite) SetupSuite() {
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

func (s *VenderDeliveryAreaPolygonGetterSuite) generateMockVendorDeliveryAreaPolygonGetModelRows(
	count int,
) (*[]domain.VendorDeliveryAreaPolygonGetModel, *sqlmock.Rows) {
	randSource := rand.NewSource(time.Now().UnixNano())
	randomer := rand.New(randSource)
	deliveryAreaPolygons := []domain.VendorDeliveryAreaPolygonGetModel{}
	for i := 0; i < count; i++ {
		deliveryAreaPolygons = append(deliveryAreaPolygons, domain.VendorDeliveryAreaPolygonGetModel{
			VendorDeliveryAreaID: randomer.Int(),
			DeliveryAreaPolygon:  &ewkb.Polygon{geom.NewPolygon(geom.XY)},
		})
	}

	rows := sqlmock.NewRows(
		[]string{
			"last_updated_on",
			"delivery_area_polygon",
			"vendor_delivery_area_id",
		},
	)
	for i := range deliveryAreaPolygons {
		rows.AddRow(
			deliveryAreaPolygons[i].LastUpdatedOn,
			deliveryAreaPolygons[i].DeliveryAreaPolygon,
			deliveryAreaPolygons[i].VendorDeliveryAreaID)
	}

	return &deliveryAreaPolygons, rows
}

func (s *VenderDeliveryAreaPolygonGetterSuite) SetupTest() {
	s.gormDBWrapper = gormtransaction.NewGormDBWrapper(s.mockDB)
	s.dapg = vendorlist.NewVenderDeliveryAreaPolygonGetter(s.gormDBWrapper)
}

func (s *VenderDeliveryAreaPolygonGetterSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.sqlMock.ExpectationsWereMet())
}

func TestVenderDeliveryAreaPolygonGetterSuite(t *testing.T) {
	suite.Run(t, new(VenderDeliveryAreaPolygonGetterSuite))
}

func (s *VenderDeliveryAreaPolygonGetterSuite) TestGet_Success() {
	models, rows := s.generateMockVendorDeliveryAreaPolygonGetModelRows(2)
	deliveryAreaIDs := []int64{1, 2}
	queryString := `
		SELECT
			vendordeliveryareaid AS vendor_delivery_area_id,
			ST_AsEWKB(deliveryareapolygon) AS delivery_area_polygon,
			lastupdatedon::timestamptz AS last_updated_on
		FROM public.tlb_vendor_deliveryarea_polygon
		WHERE vendordeliveryareaid IN`
	s.sqlMock.ExpectQuery(regexp.QuoteMeta(queryString)).WillReturnRows(rows)

	results, err := s.dapg.Get(context.Background(), &deliveryAreaIDs)

	assert.Equal(s.T(), len(*models), len(results))
	assert.ElementsMatch(s.T(), *models, results)
	assert.Nil(s.T(), err)
}
