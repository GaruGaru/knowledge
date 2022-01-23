package conf

type DatabaseType string

const (
	DatabaseTypePostgres DatabaseType = "postgres"
	DatabaseTypeSQLite   DatabaseType = "sqlite"
	DatabaseTypeMySql    DatabaseType = "mysql"
)

type Catalog struct {
	Database Database `json:"database" yaml:"database"`
}

type Database struct {
	Type   DatabaseType           `json:"type" yaml:"type"`
	Params map[string]interface{} `json:"params" yaml:"params"`
}
