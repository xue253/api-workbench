package db

import (
	"fmt"
	"log"

	"api-workbench/internal/config"
	"api-workbench/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = DB.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Environment{},
		&model.EnvironmentVariable{},
		&model.Collection{},
		&model.API{},
		&model.Assertion{},
		&model.TestCase{},
		&model.TestCaseAPI{},
		&model.TestDataSet{},
		&model.TestSuite{},
		&model.TestSuiteCase{},
		&model.ScheduledTask{},
		&model.TestRun{},
		&model.TestRunDetail{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("database initialized")
}
