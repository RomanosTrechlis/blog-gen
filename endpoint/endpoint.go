package endpoint

type Endpoint interface {
	Upload(to string) error
}
