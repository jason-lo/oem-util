package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"github.com/jason-lo/oem-util/util"
)

var (
	Stella_PB = make(map[string]string)
)

func init() {
	Stella_PB["platform"] = "Platform"
	Stella_PB["codename"] = "Code Name"
	Stella_PB["status"] = "Status"
	Stella_PB["tag"] = "LP tag"
	Stella_PB["config"] = "Config#"
	Stella_PB["m1"] = "M1"
	Stella_PB["gm"] = "Target \nGM Date"
}

var config_local = config.Config{}

func init() {
	config_local.Read()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Cmd(cmd string, shell bool) []byte {

	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out
	} else {
		out, err := exec.Command(cmd).Output()
		if err != nil {
			panic("some error found")
		}
		return out

	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
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

//func get_index_map(sheet_id string, sheet_range string) map[string]int {
func Getindex(sheet_id string, sheet_range string) map[string]int {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp0, err := srv.Spreadsheets.Values.Get(sheet_id, sheet_range).Do()

	var aString []string

	if len(resp0.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp0.Values {
			for index := range row {
				print(row[index].(string))
				print(" \n")
				aString = append(aString, row[index].(string))
			}
		}
	}

	indexof := make(map[string]int)

	for k, v := range Stella_PB {
		indexof[k] = findItemIndex(aString, v)
	}

	return indexof
}

func findItemIndex(aString []string, item string) int {

	for index, value := range aString {
		if value == item {
			return index
		}
	}
	return -1

}
