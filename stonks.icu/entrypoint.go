package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ericm/stonks/api"
	"github.com/ericm/stonks/graph"
	"github.com/oschwald/geoip2-golang"
	"github.com/piquette/finance-go/datetime"
	"github.com/spf13/viper"
)

const (
	footer     = "\nLike Stonks? Star it on GitHub: https://github.com/ericm/stonks\nstonks " + api.Version + "\n"
	geoipFile  = "GeoLite2-City.mmdb"
	geoipDload = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"
	defaultTZ  = "Europe/Dublin"
)

var db *geoip2.Reader

func main() {
	configure()
	log.Println("Downloading Geoip2 DB")
	geoData, err := http.Get(fmt.Sprintf(geoipDload, viper.GetString("geolicense")))
	if err != nil {
		log.Panic(err)
	}
	gzipReader, err := gzip.NewReader(geoData.Body)
	if err != nil {
		log.Panic(err)
	}
	tarReader := tar.NewReader(gzipReader)
	if err != nil {
		log.Panic(err)
	}
	var tarData []byte
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if filepath.Base(header.Name) == geoipFile {
			if tarData, err = ioutil.ReadAll(tarReader); err != nil {
				log.Panic(err)
			}
		}
	}
	if db, err = geoip2.FromBytes(tarData); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handleSymbol)
	fmt.Printf("Server listening on port %d\n", viper.GetInt("port"))
	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), nil)
}

func handleSymbol(w http.ResponseWriter, r *http.Request) {
	symbols := strings.Split(r.URL.Path, "/")
	plainText := false
	host := r.Header.Get("X-FORWARDED-FOR")
	if host == "" {
		var err error
		host, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	ip := net.ParseIP(host)
	record, err := db.City(ip)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tz := record.Location.TimeZone
	if tz == "" {
		tz = defaultTZ
	}
	location, err := time.LoadLocation(tz)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for _, client := range clients {
		if strings.Contains(r.Header.Get("User-Agent"), client+"/") {
			plainText = true
			break
		}
	}
	if !plainText {
		w.Header().Add("Location", "https://github.com/ericm/stonks")
		w.WriteHeader(302)
		w.Write([]byte(" "))
		return
	}
	output := ""
	num := 0
	for _, symbol := range symbols {
		if len(symbol) > 0 {
			num++
			symbol = strings.ToUpper(symbol)
			chart, err := api.GetChart(symbol, datetime.FifteenMins, nil, nil, false)
			if err != nil {
				w.WriteHeader(403)
				w.Write([]byte(err.Error() + "\n" + footer))
				return
			}
			out, err := graph.GenerateGraph(chart, 80, 12, graph.LineTheme, location)
			if err != nil {
				w.WriteHeader(403)
				w.Write([]byte(err.Error() + "\n" + footer))
				return
			}
			output += out
		}
	}
	w.WriteHeader(200)
	if num == 0 {
		w.Write([]byte("Please provide stonks in the format:\nstonks.icu/amd/intl\n" + footer))
	} else {
		w.Write([]byte(output + footer))
	}
}
