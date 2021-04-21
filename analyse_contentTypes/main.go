package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type ItemSys struct {
	Space struct {
		Sys struct {
			Type     string `json:"type"`
			Linktype string `json:"linkType"`
			ID       string `json:"id"`
		} `json:"sys"`
	} `json:"space"`
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Createdat   time.Time `json:"createdAt"`
	Updatedat   time.Time `json:"updatedAt"`
	Environment struct {
		Sys struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Linktype string `json:"linkType"`
		} `json:"sys"`
	} `json:"environment"`
	Publishedversion int       `json:"publishedVersion"`
	Publishedat      time.Time `json:"publishedAt"`
	Firstpublishedat time.Time `json:"firstPublishedAt"`
	Createdby        struct {
		Sys struct {
			Type     string `json:"type"`
			Linktype string `json:"linkType"`
			ID       string `json:"id"`
		} `json:"sys"`
	} `json:"createdBy"`
	Updatedby struct {
		Sys struct {
			Type     string `json:"type"`
			Linktype string `json:"linkType"`
			ID       string `json:"id"`
		} `json:"sys"`
	} `json:"updatedBy"`
	Publishedcounter int `json:"publishedCounter"`
	Version          int `json:"version"`
	Publishedby      struct {
		Sys struct {
			Type     string `json:"type"`
			Linktype string `json:"linkType"`
			ID       string `json:"id"`
		} `json:"sys"`
	} `json:"publishedBy"`
}

type Field struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Localized   bool   `json:"localized"`
	Required    bool   `json:"required"`
	Validations []struct {
		Linkcontenttype []string `json:"linkContentType"`
	} `json:"validations"`
	Disabled bool   `json:"disabled"`
	Omitted  bool   `json:"omitted"`
	Linktype string `json:"linkType,omitempty"`
	Items    struct {
		Type        string `json:"type"`
		Validations []struct {
			Linkcontenttype []string `json:"linkContentType"`
			Inn             []string `json:"in"`
		} `json:"validations"`
		Linktype string `json:"linkType"`
	} `json:"items,omitempty"`
}

type Item struct {
	Sys          ItemSys     `json:"sys"`
	Displayfield interface{} `json:"displayField"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Fields       []Field     `json:"fields"`
}

// Autogenerated using https://mholt.github.io/json-to-go/
type Autogenerated struct {
	Items        []Item `json:"items"`
	ContentTypes []Item `json:"contentTypes"`
}

// missingDescription expect and item and validates existance of a description
func missingDescription(item Item) (bool, error) {
	if len(item.Description) == 0 {
		return true, nil
	}
	return false, nil
}

// missingReferenceValidation expect and item and check validations for references
func missingReferenceValidation(field Field, validations map[string]string) (map[string]string, error) {

	if _, ok := validations["references"]; !ok {
		validations["references"] = ""
	}

	switch field.Type {
	case "Link":
		if len(field.Validations) == 0 && field.Linktype != "Asset" {
			if validations["references"] == "" {
				validations["references"] = fmt.Sprintf("%s", field.Name)
			} else {
				validations["references"] = fmt.Sprintf("%s,%s", validations["references"], field.Name)
			}
		}
		break
	case "Array":
		if len(field.Items.Validations) == 0 {
			if validations["references"] == "" {
				validations["references"] = fmt.Sprintf("%s", field.Name)
			} else {
				validations["references"] = fmt.Sprintf("%s,%s", validations["references"], field.Name)
			}
		}
		break
	}

	return validations, nil
}

func omittedValidation(field Field, validations map[string]string) (map[string]string, error) {

	if _, ok := validations["omitted"]; !ok {
		validations["omitted"] = ""
	}

	if field.Omitted {
		validations["omitted"] = field.Name
	}

	return validations, nil
}

func disabledValidation(field Field, validations map[string]string) (map[string]string, error) {

	if _, ok := validations["disabled"]; !ok {
		validations["disabled"] = ""
	}

	if field.Disabled {
		validations["disabled"] = field.Name
	}

	return validations, nil
}

func existVisited(visited []string, key string) bool {
	for _, v := range visited {
		if v == key {
			return true
		}
	}
	return false
}

func doReferenceTree(validationRefs map[string][]string, ToVisit []string, visited []string, contentID_to_search string, indexToCheck int) (bool, map[string][]string, []string, []string, string, int) {

	subNode, ok := validationRefs[ToVisit[indexToCheck]]

	if !ok {
		return false, validationRefs, ToVisit, visited, ToVisit[indexToCheck], indexToCheck
	} else {

		if existVisited(visited, ToVisit[indexToCheck]) {

			visited = append(visited, ToVisit[indexToCheck])
			return true, validationRefs, ToVisit, visited, ToVisit[indexToCheck], indexToCheck
		}
		visited = append(visited, ToVisit[indexToCheck])

		for ToVisitK, ToVisitV := range subNode {

			ok, a, b, visited, d, e := doReferenceTree(validationRefs, subNode, visited, ToVisitV, ToVisitK)

			if ok {
				// if true exit
				return ok, a, b, visited, d, e
			}
		}
	}

	return false, validationRefs, ToVisit, visited, contentID_to_search, indexToCheck
}

func validatereferncesLoop(obj Autogenerated) map[string]string {
	// part 0, create a map of contentypeId and item
	mapContentTypeIdName := make(map[string]Item, len(obj.Items))
	for _, vItem := range obj.Items {
		mapContentTypeIdName[vItem.Sys.ID] = vItem
	}

	// part 1.- create a map of relations: [ct1]ct2,ct3
	validationRefs := make(map[string][]string, 100)

	for _, vItem := range obj.Items {

		for _, vField := range vItem.Fields {
			switch vField.Type {
			case "Link":
				if len(vField.Validations) != 0 && vField.Linktype != "Asset" {

					if _, ok := validationRefs[vItem.Sys.ID]; !ok {
						validationRefs[vItem.Sys.ID] = make([]string, 0, 0)
					}

					for _, vValidation := range vField.Validations {
						validationRefs[vItem.Sys.ID] = append(validationRefs[vItem.Sys.ID], vValidation.Linkcontenttype...)
					}

				}
				break
			case "Array":
				if len(vField.Items.Validations) != 0 {
					if _, ok := validationRefs[vItem.Sys.ID]; !ok {
						validationRefs[vItem.Sys.ID] = make([]string, 0, 0)
					}
					for _, vValidation := range vField.Items.Validations {
						validationRefs[vItem.Sys.ID] = append(validationRefs[vItem.Sys.ID], vValidation.Linkcontenttype...)
						validationRefs[vItem.Sys.ID] = append(validationRefs[vItem.Sys.ID], vValidation.Inn...)
					}
				}
				break
			}
		}
	}

	referencesLoop := make(map[string]string, 0)
	if len(validationRefs) > 0 {

		for k, v := range validationRefs {

			visited := make([]string, 0, 0)
			visited = append(visited, k)

			// split child nodes and loop
			for index, ct_subnode := range v {
				// 	//
				if ct_subnode != "" {

					ok, _, _, visited, _, _ := doReferenceTree(validationRefs, v, visited, ct_subnode, index)
					if ok {
						// exist key
						// enrich contentType id with ContentypeName
						visitedEnriched := make([]string, 0, len(visited))
						for _, vVisited := range visited {
							vVisitedName := mapContentTypeIdName[vVisited]
							visitedEnriched = append(visitedEnriched, fmt.Sprintf("%s(%s)", vVisited, vVisitedName.Name))
						}

						referencesLoop[k] = fmt.Sprintf(" %s ", strings.Join(visitedEnriched, " -> "))
					}
					visited = visited[0:1]
				}
			}
			visited = make([]string, 0, 0)
		}
	}

	return referencesLoop
}

type errorReport struct {
	ContentTypeName string
	ContentTypeID   string
	Errors          []string
}

func main() {

	var ctFilename string
	flag.StringVar(&ctFilename, "f", "", "Specify contentType JSON file.")

	flag.Usage = func() {
		fmt.Printf("Usage of our Program: \n")
		fmt.Printf("./go-project -f /var/doc/MyContentTypeFile.json\n")
		fmt.Println("")
		fmt.Printf("How to generate it: \n")
		fmt.Println("***")
		fmt.Println("With CLI: ")
		fmt.Println("$ contentful space export --skip-content --skip-roles --skip-webhooks")
		fmt.Println("***")
		fmt.Println("***")
		fmt.Println("With CMA: https://www.contentful.com/developers/docs/references/content-management-api/#/reference/content-types/content-type-collection/get-all-content-types-of-a-space/console/curl")
		fmt.Println("***")
	}
	flag.Parse()

	if len(ctFilename) == 0 {
		fmt.Println("parameter 'f' with filename required")
		return
	} else if _, err := os.Stat(ctFilename); err != nil {
		fmt.Printf("File does not exist\n")
		return
	}

	// read file
	data, err := ioutil.ReadFile(ctFilename)
	if err != nil {
		fmt.Print(err)
	}

	// json data
	var obj Autogenerated

	// unmarshall it
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	// if in JSON, Items is empty and ContentType not, copy ContentTypes to Items
	// the customer is probably sending a JSON file done with the CLI rather than CMA
	// files are similar except that collection name is "ContentTypes" instead than "Items"
	if len(obj.Items) == 0 && len(obj.ContentTypes) > 0 {
		obj.Items = obj.ContentTypes
	}

	loopValidationErrors := validatereferncesLoop(obj)

	fmt.Printf("** Analysis Report **\n")
	fmt.Printf("Description:\n")
	fmt.Printf("[Notice] - Good practice attention.\n")
	fmt.Printf("[Warning] - Possible issue.\n")
	fmt.Printf("[Issue] - Issue, something to consider changing.\n")
	fmt.Printf("\n")
	fmt.Printf("Total ContentTypes: %d\n", len(obj.Items))
	fmt.Printf("***** REPORT ******* \n")
	fmt.Printf("\n")

	errorReports := make(map[string]errorReport, len(obj.Items))
	for _, vItem := range obj.Items {

		err := errorReport{
			ContentTypeName: vItem.Name,
			ContentTypeID:   vItem.Sys.ID,
		}

		validations := make(map[string]string, 100)

		for _, vField := range vItem.Fields {

			// add validations functions  for fields here
			validations, _ = missingReferenceValidation(vField, validations)
			validations, _ = omittedValidation(vField, validations)
			validations, _ = disabledValidation(vField, validations)
		}

		errMsg := make([]string, 0)
		if existDesc, _ := missingDescription(vItem); existDesc {
			errMsg = append(errMsg, "* [Notice] Description is missing.")
		}
		if len(vItem.Fields) == 1 {
			errMsg = append(errMsg, "* [Notice] Content type with a single fields: "+vItem.Fields[0].Name)
		}

		if ref, ok := validations["references"]; ok && len(validations["references"]) > 0 {
			errMsg = append(errMsg, "* [Warning] A reference field(s) lack validation, fields: "+ref)
		}
		if ref, ok := validations["omitted"]; ok && len(validations["omitted"]) > 0 {
			errMsg = append(errMsg, "* [Notice] Ommited from API response(ommited): "+ref)
		}

		if ref, ok := validations["disabled"]; ok && len(validations["disabled"]) > 0 {
			errMsg = append(errMsg, "* [Notice] Disabled in WebApp field(disabled): "+ref)
		}

		if msg, ok := loopValidationErrors[vItem.Sys.ID]; ok {
			errMsg = append(errMsg, "* [Issue] Infinite loop on content Type refernces: "+msg)
		}

		if len(errMsg) == 0 {
			err.Errors = []string{"No  errors."}
		} else {
			err.Errors = errMsg
		}
		errorReports[vItem.Name] = err
	}

	keys := make([]string, 0, len(errorReports))
	for k := range errorReports {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		err := errorReports[k]

		fmt.Println("ContentType name: " + err.ContentTypeName)
		fmt.Println("ContentType id: " + err.ContentTypeID)
		for _, errMsg := range err.Errors {
			fmt.Println(errMsg)
		}
		fmt.Println("")
	}
}
