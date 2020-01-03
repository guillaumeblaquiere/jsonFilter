package main

import (
	"encoding/json"
	"fmt"
	"github.com/guillaumeblaquiere/jsonFilter"
)

func main() {
	// Get the filter from a request, in query parameter for example. Define the filters field that you want
	// filters, _ := r.URL.Query()["filters"]

	//filterValue := "Key1=val1,val2:composed.SubKey=val3:Maps.entry1.key1=val5,val4"

	// Filter on simple array of string
	//filterValue := "Key1=val1,val2:arraysimple=valArray2"

	// Filter on array of struct
	//filterValue := "Key1=val1,val2:array.SubKey=valArray"

	// Filter on map of struct
	//filterValue := "Key1=val1,val2:maps.entry=val"

	// Filter on map of String
	//filterValue := "Key1=val1,val2:mapsimple.entry=val"

	// Filter on map of Array of String
	//filterValue := "Key1=val1,val2:mapofarraysimple.entry=val"

	// Filter on map of Array of Struct
	//filterValue := "Key1=val1,val2:mapofarraystruct.entry.SubKey=val"

	// Filter on map of Array of Struct
	filterValue := "Key1=val1,val2:Matrix=AA"

	filter := jsonFilter.Filter{}
	if filterValue != "" {
		err := filter.Init(filterValue, structExample{})
		if err != nil {
			//TODO error handling
			fmt.Println(err)
			return
		}
	}

	// Perform your query for getting result and match it in the struct that you want to filter
	// Example with firestore
	/*
		results := make([]structExample,0)

		iter := client.Collection("myCollection").Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Error(err)
				break
			}
			var d structExample
			doc.DataTo(&d)
			results = append(results,d)
		}
	*/
	results := getDummyExamples()

	if filterValue != "" {
		ret, err := filter.ApplyFilter(results)
		if err != nil {
			//TODO error handling
			fmt.Println(err)
			return
		}
		results = ret.([]structExample)
	}

	toPrint, _ := json.Marshal(results)
	fmt.Println(string(toPrint))

}

// Generate an example of result to filter
// Only the 2 first elements of the array are ok with the filter
func getDummyExamples() []structExample {
	return []structExample{
		{
			Key1: "val1",
			Key2: &SecondStruct{"val3"},
			Array: []SecondStruct{
				{SubKey: "valArray"},
			},
			ArraySimple: []string{"valArray2"},
			Maps: map[string]structExample{
				"entry1": {
					Key1: "val4",
				},
				"entry2": {
					Key1: "val4",
				},
			},
			MapSimple: map[string]string{
				"entry": "val",
			},

			MapOfArraySimple: map[string][]string{
				"entry": {"val"},
			},

			MapOfArrayStruct: map[string][]SecondStruct{
				"entry": {
					SecondStruct{"val"},
					SecondStruct{"val2"},
				},
				"entry2": {
					SecondStruct{"val3"},
					SecondStruct{"val4"},
				},
			},
			Matrix: [][]string{
				{"AA", "AB", "AC"},
				{"BA", "BB", "BC"},
				{"CA", "CB", "CC"},
			},
		},
		{
			Key1: "val2",
			Key2: &SecondStruct{"val3"},
			Maps: map[string]structExample{
				"entry1": {
					Key1: "val5",
				},
			},
		},
		{
			Key1: "val2",
			Key2: &SecondStruct{"val"},
		},
		{
			Key1: "val1",
			Maps: map[string]structExample{
				"entry1": {
					Key1: "val6",
				},
				"entry2": {
					Key1: "val4",
				},
			},
		},
	}
}

type SecondStruct struct {
	SubKey string `json:"skey,omitempty"`
}

type structExample struct {
	Key1             string                    `json:"key1,omitempty"`
	Key2             *SecondStruct             `json:"composed,omitempty"`
	Array            []SecondStruct            `json:"array,omitempty"`
	ArraySimple      []string                  `json:"arraysimple,omitempty"`
	Maps             map[string]structExample  `json:"maps,omitempty"`
	MapSimple        map[string]string         `json:"mapsimple,omitempty"`
	MapOfArraySimple map[string][]string       `json:"mapofarraysimple,omitempty"`
	MapOfArrayStruct map[string][]SecondStruct `json:"mapofarraystruct,omitempty"`
	Matrix           [][]string                `json:"matrix,omitempty"`
}
