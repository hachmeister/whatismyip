package router

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/dcarrillo/whatismyip/models"
	"github.com/dcarrillo/whatismyip/service"
	"github.com/gin-gonic/gin"
)

type geoDataFormatter struct {
	title  string
	format func(*models.GeoRecord) string
}

type asnDataFormatter struct {
	title  string
	format func(*models.ASNRecord) string
}

var geoOutput = map[string]geoDataFormatter{
	"country": {
		title: "Country",
		format: func(record *models.GeoRecord) string {
			return record.Country.Names["en"]
		},
	},
	"country_code": {
		title: "Country Code",
		format: func(record *models.GeoRecord) string {
			return record.Country.ISOCode
		},
	},
	"city": {
		title: "City",
		format: func(record *models.GeoRecord) string {
			return record.City.Names["en"]
		},
	},
	"latitude": {
		title: "Latitude",
		format: func(record *models.GeoRecord) string {
			return fmt.Sprintf("%f", record.Location.Latitude)
		},
	},
	"longitude": {
		title: "Longitude",
		format: func(record *models.GeoRecord) string {
			return fmt.Sprintf("%f", record.Location.Longitude)
		},
	},
	"postal_code": {
		title: "Postal Code",
		format: func(record *models.GeoRecord) string {
			return record.Postal.Code
		},
	},
	"time_zone": {
		title: "Time Zone",
		format: func(record *models.GeoRecord) string {
			return record.Location.TimeZone
		},
	},
}

var asnOutput = map[string]asnDataFormatter{
	"number": {
		title: "ASN Number",
		format: func(record *models.ASNRecord) string {
			return fmt.Sprintf("%d", record.AutonomousSystemNumber)
		},
	},
	"organization": {
		title: "ASN Organization",
		format: func(record *models.ASNRecord) string {
			return record.AutonomousSystemOrganization
		},
	},
}

func getGeoAsString(ctx *gin.Context) {
	field := strings.ToLower(ctx.Params.ByName("field"))
	ip := service.Geo{IP: net.ParseIP(ctx.ClientIP())}
	record := ip.LookUpCity()

	if field == "" {
		ctx.String(http.StatusOK, geoCityRecordToString(record))
	} else if g, ok := geoOutput[field]; ok {
		ctx.String(http.StatusOK, g.format(record)+"\n")
	} else {
		ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound)+"\n")
	}
}

func getASNAsString(ctx *gin.Context) {
	field := strings.ToLower(ctx.Params.ByName("field"))
	ip := service.Geo{IP: net.ParseIP(ctx.ClientIP())}
	record := ip.LookUpASN()

	if field == "" {
		ctx.String(http.StatusOK, geoASNRecordToString(record))
	} else if g, ok := asnOutput[field]; ok {
		ctx.String(http.StatusOK, g.format(record)+"\n")
	} else {
		ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound)+"\n")
	}
}

func geoCityRecordToString(record *models.GeoRecord) string {
	var output string

	keys := make([]string, 0, len(geoOutput))
	for k := range geoOutput {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		output += fmt.Sprintf("%s: %v\n", geoOutput[k].title, geoOutput[k].format(record))
	}

	return output
}

func geoASNRecordToString(record *models.ASNRecord) string {
	var output string

	keys := make([]string, 0, len(asnOutput))
	for k := range asnOutput {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		output += fmt.Sprintf("%s: %v\n", asnOutput[k].title, asnOutput[k].format(record))
	}

	return output
}
