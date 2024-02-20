package db

import (
	"bloom/config"
	"bloom/db/models/taxon"
	"bloom/utils"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"slices"

	"github.com/ztrue/tracerr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), config.MongoDBConnTimeout)
	defer cancel()
	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDBURI))
	err = tracerr.Wrap(err)
	utils.CheckError(err)
	collectionNames, err := DB().ListCollectionNames(context.Background(), bson.D{})
	err = tracerr.Wrap(err)
	utils.CheckError(err)
	if !slices.Contains(collectionNames, "taxon") {
		fmt.Println("Rebuilding taxon collection")
		Build()
	}
}

func DB() *mongo.Database {
	return Client.Database(config.MongoDBName)
}

func Build() {
	taxonColl := DB().Collection("taxon")
	f, err := os.Open("sdata/WFOTaxonomicBackbone/classification.csv")
	err = tracerr.Wrap(err)
	utils.CheckError(err)
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true
	data, err := csvReader.ReadAll()
	utils.CheckError(err)
	defer f.Close()

	for i, line := range data {
		if i > 0 { // omit header line
			entry := taxon.Entry{
				ScientificName:       line[3],
				TaxonRank:            line[4],
				Family:               line[7],
				Subfamily:            line[8],
				Tribe:                line[9],
				Subtribe:             line[10],
				Genus:                line[11],
				Subgenus:             line[12],
				SpecificEpithet:      line[13],
				InfraspecificEpithet: line[14],
				TaxonomicStatus:      line[18],
			}
			_, err := taxonColl.InsertOne(context.Background(), entry)
			err = tracerr.Wrap(err)
			utils.CheckError(err)
		}
	}
}
