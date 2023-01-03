package migrationservice

import (
	"context"

	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
	polygonService "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/service/delivery_area_polygon"
	"github.com/aydinmuzaffer/migration-tool-service/src/utils/collectionutils"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const DEFAULT_CHUNK_SIZE int = 500

type MigrationService interface {
	Migrate(ctx context.Context) error
}

type migrationService struct {
	chunkSize                 int
	deliveryAreaGetter        domain.VendorDeliveryAreaGetter
	deliveryAreaSetter        domain.VenderDeliveryAreaSetter
	deliveryAreaPolygonGetter domain.VenderDeliveryAreaPolygonGetter
	deliveryAreaPolygonSetter domain.VenderDeliveryAreaPolygonSetter
	polygonService            polygonService.PolygonService
}

func (m *migrationService) Migrate(ctx context.Context) error {
	areaIds, err := m.deliveryAreaGetter.GetIds(ctx)
	if err != nil {
		return errors.Wrap(err, "[MigrationService].[Migrate] error on  deliveryAreaGetter.GetIds")
	}

	iteration := len(*areaIds) / m.chunkSize
	var g errgroup.Group
	for i := 0; i <= iteration; i++ {
		areaIdsChunk := collectionutils.GetChunk(areaIds, m.chunkSize, i)
		func(areaIds *[]int64) {
			g.Go(func() error {
				areas, err := m.deliveryAreaGetter.Get(ctx, areaIds)
				if err != nil {
					return errors.Wrap(err, "[MigrationService].[Migrate] error on  deliveryAreaGetter.Get")
				}

				if err := m.deliveryAreaSetter.Insert(ctx, areas); err != nil {
					return errors.Wrap(err, "[MigrationService].[Migrate] error on  deliveryAreaGetter.Insert")
				}

				polygons, err := m.deliveryAreaPolygonGetter.Get(ctx, areaIds)
				if err != nil {
					return errors.Wrap(err, "[MigrationService].[Migrate] error on  deliveryAreaPolygonGetter.Get")
				}

				polygonsToInsert, err := m.polygonService.For(polygons).GetPricingPolygons()
				if err != nil {
					return errors.Wrap(err, "[MigrationService].[Migrate] error on  polygonService.GetPricingPolygons")
				}

				if err := m.deliveryAreaPolygonSetter.Insert(ctx, polygonsToInsert); err != nil {
					return errors.Wrap(err, "[MigrationService].[Migrate] error on  deliveryAreaPolygonGetter.Get")
				}
				return nil
			})
		}(areaIdsChunk)
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func NewMigrationService(deliveryAreaGetter domain.VendorDeliveryAreaGetter, deliveryAreaSetter domain.VenderDeliveryAreaSetter,
	deliveryAreaPolygonGetter domain.VenderDeliveryAreaPolygonGetter, deliveryAreaPolygonSetter domain.VenderDeliveryAreaPolygonSetter,
	polygonService polygonService.PolygonService, chunkSize int) MigrationService {
	return &migrationService{
		deliveryAreaGetter:        deliveryAreaGetter,
		deliveryAreaSetter:        deliveryAreaSetter,
		deliveryAreaPolygonGetter: deliveryAreaPolygonGetter,
		deliveryAreaPolygonSetter: deliveryAreaPolygonSetter,
		polygonService:            polygonService,
		chunkSize:                 chunkSize,
	}
}
