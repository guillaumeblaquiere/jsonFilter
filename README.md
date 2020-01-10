# Overview

[![Build Status](https://github.com/guillaumeblaquiere/jsonFilter/workflows/build/badge.svg)](https://github.com/guillaumeblaquiere/jsonFilter/actions?query=workflow%3Abuild)
[![Test Status](https://github.com/guillaumeblaquiere/jsonFilter/workflows/test/badge.svg)](https://github.com/guillaumeblaquiere/jsonFilter/actions?query=workflow%3Atest)
[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/guillaumeblaquiere/jsonFilter)

This library allows you to apply **post processing filters** to the Datastore/Firestore results. 

The filter is API oriented and designed to be provided by an API consumer in param to your its request.

# Why to use this additional filter

Firestore and datastore have several [query limitation for Firestore](https://firebase.google.com/docs/firestore/query-data/queries#query_limitations)
and [for Datastore](https://cloud.google.com/datastore/docs/concepts/queries#restrictions_on_queries)

- Only one `array-contains-any` or `IN` is allowed per query
- Only 10 elements are allowed in a `IN` or `array-contains-any` clause
- When you filter on range AND a field, you need to create a composite index beforehand
- You can search in the arrays that only contain values, and not object structure.

The library allows you to unlock these limitation:

- Allow to filter on any type and nested type (map, array, map of array/object, array of map/object,...)
- Allow to use several filters `IN` on the set of data
- Use more than 10 elements in a `IN` condition
- Don't required any composite index creation 

## limitation
 
There is the known limitation of this library. These can be implemented -> Open a feature request!
 
 - No wildcard like * to replace any JSON field name.
 - No wildcard like * or regex to filter on values
 - Only equals comparison is available. Others (!= < >) aren't available

## Performance concern

The filters should be applied only on a small array of results and the filtering overhead is very 
small

Indeed, the documents have to be read for being filter. If you read thousand of document, 
you will pay a lot for nothing
 
In addition, your API response time will take more time because of the high number of documents 
to recover
 and the filtering duration. 

# Way of working

This library work with Go app and use reflection. It performs 3 things

- Check if the provided filter is valid 
- Compile the filter according with the data structure to filter -> Validate the filter against
 the structure to filter
- Apply the filter to the array of structure

See [example](https://github.com/guillaumeblaquiere/jsonFilter/blob/master/examples/example.go) for a practical implementation.

# Filter format

The default filter format is the following
```
key1=val1,val2:key2.subkey=val3
```
**Where:**

-key1 is the JSON field name to filter. You can use composed filter to browse your JSON tree, 
like key2.subkey
- Val1, val2, val3 are the values to compare
- The tuple key + value(s) is named Filter

**Behavior:**

- A filter is considered as OK if at least one of values matches (equivalent of OR or IN
 conditions between values)
- When the filter is applied to an array of structs, the struct is kept in the result if 
all the filters are OK (AND condition between filters)

## Customize filter format

The default filter format use these caracter

- Keys and values are separated by equal sign = by default
- Filters are separated by colon : by default
- Values are separated by comma , by default
- Different fields value of a composed key is dot . by default

You can set an Options structure on filter to customize your filter like this

```
	o := &jsonFilter.Options{
		MaxDepth:             4,
		KeyValueSeparator:    "=",
		ValueSeparator:       ",",
		KeysSeparator:        ":",
		ComposedKeySeparator: "->",
	}
	
	filter.SetOptions(o)
```

If you don't define a part of the option, the default value is used for this part (a log 
message display this)

## Max depth

You can also define the max depth of composed key. By default, this value is set to 0, 
which means infinite. You can override this value in the option structure.

# Filter value type

You can filter on these simple types

- string
- int
- float
- bool

Complex type are supported
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
 
## Special filter on map
In JSON, the map representation is the following
```
{
    "mapsSimple":{
        "entryMap1":"value1",
        "entryMap2":"value2"
    },
    "mapsStruct":{
        "entryMap1": {
            "fieldName":"value1"
        },
        "entryMap2":{
            "fieldName":"value2"
        },
    },
    "mapsArray":{
        "entryMap1": [
            {
                "fieldName":"value1"
            }
        ],
        "entryMap2":[
            {
                "fieldName":"value2"
            },
        ]
    },

}
```

The filter key will be the following

- `mapsSimple.entryMap1` if it's a simple map
- `mapsStruct.entryMap1.fieldName` if it's a map of structure
- `mapsArray.entryMap1.fieldName` if it's a map of Array. The array is invisible in the processing

# Licence

This library is licensed under Apache 2.0. Full license text is available in [LICENSE.](https://github.com/guillaumeblaquiere/jsonFilter/blob/master/LICENSE)