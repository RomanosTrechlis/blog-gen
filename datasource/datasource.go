package datasource

// DataSource fetches data from an endpoint
type DataSource interface {
	Fetch(from, to string) ([]string, error)
}
