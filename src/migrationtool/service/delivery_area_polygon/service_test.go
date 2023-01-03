package deliveryareapolygonservice_test

import (
	"testing"
	"time"

	mocks "github.com/aydinmuzaffer/migration-tool-service/mocks/service/delivery_area_polygon"
	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
	deliveryareapolygonservice "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/service/delivery_area_polygon"
	"github.com/stretchr/testify/suite"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
)

type PolygonServiceTestSuite struct {
	suite.Suite
	service            deliveryareapolygonservice.PolygonService
	polygonBuilderMock *mocks.PolygonBuilder
}

func (s *PolygonServiceTestSuite) SetupTest() {
	s.polygonBuilderMock = &mocks.PolygonBuilder{}
	s.service = deliveryareapolygonservice.NewPolygonService(s.polygonBuilderMock)
}

func TestCashPaymentService(t *testing.T) {
	suite.Run(t, new(PolygonServiceTestSuite))
}

func (s *PolygonServiceTestSuite) TestGetPricingPolygons() {
	polygon1 := geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
		{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
	})
	polygon2 := geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{
		{{0, 1}, {1, 2}, {1, 1}, {0, 1}, {0, 0}},
	})
	vendorListPolygons := []domain.VendorDeliveryAreaPolygonGetModel{
		{
			VendorDeliveryAreaID: 1,
			DeliveryAreaPolygon:  &ewkb.Polygon{polygon1},
			LastUpdatedOn:        time.Now(),
		},
		{
			VendorDeliveryAreaID: 2,
			DeliveryAreaPolygon:  &ewkb.Polygon{polygon2},
			LastUpdatedOn:        time.Now(),
		},
	}

	expectedPricingPolygons := []domain.VendorDeliveryAreaPolygonInsertModel{
		{
			VendorDeliveryAreaID:            vendorListPolygons[0].VendorDeliveryAreaID,
			DeliveryAreaPolygon:             vendorListPolygons[0].DeliveryAreaPolygon,
			DeliveryAreaHash:                "Hash",
			SRID:                            "1234",
			DeliveryAreaPolygonInsertString: "InsertString",
			LastUpdatedOn:                   vendorListPolygons[0].LastUpdatedOn,
		},
		{
			VendorDeliveryAreaID:            vendorListPolygons[1].VendorDeliveryAreaID,
			DeliveryAreaPolygon:             vendorListPolygons[1].DeliveryAreaPolygon,
			DeliveryAreaHash:                "Hash",
			SRID:                            "1234",
			DeliveryAreaPolygonInsertString: "InsertString",
			LastUpdatedOn:                   vendorListPolygons[1].LastUpdatedOn,
		},
	}

	s.polygonBuilderMock.On("SetVendorListPolygons", vendorListPolygons).Return(nil)
	s.polygonBuilderMock.On("PopulateInsertModels").Return(nil)
	s.polygonBuilderMock.On("SetPolygonHashes").Return(nil)
	s.polygonBuilderMock.On("SetPolygonInsertString").Return(nil)
	s.polygonBuilderMock.On("Build").Return(&expectedPricingPolygons)

	pricingPolygons, err := s.service.For(vendorListPolygons).GetPricingPolygons()

	s.Nil(err)
	s.NotNil(pricingPolygons)
	s.Equal(*pricingPolygons, expectedPricingPolygons)
	s.polygonBuilderMock.AssertCalled(s.T(), "PopulateInsertModels")
	s.polygonBuilderMock.AssertCalled(s.T(), "SetPolygonHashes")
	s.polygonBuilderMock.AssertCalled(s.T(), "SetPolygonInsertString")
	s.polygonBuilderMock.AssertCalled(s.T(), "Build")

}
