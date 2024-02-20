package users

type Entry struct {
	Username    string
	Email       string
	Password    string
	Permissions int
	DateCreated int64
	Lat         float64
	Long        float64
	Garden      []string
}
