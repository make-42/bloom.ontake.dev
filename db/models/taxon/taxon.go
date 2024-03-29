package taxon

type Entry struct {
	ScientificName       string
	TaxonRank            string
	Family               string
	Subfamily            string
	Tribe                string
	Subtribe             string
	Genus                string
	Subgenus             string
	SpecificEpithet      string
	InfraspecificEpithet string
	TaxonomicStatus      string
}

type EntryIDED struct {
	ScientificName       string
	TaxonRank            string
	Family               string
	Subfamily            string
	Tribe                string
	Subtribe             string
	Genus                string
	Subgenus             string
	SpecificEpithet      string
	InfraspecificEpithet string
	TaxonomicStatus      string
	ID                   string `json:"id" bson:"_id"`
}
