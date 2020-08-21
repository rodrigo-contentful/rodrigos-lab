package main

/*
todo: upsertcontentypes need to including validations and field options
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	// unsuported golang SDK :/
	//https://github.com/contentful-labs/contentful-go
	contentful "github.com/contentful-labs/contentful-go"
)

// Config defines a script configuration
type Config struct {
	OrgID       string   `json:"organizationID"`
	CmaToken    string   `json:"cmaToken"`
	SpaceOrigID string   `json:"spaceOriginID"`
	SpaceDestID []string `json:"spaceDestID"`
}

//SysInner ...
type SysInner struct {
	ID       string `json:"id"`
	LinkType string `json:"linkType"`
	Type     string `json:"type"`
}

//ContentType ...
type ContentType struct {
	Sys SysInner `json:"sys"`
}

// Environment ...
type Environment struct {
	Sys SysInner `json:"sys"`
}

// Space ..
type Space struct {
	Sys SysInner `json:"sys"`
}

// Tags ...
type Tags struct {
	Tags interface{} `json:"tags"`
}

// Sys define sthe sys part of the payload
type Sys struct {
	ID        string      `json:"id"`
	Revision  int         `json:"revision"`
	Type      string      `json:"type"`
	Space     Space       `json:"space"`
	CT        ContentType `json:"contentType"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	Env       Environment `json:"environment"`
}

// Fields defines the example ContentType fields
type Fields struct {
	Test  interface{}       `json:"test"`
	Title map[string]string `json:"title"`
	Turn  map[string]bool   `json:"turn"`
}

// Payload defines the JSON payload from getContentType
type Payload struct {
	Fields   Fields `json:"fields"`
	Metadata Tags   `json:"metadata"`
	Sys      Sys    `json:"sys"`
}

func main() {

	// json data
	var config Config
	// load config
	fmt.Println("Loading config")
	if !loadConfig(&config) {
		return
	}
	fmt.Println("Success loading config, OrgID: ", config.OrgID)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var t Payload
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		// get locales keys
		var keys string
		for k := range t.Fields.Title {
			keys = k
		}

		migrate(t, keys, &config)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}

// loadConfig returns the json configuration as structure
func loadConfig(config *Config) bool {
	// read file
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Print(err)
	}

	// unmarshall it
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	return true
}

// migrate will upsert a content type from space A to space B, later will upsert an entry from that content type from space A to space B
func migrate(t Payload, keys string, config *Config) {

	orgID := config.OrgID
	cmaToken := config.CmaToken   // observe your CMA token from Contentful's web page
	spaceID := config.SpaceOrigID // space origin
	cma := contentful.NewCMA(cmaToken).SetOrganization(orgID)
	cma.Debug = true

	if strings.ToLower(t.Sys.Type) == "entry" {
		fmt.Println("Webhook is for Entry")
	}

	// copy CT
	// GetContentType from space A
	ct, e := cma.ContentTypes.Get(spaceID, t.Sys.CT.Sys.ID)
	if e != nil {
		fmt.Printf("Errro gettong CT: %+v\n", e)
	}

	ctNew := &contentful.ContentType{
		Name:   ct.Name,
		Fields: ct.Fields,
	}

	// upserting ContentTypes pero space
	for _, desSpaceID := range config.SpaceDestID {
		e := upsertContentType(desSpaceID, cma, ctNew)
		if e != nil {
			fmt.Printf("Error: %s", e.Error())
		}
	}

	//copy entry
	// Get Entry from space A
	entry, e := cma.Entries.Get(spaceID, t.Sys.ID)
	if e != nil {
		fmt.Printf("Error getting Entry: %+v\n", e)
	}

	//Upserting Entries per space
	for _, desSpaceID := range config.SpaceDestID {
		e := upsertEntry(desSpaceID, cma, ctNew, entry)
		if e != nil {
			fmt.Printf("%s", e.Error())
		}
	}
}

// upsertContentType expects a space destination, cma client and conten type and will upsert a ContentType in destantion space
func upsertContentType(spaceCopyID string, cma *contentful.Client, ctNew *contentful.ContentType) error {

	/*
	   {
	      "name": "Blog Post",
	      "fields": [
	        {
	          "id": "title",
	          "name": "Title",
	          "required": true,
	          "localized": true,
	          "type": "Text"
	        },
	        {
	          "id": "body",
	          "name": "Body",
	          "required": true,
	          "localized": true,
	          "type": "Text"
	        }
	      ]
	    }
	*/

	ctsCollection := cma.ContentTypes.List(spaceCopyID)
	ctsCollection, err := ctsCollection.Next() // makes the actual api call
	if err != nil {
		log.Fatal(err)
		return err
	}

	cts := ctsCollection.ToContentType() // make the type assertion
	for _, c := range cts {
		if c.Name == ctNew.Name {
			// ct exist, update version
			ctNew.Sys = c.Sys
			break
		}
	}

	// Upsert contentype to space B
	err = cma.ContentTypes.Upsert(spaceCopyID, ctNew)
	if err != nil {
		fmt.Printf("Error copying CT: %+v\n", err)
		return err
	}

	// ContentTypes are copied as draft, here we can publish the contentType on space B programatically
	err = cma.ContentTypes.Activate(spaceCopyID, ctNew)
	if err != nil {
		fmt.Printf("Error activating CT: %+v\n", err)
		return err
	}
	return nil
}

// upsertEntry expects a space destination, cma client and content type, and entry and will upsert an Entry in destantion space
func upsertEntry(spaceCopyID string, cma *contentful.Client, ctNew *contentful.ContentType, entry *contentful.Entry) error {

	/*
	   {
	      "fields": {
	        "title": {
	          "en-US": "Hello, World!"
	        },
	        "body": {
	          "en-US": "Bacon is healthy!"
	        }
	      }
	    }
	*/
	entt := &contentful.Entry{
		Sys: &contentful.Sys{
			ID: entry.Sys.ID,
			ContentType: &contentful.ContentType{
				Sys: &contentful.Sys{
					ID: ctNew.Sys.ID,
				},
			},
		},
		Fields: entry.Fields,
	}

	// Upsert entry on space B
	err := cma.Entries.Upsert(spaceCopyID, entt)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Entry is copied as draft, here we can publish it programatically
	err = cma.Entries.Publish(spaceCopyID, entt)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
