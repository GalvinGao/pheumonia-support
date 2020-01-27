package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	cookie = "xxx" // input your cookie here
)

var (
	client = http.Client{
		Timeout: time.Minute,
	}
	Logger     *log.Logger
	folder     = "files/"
	sheetIndex int
)

type Response struct {
	DownloadURL string `json:"downloadUrl"`
}

func get(name string, guid string) {
	Logger.Printf("fetching download url from api...")

	val := url.Values{}
	val.Add("type", "xlsx")
	val.Add("file", guid)
	val.Add("returnJson", "1")
	val.Add("name", name)
	val.Add("isAsync", "0")
	query := val.Encode()

	u := fmt.Sprintf("https://shimo.im/lizard-api/files/%s/export?%s", guid, query)
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		Logger.Printf("failed to fetch http resource: %v", err)
		panic(err)
	}
	request.Header.Set("Cookie", cookie)
	request.Header.Set("Referer", fmt.Sprintf("https://shimo.im/sheets/%s/MODOC", guid))

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		Logger.Printf("failed to fetch http resource: status code %v", response.StatusCode)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var resp Response
	err = json.Unmarshal(responseBytes, &resp)
	if err != nil {
		panic(err)
	}

	Logger.Printf("got download url: %s", resp.DownloadURL)
	Logger.Printf("downloading file...")

	response, err = client.Get(resp.DownloadURL)
	if err != nil {
		panic(err)
	}

	parse(name, response)
}

func parse(name string, response *http.Response) {
	Logger.Printf("got file. converting...")
	xlsxBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	xlsxFile, err := xlsx.OpenBinary(xlsxBytes)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(filepath.Join(folder, name+".csv"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	cw := csv.NewWriter(f)
	sheet := xlsxFile.Sheets[sheetIndex]
	var vals []string
	for _, row := range sheet.Rows {
		if row != nil {
			vals = vals[:0]
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					vals = append(vals, err.Error())
				}
				vals = append(vals, str)
			}
		}
		cw.Write(vals)
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 + 1 {
		panic("not enough argument. usage: ./executable [guid] [name] (sheetId)")
	}
	guid := os.Args[1]
	name := os.Args[2]
	var sheetId string
	if len(os.Args) > 3 {
		sheetId = os.Args[3]
	} else {
		sheetId = ""
	}

	Logger = log.New(os.Stdout, "[main] ", log.LstdFlags)

	if guid == "" || name == "" {
		Logger.Printf("need guid and name. usage: ./executable [guid] [name] (sheetId)")
		panic("missing required params")
	}
	if sheetId == "" {
		Logger.Printf("not providing sheetIndex id. using 0 (the first sheetIndex)")
	} else {
		s, err := strconv.ParseInt(sheetId, 10, 32)
		if err != nil {
			panic(err)
		}
		sheetIndex = int(s)
		Logger.Printf("set sheetIndex id to %v", sheetIndex)
	}
	Logger.Printf("converting shimo.im document, guid %s, name %s", guid, name)

	err := os.MkdirAll(folder, 0755)
	if err != nil {
		Logger.Printf("failed to create folder: %v", err)
	}
	get(name, guid)
	Logger.Printf("converted file and saved at %s", filepath.Join(folder, name + ".csv"))
}
