package dbconn

import (
	"embed"
	"fmt"
	"github.com/almartino/entando-go/datasource"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New instantiates a database connection by using the Entando database configuration.
func New(fs embed.FS) *gorm.DB {
	config := datasource.NewConfiguration()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Rome", config.Host, config.Username, config.Password, config.Name, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := migrateDb(&config, fs); err != nil {
		panic(err)
	}
	return db
}

func migrateDb(config *datasource.Configuration, fs embed.FS) error {
	userPwd := fmt.Sprintf("%s:%s", config.Username, config.Password)
	url := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable", userPwd, config.Host, config.Port, config.Name)
	source, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
