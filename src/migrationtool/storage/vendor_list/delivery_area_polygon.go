package vendorlist

import (
	"context"

	"github.com/aydinmuzaffer/migration-tool-service/src/db/gormtransaction"
	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
)

var _ domain.VenderDeliveryAreaPolygonGetter = (*venderDeliveryAreaPolygonGetter)(nil) // compile time proof

type venderDeliveryAreaPolygonGetter struct {
	dw gormtransaction.DBWrapper
}

func (g *venderDeliveryAreaPolygonGetter) Get(ctx context.Context, ids *[]int64) ([]domain.VendorDeliveryAreaPolygonGetModel, error) {
	queryString := `
		SELECT
			vendordeliveryareaid AS vendor_delivery_area_id,
			ST_AsEWKB(deliveryareapolygon) AS delivery_area_polygon,
			lastupdatedon::timestamptz AS last_updated_on
		FROM public.tlb_vendor_deliveryarea_polygon
		WHERE vendordeliveryareaid IN (?)`
	var deliveryAreaPolygons []domain.VendorDeliveryAreaPolygonGetModel
	db := g.dw.GetDB(ctx).Raw(queryString, *ids).Scan(&deliveryAreaPolygons)
	return deliveryAreaPolygons, db.Error
}

func NewVenderDeliveryAreaPolygonGetter(
	dw gormtransaction.DBWrapper,
) domain.VenderDeliveryAreaPolygonGetter {
	return &venderDeliveryAreaPolygonGetter{
		dw: dw,
	}
}
