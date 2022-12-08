package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"gorm.io/driver/postgres"
        "gorm.io/gorm"
        "gorm.io/gorm/logger"
)

type Ips struct {
        gorm.Model
        IP       string
        Hostname string
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", fromAgent)
	e.GET("/", fromCppm)

	e.Logger.Fatal(e.Start(":5678"))
}

func fromAgent(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusNotFound, "Failed")
	}
	host := string(body)

	dsn := "host=10.2.13.132 user=admin password=admin dbname=postgresdb port=5432"
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
                Logger: logger.Default.LogMode(logger.Info),
        })
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

        if err != nil {
                panic(err)
        }
        db.AutoMigrate(&Ips{})

	var ips Ips

	ip := strings.Split(host, "\t")[0]
	hostname := strings.Split(host, "\t")[1]

	result := db.First(&ips, "Hostname = ?", hostname)

	if result.RowsAffected == 0 {
		db.Create(&Ips{IP: ip, Hostname: hostname})
	}

	return c.String(http.StatusOK, string(body))
}

func fromCppm(c echo.Context) error {
	dsn := "host=10.2.13.132 user=admin password=admin dbname=postgresdb port=5432"
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
                Logger: logger.Default.LogMode(logger.Info),
        })
        if err != nil {
                panic(err)
        }
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	db.AutoMigrate(&Ips{})

	iptable := []Ips{}

        result := db.Find(&iptable)

	var body string

        if result.RowsAffected != 0 {
	        for _, ip := range iptable {
               		body += ip.IP + " "
			body += ip.Hostname + "\n"
        	}
	}
	return c.String(http.StatusOK, string(body))
}
