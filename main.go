package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func createStore(filePath string) map[string]string {
	fi, _ := os.Stat(filePath)
	return make(map[string]string, fi.Size())
}

func createRecord(keys []string, values []string) string {
	var b strings.Builder
	b.WriteByte('{')
	i := 0
	for ; i < len(keys)-1; i++ {
		b.WriteString(fmt.Sprintf("\"%s\": \"%s\",", keys[i], values[i]))
	}
	b.WriteString(fmt.Sprintf("\"%s\": \"%s\"}", keys[i], values[i]))
	return b.String()
}

func readCsvFile(filePath string, recordsStore map[string]string) {
	f, _ := os.Open(filePath)
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()
	for i := 1; i < len(records); i++ {
		recordsStore[records[i][1]] = createRecord(records[0], records[i])
	}
}

func createArrRecords(ids []string, recordsStore map[string]string) string {
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < len(ids); {
		record, ok := recordsStore[ids[i]]
		if ok {
			b.WriteString(fmt.Sprintf("%s", record))
		}
		i++
		for ; i < len(ids); i++ {
			if _, ok := recordsStore[ids[i]]; ok {
				b.WriteByte(',')
				break
			}
		}
	}
	b.WriteByte(']')
	return b.String()
}

func getItems(recordsStore map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query()["id"]
		w.Header().Set("Content-Type", "application/json")
		items := createArrRecords(id, recordsStore)
		w.Write([]byte(items))
	}
}

func main() {
	filename := "./ueba.csv"
	recordsStore := createStore(filename)
	readCsvFile(filename, recordsStore)
	http.HandleFunc("/get-items", getItems(recordsStore))
	http.ListenAndServe(":3333", nil)
}
