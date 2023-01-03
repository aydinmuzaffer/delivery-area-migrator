package vendorlist

import (
	"context"

	"github.com/aydinmuzaffer/migration-tool-service/src/db/gormtransaction"
	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
)

var _ domain.VendorDeliveryAreaGetter = (*vendorDeliveryAreaGetter)(nil) // compile time proof

type vendorDeliveryAreaGetter struct {
	dw gormtransaction.DBWrapper
}

func NewVendorDeliveryAreaGetter(dw gormtransaction.DBWrapper) domain.VendorDeliveryAreaGetter {
	return &vendorDeliveryAreaGetter{
		dw: dw,
	}
}

func (g *vendorDeliveryAreaGetter) Get(ctx context.Context, ids *[]int64) (*[]domain.VendorDeliveryAreaModel, error) {
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
		WHERE id IN (?)`
	var vendorAreas []domain.VendorDeliveryAreaModel
	db := g.dw.GetDB(ctx).Raw(queryString, *ids).Scan(&vendorAreas)
	return &vendorAreas, db.Error
}

func (g *vendorDeliveryAreaGetter) GetIds(ctx context.Context) (*[]int64, error) {
	var vendorIds []int64
	db := g.dw.GetDB(ctx).Raw("SELECT id FROM public.tlb_vendor_deliveryarea WHERE isdeleted is false").Scan(&vendorIds)
	return &vendorIds, db.Error
}
