package schemas

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

type Deploy struct {
	Port int
}

type Config struct {
	DB     DatabaseConfig
	Deploy Deploy
}
