package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Address struct {
	City    string `csv:"kota"`
	Country string `csv:"negara"`
}

type User struct {
	Name string `csv:"nama"`
	Address
	Age       int       `csv:"umur,omitempty"`
	CreatedAt time.Time `csv:"createdAt"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/user", UserFunc)
	e.GET("/download", Download)
	e.GET("/download/v2", DownloadUpd)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":9090"))
}

func UserFunc(c echo.Context) error {
	u := &User{
		Name:      "Fajar",
		Address:   Address{"Jakarta", "Indonesia"},
		Age:       20,
		CreatedAt: time.Now().Local().Add(time.Hour * 1),
	}
	return c.JSON(http.StatusOK, u)

}

func Download(c echo.Context) error {
	users := []User{
		{
			Name:      "Fajar",
			Address:   Address{"Jakarta", "Indonesia"},
			Age:       20,
			CreatedAt: time.Now().Local().Add(time.Hour * 1),
		},
		{
			Name:      "Ahmad",
			Address:   Address{"Bandung", "Indonesia"},
			Age:       21,
			CreatedAt: time.Now().Local().Add(time.Hour * 2),
		},
		{
			Name:      "Islami",
			Address:   Address{"Bengkulu", "Indonesia"},
			Age:       22,
			CreatedAt: time.Now().Local().Add(time.Hour * 3),
		},
	}

	b, err := csvutil.Marshal(users)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Write(b)
	return c.Attachment("file.csv", "file.csv")

}

func DownloadUpd(c echo.Context) error {
	users := []User{
		{
			Name:      "Fajar",
			Address:   Address{"Jakarta", "Indonesia"},
			Age:       20,
			CreatedAt: time.Now().Local().Add(time.Hour * 1),
		},
		{
			Name:      "Ahmad",
			Address:   Address{"Bandung", "Indonesia"},
			Age:       21,
			CreatedAt: time.Now().Local().Add(time.Hour * 2),
		},
		{
			Name:      "Islami",
			Address:   Address{"Bengkulu", "Indonesia"},
			Age:       22,
			CreatedAt: time.Now().Local().Add(time.Hour * 3),
		},
	}

	filename := "./file/report" + time.Now().Format("2006-01-02") + ".csv"
	csvFile, err := os.Create(filename)

	if err != nil {
		log.Println(err)
	}

	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	if err := csvutil.NewEncoder(w).Encode(users); err != nil {
		fmt.Println("error:", err)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Println("error:", err)
	}
	return c.Attachment(filename, "tes.csv")

}
