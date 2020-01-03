# Overview

This library allow you to apply a post processing filters to the Datastore/Firestore results. 

The filter is API oriented and is designed to be provided by a consumer in param to your its request.

# Way of working

This library work with Go app and use reflection. It performs 3 things

- Check if the provided filter is valid 
- Compile the filter according with the data structure to filter -> Validate the filter against the structure to filter
- Apply the filter to the array of structure

# Why to use this additional filter

# Filter format

The default filter format is the following
```
key1=val1,val2:key2.subkey=val3
```
**Where:**

-key1 is the JSON field name to filter. You can use composed filter to browse your JSON tree, like key2.subkey
- Val1, val2, val3 are the values to compare
- The tuple key + value(s) is named Filter

**Behavior:**

- A filter is considered as OK if at least one of values matches (equivalent of OR or IN conditions between values)
- When the filter is applied to an array of structs, the struct is kept in the result if all the filters are OK (AND condition between filters)

## Customize filter format

The default filter format use these caracter

- Keys and values are separated by equal sign = by default
- Filters are separated by colon : by default
- Values are separated by comma , by default
- Different fields value of a composed key is dot . by default

You can set an Options structure on filter to customize your filter like this

```

```

If you don't define a part of the option, the default value is used for this part (a log message display this)

## Max depth

You can also define the max depth of composed key. By default, this value is set to 0, which means infinite. You can override this value in the option structure.

# Filter value type

You can filter on these simple types

- string
- int
- float
- bool

Complex type are supported
- pointer (invisible in JSON result but your structure can include filters)
- array
- struct
- map

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

}
```

The filter key will be the following

- `mapsSimple.entryMap1` if it's a simple map
- `mapsStruct.entryMap1.fieldName` if it's a map of structure

# limitation
No wildcard like * to replace any JSON field name.
Only equals comparison is available. Others (!= < >) can be implemented -> Open a feature request!
# Licence