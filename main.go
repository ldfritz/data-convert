// This package manages converting files to and from JSON and CSV.
// Unless specified the source format is based off the file extension.
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	source := os.Args[1]
	ext := filepath.Ext(source)
	switch {
	case strings.EqualFold(ext, ".csv"):
		fmt.Println(ExportJSON(ImportCSV(source)))
	case strings.EqualFold(ext, ".json"):
		fmt.Println(ExportCSV(ImportJSON(source)))
	default:
		fmt.Printf("error: unexpected file extension: %s\n", ext)
		return
	}
}

func ImportCSV(file string) []map[string]interface{} {
	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	header, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	var result []map[string]interface{}
	for _, l := range lines {
		m := make(map[string]interface{})
		for i, h := range header {
			m[h] = l[i]
		}
		result = append(result, m)
	}
	return result
}

func ExportCSV(data []map[string]interface{}) string {
	var keys []string
	for k := range data[0] {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	records := [][]string{keys}
	for _, v := range data {
		var row []string
		for _, k := range keys {
			row = append(row, v[k].(string))
		}
		records = append(records, row)
	}
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	w.WriteAll(records)
	return b.String()
}

func ImportJSON(filename string) []map[string]interface{} {
	var result []map[string]interface{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(content, &result)
	return result
}

func ExportJSON(data []map[string]interface{}) string {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}
