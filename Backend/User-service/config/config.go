package config

import "os"


type Config struct {
DatabaseDSN string
JWTSecret string
Port string
}


func Load() *Config {
return &Config{
DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=users port=5432 sslmode=disable TimeZone=UTC"),
JWTSecret: getEnv("JWT_SECRET", "supersecret"),
Port: getEnv("PORT", "8080"),
}
}


func getEnv(key, fallback string) string {
if v := os.Getenv(key); v != "" {
return v
}
return fallback
}