package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aydinmuzaffer/migration-tool-service/src/db/gormtransaction"
	deliveryareapolygonservice "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/service/delivery_area_polygon"
	migrationservice "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/service/migration"
	pricing_storage "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/storage/pricing"
	vendorlist "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/storage/vendor_list"
	"github.com/aydinmuzaffer/migration-tool-service/src/utils/envutils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getVendorDb(address string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(address), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getPricingDb(address string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(address), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var elapsed time.Duration
	start := time.Now()
	areaBulkInsertChunkSize, err := strconv.Atoi(envutils.GetENV("AREA_BULK_INSERT_CHUNK_SIZE", "999"))
	if err != nil {
		log.Fatal(err)
	}

	pricingDbAddressDefault := "host=localhost port=5433 user=pgsrv_pricing_owner password=pricing dbname=pricing_local sslmode=disable"
	pricingDbAddress := envutils.GetENV("PRICING_DB_ADDRESS", pricingDbAddressDefault)
	pricingDb, err := getPricingDb(pricingDbAddress)
	if err != nil {
		log.Fatal(err)
	}

	pdb, err := pricingDb.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := pdb.Ping(); err != nil {
		log.Fatal(err)
	}
	defer pdb.Close()

	vendorDbAddressDefault := "host=localhost port=8080 user=pgsrv_vendor_list_qa_app password=password dbname=vendor_list_qa"
	vendorDbAddress := envutils.GetENV("VENDOR_DB_ADDRESS", vendorDbAddressDefault)
	vendorDb, err := getVendorDb(vendorDbAddress)
	if err != nil {
		log.Fatal(err)
	}

	vdb, err := vendorDb.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := vdb.Ping(); err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()
	vendorDbWrapper := gormtransaction.NewGormDBWrapper(vendorDb)
	deliveryAreaGetter := vendorlist.NewVendorDeliveryAreaGetter(vendorDbWrapper)
	deliveryAreaPolygonGetter := vendorlist.NewVenderDeliveryAreaPolygonGetter(vendorDbWrapper)

	pricingDbWrapper := gormtransaction.NewGormDBWrapper(pricingDb)
	deliveryAreaSetter := pricing_storage.NewVenderDeliveryAreaSetter(pricingDbWrapper)
	deliveryAreaPolygonSetter := pricing_storage.NewVenderDeliveryAreaPolygonSetter(pricingDbWrapper)

	polygonBuilder := deliveryareapolygonservice.NewPolygonBuilder()
	polygonService := deliveryareapolygonservice.NewPolygonService(polygonBuilder)
	migrationservice := migrationservice.NewMigrationService(deliveryAreaGetter, deliveryAreaSetter,
		deliveryAreaPolygonGetter, deliveryAreaPolygonSetter, polygonService, areaBulkInsertChunkSize)

	pricingTransactor := gormtransaction.NewTransactor(pricingDb)

	trErr := pricingTransactor.WithinTransaction(ctx,
		func(txCtx context.Context) error {
			return migrationservice.Migrate(ctx)
		})
	if trErr != nil {
		log.Fatal(err)
	}
	elapsed = time.Since(start)
	log.Printf("Whole Migration took %s", elapsed)
}
