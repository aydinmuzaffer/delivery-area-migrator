package domain

import (
	"context"
	"time"
)

type VendorDeliveryAreaModel struct {
	ID                    int64
	VendorID              int
	PorygonDeliveryAreaID string
	DeliveryAreaStatus    int
	DeliveryFeeType       int16
	DeliveryFeeValue      float64
	MinimumOrderValue     float64
	MunicipalityTaxType   int16
	MunicipalityTaxValue  float64
	TouristTaxType        int16
	TouristTaxValue       float64
	LastUpdatedOn         time.Time
}

type VendorDeliveryAreaGetter interface {
	Get(ctx context.Context, ids *[]int64) (*[]VendorDeliveryAreaModel, error)
	GetIds(ctx context.Context) (*[]int64, error)
}

type VenderDeliveryAreaSetter interface {
	Insert(ctx context.Context, areas *[]VendorDeliveryAreaModel) error
}
