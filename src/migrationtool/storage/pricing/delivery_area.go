package pricing

import (
	"context"
	"fmt"
	"strings"

	"github.com/aydinmuzaffer/migration-tool-service/src/db/gormtransaction"
	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
)

var _ domain.VenderDeliveryAreaSetter = (*venderDeliveryAreaSetter)(nil) // compile time proof

type venderDeliveryAreaSetter struct {
	dw gormtransaction.DBWrapper
}

func (s *venderDeliveryAreaSetter) Insert(ctx context.Context, areas *[]domain.VendorDeliveryAreaModel) error {
	insertScript := `INSERT INTO dbo."vendor_delivery_area" ("id", "vendor_id", "porygon_delivery_area_id", "delivery_area_status",
	"delivery_fee_type", "delivery_fee_value", "minimum_order_value", "municipality_tax_type", "municipality_tax_value", "tourist_tax_type", "tourist_tax_value") 
   	VALUES %s`

	columnNumber := 11
	valueStrings := make([]string, len(*areas))
	valueArgs := make([]interface{}, len(*areas)*columnNumber)
	for index, post := range *areas {
		valueStrings[index] = fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			index*columnNumber+1,
			index*columnNumber+2,
			index*columnNumber+3,
			index*columnNumber+4,
			index*columnNumber+5,
			index*columnNumber+6,
			index*columnNumber+7,
			index*columnNumber+8,
			index*columnNumber+9,
			index*columnNumber+10,
			index*columnNumber+11,
		)
		valueArgs[index*columnNumber+0] = post.ID
		valueArgs[index*columnNumber+1] = post.VendorID
		valueArgs[index*columnNumber+2] = post.PorygonDeliveryAreaID
		valueArgs[index*columnNumber+3] = post.DeliveryAreaStatus
		valueArgs[index*columnNumber+4] = fmt.Sprintf("%v", post.DeliveryFeeType)
		valueArgs[index*columnNumber+5] = post.DeliveryFeeValue
		valueArgs[index*columnNumber+6] = post.MinimumOrderValue
		valueArgs[index*columnNumber+7] = fmt.Sprintf("%v", post.MunicipalityTaxType)
		valueArgs[index*columnNumber+8] = post.MunicipalityTaxValue
		valueArgs[index*columnNumber+9] = fmt.Sprintf("%v", post.TouristTaxType)
		valueArgs[index*columnNumber+10] = post.TouristTaxValue
	}
	stmt := fmt.Sprintf(insertScript, strings.Join(valueStrings, ","))
	db := s.dw.GetDB(ctx).Exec(stmt, valueArgs...)
	return db.Error
}

func NewVenderDeliveryAreaSetter(
	dw gormtransaction.DBWrapper,
) domain.VenderDeliveryAreaSetter {
	return &venderDeliveryAreaSetter{
		dw: dw,
	}
}
