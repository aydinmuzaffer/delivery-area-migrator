package deliveryareapolygonservice

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
	"github.com/pkg/errors"
	"github.com/twpayne/go-geom"
	"golang.org/x/sync/errgroup"
)

type PolygonService interface {
	For([]domain.VendorDeliveryAreaPolygonGetModel) PolygonService
	GetPricingPolygons() (*[]domain.VendorDeliveryAreaPolygonInsertModel, error)
}

type polygonService struct {
	wg             errgroup.Group
	polygonBuilder PolygonBuilder
}

func NewPolygonService(polygonBuilder PolygonBuilder) PolygonService {
	return &polygonService{
		polygonBuilder: polygonBuilder,
	}
}

func (s *polygonService) For(polygons []domain.VendorDeliveryAreaPolygonGetModel) PolygonService {
	s.polygonBuilder.SetVendorListPolygons(polygons)
	return s
}

func (s *polygonService) GetPricingPolygons() (*[]domain.VendorDeliveryAreaPolygonInsertModel, error) {
	s.polygonBuilder.PopulateInsertModels()

	// func(b PolygonBuilder) {
	// 	s.wg.Go(func() error {
	// 		if _, err := b.SetPolygonHashes(); err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	})
	// }(s.polygonBuilder)

	s.wg.Go(func() error {
		if err := s.polygonBuilder.SetPolygonHashes(); err != nil {
			return err
		}
		return nil
	})

	// func(b PolygonBuilder) {
	// 	s.wg.Go(func() error {
	// 		b.SetPolygonInsertString()
	// 		return nil
	// 	})
	// }(s.polygonBuilder)

	s.wg.Go(func() error {
		s.polygonBuilder.SetPolygonInsertString()
		return nil
	})

	if err := s.wg.Wait(); err != nil {
		return nil, err
	}

	return s.polygonBuilder.Build(), nil
}

type PolygonBuilder interface {
	SetVendorListPolygons(polygons []domain.VendorDeliveryAreaPolygonGetModel)
	PopulateInsertModels()
	SetPolygonHashes() error
	SetPolygonInsertString()
	Build() *[]domain.VendorDeliveryAreaPolygonInsertModel
}

type polygonBuilder struct {
	polygons         []domain.VendorDeliveryAreaPolygonGetModel
	polygonsToInsert []domain.VendorDeliveryAreaPolygonInsertModel
}

func NewPolygonBuilder() PolygonBuilder {
	return &polygonBuilder{}
}

func (b *polygonBuilder) SetVendorListPolygons(polygons []domain.VendorDeliveryAreaPolygonGetModel) {
	b.polygons = polygons
}

func (b *polygonBuilder) PopulateInsertModels() {
	polygonsToInsert := make([]domain.VendorDeliveryAreaPolygonInsertModel, len(b.polygons))
	for i := range b.polygons {
		p := &b.polygons[i]
		polygonToInsert := domain.VendorDeliveryAreaPolygonInsertModel{
			VendorDeliveryAreaID: p.VendorDeliveryAreaID,
			DeliveryAreaPolygon:  p.DeliveryAreaPolygon,
			SRID:                 strconv.Itoa(p.DeliveryAreaPolygon.SRID()),
		}
		polygonsToInsert[i] = polygonToInsert
	}
	b.polygonsToInsert = polygonsToInsert
}

func (b *polygonBuilder) SetPolygonHashes() error {
	for i := range b.polygonsToInsert {
		p := &b.polygonsToInsert[i]
		hash, err := generateHashForPolygon(p.DeliveryAreaPolygon.Coords())
		if err != nil {
			errors.Wrap(err, "[PolygonBuilder].SetPolygonHashes")
		}
		p.DeliveryAreaHash = *hash
	}
	return nil
}

func (b *polygonBuilder) SetPolygonInsertString() {
	for i := range b.polygonsToInsert {
		p := &b.polygonsToInsert[i]
		deliveryAreaPolygonInsertString := generatePolygonString(p.DeliveryAreaPolygon.Coords())
		p.DeliveryAreaPolygonInsertString = fmt.Sprintf("SRID=%s;%s", p.SRID, deliveryAreaPolygonInsertString)
	}
}

func (b *polygonBuilder) Build() *[]domain.VendorDeliveryAreaPolygonInsertModel {
	return &b.polygonsToInsert
}

func generatePolygonString(coords [][]geom.Coord) string {
	insertString := "POLYGON ((%s))"
	coordinates := []string{}
	for _, coord := range coords {
		for _, c := range coord {
			xy := []string{}
			xy = append(xy, fmt.Sprintf("%v", c.X()))
			xy = append(xy, fmt.Sprintf("%v", c.Y()))
			coordinates = append(coordinates, strings.Join(xy, " "))
		}
	}
	//prevent unclosed rings for Polygons.
	//https://stackoverflow.com/questions/35128026/i-am-trying-to-insert-polygon-data-in-table-then-i-get-an-error
	first := coordinates[0]
	last := coordinates[len(coordinates)-1]
	if first != last {
		coordinates = append(coordinates, first)
	}
	return fmt.Sprintf(insertString, strings.Join(coordinates, ","))
}

func generateHashForPolygon(coordinates [][]geom.Coord) (*string, error) {
	bytes, err := json.Marshal(coordinates)
	if err != nil {
		return nil, fmt.Errorf("can't seriliaze coordinates: %v", err)
	}
	h := sha512.New384()
	if _, err := h.Write(bytes); err != nil {
		return nil, fmt.Errorf("can't hash coordinates: %v", err)
	}
	hashedCoordinates := hex.EncodeToString(h.Sum(nil))
	return &hashedCoordinates, nil
}
