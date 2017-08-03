package datasource

import (
	"fmt"
)

// DataSource fetches data from an endpoint
type DataSource interface {
	Fetch(from, to string) ([]string, error)
}

// New is a data source factory
func New(dsType string) (ds DataSource, err error) {
	switch dsType {
	case "git":
		ds = newGitDataSource()
	case "local":
		ds = newLocalDataSource()
	case "":
		err = fmt.Errorf("please provide a datasource in the configuration file")
	}
	return ds, err
}
