package database

import (
    "log"

    "github.com/oschwald/maxminddb-golang"
)

type ASNRecord struct {
    ASN    string `maxminddb:"asn"`
    Domain string `maxminddb:"domain"`
    Name   string `maxminddb:"name"`
}

type CountryRecord struct {
    Continent    string `maxminddb:"continent"`
    Continent_name string `maxminddb:"continent_name"`
    Country   string `maxminddb:"country"`
    Country_name   string `maxminddb:"country_name"`
}

func OpenASNDB(path string) (*maxminddb.Reader, error) {
    db, err := maxminddb.Open(path)
    if err != nil {
        log.Fatal("Error opening ASN database:", err)
    }
    return db, err
}

func OpenCountryDB(path string) (*maxminddb.Reader, error) {
    db, err := maxminddb.Open(path)
    if err != nil {
        log.Fatal("Error opening country database:", err)
    }
    return db, err
}
