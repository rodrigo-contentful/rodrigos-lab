package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	// unsuported golang SDK :/
	//https://github.com/contentful-labs/contentful-go
	contentful "github.com/contentful-labs/contentful-go"
)

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

		publish(t, keys)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}

// publish will upsert a content type from space A to space B, later will upsert an entry from that content type from space A to space B
func publish(t Payload, keys string) {

	orgID := "SOME_ORG_ID"
	cmaToken := "SOME_CMA_TOKEN" // observe your CMA token from Contentful's web page
	spaceID := "SPACE_ID_A"      // space origin
	spaceCopyID := "SPACE_ID_B"  // space destination
	cma := contentful.NewCMA(cmaToken)
	cma.SetOrganization(orgID)
	cma.Debug = true

	// e, spaces := cma.Spaces.Get(spaceID)
	// if e != nil {
	// 	fmt.Printf("%+v", e)
	// }
	// fmt.Printf("%+v", spaces)

	if strings.ToLower(t.Sys.Type) == "entry" {
		fmt.Println("Webhook is for Entry")
	}

	// copy CT
	// GetContentType from space A
	ct, e := cma.ContentTypes.Get(spaceID, t.Sys.CT.Sys.ID)
	if e != nil {
		fmt.Printf("Errro gettong CT: %+v\n", e)
	}

	// Upsert contentype to space B
	e = cma.ContentTypes.Upsert(spaceCopyID, ct)
	if e != nil {
		fmt.Printf("Error copying CT: %+v\n", e)
	}

	// ContentTypes are copied as draft, here we can publish the contentType on space B programatically
	e = cma.ContentTypes.Activate(spaceCopyID, ct)
	if e != nil {
		fmt.Printf("Error activating CT: %+v\n", e)
	}

	//copy entry
	// Get Entry from space A
	entry, e := cma.Entries.Get(spaceID, t.Sys.ID)
	if e != nil {
		fmt.Printf("Error getting Entry: %+v\n", e)
	}

	// Upsert entry on space B
	err := cma.Entries.Upsert(spaceCopyID, entry)
	if err != nil {
		log.Fatal(err)
	}

	// Entry is copied as draft, here we can publish it programatically
	err = cma.Entries.Publish(spaceCopyID, entry)
	if err != nil {
		log.Fatal(err)
	}

}
