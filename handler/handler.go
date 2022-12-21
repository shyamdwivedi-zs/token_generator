package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/gmail/v1"
)

type tokenHandler struct {
	authCode chan string
}

func New(authCode chan string) *tokenHandler {
	return &tokenHandler{authCode: authCode}
}

func (th tokenHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("state is: ", r.FormValue("state"))

	th.authCode <- r.FormValue("code")

	fmt.Fprintf(w, "Success...!!")
}

func (th tokenHandler) GenerateToken(authCode chan string) {
	apiScopes := map[string][]string{
		"calendar": {calendar.CalendarScope, calendar.CalendarReadonlyScope, calendar.CalendarEventsScope, calendar.CalendarEventsReadonlyScope},
		"drive":    {drive.DriveScope, drive.DriveFileScope, drive.DriveAppdataScope, drive.DriveMetadataScope, drive.DriveScriptsScope, drive.DriveReadonlyScope, drive.DriveMetadataReadonlyScope, drive.DrivePhotosReadonlyScope},
		"mail":     {gmail.MailGoogleComScope, gmail.GmailModifyScope, gmail.GmailReadonlyScope, gmail.GmailComposeScope, gmail.GmailSendScope},
	}

	tokens := make(map[string]*oauth2.Token)

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	for api, scope := range apiScopes {
		config, err := google.ConfigFromJSON(b, scope...)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}

		tokens[api] = getTokenFromWeb(config, authCode)
	}

	err = saveToken(tokens)
	if err != nil {
		log.Fatalf("Unable to create token file: %v", err)
	}

	fmt.Println("\n\ntoken save to token.json")
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config, authCode chan string) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser and authorize: \n%v\n", authURL)

	code := <-authCode

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return token
}

func saveToken(tokens map[string]*oauth2.Token) error {
	f, err := os.Create("token.json")
	if err != nil {
		return err
	}

	err = json.NewEncoder(f).Encode(tokens)
	if err != nil {
		return err
	}

	return nil
}
