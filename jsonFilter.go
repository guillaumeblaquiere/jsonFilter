/*
Apply a post processing filters to the Datastore/Firestore results mapped in struct with json tag or not.

The filter format is designed to be passed in API param (query or path). The filters can express compound operation.

During the processing the values to filter (an array) is passed to the library to apply the filters.
A filter can be composed to several part:
  - Several filter elements
  - Each filter element have  a tree path into the JSON, name the key, an operator and the value(s) to compare

Each filter element must return OK for keeping the entry value. For this, 4 operators are allowed:
  - The equality, comparable to IN sql clause: at least one value must matches. Default operator is `=`
  - The not equality, comparable to NOT IN sql clause: all values mustn't match. Default operator is `!=`
  - The Greater Than: only one numeric can be compared. Default operator is `>`
  - The Lower Than: only one numeric can be compared. Default operator is `<`

It's possible to combine operators on the same key, for example k1 < 10 && k1 != 2.
The same operator on the same key will raise an error.

The filters are applicable on this list types and structures (and combination possibles):
  - simple types
	- string
	- int
	- float
	- bool
  - Complex type
	- pointer (invisible in JSON result but your structure can include filters)
	- struct
	- array
	  - of simple types
	  - of map
	  - of array
	  - of pointer
	- map
	  - of simple types
	  - of map
	  - of array
	  - of pointer

This library works with Go app and use reflection. It performs 3 things
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
	"strconv"
	"strings"
)

/*
Structure to define the option of the Filter.

You can customize it if you want. Else the default values are applied
*/
type Options struct {
	// Limit the depth of the key search. In case of complex object, can limit the compute resources. 0 means infinite. Default is '0'
	MaxDepth int
	// Character(s) to separate key (filter name)  from values (value to compare) for an equal comparison. Default is '='
	EqualKeyValueSeparator string
	// Character(s) to separate key (filter name)  from values (value to compare) for a greater than comparison. Default is '>'
	GreaterThanKeyValueSeparator string
	// Character(s) to separate key (filter name)  from values (value to compare) for a lower than comparison. Default is '<'
	LowerThanKeyValueSeparator string
	// Character(s) to separate key (filter name)  from values (value to compare) for a not equal comparison. Default is '!='
	NotEqualKeyValueSeparator string
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
	filter  []kov
}

type kov struct {
	Key      string
	Operator string
	Values   []string
}

// Default option used in case of no specific set.
var defaultOption = &Options{
	MaxDepth:                     0,
	EqualKeyValueSeparator:       "=",
	GreaterThanKeyValueSeparator: ">",
	LowerThanKeyValueSeparator:   "<",
	NotEqualKeyValueSeparator:    "!=",
	ValueSeparator:               ",",
	KeysSeparator:                ":",
	ComposedKeySeparator:         ".",
}

/*
Set the option to the filter.

If the option is nil, the default option will be used.

If there is some missing or incorrect value to the defined option, a warning message is displayed and the erroneous part
is replace by the default ones.

To set option:
	filter := jsonFilter.Filter{}

	o := &jsonFilter.Options{
		MaxDepth:             			4,
		EqualKeyValueSeparator:    		"=",
  		GreaterThanKeyValueSeparator: 	">",
		LowerThanKeyValueSeparator:   	"<",
		NotEqualKeyValueSeparator:    	"!=",
		ValueSeparator:       			",",
		KeysSeparator:        			":",
		ComposedKeySeparator: 			"->",
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
	if o.EqualKeyValueSeparator == "" {
		o.EqualKeyValueSeparator = defaultOption.EqualKeyValueSeparator
		log.Warnf("EqualKeyValueSeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.EqualKeyValueSeparator)
	}
	if o.GreaterThanKeyValueSeparator == "" {
		o.GreaterThanKeyValueSeparator = defaultOption.GreaterThanKeyValueSeparator
		log.Warnf("GreaterThanKeyValueSeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.GreaterThanKeyValueSeparator)
	}
	if o.LowerThanKeyValueSeparator == "" {
		o.LowerThanKeyValueSeparator = defaultOption.LowerThanKeyValueSeparator
		log.Warnf("LowerThanKeyValueSeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.LowerThanKeyValueSeparator)
	}
	if o.NotEqualKeyValueSeparator == "" {
		o.NotEqualKeyValueSeparator = defaultOption.NotEqualKeyValueSeparator
		log.Warnf("NotEqualKeyValueSeparator can't be empty. Option entry ignored, default used %q \n", defaultOption.NotEqualKeyValueSeparator)
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
  - Duplicated entry in the filter key name for the same operator
  - Violation of filter format:
    - No values for a key
    - No key for a filter
    - More than 1 value for Greater Than and Lower than operator
    - Not a numeric (float compliant) value for Greater Than and Lower than operator
  - Filter key not exist in the provided interface
    - Struct field name not match the filter key
    - Struct json tag not match the filter key

*/
func (f *Filter) Init(v string, i interface{}) (err error) {
	if f.options == nil {
		f.options = defaultOption
	}
	fts, err := f.parseFilter(v)
	if err != nil {
		return
	}
	err = f.compileFilter(fts, reflect.TypeOf(i))
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
func (f *Filter) ApplyFilter(e interface{}) (interface{}, error) {

	eav := reflect.ValueOf(e) //entry array value
	// Check if the entry is an array
	if eav.Kind() != reflect.Slice {
		log.Errorf("The entries is not of type Array but of type %s. Filter can be applied only on an array", eav.Type())
		return nil, errors.New("internal error")
	}

	//Init ret with the max possible length
	ret := reflect.MakeSlice(eav.Type(), 0, eav.Len())

	// Iterate on all e
	for i := 0; i < eav.Len(); i++ {
		evs := eav.Index(i)
		// Flag for keeping or not the entry in the result set
		keepV := true

		// Apply all the filters
		for _, kov := range f.filter {
			k := kov.Key
			// Get the values of the entry for this Filter
			// Find all possible values per entry in case of composite key
			evl := f.findValueInComposedKey(k, evs) //entry value list

			// Flag to know if at least one Filter values matches the entry value
			m := false // match

			// Select the correct operator
			switch kov.Operator {
			case f.options.EqualKeyValueSeparator:
				// Iterate over the matching value of the entry
				for _, ev := range evl {
					// Iterate over the filter possible value.
					for _, v := range kov.Values {
						// If only one matches, the IN operator is valid
						if fmt.Sprint(ev) == v {
							m = true
							break
						}
					}
					// If the field match stop the loop: At least 1 of the Filter Values have to match (OR condition)
					if m {
						break
					}
				}
			case f.options.NotEqualKeyValueSeparator:
				m = true
				// Iterate over the matching value of the entry
				for _, ev := range evl {
					// Iterate over the filter possible value.
					for _, v := range kov.Values {
						// If only one value matches, the NOT IN operator doesn't match: All values must not be in
						if fmt.Sprint(ev) == v {
							m = false
							break
						}
					}
					if !m {
						break
					}
				}
			case f.options.GreaterThanKeyValueSeparator:
				for _, ev := range evl {
					v := kov.Values[0] // always 1 values for greater than operator
					//Compare only numeric values
					vf, _ := strconv.ParseFloat(v, 10) // assume that possible thanks to parser check
					evf, err := strconv.ParseFloat(fmt.Sprint(ev), 10)
					if err == nil && evf > vf {
						m = true
						break
					}
				}
			case f.options.LowerThanKeyValueSeparator:
				for _, ev := range evl {
					v := kov.Values[0] // always 1 values for greater than operator
					//Compare only numeric values
					vf, _ := strconv.ParseFloat(v, 10) // assume that possible thanks to parser check
					evf, err := strconv.ParseFloat(fmt.Sprint(ev), 10)
					if err == nil && evf < vf {
						m = true
						break
					}
				}
			}

			//If any the Filter value matches the entry field value, we don't keep it in the result set
			// and break the loop because all fields must match. If one fail, stop here
			if !m {
				keepV = false
				break
			}
		}
		// If all fields matches, keep the entry in the result set
		if keepV {
			ret = reflect.Append(ret, evs)
		}
	}
	return ret.Interface(), nil
}

// Find all values (leaf value) associated with a composed key (filter name).
// Return always an array of values in case of search in sub elements which are an array of structs
func (f *Filter) findValueInComposedKey(k string, evs reflect.Value) []reflect.Value {
	kp := strings.Split(k, f.options.ComposedKeySeparator) //key p
	vs := []reflect.Value{evs}                             // values

	//Scan all p of the composed key, going deeper and deeper
	for _, p := range kp {
		r := make([]reflect.Value, 0) //result

		// Scan recursively all sub values found
		for _, v := range vs {

			var res reflect.Value
			// If the current element is a map
			if v.Kind() == reflect.Map {
				// search the matching key in the value list
				foundEntry := false
				for _, val := range v.MapKeys() {
					if fmt.Sprint(val) == p {
						res = v.MapIndex(val)
						foundEntry = true
						break //only one entry in the map key list
					}
				}
				if !foundEntry {
					//If no entry match the key of the map key list, continue to the next value, forget this p of the tree
					continue
				}
			} else { // if not, scan the structure
				res = v.FieldByName(p)
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
				r = extractValueFromSlice(r, res)
			} else {
				r = append(r, res)
			}

		}
		vs = r
	}
	return vs
}

//Recursive loop for getting all the values from a Tensor (array of N dimension)
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

// Return error if there is 2 times the same key with the same operator
// return error if the composed filter depth is higher than this defined in options (0 = infinite)
// Filter default option pattern is key1=value1,value2:key1!=value:key2=value3,value4
func (f *Filter) parseFilter(filterInput string) (kovs []kov, err error) {
	kovs = []kov{}
	fts := strings.Split(filterInput, f.options.KeysSeparator)

	// Parse all fts found
	for _, ft := range fts {

		kv, op := f.getFilterAndValue(ft)

		// If there isn't values part, it's an error
		if !isKeyValuesValidPair(kv) {
			return nil, errors.New(fmt.Sprintf("No values defined for the key %s ft", ft))
		}

		k := kv[0]

		// Check if key is not empty
		if k == "" {
			return nil, errors.New("No filter key")
		}

		// Check the max depth
		if f.options.MaxDepth > 0 && len(strings.Split(k, f.options.ComposedKeySeparator)) > f.options.MaxDepth {
			return nil, errors.New(fmt.Sprintf("The Filter key %s doen't match the max depth key set to %d", k, f.options.MaxDepth))
		}

		// Check if the key with the same operator has been already set in the map
		for _, kov := range kovs {
			if kov.Key == k && kov.Operator == op {
				return nil, errors.New(fmt.Sprintf("The key %s already for the operator %s exist in the Filter field", k, op))
			}
		}

		// extract the values and set them to the map
		v := strings.Split(kv[1], f.options.ValueSeparator)

		if op == f.options.GreaterThanKeyValueSeparator || op == f.options.LowerThanKeyValueSeparator {
			if len(v) > 1 {
				return nil, errors.New("the Filter 'greater than' and 'lower than' must have exactly 1 value")
			}
			if _, err := strconv.ParseFloat(v[0], 10); err != nil {
				return nil, errors.New(fmt.Sprintf("the Filter 'greater than' and 'lower than' must have a numeric value. Here %v", filterInput[0]))
			}
		}

		kovs = append(kovs, kov{
			Key:      k,
			Operator: op,
			Values:   v,
		})
	}
	return
}

// Get the key and the values of the filters by testing possible operators
// return an empty array if any separator matches
// In case of 2 separators work when splitting the filter, we keep only the longest separator
// Example: in case of = and != both will split on =, but we only keep != because it's the longest
func (f *Filter) getFilterAndValue(filter string) (fkvs []string, op string) {

	if fkv := strings.Split(filter, f.options.EqualKeyValueSeparator); isKeyValuesValidPair(fkv) && len(op) < len(f.options.EqualKeyValueSeparator) {
		fkvs = fkv
		op = f.options.EqualKeyValueSeparator
	}

	if fkv := strings.Split(filter, f.options.NotEqualKeyValueSeparator); isKeyValuesValidPair(fkv) && len(op) < len(f.options.NotEqualKeyValueSeparator) {
		fkvs = fkv
		op = f.options.NotEqualKeyValueSeparator
	}

	if fkv := strings.Split(filter, f.options.GreaterThanKeyValueSeparator); isKeyValuesValidPair(fkv) && len(op) < len(f.options.GreaterThanKeyValueSeparator) {
		fkvs = fkv
		op = f.options.GreaterThanKeyValueSeparator
	}

	if fkv := strings.Split(filter, f.options.LowerThanKeyValueSeparator); isKeyValuesValidPair(fkv) && len(op) < len(f.options.LowerThanKeyValueSeparator) {
		fkvs = fkv
		op = f.options.LowerThanKeyValueSeparator
	}

	return
}

func isKeyValuesValidPair(fkvs []string) bool {
	return len(fkvs) == 2
}

// Find the struct field name in relation with the Filter name provided in the query
// The search is performed in the json tag of the struct field and on the struct field name in case of missing tag;
func (f *Filter) compileFilter(kovs []kov, t reflect.Type) (err error) {
	f.filter = []kov{}

	//for all  filters, search is a struct field name match with it
	for _, kov := range kovs {
		k := kov.Key
		ck := "" // composed key
		ckp := strings.Split(k, f.options.ComposedKeySeparator)
		ct := t //current type

		// validate the struct field name according with the key composition. Going deeper and deeper
		for i, p := range ckp {
			var cp string //composed part

			// If array, loop on it to get the element contained in the tensor (Array of N dimension)
			for ct.Kind() == reflect.Slice {
				ct = ct.Elem()
				//In case of ptr
				if ct.Kind() == reflect.Ptr {
					ct = ct.Elem()
				}
			}
			//if map, keep the key as is
			if ct.Kind() == reflect.Map {
				cp = p
				ct = ct.Elem()
				//In case of ptr
				if ct.Kind() == reflect.Ptr {
					ct = ct.Elem()
				}
			} else { // look into the structure

				fs := foundFieldInStruct(p, ct)
				// If no match found, raise an error
				if fs == nil {
					log.Debugf("The Filter key %s not exist in the type %s", p, t.Name())
					return errors.New(fmt.Sprintf("The Filter key %s not exist in the returned object", ck+" "+p))
				}
				ct = fs.Type

				// If it's not the root element of the composed key, add a separator the the filter name
				cp = fs.Name
			}
			//Add a composed separator if it's not the root p
			if i != 0 {
				ck += f.options.ComposedKeySeparator
			}
			ck += cp
		}
		kov.Key = ck
		f.filter = append(f.filter, kov)
	}
	return
}

// Return the structField found according with the filter key name and the type to scan.
// Return nil if nothing found in the type.
func foundFieldInStruct(k string, t reflect.Type) *reflect.StructField {
	ct := t // current type
	//In case of ptr
	if t.Kind() == reflect.Ptr {
		ct = t.Elem()
	}

	for i := 0; i < ct.NumField(); i++ {
		f := ct.Field(i)

		// on each fields, check if the Filter can be applied
		if strings.Contains(f.Tag.Get("json"), k) ||
			f.Name == k {
			// When found,add it to the map and go to the next Filter f
			return &f
		}
	}
	return nil
}
