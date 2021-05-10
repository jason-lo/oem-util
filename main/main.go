package main

import (
	"fmt"
	"github.com/jason-lo/oem-util/sheet"
)

func main() {
	//https://docs.google.com/spreadsheets/d/1R8S1fwJTiPTYTqdLu85fZ4kCseUrEsvPl09aoM0AZKU/edit#gid=1623281831
	spreadsheetId := "1R8S1fwJTiPTYTqdLu85fZ4kCseUrEsvPl09aoM0AZKU"
	spreadsheetRange := "CMIT/WKS 2021!A2:BF2"
	//indexof := get_index_map(spreadsheetId, spreadsheetRange)
	indexof := sheet.Getindex(spreadsheetId, spreadsheetRange)

	for k, v := range indexof {
		fmt.Printf("key[%s] value[%d]\n", k, v)
	}
}
