package observations

type Entry struct {
	UserID         string
	TaxonID        string
	Lat            float64
	Long           float64
	BloomStartDate int64
	BloomPeakDate  int64
	BloomEndDate   int64
	DateCreated    int64
	DateModified   int64
}
