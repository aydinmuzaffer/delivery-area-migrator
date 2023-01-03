package domain

import (
	"context"
	"time"

	"github.com/twpayne/go-geom/encoding/ewkb"
)

type VendorDeliveryAreaPolygonGetModel struct {
	VendorDeliveryAreaID int
	DeliveryAreaPolygon  *ewkb.Polygon
	LastUpdatedOn        time.Time
}

type VendorDeliveryAreaPolygonInsertModel struct {
	VendorDeliveryAreaID            int
	DeliveryAreaPolygon             *ewkb.Polygon
	DeliveryAreaHash                string
	SRID                            string
	DeliveryAreaPolygonInsertString string
	LastUpdatedOn                   time.Time
}

type VenderDeliveryAreaPolygonGetter interface {
	Get(ctx context.Context, ids *[]int64) ([]VendorDeliveryAreaPolygonGetModel, error)
}

type VenderDeliveryAreaPolygonSetter interface {
	Insert(ctx context.Context, deliveryAreaPolygons *[]VendorDeliveryAreaPolygonInsertModel) error
}
