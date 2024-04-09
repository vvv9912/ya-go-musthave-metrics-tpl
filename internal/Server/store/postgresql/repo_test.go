package postgresql

import (
	"context"
	"github.com/stretchr/testify/suite"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store"
	"testing"
	"time"
)

type testStorager interface {
	store.Storager
	clean(ctx context.Context) error
}
type Config struct {
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration
	Username       string
	Password       string
	DBName         string
	MigrationVer   int

	Host string
	Port int
}
type PostrgresTestSuite struct {
	suite.Suite
	testStorager

	tc  *tcpostgres.PostgresContainer
	cfg *Config
}

func (ts *PostrgresTestSuite) SetupSuite() {
	//todo Add Function Migrate
	//cfg := &Config{
	//	ConnectTimeout: 5 * time.Second,
	//	QueryTimeout:   5 * time.Second,
	//	Username:       "postgres",
	//	Password:       "test",
	//	DBName:         "postgres",
	//	MigrationVer:   1,
	//}
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()
	//
	//pgc, err := tcpostgres.RunContainer(ctx,
	//	testcontainers.WithImage("postgres:latest"),
	//	tcpostgres.WithDatabase(cfg.DBName),
	//	tcpostgres.WithUsername(cfg.Username),
	//	tcpostgres.WithPassword(cfg.Password),
	//	testcontainers.WithWaitStrategy(
	//		wait.ForLog("database system is ready to accept connection").
	//			WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	//)
	//
	//require.NoError(ts.T(), err)
	//
	//cfg.Host, err = pgc.Host(ctx)
	//require.NoError(ts.T(), err)
	//
	//port, err := pgc.MappedPort(ctx, "5432")
	//require.NoError(ts.T(), err)
	//
	//cfg.Port = port.Int()
	//
	//ts.tc = pgc
	//ts.cfg = cfg
	//
	//database_dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	//log.Println(database_dsn)
	//
	//db, err := sqlx.Connect("postgres", database_dsn)
	//
	//storage := NewDatabase(db)
	//
	//ts.testStorager = storage
	//
	//err = Migrate(db) //todo add funct migrate
	//require.NoError(ts.T(), err)
	//
	//ts.T().Logf("stared postgres at %s:%d", cfg.Host, cfg.Port)

}
func TestDatabase_getAllCounter(t *testing.T) {

}

func TestDatabase_getAllGauge(t *testing.T) {

}

func TestDatabase_getCounter(t *testing.T) {

}

func TestDatabase_getGauge(t *testing.T) {

}

func TestDatabase_updateCounter(t *testing.T) {

}

func TestDatabase_updateGauge(t *testing.T) {

}

func TestDatabase_updateMetricsBatch(t *testing.T) {

}

func TestNewDatabase(t *testing.T) {
}
