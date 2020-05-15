package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	flags "github.com/jessevdk/go-flags"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	sheetAPIScope  = "https://www.googleapis.com/auth/spreadsheets"
	timeFormat     = "2006/01/02"
	spreadsheetID  = "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	sheetName      = "Class Data"
	readRange      = "!A2:E"
	tokFile        = "token.json"
	credentialFile = "credentials.json"
)

type options struct {
	Column string `short:"c" long:"column" description:"a column to update"`
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func increment(n int) int {
	return n + 1
}

func (o *options) toColumnNum() int {
	switch o.Column {
	case "B":
		return 1
	case "C":
		return 2
	default:
		log.Fatalf("invalid column: %v", o.Column)
	}

	return 100
}

func setup() *sheets.Service {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheetAPIScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return srv
}

func cellIncrement(
	opts options,
	valueRange *sheets.ValueRange,
) (string, *sheets.ValueRange) {
	now := time.Now().Local()
	var rowNum string
	var newVal int
	if len(valueRange.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, row := range valueRange.Values {
			if len(row) < 3 {
				log.Printf("length is less than 3. row: %v, length: %v", i+1, len(row))
				continue
			}
			str, ok := row[0].(string)
			if ok {
				if now.Format(timeFormat) == str {
					// 更新対象となるセルの既存値を取得する
					val, err := strconv.Atoi(row[opts.toColumnNum()].(string))
					if err != nil {
						log.Fatalf("Unable to retrieve data from sheet: %v", err)
					}
					// 既存の値をincrementした値を格納する
					newVal = increment(val)
					// セルの行を取得
					rowNum = strconv.Itoa(i + 1)
				}
			}
		}
	}

	writeRange := sheetName + "!" + opts.Column + rowNum + ":" + opts.Column + rowNum
	return writeRange, &sheets.ValueRange{
		Values: [][]interface{}{
			{newVal},
		},
	}
}

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalf("Unable to parse flag: %v", err)
	}

	srv := setup()

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, sheetName+readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	writeRange, newValueRange := cellIncrement(opts, resp)
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, newValueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
}
