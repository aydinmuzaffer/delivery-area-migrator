package pricing

import (
	"context"
	"fmt"
	"strings"

	"github.com/aydinmuzaffer/migration-tool-service/src/db/gormtransaction"
	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
)

var _ domain.VenderDeliveryAreaPolygonSetter = (*venderDeliveryAreaPolygonSetter)(nil) // compile time proof

type venderDeliveryAreaPolygonSetter struct {
	dw gormtransaction.DBWrapper
}

func (s *venderDeliveryAreaPolygonSetter) Insert(ctx context.Context, deliveryAreaPolygons *[]domain.VendorDeliveryAreaPolygonInsertModel) error {
	insertScript := `INSERT INTO dbo.vendor_delivery_area_polygon (vendor_delivery_area_id, delivery_area_hash, delivery_area_polygon)
					 VALUES %s`

	columnNumber := 3
	valueStrings := make([]string, len(*deliveryAreaPolygons))
	valueArgs := make([]interface{}, len(*deliveryAreaPolygons)*columnNumber)
	for index, p := range *deliveryAreaPolygons {
		valueStrings[index] = fmt.Sprintf("($%d,$%d,$%d)",
			index*columnNumber+1,
			index*columnNumber+2,
			index*columnNumber+3,
		)
		valueArgs[index*columnNumber+0] = p.VendorDeliveryAreaID
		valueArgs[index*columnNumber+1] = p.DeliveryAreaHash
		valueArgs[index*columnNumber+2] = p.DeliveryAreaPolygonInsertString
	}
	stmt := fmt.Sprintf(insertScript, strings.Join(valueStrings, ","))
	db := s.dw.GetDB(ctx).Exec(stmt, valueArgs...)
	return db.Error
}

func NewVenderDeliveryAreaPolygonSetter(
	dw gormtransaction.DBWrapper,
) domain.VenderDeliveryAreaPolygonSetter {
	return &venderDeliveryAreaPolygonSetter{
		dw: dw,
	}
}
