package pkg

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var apiDateTimeFormat = "2006-01-02T15:04:05.999999999Z"

func renderTable(header *table.Row, rows *[]table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(*header)
	t.AppendRows(*rows)
	t.Render()
}

func getLocation() *time.Location {
	location, err := time.LoadLocation(viper.GetString("location"))
	if err != nil {
		log.Fatal(err)
	}
	return location
}

func localDateTimeString(apiDate string) string {
	t, err := time.Parse(apiDateTimeFormat, apiDate)
	if err != nil {
		return ""
	}
	return t.In(getLocation()).Format("2006-01-02 15:04:05")
}
