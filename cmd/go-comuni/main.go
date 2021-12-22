package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"

	"github.com/panta/go-comuni-italiani"
)

const (
	CODICI_ISTAT_COMUNI_CSV_URL = "https://www.istat.it/storage/codici-unita-amministrative/Elenco-comuni-italiani.csv"
)

type csvFieldEntry struct {
	FieldName         string
	FieldTypeName     string
	TagName           string
	HeaderColumnName  string
	HeaderColumnIndex int
}

func (entry *csvFieldEntry) getField(obj interface{}) reflect.Value {
	pointToStruct := reflect.ValueOf(obj)
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}
	curField := curStruct.FieldByName(entry.FieldName) // type: reflect.Value
	if !curField.IsValid() {
		panic("not found/valid:" + entry.FieldName)
	}
	return curField
}

type csvFields struct {
	Type               reflect.Type
	TagToEntry         map[string]*csvFieldEntry
	ColumnNameToEntry  map[string]*csvFieldEntry
	ColumnIndexToEntry map[int]*csvFieldEntry
}

func parseStructTags(header []string) csvFields {
	headerIndexMap := map[string]int{}

	for headerIndex, headerColumnName := range header {
		headerIndexMap[headerColumnName] = headerIndex
	}

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(comuni.Comune{})

	csvFields := csvFields{
		Type:               t,
		TagToEntry:         make(map[string]*csvFieldEntry),
		ColumnNameToEntry:  make(map[string]*csvFieldEntry),
		ColumnIndexToEntry: make(map[int]*csvFieldEntry),
	}

	for iterHeaderColumnName, iterHeaderIndex := range headerIndexMap {
		entry := csvFieldEntry{
			TagName:           "",
			HeaderColumnName:  iterHeaderColumnName,
			HeaderColumnIndex: iterHeaderIndex,
		}
		csvFields.ColumnNameToEntry[iterHeaderColumnName] = &entry
		csvFields.ColumnIndexToEntry[iterHeaderIndex] = &entry
	}

	tagName := "csv"

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get(tagName)

		var entry *csvFieldEntry
		headerColumnName := ""
		headerIndex := -1
		for iterHeaderColumnName, iterHeaderIndex := range headerIndexMap {
			iterHeaderColumnName = strings.TrimSpace(iterHeaderColumnName)
			if strings.HasSuffix(tag, "...") {
				prefix := tag[:len(tag)-3]
				prefixLen := len(prefix)
				if len(iterHeaderColumnName) >= prefixLen &&
					strings.EqualFold(strings.TrimSpace(iterHeaderColumnName[:prefixLen]), strings.ToLower(strings.TrimSpace(prefix))) {
					headerColumnName = iterHeaderColumnName
					headerIndex = iterHeaderIndex
					entry = csvFields.ColumnIndexToEntry[iterHeaderIndex]
					break
				}
			} else if strings.EqualFold(strings.TrimSpace(iterHeaderColumnName), strings.ToLower(strings.TrimSpace(tag))) {
				headerColumnName = iterHeaderColumnName
				headerIndex = iterHeaderIndex
				entry = csvFields.ColumnIndexToEntry[iterHeaderIndex]
				break
			}
		}

		if entry != nil {
			entry.TagName = tag
			entry.FieldName = field.Name
			entry.FieldTypeName = field.Type.Name()
		} else {
			entry = &csvFieldEntry{
				TagName:           tag,
				FieldName:         field.Name,
				FieldTypeName:     field.Type.Name(),
				HeaderColumnName:  headerColumnName,
				HeaderColumnIndex: headerIndex,
			}
		}
		csvFields.TagToEntry[tag] = entry
	}

	return csvFields
}

func DetermineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	return
}

func parseBoolField(fieldName string, fieldStringValue string) (bool, error) {
	value, err := strconv.ParseBool(fieldStringValue)
	if err != nil {
		return false, fmt.Errorf("can't parse boolean field '%s' ('%s'): %w", fieldName, fieldStringValue, err)
	}
	return value, nil
}

func parseInt32Field(fieldName string, fieldStringValue string) (int32, error) {
	value, err := strconv.ParseInt(fieldStringValue, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("can't parse field '%s' ('%s'): %w", fieldName, fieldStringValue, err)
	}
	return int32(value), nil
}

func parseInt64Field(fieldName string, fieldStringValue string) (int64, error) {
	value, err := strconv.ParseInt(fieldStringValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse field '%s' ('%s'): %w", fieldName, fieldStringValue, err)
	}
	return int64(value), nil
}

func doConvert(source io.Reader, output string) error {
	allData, err := ioutil.ReadAll(source)
	if err != nil {
		return fmt.Errorf("can't read from source: %w", err)
	}

	encodingReader := bytes.NewReader(allData)

	encoding_, _, _, err := DetermineEncodingFromReader(encodingReader)
	if err != nil {
		panic(err)
	}

	dataReader := bytes.NewReader(allData)
	// map the CSV encoding to UTF-8
	// encoding_ = charmap.Windows1252
	reEncoder := encoding_.NewDecoder().Reader(dataReader)
	r := csv.NewReader(reEncoder)
	r.Comma = ';'
	r.Comment = '#'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	header := records[0]
	records = records[1:]

	csvFields := parseStructTags(header)

	comuniList := []comuni.Comune{}
	for _, record := range records {
		comune := comuni.Comune{}

		for columnIndex, columnValue := range record {
			entry := csvFields.ColumnIndexToEntry[columnIndex]
			if (entry == nil) || (entry.FieldName == "") {
				continue
			}
			field := entry.getField(&comune)
			switch entry.FieldTypeName {
			case "string":
				field.SetString(strings.TrimSpace(columnValue))

			case "bool":
				boolValue, err := parseBoolField(entry.FieldName, columnValue)
				if err != nil {
					log.Fatalf("can't convert bool csv field '%s' for record '%s' column: %v", entry.HeaderColumnName, entry.FieldName, err)
				}
				field.SetBool(boolValue)

			case "int32":
				intValue, err := parseInt32Field(entry.FieldName, columnValue)
				if err != nil {
					log.Fatalf("can't convert int32 csv field '%s' for record '%s' column: %v", entry.HeaderColumnName, entry.FieldName, err)
				}
				field.SetInt(int64(intValue))

			case "int":
				fallthrough
			case "int64":
				intValue, err := parseInt64Field(entry.FieldName, columnValue)
				if err != nil {
					log.Fatalf("can't convert int64 csv field '%s' for record '%s' column: %v", entry.HeaderColumnName, entry.FieldName, err)
				}
				field.SetInt(int64(intValue))

			default:
				log.Fatalf("csv field type '%s' not handled - csv field '%s' for record '%s' column: %v", entry.FieldTypeName, entry.HeaderColumnName, entry.FieldName, err)
			}
		}

		comuniList = append(comuniList, comune)
	}

	data, err := json.Marshal(&comuniList)
	if err != nil {
		return fmt.Errorf("can't marshal to JSON: %w", err)
	}

	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("can't create output file '%s': %w", output, err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(data)
	if err != nil {
		return fmt.Errorf("can't write JSON to '%s': %w", output, err)
	}

	return nil
}

func handleCmdConvert(source string, output string, allowInsecure bool, args []string) {
	if strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "http://") {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: allowInsecure}
		httpClient := &http.Client{Transport: customTransport}

		resp, err := httpClient.Get(source)
		if err != nil {
			log.Fatalf("can't fetch from '%s': %v", source, err)
		}

		defer resp.Body.Close()

		if err := doConvert(resp.Body, output); err != nil {
			log.Fatalf("can't process '%s': %v", source, err)
		}
	} else {
		comuniFile, err := os.OpenFile(source, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatalf("can't open '%s': %v", source, err)
		}
		defer comuniFile.Close()

		if err := doConvert(comuniFile, output); err != nil {
			log.Fatalf("can't process '%s': %v", source, err)
		}
	}
}

func handleCmdDownload(source string, output string, allowInsecure bool, args []string) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: allowInsecure}
	httpClient := &http.Client{Transport: customTransport}

	resp, err := httpClient.Get(source)
	if err != nil {
		log.Fatalf("can't fetch from '%s': %v", source, err)
	}

	defer resp.Body.Close()

	outputFile, err := os.Create(output)
	if err != nil {
		log.Fatalf("can't create output file '%s': %v", output, err)
	}
	defer outputFile.Close()

	_, err = outputFile.ReadFrom(resp.Body)
	if err != nil {
		log.Fatalf("can't write remote source to '%s': %v", output, err)
	}
}

func main() {
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	convertAllowInsecure := convertCmd.Bool("allow-insecure", false, "allow insecure")
	convertSourceName := convertCmd.String("source", CODICI_ISTAT_COMUNI_CSV_URL, "source")
	convertOutputName := convertCmd.String("output", "comuni.json", "output file")

	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	downloadAllowInsecure := downloadCmd.Bool("allow-insecure", false, "allow insecure")
	downloadSourceName := downloadCmd.String("source", CODICI_ISTAT_COMUNI_CSV_URL, "source")
	downloadOutputName := downloadCmd.String("output", "comuni.csv", "output file")

	if len(os.Args) < 2 {
		log.Fatalf("expected a subcommand (eg 'convert' or 'download')")
	}

	switch os.Args[1] {
	case "convert":
		if err := convertCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("can't parse command arguments: %v", err)
		}
		handleCmdConvert(*convertSourceName, *convertOutputName, *convertAllowInsecure, convertCmd.Args())
	case "download":
		if err := downloadCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("can't parse command arguments: %v", err)
		}
		handleCmdDownload(*downloadSourceName, *downloadOutputName, *downloadAllowInsecure, convertCmd.Args())
	default:
		log.Fatalf("expected 'convert' or 'download' subcommand")
	}
}
