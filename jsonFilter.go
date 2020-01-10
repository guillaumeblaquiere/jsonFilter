/*
Apply a post processing filters to the Datastore/Firestore results mapped in struct with json tag or not.

The filter is API oriented and designed to be provided by an API consumer in param to your its request.
This library work with Go app and use reflection. It performs 3 things
  - Check if the provided filter is valid.
  - Compile the filter according with the data structure to filter -> Validate the filter against the structure to filter.
  - Apply the filter to the array of structure.
*/
package jsonFilter

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

/*
Structure to define the option of the Filter.

You can customize it if you want. Else the default values are applied
*/
type Options struct {
	// Limit the depth of the key search. In case of complex object, can limit the compute resources. 0 means infinite. Default is '0'
	MaxDepth int
	// Character(s) to separate key (filter name)  from values (value to compare). Default is '='
	KeyValueSeparator string
	//  Character(s) to separate values (value to compare). Default is ','
	ValueSeparator string
	// Character(s) to separate keys (filters name). Default is ':'
	KeysSeparator string
	// Character(s) to separate key part in case of composed key (filter.subfilter) . Default is '.'
	ComposedKeySeparator string
}

/*
Filter structure to use for filtering. Init the default value like this
  filter := jsonFilter.Filter{}

*/
type Filter struct {
	// Only private fields
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

/*
Set the option to the filter.

If the option is nil, the default option will be used.

If there is some missing or incorrect value to the defined option, a warning message is displayed and the erroneous part
is replace by the default ones.

To set option:
	filter := jsonFilter.Filter{}

	o := &jsonFilter.Options{
		MaxDepth:             4,
		KeyValueSeparator:    "=",
		ValueSeparator:       ",",
		KeysSeparator:        ":",
		ComposedKeySeparator: "->",
	}

	filter.SetOptions(o)

*/
func (f *Filter) SetOptions(o *Options) {
	if o == nil {
		o = defaultOption
		log.Warn("options can't be nil. Options ignored, default used")
	}
	if o.MaxDepth < 0 {
		o.MaxDepth = defaultOption.MaxDepth
		log.Warnf("MaxDepth must be positive. 0 means infinite depth. Option entry ignored, default used %q \n", defaultOption.MaxDepth)
	}
	if o.KeyValueSeparator == "" {
		o.KeyValueSeparator = defaultOption.KeyValueSeparator
		log.Warnf("KeyValueSeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.KeyValueSeparator)
	}
	if o.ValueSeparator == "" {
		o.ValueSeparator = defaultOption.ValueSeparator
		log.Warnf("ValueSeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.ValueSeparator)
	}
	if o.KeysSeparator == "" {
		o.KeysSeparator = defaultOption.KeysSeparator
		log.Warnf("KeysSeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.KeysSeparator)
	}
	if o.ComposedKeySeparator == "" {
		o.ComposedKeySeparator = defaultOption.ComposedKeySeparator
		log.Warnf("ComposedKeySeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.ComposedKeySeparator)
	}
	f.options = o
}

/*
Initialize the filter with the requested filter and the struct on which to apply later the filter

The filter parsing and compilation are saved in the Filter struct.

Errors are returned in case of:
  - Duplicated entry in the filter keyname
  - Violation of filter format:
    - No values for a key
    - No key for a filter
  - Filter key not exist in the provided interface
    - Struct field name not match the filter key
    - Struct json tag not match the filter key

*/
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
Apply the initialized Filter to a list (array) of struct. The type of array elements is the same as this one provided
in the Init method. The entries must be an array.

Return an array with only the matching entries, else an error is returned.

Cast the return in the array type like this:
	ret, err := filter.ApplyFilter(results)
	if err != nil {
		// Perform error handling
		fmt.Println(err)
		return
	}
	// Cast here the return
	results = ret.([]structExample)

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
			val := valueToScan

			var res reflect.Value
			// If the current element is a map
			if val.Kind() == reflect.Map {
				// search the matching key in the value list
				foundEntry := false
				for _, v := range val.MapKeys() {
					if fmt.Sprint(v) == part {
						res = val.MapIndex(v)
						foundEntry = true
						break //only one entry in the map key list
					}
				}
				if !foundEntry {
					//If no entry match the key of the map key list, continue to the next value, forget this part of the tree
					continue
				}
			} else { // if not, scan the structure
				res = val.FieldByName(part)
			}

			//In case of pointer
			if res.Kind() == reflect.Ptr {
				//If the pointer lead to nil value
				if res.Pointer() == 0 {
					continue
				}
				res = res.Elem()
			}

			// In case of array found, add all the matching values to the result (or next value th scan if not the leaf)
			if res.Kind() == reflect.Slice {
				scanResult = extractValueFromSlice(scanResult, res)
			} else {
				scanResult = append(scanResult, res)
			}

		}
		valuesToScan = scanResult
	}
	return valuesToScan
}

/*
Recursive loop for getting all the values from a Tensor (array of N dimension)
*/
func extractValueFromSlice(r []reflect.Value, v reflect.Value) []reflect.Value {
	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Kind() == reflect.Slice {
			r = extractValueFromSlice(r, v.Index(i))
		} else {
			curVal := v.Index(i)
			//In case of pointer
			if curVal.Kind() == reflect.Ptr {
				//If the pointer lead to nil value
				if curVal.Pointer() == 0 {
					continue
				}
				curVal = curVal.Elem()
			}
			r = append(r, curVal)
		}
	}
	return r
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

		// Check if key is not empty
		if key == "" {
			return nil, errors.New("No filter key")
		}

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
			var composedPart string

			// If array, loop on it to get the element contained in the tensor (Array of N dimension)
			for objectToInspect.Kind() == reflect.Slice {
				objectToInspect = objectToInspect.Elem()
				//In case of ptr
				if objectToInspect.Kind() == reflect.Ptr {
					objectToInspect = objectToInspect.Elem()
				}
			}
			//if map, keep the key as is
			if objectToInspect.Kind() == reflect.Map {
				composedPart = part
				objectToInspect = objectToInspect.Elem()
				//In case of ptr
				if objectToInspect.Kind() == reflect.Ptr {
					objectToInspect = objectToInspect.Elem()
				}
			} else { // look into the structure

				fieldStruct := foundFieldInStruct(part, objectToInspect)
				// If no match found, raise an error
				if fieldStruct == nil {
					log.Debugf("The Filter key %s not exist in the type %s", part, t.Name())
					return errors.New(fmt.Sprintf("The Filter key %s not exist in the returned object", composedFilterKey+" "+part))
				}
				objectToInspect = fieldStruct.Type

				// If it's not the root element of the composed key, add a separator the the filter name
				composedPart = fieldStruct.Name
			}
			//Add a composed separator if it's not the root part
			if i != 0 {
				composedFilterKey += f.options.ComposedKeySeparator
			}
			composedFilterKey += composedPart
		}
		f.filter[composedFilterKey] = filterValues
	}
	return
}

// Return the structField found according with the filter key name and the type to scan.
// Return nil if nothing found in the type.
func foundFieldInStruct(filterKey string, t reflect.Type) *reflect.StructField {
	internalType := t
	//In case of ptr
	if t.Kind() == reflect.Ptr {
		internalType = t.Elem()
	}

	for i := 0; i < internalType.NumField(); i++ {
		field := internalType.Field(i)

		// on each fields, check if the Filter can be applied
		if strings.Contains(field.Tag.Get("json"), filterKey) ||
			field.Name == filterKey {
			// When found,add it to the map and go to the next Filter field
			return &field
		}
	}
	return nil
}
