package main

//使用database模块示范

import (
    ...
    "your_project_name/database"
    ...
)

func main() {
    // 打开ASN数据库
    asnDB, err := database.OpenASNDB("/data/ipinfo/db/asn.mmdb")
    ...
    // 打开国家数据库
    countryDB, err := database.OpenCountryDB("/data/ipinfo/db/country.mmdb")
    ...
}

...
func ipLookupHandler(w http.ResponseWriter, r *http.Request) {
    ...
    // 查询ASN记录
    var asn database.ASNRecord
    err := asnDB.Lookup(ip, &asn)
    ...
    // 查询国家记录
    var country database.CountryRecord
    err = countryDB.Lookup(ip, &country)
    ...
}
