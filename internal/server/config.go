package server

type PostgresConfig struct {
	ConnectionString string `yaml:"conn"`
}

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
}
