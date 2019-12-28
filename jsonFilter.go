package jsonFilter

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

/*
Structure to define the option of the Filter

	MaxDepth             Limit the depth of the key search. in case of complex object, can limit the compute resources
	KeyValueSeparator    Character(s) to separate key (filter name)  from values (value to compare). Default is '='
	ValueSeparator       Character(s) to separate values (value to compare). Default is ','
	KeysSeparator        Character(s) to separate keys (filters name). Default is ':'
	ComposedKeySeparator Character(s) to separate key part in case of composed key (filter.subfilter) . Default is '.'
*/
type Options struct {
	MaxDepth             int
	KeyValueSeparator    string
	ValueSeparator       string
	KeysSeparator        string
	ComposedKeySeparator string
}

type Filter struct {
	options *Options
	filter  map[string][]string
}

// Default option used in case of no specific set.
var defaultOption = &Options{
	MaxDepth:             0,
	KeyValueSeparator:    "=",
	ValueSeparator:       ",",
	KeysSeparator:        ":",
	ComposedKeySeparator: ".",
}

// If options are set, all the fields have to be defined. No empty string allowed, no negative max depth allowed
// In case of error, no options are applied and the default option is used instead.
func (f *Filter) SetOptions(o *Options) error {
	if o == nil {
		errorMessage := "options can't be nil. Options ignored, default used"
		f.options = defaultOption
		log.Error(errorMessage)
		return errors.New(errorMessage)
	}
	if o.MaxDepth < 0 {
		errorMessage := "MaxDepth must be positive. 0 min infinite depth. Options ignored, default used"
		f.options = defaultOption
		log.Error(errorMessage)
		return errors.New(errorMessage)
	}
	if o.KeyValueSeparator == "" {
		errorMessage := "KeyValueSeparator can't be empty. Options ignored, default used"
		f.options = defaultOption
		log.Error(errorMessage)
		return errors.New(errorMessage)
	}
	if o.ValueSeparator == "" {
		errorMessage := "ValueSeparator can't be empty. Options ignored, default used"
		f.options = defaultOption
		log.Error(errorMessage)
		return errors.New(errorMessage)
	}
	if o.KeysSeparator == "" {
		errorMessage := "KeysSeparator can't be empty. Options ignored, default used"
		f.options = defaultOption
		log.Error(errorMessage)
		return errors.New(errorMessage)
	}
	if o.ComposedKeySeparator == "" {
		errorMessage := "ComposedKeySeparator can't be empty. Options ignored, default used"
		f.options = defaultOption
		log.Error(errorMessage)
		return errors.New(errorMessage)
	}
	f.options = o
	return nil
}

// Initialize the filter with the requested filter and the struct on which to apply later the filter
// The filter result is saved in the Filter struct
func (f *Filter) Init(filterValue string, i interface{}) (err error) {
	if f.options == nil {
		f.options = defaultOption
	}
	filters, err := f.parseFilter(filterValue)
	if err != nil {
		return
	}
	err = f.compileFilter(filters, reflect.TypeOf(i))
	if err != nil {
		return
	}
	return
}

/*
Apply the typed Filter to the entries. The entries must be an array.
Return an array with only the matching entries

Matching the filters means

 * match all defined filters (AND condition)
 * On one Filter, match at least 1 value of the values list defined on the Filter (OR condition)
*/
func (f *Filter) ApplyFilter(entries interface{}) (interface{}, error) {

	entriesReflected := reflect.ValueOf(entries)
	// Check if the entry is an array
	if entriesReflected.Kind() != reflect.Slice {
		log.Errorf("The entries is not of type Array but of type %s. Filter can be applied only on an array", entriesReflected.Type())
		return nil, errors.New("internal error")
	}

	//Init ret with the max possible length
	ret := reflect.MakeSlice(entriesReflected.Type(), 0, entriesReflected.Len())

	// Iterate on all entries
	for i := 0; i < entriesReflected.Len(); i++ {
		entryValues := entriesReflected.Index(i)
		// Flag for keeping or not the entry in the result set
		keepValue := true

		// Apply all the filters
		for filterKey, filterValues := range f.filter {
			// Get the values of the entry for this Filter
			// Find all possible values per entry in case of composite key
			entryValueList := f.findValueInComposedKey(filterKey, entryValues)

			// Flag to know if at least one Filter values matches the entry value
			filterMatch := false

			// Loop on the Filter values and test all against the entry field value
			for _, filterValue := range filterValues {
				for _, entryValue := range entryValueList {
					//Use Sprint for converting the entry value to String
					if fmt.Sprint(entryValue) == filterValue {
						// If the field match, flag it and stop the loop: At least 1 of the Filter Values have to match (OR condition)
						filterMatch = true
						break
					}
				}
				// Key the object if at least one of the leaf match
				if filterMatch {
					break
				}
			}
			//If any the Filter value matches the entry field value, we don't keep it in the result set
			// and break the loop because all fields must match. If one fail, stop here
			if !filterMatch {
				keepValue = false
				break
			}
		}
		// If all fields matches, keep the entry in the result set
		if keepValue {
			ret = reflect.Append(ret, entryValues)
		}
	}
	return ret.Interface(), nil
}

// Find all values (leaf value) associated with a composed key (filter name).
// Return always an array of values in case of search in sub elements which are an array of structs
func (f *Filter) findValueInComposedKey(filterKey string, entryValues reflect.Value) []reflect.Value {
	filterKeyPart := strings.Split(filterKey, f.options.ComposedKeySeparator)
	valuesToScan := []reflect.Value{entryValues}

	//Scan all part of the composed key, going deeper and deeper
	for _, part := range filterKeyPart {
		scanResult := make([]reflect.Value, 0)

		// Scan recursively all sub values found
		for _, valueToScan := range valuesToScan {
			res := valueToScan.FieldByName(part)

			// In case of array found, add all the matching values to the result (or next value th scan if not the leaf)
			if res.Kind() == reflect.Slice {
				for i := 0; i < res.Len(); i++ {
					scanResult = append(scanResult, res.Index(i))
				}
			} else {
				scanResult = append(scanResult, res)
			}
		}
		valuesToScan = scanResult
	}
	return valuesToScan
}

// Return error if 2 time the same key in the Filter
// return error if the composed filter depth is higher than this defined in options (0 = infinite)
// Filter default option pattern is key1=value1,value2:key2=value3,value4
func (f *Filter) parseFilter(filterValue string) (filterMap map[string][]string, err error) {
	filterMap = map[string][]string{}
	filters := strings.Split(filterValue, f.options.KeysSeparator)

	// Parse all filters found
	for _, filter := range filters {
		filterKeyValue := strings.Split(filter, f.options.KeyValueSeparator) // index 0 = key, index 1 = value(s)

		// If there isn't values part, it's an error
		if len(filterKeyValue) < 2 {
			return nil, errors.New(fmt.Sprintf("No values defined for the key %s filter", filter))
		}

		key := filterKeyValue[0]
		// Check the max depth
		if f.options.MaxDepth > 0 && len(strings.Split(key, f.options.ComposedKeySeparator)) > f.options.MaxDepth {
			return nil, errors.New(fmt.Sprintf("The Filter key %s doen't match the max depth key set to %d", key, f.options.MaxDepth))
		}

		// Check if the key has been already set in the map
		_, here := filterMap[key]
		if here {
			return nil, errors.New(fmt.Sprintf("The key %s already exist in the Filter field", filter))
		}

		// extract the values and set them to the map
		filterValues := strings.Split(filterKeyValue[1], f.options.ValueSeparator)
		filterMap[key] = filterValues
	}
	return
}

// Find the struct field name in relation with the Filter name provided in the query
// The search is performed in the json tag of the struct field and on the struct field name in case of missing tag;
func (f *Filter) compileFilter(filterMap map[string][]string, t reflect.Type) (err error) {
	f.filter = map[string][]string{}

	//for all  filters, search is a struct field name match with it
	for filterKey, filterValues := range filterMap {
		composedFilterKey := ""
		filterKeyPart := strings.Split(filterKey, f.options.ComposedKeySeparator)
		objectToInspect := t

		// validate the struct field name according with the key composition. Going deeper and deeper
		for i, part := range filterKeyPart {
			fieldStruct := foundFieldInStruct(part, objectToInspect)
			// If no match found, raise an error
			if fieldStruct == nil {
				log.Debugf("The Filter key %s not exist in the type %s", part, t.Name())
				return errors.New(fmt.Sprintf("The Filter key %s not exist in the returned object", composedFilterKey+" "+part))
			}
			objectToInspect = fieldStruct.Type

			// If slice, get only the elements type
			if objectToInspect.Kind() == reflect.Slice {
				objectToInspect = objectToInspect.Elem()
			}

			// If it's not the root element of the composed key, add a separator the the filter name
			if i != 0 {
				composedFilterKey += f.options.ComposedKeySeparator
			}
			composedFilterKey += fieldStruct.Name
		}
		f.filter[composedFilterKey] = filterValues
	}
	return
}

// Return the structField found according with the filter key name and the type to scan.
// Return nil if nothing found in the type.
func foundFieldInStruct(filterKey string, t reflect.Type) *reflect.StructField {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// on each fields, check if the Filter can be applied
		if field.Tag.Get("json") == filterKey ||
			field.Name == filterKey {
			// When found,add it to the map and go to the next Filter field
			return &field
		}
	}
	return nil
}