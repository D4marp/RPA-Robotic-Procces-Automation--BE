package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"rpa-backend/internal/config"
)

func TestGetEnvReturnsValueWhenSet(t *testing.T) {
	t.Setenv("TEST_KEY", "hello")
	assert.Equal(t, "hello", config.GetEnv("TEST_KEY", "default"))
}

func TestGetEnvReturnsFallbackWhenUnset(t *testing.T) {
	os.Unsetenv("UNSET_KEY")
	assert.Equal(t, "default", config.GetEnv("UNSET_KEY", "default"))
}

func TestMySQLDSNFromEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_USER", "rpa")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_NAME", "rpa_db")

	dsn := config.MySQLDSN()
	assert.Contains(t, dsn, "rpa:secret@tcp(db.example.com:3307)/rpa_db")
	assert.Contains(t, dsn, "charset=utf8mb4")
}

func TestMySQLDSNFromDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local")
	assert.Equal(t, "user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local", config.MySQLDSN())
}

func TestMySQLDSNFromStandardDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "mysql://myuser:mypass@db-host:3306/rpa_db?charset=utf8mb4")
	dsn := config.MySQLDSN()
	assert.Equal(t, "myuser:mypass@tcp(db-host:3306)/rpa_db?charset=utf8mb4&parseTime=True&loc=Local", dsn)
}

func TestMySQLDSNFromStandardDatabaseURLNoProtocolPrefix(t *testing.T) {
	t.Setenv("DATABASE_URL", "myuser:mypass@db-host.example.com:3306/rpa_db")
	dsn := config.MySQLDSN()
	assert.Equal(t, "myuser:mypass@tcp(db-host.example.com:3306)/rpa_db?charset=utf8mb4&parseTime=True&loc=Local", dsn)
}
