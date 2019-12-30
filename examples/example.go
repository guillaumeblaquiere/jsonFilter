package main

import (
	"encoding/json"
	"fmt"
	"github.com/guillaumeblaquiere/jsonFilter"
)

func main() {
	// Get the filter from a request, in query parameter for example. Define the filters field that you want
	// filters, _ := r.URL.Query()["filters"]

	filterValue := "Key1=val1,val2:composed.SubKey=val3"

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
			SecondKey: struct {
				SubKey string `json:"skey"`
			}{"val3"},
		},
		{
			Key1: "val2",
			SecondKey: struct {
				SubKey string `json:"skey"`
			}{"val3"},
		},
		{
			Key1: "val2",
			SecondKey: struct {
				SubKey string `json:"skey"`
			}{"val"},
		},
		{
			Key1: "val1",
		},
	}
}

type structExample struct {
	Key1      string `json:"key1"`
	SecondKey struct {
		SubKey string `json:"skey"`
	} `json:"composed"`
}
