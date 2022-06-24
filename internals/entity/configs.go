package entity

type PSQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func DefaultPSQLConfig() *PSQLConfig {
	return &PSQLConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

}
