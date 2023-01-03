package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"golang.org/x/sync/errgroup"
)

func getVendorDb() (*sqlx.DB, error) {
	host := "localhost"
	port := 8080
	user := "pgsrv_vendor_list_qa_app"
	password := "password"
	dbname := "vendor_list_qa"
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getPricingDb() (*sqlx.DB, error) {
	host := "localhost"
	port := 5433
	user := "pgsrv_pricing_owner"
	password := "pricing"
	dbname := "pricing_local"
	sslMode := "disable"
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getAreaIds(vendorDb *sqlx.DB) (*[]int64, error) {
	var vendorIds []int64
	rows, err := vendorDb.Query("SELECT id FROM public.tlb_vendor_deliveryarea WHERE isdeleted is false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var vendorId int64
		if err := rows.Scan(&vendorId); err != nil {
			return nil, err
		}
		vendorIds = append(vendorIds, vendorId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &vendorIds, nil
}

type VendorArea struct {
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

type DeliveryAreaPolygon struct {
	VendorDeliveryAreaID            int
	DeliveryAreaHash                string
	DeliveryAreaPolygon             *ewkb.Polygon
	DeliveryAreaPolygonInsertString string
	LastUpdatedOn                   time.Time
}

func getAreas(db *sqlx.DB, ids *[]int64) (*[]VendorArea, error) {
	queryString := `
		SELECT
			id AS id,
			vendorid AS vendorId,
			porygondeliveryareaid AS porygonDeliveryAreaId,
			deliveryareastatus AS deliveryAreaStatus,
			deliveryfeetype AS deliveryFeeType,
			deliveryfeevalue AS deliveryFeeValue,
			minimumordervalue::numeric(9,2) AS minimumOrderValue,
			municipalitytaxtype AS municipalityTaxType,
			municipalitytaxvalue AS municipalityTaxValue,
			touristtaxtype AS touristTaxType,
			touristtaxvalue AS touristTaxValue,
			lastupdatedon::timestamptz AS lastUpdatedOn
	FROM public.tlb_vendor_deliveryarea
	WHERE id IN (?)
	`

	query, args, err := sqlx.In(queryString, *ids)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query) // sqlx.In returns queries with the `?` bindvar, rebind it here for matching the database in used (e.g. postgre, oracle etc, can skip it if you use mysql)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendorAreas []VendorArea

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var vendorArea VendorArea
		if err := rows.Scan(&vendorArea.ID, &vendorArea.VendorID, &vendorArea.PorygonDeliveryAreaID, &vendorArea.DeliveryAreaStatus,
			&vendorArea.DeliveryFeeType, &vendorArea.DeliveryFeeValue, &vendorArea.MinimumOrderValue, &vendorArea.MunicipalityTaxType,
			&vendorArea.MunicipalityTaxValue, &vendorArea.TouristTaxType, &vendorArea.TouristTaxValue, &vendorArea.LastUpdatedOn); err != nil {
			return nil, err
		}
		vendorAreas = append(vendorAreas, vendorArea)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &vendorAreas, nil
}

func getPolygons(db *sqlx.DB, ids *[]int64) ([]*DeliveryAreaPolygon, error) {
	queryString := `
		SELECT
			vendordeliveryareaid AS vendorDeliveryAreaId,
			ST_AsEWKB(deliveryareapolygon) AS deliveryAreaPolygon,
			lastupdatedon::timestamptz AS lastUpdatedOn
		FROM public.tlb_vendor_deliveryarea_polygon
		WHERE vendordeliveryareaid IN (?)`

	query, args, err := sqlx.In(queryString, *ids)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query) // sqlx.In returns queries with the `?` bindvar, rebind it here for matching the database in used (e.g. postgre, oracle etc, can skip it if you use mysql)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveryAreaPolygons []*DeliveryAreaPolygon

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var polygon DeliveryAreaPolygon
		if err := rows.Scan(&polygon.VendorDeliveryAreaID, &polygon.DeliveryAreaPolygon, &polygon.LastUpdatedOn); err != nil {
			return nil, err
		}
		deliveryAreaPolygons = append(deliveryAreaPolygons, &polygon)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return deliveryAreaPolygons, nil
}

func insertAreas(pricingDb *sqlx.DB, areas *[]VendorArea) error {
	insertScript := `INSERT INTO dbo."vendor_delivery_area" ("id", "vendor_id", "porygon_delivery_area_id", "delivery_area_status",
	"delivery_fee_type", "delivery_fee_value", "minimum_order_value", "municipality_tax_type", "municipality_tax_value", "tourist_tax_type", "tourist_tax_value") 
   	VALUES %s`

	columnNumber := 11
	valueStrings := make([]string, len(*areas))
	valueArgs := make([]interface{}, len(*areas)*columnNumber)
	for index, post := range *areas {
		valueStrings[index*columnNumber] = fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
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
		valueArgs[index*columnNumber+4] = post.DeliveryFeeType
		valueArgs[index*columnNumber+5] = post.DeliveryFeeValue
		valueArgs[index*columnNumber+6] = post.MinimumOrderValue
		valueArgs[index*columnNumber+7] = post.MunicipalityTaxType
		valueArgs[index*columnNumber+8] = post.MunicipalityTaxValue
		valueArgs[index*columnNumber+9] = post.TouristTaxType
		valueArgs[index*columnNumber+10] = post.TouristTaxValue
	}
	stmt := fmt.Sprintf(insertScript, strings.Join(valueStrings, ","))
	_, err := pricingDb.Exec(stmt, valueArgs...)
	return err
}

func insertPolygons(pricingDb *sqlx.DB, deliveryAreaPolygons []*DeliveryAreaPolygon) error {
	insertScript := `INSERT INTO dbo.vendor_delivery_area_polygon (vendor_delivery_area_id, delivery_area_hash, delivery_area_polygon)
					 VALUES %s`

	columnNumber := 3
	valueStrings := make([]string, len(deliveryAreaPolygons))
	valueArgs := make([]interface{}, len(deliveryAreaPolygons)*columnNumber)
	for index, p := range deliveryAreaPolygons {
		valueStrings[index*columnNumber] = fmt.Sprintf("($%d,$%d,$%d)",
			index*columnNumber+1,
			index*columnNumber+2,
			index*columnNumber+3,
		)
		deliveryAreaPolygon := fmt.Sprintf("SRID=%d;%s", p.DeliveryAreaPolygon.SRID(), p.DeliveryAreaPolygonInsertString)
		valueArgs[index*columnNumber+0] = p.VendorDeliveryAreaID
		valueArgs[index*columnNumber+1] = p.DeliveryAreaHash
		valueArgs[index*columnNumber+2] = deliveryAreaPolygon
	}
	stmt := fmt.Sprintf(insertScript, strings.Join(valueStrings, ","))
	_, err := pricingDb.Exec(stmt, valueArgs...)
	return err
}

func bulkMigrate(vendorDb *sqlx.DB, pricingDb *sqlx.DB, areaIds *[]int64, chunkSize int) {
	iteration := len(*areaIds) / chunkSize
	var wg sync.WaitGroup
	for i := 0; i <= iteration; i++ {
		wg.Add(1)
		areaIdsChunk := getChunk(areaIds, chunkSize, i)
		go func(wg *sync.WaitGroup, vendorDb *sqlx.DB, pricingDb *sqlx.DB, areaIds *[]int64) {
			defer wg.Done()
			areas, _ := getAreas(vendorDb, areaIdsChunk)
			if areas != nil {
				_ = insertAreas(pricingDb, areas)
			}
			// areas, err := getAreas(vendorDb, areaIdsChunk)
			// if err != nil {
			// 	errChan <- err
			// } else if err := insertAreas(pricingDb, areas); err != nil {
			// 	errChan <- err
			// }
		}(&wg, vendorDb, pricingDb, areaIdsChunk)
	}
	wg.Wait()
}

func bulkMigrateV3(vendorDb *sqlx.DB, pricingDb *sqlx.DB, areaIds *[]int64, chunkSize int) error {
	iteration := len(*areaIds) / chunkSize
	var g errgroup.Group
	for i := 0; i <= iteration; i++ {
		areaIdsChunk := getChunk(areaIds, chunkSize, i)
		func(vendorDb *sqlx.DB, pricingDb *sqlx.DB, areaIds *[]int64) {
			g.Go(func() error {
				areas, err := getAreas(vendorDb, areaIdsChunk)
				if err != nil {
					fmt.Println(fmt.Errorf("error getting areas in iteration:%d %v", i, err))
				}

				if err := insertAreas(pricingDb, areas); err != nil {
					fmt.Println(fmt.Errorf("error inserting areas in iteration:%d %v", i, err))
				}

				polygons, err := getPolygons(vendorDb, areaIdsChunk)
				if err != nil {
					fmt.Println(fmt.Errorf("error getting polygons in iteration:%d %v", i, err))
				}

				preparePolygonsForInsert(polygons)
				if err := insertPolygons(pricingDb, polygons); err != nil {
					fmt.Println(fmt.Errorf("error inserting polygons in iteration:%d %v", i, err))
				}
				return nil
			})
		}(vendorDb, pricingDb, areaIdsChunk)
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func preparePolygonsForInsert(polygons []*DeliveryAreaPolygon) {
	for _, p := range polygons {
		coordinates := p.DeliveryAreaPolygon.Coords()
		hash, err := generateHashForPolygon(coordinates)
		if err != nil {
			fmt.Printf("Error while polygon hashing %v", err)
			continue
		}
		p.DeliveryAreaHash = *hash
		p.DeliveryAreaPolygonInsertString = generatePolygonString(coordinates)
	}
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

func bulkMigrate2(vendorDb *sqlx.DB, pricingDb *sqlx.DB, areaIds *[]int64, chunkSize int) (chan error, chan int) {
	iteration := len(*areaIds) / chunkSize
	errOuterChan := make(chan error, iteration*2)
	quit := make(chan int)
	go func(errChan chan error, quit chan int) {
		errInnerChan := make(chan error, iteration*2)
		var wg *sync.WaitGroup
		for i := 0; i <= iteration; i++ {
			select {
			case err := <-errInnerChan:
				errOuterChan <- err
			default:
				wg.Add(1)
				areaIdsChunk := getChunk(areaIds, chunkSize, i)
				go func(wg *sync.WaitGroup, errChan chan error, vendorDb *sqlx.DB, pricingDb *sqlx.DB, areaIds *[]int64) {
					defer wg.Done()
					areas, err := getAreas(vendorDb, areaIdsChunk)
					if err != nil {
						errChan <- err
					} else if err := insertAreas(pricingDb, areas); err != nil {
						errChan <- err
					}
				}(wg, errInnerChan, vendorDb, pricingDb, areaIdsChunk)
			}
		}
		wg.Wait()
		quit <- 1
	}(errOuterChan, quit)

	return errOuterChan, quit
}

func getChunk[T any](collection *[]T, chunkSize, iteration int) *[]T {
	var howManyToPick = chunkSize
	if len(*collection) < (iteration*chunkSize + chunkSize) {
		howManyToPick = len(*collection) - (iteration * chunkSize)
	}
	from := iteration * chunkSize
	to := from + howManyToPick
	chunk := (*collection)[from:to]
	return &chunk
}

func getENV(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func main() {
	var elapsed time.Duration
	start := time.Now()
	areaBulkInsertChunkSize, err := strconv.Atoi(getENV("AREA_BULK_INSERT_CHUNK_SIZE", "999"))
	if err != nil {
		log.Fatal(err)
	}

	pricingDb, err := getPricingDb()
	if err != nil {
		log.Fatal(err)
	}
	defer pricingDb.Close()

	vendorDb, err := getVendorDb()
	if err != nil {
		log.Fatal(err)
	}
	defer vendorDb.Close()

	areaIds, err := getAreaIds(vendorDb)
	if err != nil {
		log.Fatal(err)
	}
	//areaIds := &[]int64{92711}
	elapsed = time.Since(start)
	fmt.Printf("Pulled %d many area Ids and took %s", len(*areaIds), elapsed)

	bulkStarted := time.Now()

	if err := bulkMigrateV3(vendorDb, pricingDb, areaIds, areaBulkInsertChunkSize); err != nil {
		log.Fatal(err)
	}

	elapsed = time.Since(bulkStarted)
	log.Printf("Bulk insert took %s", elapsed)

	elapsed = time.Since(start)
	log.Printf("Whole Migration took %s", elapsed)
}
