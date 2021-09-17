package main

/*
Tool to analyse content types JSON programatically
*/
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"strconv"
	"time"
	// "bytes"
	// "encoding/gob"
	"math"
)

// lower limit for similarity of type names for similitud comparison
const simLowerLimitPrecentageName = .75
const simLowerLimitPrecentageFields = .85

// ItemSys describes sys element of a content type
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

// Field describes a field of a cotent type
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

// Item describes a conten type
type Item struct {
	Sys          ItemSys     `json:"sys"`
	Displayfield interface{} `json:"displayField"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Fields       []Field     `json:"fields"`
}

// ModelItems using https://mholt.github.io/json-to-go/
type ModelItems struct {
	Items        []Item `json:"items"`
	ContentTypes []Item `json:"contentTypes"`
}

type spaceParsed struct {
	Name   string
	ItemsCount []contentTypeParsed
}

type contentTypeParsed struct {
	Name   string
	FieldsCount []fieldItemsHashCheck
	Hash []byte
}

type fieldItemsHashCheck struct {
	Type        string `json:"type"`
	Localized   bool   `json:"localized"`
	Required    bool   `json:"required"`
	Disabled bool   `json:"disabled"`
	Omitted  bool   `json:"omitted"`
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

// textFieldValidationByName Based on a field name match a contentful field validation
func textFieldValidationByName(field Field, validations map[string]string) (map[string]string, error) {

	if _, ok := validations["fieldValidationByName"]; !ok {
		validations["fieldValidationByName"] = ""
	}

	var fieldValidNames = []string{"email","e-mail", "phone", "telephone", "mobile", "url", "link", "date", "time"}

	for _, fieldValidName := range fieldValidNames {
		nameTmp := strings.Trim(field.Name, " ")
		nameTmp = strings.Trim(nameTmp, "_")
		nameTmp = strings.Trim(nameTmp, "-")
		if strings.Contains(strings.ToLower(nameTmp), fieldValidName) {
			if len(field.Validations) == 0 && field.Type == "Symbol" {
				validations["fieldValidationByName"] = fmt.Sprintf("Possible validation missing: field name '%s' matches a Contentful text validation '%s'", field.Name, fieldValidName)
				return validations, nil
			}
		}
	}
	return validations, nil
}

// omittedValidation validated if field has the omitted flag active
func omittedValidation(field Field, validations map[string]string) (map[string]string, error) {

	if _, ok := validations["omitted"]; !ok {
		validations["omitted"] = ""
	}

	if field.Omitted {
		validations["omitted"] = field.Name
	}

	return validations, nil
}

// disabledValidation validated if field has the disable flag active
func disabledValidation(field Field, validations map[string]string) (map[string]string, error) {

	if _, ok := validations["disabled"]; !ok {
		validations["disabled"] = ""
	}

	if field.Disabled {
		validations["disabled"] = field.Name
	}

	return validations, nil
}

var htmlElements = []string{
	"form",
	"input",
	"datalist",
	"fieldset",
	"keygen",
	"legend",
	"optgroup",
	"option",
	"output",
	"select",
	"textarea",
	"input",
	"output",
	"select",
	"button",
	"option",
	"textarea",
	"optgroup",
	"fieldset",
}

// fieldNameAsHTMLElement compare the name of a field with a html element, possible micropy
func fieldNameAsHTMLElement(field Field, validations map[string]string) (map[string]string, error) {

	for _, elementName := range htmlElements {
		if strings.Contains(strings.ToLower(field.Name), elementName) {
			validations["htmlName"] = fmt.Sprintf("Possible microcopy: field name '%s' includes html element name '%s'", field.Name, elementName)
			return validations, nil
		}
	}

	return validations, nil
}

// fieldCTA Based on field name find the text "button" to find a possible microcpy
func fieldCTA(field Field, validations map[string]string) (map[string]string, error) {

	if strings.Contains(strings.ToLower(field.Name), "button") {
		validations["htmlName"] = fmt.Sprintf("Possible CTA(call to action): field name '%s' includes the text 'button'", field.Name)
		return validations, nil
	}

	return validations, nil
}

// fieldNotResponsive Based on field name find the text "desktop", "mobile", "tablet" to find a possible responsive field.
func fieldNotResponsive(field Field, validations map[string]string) (map[string]string, error) {

	responsiveLabels := []string{"desktop", "mobile", "tablet"}
	for _, elementName := range responsiveLabels {
		if strings.Contains(strings.ToLower(field.Name), elementName) {
			validations["htmlName"] = fmt.Sprintf("%s Possible NON responsive field: field name '%s' includes the text '%s' \n", validations["htmlName"], field.Name, elementName)
			return validations, nil
		}
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

// doReferenceTree navigates each reference field and create their path
func doReferenceTree(validationRefs map[string][]string, ToVisit []string, visited []string, indexToCheck int) (bool, []string) {
	subNode, ok := validationRefs[ToVisit[indexToCheck]]
	if !ok {
		return false, visited
	} else {

		if existVisited(visited, ToVisit[indexToCheck]) {

			visited = append(visited, ToVisit[indexToCheck])

			return true, visited
		}
		visited = append(visited, ToVisit[indexToCheck])

		for ToVisitK := range subNode {

			ok, visited := doReferenceTree(validationRefs, subNode, visited, ToVisitK)

			if ok {
				// if true exit
				return ok, visited
			}
		}
	}

	return false, visited
}

// noticeLog print log Notice
func noticeLog(msg string) string {
	return fmt.Sprintf("[Notice] 💡 - %s", msg)
}

// attentionLog print log Attention
func attentionLog(msg string) string {
	return fmt.Sprintf("[Attention] 🔍 - %s", msg)
}

// warningLog print log Warning
func warningLog(msg string) string {
	return fmt.Sprintf("[Warning] 🏮 - %s", msg)
}

// issueLog print log Issue
func issueLog(msg string) string {
	return fmt.Sprintf("[Issue] ⛔ - %s", msg)
}

// validatereferncesLoop inspect the references to find loops
func validatereferncesLoop(obj ModelItems) (map[string]string, map[string]string) {
	nonOrphanContentTypes := make(map[string]string, len(obj.Items))
	// part 0, create a map of contentypeId and item
	mapContentTypeIDName := make(map[string]Item, len(obj.Items))
	for _, vItem := range obj.Items {
		mapContentTypeIDName[vItem.Sys.ID] = vItem
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

		// create a map of contentTypes ids that are related TO or FROM
		for contentTypeFrom, contentTypeRelations := range validationRefs {
			nonOrphanContentTypes[contentTypeFrom] = contentTypeFrom
			for _, contentTypeChild := range contentTypeRelations {
				nonOrphanContentTypes[contentTypeChild] = contentTypeChild
			}
		}

		for k, v := range validationRefs {

			visited := make([]string, 0, 0)
			visited = append(visited, k)

			// split child nodes and loop
			for index, ctSubnode := range v {
				//
				if ctSubnode != "" {

					ok, visited := doReferenceTree(validationRefs, v, visited, index)
					if ok {
						// exist key
						// enrich contentType id with ContentypeName
						visitedEnriched := make([]string, 0, len(visited))
						for _, vVisited := range visited {
							vVisitedName := mapContentTypeIDName[vVisited]
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

	return referencesLoop, nonOrphanContentTypes
}

type errorReport struct {
	ContentTypeName string
	ContentTypeID   string
	Errors          []string
}

type fieldValidation struct {
	ContentTypeName string
	FieldName       string
	FieldID         string
	FieldType       string
	HideDefault     bool
}

func iterateDirectory(path, ctDirectory string) []string {
	res := make([]string, 0, 10)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(fmt.Sprintf("%s/%s", path, info.Name())) == ".json" {
			res = append(res, fmt.Sprintf("%s/%s", ctDirectory, info.Name()))
		}
		return nil
	})
	return res
}

func fieldDuplicated(f Field, duplications map[string][]fieldValidation) []string {
	res := make([]string, 0, 10)
	if findDuplication, ok := duplications[f.ID]; ok {
		if len(findDuplication) > 1 {
			for _, vv := range findDuplication {
				// filter fields set as default from types
				// need to filter some fields liek slug, title etc, find a smart way of doing it
				if !vv.HideDefault {
					res = append(res, fmt.Sprintf("** ContentTypeID: '%s'", vv.ContentTypeName))
				}
			}
		}
	}
	return res
}

// func hash(s []fieldItemsHashCheck) []byte {
//     var b bytes.Buffer
//     gob.NewEncoder(&b).Encode(s)
//     return b.Bytes()
// }


func processJSONFile(ctFilename string, contentTypesParsed []spaceParsed, nameThreshhold, fieldsThreshold float64) []spaceParsed {
// func processJSONFile(ctFilename string, contentTypesParsed map[string][]contentTypeParsed) {
	// read file
	data, err := ioutil.ReadFile(ctFilename)
	if err != nil {
		fmt.Print(err)
	}

	// json data
	var obj ModelItems

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
	loopValidationErrors, nonOrphanContentTypes := validatereferncesLoop(obj)

	fmt.Println()
	fmt.Println()
	fmt.Printf("***************************************************************\n")
	fmt.Println()
	fmt.Printf("** Analysis Report Space '%s' **\n", obj.Items[0].Sys.Space.Sys.ID)
	fmt.Printf("Description:\n")
	fmt.Printf(noticeLog("Good practice/recomendation.\n"))
	fmt.Printf(attentionLog("Something to have a look at.\n"))
	fmt.Printf(warningLog("Possible issue.\n"))
	fmt.Printf(issueLog("Something to consider changing.\n"))
	fmt.Printf("\n")
	fmt.Printf("Total ContentTypes: %d\n", len(obj.Items))
	fmt.Printf("***** REPORT ******* \n")
	fmt.Printf("\n")

	// create aray of contenttypes of space to multispace checkup
	ctItems := make([]contentTypeParsed, 0, 10)
	errorReports := make(map[string]errorReport, len(obj.Items))

	fieldDuplication := make(map[string][]fieldValidation, len(obj.Items))
	for _, vItem := range obj.Items {

		ctp := contentTypeParsed{
			Name: vItem.Name,
		}
		// ctp.FieldsCount = make(map[string]int, len(vItem.Fields))
		// ctItems = append(ctItems, ctp)

		err := errorReport{
			ContentTypeName: vItem.Name,
			ContentTypeID:   vItem.Sys.ID,
		}

		// find the display field, if not existante is null or empty
		// that is why is needed to cast to string
		vItemDisplayfield := ""
		switch vItem.Displayfield.(type) { // the switch uses the type of the interface
		case string:
			vItemDisplayfield = vItem.Displayfield.(string)
		}

		validations := make(map[string]string, 100)

		/*
			Steps

			1.- Identify microcpy - HTML elements labels
			2.- Identify CTA
			3.- Duplicated medias
			4.- Standard Naming
		*/
		for _, vField := range vItem.Fields {

			ctp.FieldsCount = append(ctp.FieldsCount, fieldItemsHashCheck{
				Type:vField.Type,
				Localized:vField.Localized,
				Required:vField.Required,
				Disabled:vField.Disabled,
				Omitted:vField.Omitted,
			})
				// add validations functions  for fields here
			validations, _ = missingReferenceValidation(vField, validations)
			validations, _ = omittedValidation(vField, validations)
			validations, _ = disabledValidation(vField, validations)
			validations, _ = fieldNameAsHTMLElement(vField, validations)
			validations, _ = fieldCTA(vField, validations)
			validations, _ = fieldNotResponsive(vField, validations)
			validations, _ = textFieldValidationByName(vField, validations)

			// Add contentYtpes and fields to map of duplicated values
			nFielValidation := fieldValidation{
				ContentTypeName: vItem.Name,
				FieldName:       vField.Name + " " + vItemDisplayfield,
				FieldID:         vField.ID,
				FieldType:       vField.Type,
				HideDefault:     false,
			}

			if vItemDisplayfield == nFielValidation.FieldID {
				nFielValidation.HideDefault = true
			}

			if fdItem, ok := fieldDuplication[vField.ID]; ok {
				fdItem = append(fdItem, nFielValidation)
				fieldDuplication[vField.ID] = fdItem
			} else {
				fieldDuplication[vField.ID] = append(fieldDuplication[vField.ID], nFielValidation)
			}
		}

		// hash fields to check later
		// ctp.Hash=Hash(ctp.FieldsCount)

		// add the the content type name and fields to validation array of space
		ctItems = append(ctItems, ctp)

		errMsg := make([]string, 0)
		if existDesc, _ := missingDescription(vItem); existDesc {
			//errMsg = append(errMsg, "* 💡 [Notice] Description is missing.")
			errMsg = append(errMsg, noticeLog("Description is missing."))
		}
		if len(vItem.Fields) == 1 {
			//errMsg = append(errMsg, "* 💡 [Notice] Content type with a single fields: "+vItem.Fields[0].Name)
			errMsg = append(errMsg, noticeLog("Content type with a single fields: "+vItem.Fields[0].Name))
		}

		if ref, ok := validations["references"]; ok && len(validations["references"]) > 0 {
			//errMsg = append(errMsg, "* 🔍 [Warning] A reference field(s) lack validation, fields: "+ref)
			errMsg = append(errMsg, warningLog("A reference field(s) lack validation, fields: "+ref))
		}
		if ref, ok := validations["omitted"]; ok && len(validations["omitted"]) > 0 {
			//errMsg = append(errMsg, "* 💡 [Notice] Ommited from API response(ommited): "+ref)
			errMsg = append(errMsg, noticeLog("Ommited from API response(ommited): "+ref))
		}

		if ref, ok := validations["disabled"]; ok && len(validations["disabled"]) > 0 {
			//errMsg = append(errMsg, "* 💡 [Notice] Disabled in WebApp field(disabled): "+ref)
			errMsg = append(errMsg, noticeLog("Disabled in WebApp field(disabled): "+ref))
		}

		if ref, ok := validations["htmlName"]; ok && len(validations["htmlName"]) > 0 {
			//errMsg = append(errMsg, "* 💡 [Notice] "+ref)
			errMsg = append(errMsg, attentionLog(ref))
		}

		if ref, ok := validations["fieldValidationByName"]; ok && len(validations["fieldValidationByName"]) > 0 {
			//errMsg = append(errMsg, "* 💡 [Notice] "+ref)
			errMsg = append(errMsg, attentionLog(ref))
		}

		if msg, ok := loopValidationErrors[vItem.Sys.ID]; ok {
			//errMsg = append(errMsg, "* 🔍 [Warning] Infinite loop on content Type refernces: "+msg)
			errMsg = append(errMsg, warningLog("Infinite loop on content Type refernces: "+msg))
		}

		// is orphan
		if _, ok := nonOrphanContentTypes[vItem.Sys.ID]; !ok {
			errMsg = append(errMsg, warningLog("Content type not referenced(Orphan)."))
		}

		// no display field selected
		if len(vItemDisplayfield) == 0 {
			errMsg = append(errMsg, issueLog("Content type ha sno title field."))
		}

		if len(errMsg) == 0 {
			err.Errors = []string{"* 🥇 No  errors."}
		} else {
			err.Errors = errMsg
		}
		errorReports[vItem.Name] = err
	}

	// add contentTypes items to map of spaces
	// to multispace validation
	contentTypesParsed = append(contentTypesParsed,spaceParsed{
		Name:obj.Items[0].Sys.Space.Sys.ID,
		ItemsCount:ctItems,
	})
	// contentTypesParsed[obj.Items[0].Sys.Space.Sys.ID] = ctItems

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

	fmt.Println("")
	fmt.Println("")
	fmt.Println("Similar Content Types found:")
	fmt.Println("")
	for itemIndex,itemValue := range ctItems{

		for _,itemStep := range ctItems[(itemIndex+1):]{	
			SimilarContentTypesFormat(itemValue,itemStep,nameThreshhold, fieldsThreshold,"","")
				
		}	
	}

	return contentTypesParsed
}

func multiSpaceValidations(contentTypesParsed []spaceParsed, nameThreshhold, fieldsThreshold float64 ) {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("*** Validating multispace reusabaility ** ")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Content type naming repetetition")

	for firstSpaceIndex,firstSpaceValue := range contentTypesParsed{

		for _,SecondSpaceValue := range contentTypesParsed[(firstSpaceIndex+1):]{	

			for _, firstSpaceValueType := range firstSpaceValue.ItemsCount {
				for _, SecondSpaceValueType := range SecondSpaceValue.ItemsCount {
					SimilarContentTypesFormat(firstSpaceValueType,SecondSpaceValueType,nameThreshhold, fieldsThreshold,firstSpaceValue.Name,SecondSpaceValue.Name)
				}
			}

		}	
	}
}

func splitChars(s string) []string {
	chars := make([]string, 0, len(s))
	// Assume ASCII inputs
	for i := 0; i != len(s); i++ {
		chars = append(chars, string(s[i]))
	}
	return chars
}


func SimilarContentTypesFormat(itemA,itemB contentTypeParsed, nameThreshhold, fieldsThreshold float64 ,spaceIdA, spaceIdB string){

	// do text for itemA
	testA := make([]string,0,len(itemA.FieldsCount)+1)
	for _,vItemA := range itemA.FieldsCount	{
		rr22,_:=json.Marshal(&vItemA)
		testA=append(testA,string(rr22))
	}
	
	// do text for itemA
	testB := make([]string,0,len(itemB.FieldsCount)+1)
	for _,vItemB := range itemB.FieldsCount	{
		rr22,_:=json.Marshal(&vItemB)
		testB=append(testB,string(rr22))
	}

	testAA := strings.Join(testA,"\n")
	testBB := strings.Join(testB,"\n")
	
	// checking first if name has similarity
	seqMatcherA := NewMatcher(splitChars(itemA.Name),splitChars(itemB.Name))
	// if (seqMatcherA.QuickRatio() < nameThreshhold){
	if (seqMatcherA.Ratio() < nameThreshhold){
		return
	}

	// now check similarity of fields
	seqMatcher := NewMatcher(SplitLines(testAA),SplitLines(testBB))
	// if (seqMatcher.QuickRatio() >= fieldsThreshold){
	if (seqMatcher.Ratio() >= fieldsThreshold){

		if len(spaceIdA)!=0 || len(spaceIdB)!=0{
			fmt.Printf("Space Id's: '%s' , '%s'\n",spaceIdA, spaceIdB)
		}

		fmt.Printf("Content Types: '%s' , '%s'\n",itemA.Name, itemB.Name)
		fmt.Printf("Similarity of name: %.0f %s\n",math.Ceil(seqMatcherA.Ratio()*100),"%")
		fmt.Printf("Similarity of fields: %.0f %s\n",math.Ceil(seqMatcher.Ratio()*100),"%")
		fmt.Printf("Fields: \n")
		for fieldIndex,fieldName := range itemA.FieldsCount{
			fmt.Printf("A%d .- Type: %s\n",(fieldIndex+1),fieldName.Type)
			fmt.Printf("     Required: %s\n",strconv.FormatBool(fieldName.Required))
			fmt.Printf("     Omitted: %s\n",strconv.FormatBool(fieldName.Omitted))
			fmt.Printf("     Disabled: %s\n",strconv.FormatBool(fieldName.Disabled))
			fmt.Printf("     Localized: %s\n",strconv.FormatBool(fieldName.Localized))
			fmt.Println("")
		}
		fmt.Println("")
		for fieldIndex,fieldName := range itemB.FieldsCount{
			fmt.Printf("B%d .- Type: %s\n",(fieldIndex+1),fieldName.Type)
			fmt.Printf("     Required: %s\n",strconv.FormatBool(fieldName.Required))
			fmt.Printf("     Omitted: %s\n",strconv.FormatBool(fieldName.Omitted))
			fmt.Printf("     Disabled: %s\n",strconv.FormatBool(fieldName.Disabled))
			fmt.Printf("     Localized: %s\n",strconv.FormatBool(fieldName.Localized))
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func main() {

	var ctFilename, ctDirectory string
	var treshholdName,treshholdFields float64
	flag.StringVar(&ctFilename, "f", "", "Specify contentType JSON file.")
	flag.StringVar(&ctDirectory, "d", "", "Specify directory of contentType JSON files.")
	flag.Float64Var(&treshholdName, "tn", simLowerLimitPrecentageName, "Precentage treshhold for naming matching, default 75%")
	flag.Float64Var(&treshholdFields, "tf", simLowerLimitPrecentageFields, "Precentage treshhold for fields matching, default 85%")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -f  Content type json to parse.")
		fmt.Println("  -d  Folder of Content types json to parse.(multispace)")
		fmt.Println("  -tn  Precentage treshhold for naming matching, default 75%")
		fmt.Println("  -tf  Precentage treshhold for fields matching, default 85%")
		fmt.Println("")
		fmt.Printf("Usage of our Program: \n")
		fmt.Printf("./go-project -f /var/doc/MyContentTypeFile.json\n")
		fmt.Println("")
		fmt.Printf("How to generate it: \n")
		fmt.Println("***")
	}
	flag.Parse()

	if len(ctFilename) == 0 && len(ctDirectory) == 0 {
		fmt.Println("parameter 'f' with filename or 'd' with directory required")
		return
	}

	filesToParse := make([]string, 0, 10)
	if len(ctFilename) != 0 {
		if _, err := os.Stat(ctFilename); err != nil {
			fmt.Printf("File '%s' does not exist \n", ctFilename)
			return
		}
		filesToParse = append(filesToParse, ctFilename)
	} else if len(ctDirectory) != 0 {
		if _, err := os.Stat(ctDirectory); err != nil {
			fmt.Printf("Direcotry '%s' does not exist \n", ctDirectory)
			return
		}

		currentDirectory, err := os.Getwd()
		if err != nil {
			fmt.Printf("%+v", err.Error())
		}
		currentDirectory = fmt.Sprintf("%s/", ctDirectory)
		filesToParse = iterateDirectory(currentDirectory, ctDirectory)

	}

	contentTypesParsed := make([]spaceParsed,0, 100)
	for _, fileName := range filesToParse {
		contentTypesParsed=processJSONFile(fileName, contentTypesParsed,treshholdName,treshholdFields)
	}

	multiSpaceValidations(contentTypesParsed,treshholdName,treshholdFields)

}


