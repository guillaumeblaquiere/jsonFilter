package jsonFilter

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFilter_ApplyFilter(t *testing.T) {
	type fields struct {
		options *Options
		filter  map[string][]string
	}
	type args struct {
		entries interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Minimal filter",
			fields: fields{
				options: defaultOption,
				filter: map[string][]string{
					"FirstLevelString": {"value1"},
				},
			},
			args: args{entries: []testStruct{
				{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{},
				},
			}},
			want: []testStruct{
				{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{},
				},
			},
			wantErr: false,
		},
		{
			name: "Composite filter",
			fields: fields{
				options: defaultOption,
				filter: map[string][]string{
					"FirstLevelArray.LevelString": {"string1"},
				},
			},
			args: args{entries: []testStruct{
				{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray: []secondLevelStruct{
						{LevelString: "string1"},
					},
					FirstLevelStruct: secondLevelStruct{},
				},
				{
					FirstLevelString: "value2",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray: []secondLevelStruct{
						{LevelString: "string2"},
					},
					FirstLevelStruct: secondLevelStruct{},
				},
				{
					FirstLevelString: "value3",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray: []secondLevelStruct{
						{LevelString: "string1"},
						{LevelString: "string2"},
					},
					FirstLevelStruct: secondLevelStruct{},
				},
			}},
			want: []testStruct{
				{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray: []secondLevelStruct{
						{LevelString: "string1"},
					},
					FirstLevelStruct: secondLevelStruct{},
				},
				{
					FirstLevelString: "value3",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray: []secondLevelStruct{
						{LevelString: "string1"},
						{LevelString: "string2"},
					},
					FirstLevelStruct: secondLevelStruct{},
				},
			},
			wantErr: false,
		},
		{
			name: "Minimal filter",
			fields: fields{
				options: defaultOption,
				filter: map[string][]string{
					"FirstLevelString": {"value1"},
				},
			},
			args: args{entries: testStruct{
				FirstLevelString: "value1",
				FirstLevelInt:    0,
				FirstLevelFloat:  0,
				FirstLevelBool:   false,
				FirstLevelArray:  nil,
				FirstLevelStruct: secondLevelStruct{},
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			got, err := f.ApplyFilter(tt.args.entries)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_SetOptions(t *testing.T) {
	type fields struct {
		options *Options
		filter  map[string][]string
	}
	type args struct {
		o *Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Correct Option",
			fields:  fields{},
			args:    args{o: defaultOption},
			wantErr: false,
		},
		{
			name:    "nil option",
			fields:  fields{},
			args:    args{o: nil},
			wantErr: true,
		},
		{
			name:   "negative max depth",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:             -1,
				KeyValueSeparator:    "",
				ValueSeparator:       "",
				KeysSeparator:        "",
				ComposedKeySeparator: "",
			}},
			wantErr: true,
		},
		{
			name:   "empty KeyValueSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:             1,
				KeyValueSeparator:    "",
				ValueSeparator:       ",",
				KeysSeparator:        ":",
				ComposedKeySeparator: ".",
			}},
			wantErr: true,
		},
		{
			name:   "empty ValueSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:             1,
				KeyValueSeparator:    "=",
				ValueSeparator:       "",
				KeysSeparator:        ":",
				ComposedKeySeparator: ".",
			}},
			wantErr: true,
		},
		{
			name:   "empty KeysSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:             1,
				KeyValueSeparator:    "=",
				ValueSeparator:       ",",
				KeysSeparator:        "",
				ComposedKeySeparator: ".",
			}},
			wantErr: true,
		},
		{
			name:   "empty ComposedKeySeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:             1,
				KeyValueSeparator:    "=",
				ValueSeparator:       ",",
				KeysSeparator:        ":",
				ComposedKeySeparator: "",
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			if err := f.SetOptions(tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("SetOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilter_Init(t *testing.T) {
	type fields struct {
		options *Options
		filter  map[string][]string
	}
	type args struct {
		filterValue string
		i           interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "minimal test",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterValue: "FirstLevelString=val2",
				i:           testStruct{},
			},
			wantErr: false,
		},
		{
			name: "minimal test without option set",
			fields: fields{
				options: nil,
				filter:  nil,
			},
			args: args{
				filterValue: "FirstLevelString=val2",
				i:           testStruct{},
			},
			wantErr: false,
		},
		{
			name: "error double key",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterValue: "FirstLevelString=val2:FirstLevelString=val1",
				i:           testStruct{},
			},
			wantErr: true,
		},
		{
			name: "error unknown key",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterValue: "unknownKey=val1",
				i:           testStruct{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			if err := f.Init(tt.args.filterValue, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilter_compileFilter(t *testing.T) {
	type fields struct {
		options *Options
		filter  map[string][]string
	}
	type args struct {
		filterMap map[string][]string
		t         reflect.Type
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFilter map[string][]string
		wantErr    bool
	}{
		{
			name: "minimal test on Name",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"stringFirstLevel": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"FirstLevelString": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "minimal test on Tag",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"stringFirstLevel": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"FirstLevelString": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "Nokey",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"noKey": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{},
			wantErr:    true,
		},
		{
			name: "filter in struct on Name",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"FirstLevelStruct.LevelString": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"FirstLevelStruct.LevelString": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "filter in struct on Tag",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"structFirstLevel.stringSecondLevel": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"FirstLevelStruct.LevelString": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "filter in struct on mixed tag/name",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"FirstLevelStruct.stringSecondLevel": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"FirstLevelStruct.LevelString": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "Filter in array values on Name",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"FirstLevelArray": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"FirstLevelArray": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "Filter in array values on Tag",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"arrayFirstLevel": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"FirstLevelArray": {"val1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			if err := f.compileFilter(tt.args.filterMap, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("compileFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(f.filter, tt.wantFilter) {
				t.Errorf("compileFilter() gotFilterMap = %v, want %v", f.filter, tt.wantFilter)
			}
		})
	}
}

func TestFilter_findValueInComposedKey(t *testing.T) {
	type fields struct {
		options *Options
		filter  map[string][]string
	}
	type args struct {
		filterKey   string
		entryValues reflect.Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []reflect.Value
	}{
		{
			name: "Ok String",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "FirstLevelString",
				entryValues: reflect.ValueOf(testStruct{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf("value1"),
			},
		},
		{
			name: "Ok int",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "FirstLevelInt",
				entryValues: reflect.ValueOf(testStruct{
					FirstLevelString: "value1",
					FirstLevelInt:    2,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(2),
			},
		},
		{
			name: "Ok float",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "FirstLevelFloat",
				entryValues: reflect.ValueOf(testStruct{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  -1.2,
					FirstLevelBool:   false,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(-1.2),
			},
		},
		{
			name: "Ok bool",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "FirstLevelBool",
				entryValues: reflect.ValueOf(testStruct{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   true,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(true),
			},
		},
		{
			name: "Ok Array",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "FirstLevelArray",
				entryValues: reflect.ValueOf(testStruct{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray: []secondLevelStruct{
						{LevelString: "string1"},
						{LevelString: "string2"},
					},
					FirstLevelStruct: secondLevelStruct{},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(secondLevelStruct{"string1"}),
				reflect.ValueOf(secondLevelStruct{"string2"}),
			},
		},
		{
			name: "Ok Struct",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "FirstLevelStruct",
				entryValues: reflect.ValueOf(testStruct{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{"string3"},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(secondLevelStruct{"string3"}),
			},
		},
		{
			name: "noKey", // This case should never occur
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "noKey",
				entryValues: reflect.ValueOf(testStruct{
					FirstLevelString: "value1",
					FirstLevelInt:    0,
					FirstLevelFloat:  0,
					FirstLevelBool:   false,
					FirstLevelArray:  nil,
					FirstLevelStruct: secondLevelStruct{"string3"},
				}),
			},
			want: []reflect.Value{reflect.ValueOf(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			got := f.findValueInComposedKey(tt.args.filterKey, tt.args.entryValues)
			for i := range got {
				if !reflect.DeepEqual(fmt.Sprint(got[i]), fmt.Sprint(tt.want[i])) {
					t.Errorf("findValueInComposedKey() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestFilter_parseFilter(t *testing.T) {
	type fields struct {
		options *Options
		filter  map[string][]string
	}
	type args struct {
		filterValue string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantFilterMap map[string][]string
		wantErr       bool
	}{
		{
			name: "minimal filter",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1=val1"},
			wantFilterMap: map[string][]string{
				"key1": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "filter multi values",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1=val1"},
			wantFilterMap: map[string][]string{
				"key1": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "multi filter single value",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1=val1"},
			wantFilterMap: map[string][]string{
				"key1": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "complex filter",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1=val1,val2,val3:key2=val4:key3=val5,val6"},
			wantFilterMap: map[string][]string{
				"key1": {"val1", "val2", "val3"},
				"key2": {"val4"},
				"key3": {"val5", "val6"},
			},
			wantErr: false,
		},
		{
			name: "Wrong filter: no values",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "key1"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "Wrong filter: multi keys",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "key1=val1:key1=val2"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "Wrong filter: too deep",
			fields: fields{
				options: &Options{
					MaxDepth:             1,
					KeyValueSeparator:    "=",
					ValueSeparator:       ",",
					KeysSeparator:        ":",
					ComposedKeySeparator: ".",
				},
				filter: nil, //always null at parsing time
			},
			args:          args{filterValue: "key1.tooDeep=val1"},
			wantFilterMap: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			gotFilterMap, err := f.parseFilter(tt.args.filterValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFilterMap, tt.wantFilterMap) {
				t.Errorf("parseFilter() gotFilterMap = %v, want %v", gotFilterMap, tt.wantFilterMap)
			}
		})
	}
}

func Test_foundFieldInStruct(t *testing.T) {
	type args struct {
		filterKey string
		t         reflect.Type
	}
	tests := []struct {
		name          string
		args          args
		wantFieldName *reflect.StructField
	}{
		{
			name: "FirstLevelString by name",
			args: args{
				filterKey: "FirstLevelString",
				t:         reflect.TypeOf(testStruct{}),
			},
			wantFieldName: getField(reflect.TypeOf(testStruct{}), 0),
		},
		{
			name: "FirstLevelString by tag",
			args: args{
				filterKey: "stringFirstLevel",
				t:         reflect.TypeOf(testStruct{}),
			},
			wantFieldName: getField(reflect.TypeOf(testStruct{}), 0),
		},
		{
			name: "not found",
			args: args{
				filterKey: "not found",
				t:         reflect.TypeOf(testStruct{}),
			},
			wantFieldName: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFieldName := foundFieldInStruct(tt.args.filterKey, tt.args.t); !reflect.DeepEqual(gotFieldName, tt.wantFieldName) {
				t.Errorf("foundFieldInStruct() = %v, want %v", gotFieldName, tt.wantFieldName)
			}
		})
	}
}

// Impossible to get the address of a struct field, but only of var.
// This function is mandatory for test
func getField(t reflect.Type, i int) *reflect.StructField {
	field := t.Field(i)
	return &field
}

type secondLevelStruct struct {
	LevelString string `json:"stringSecondLevel"`
}

type testStruct struct {
	FirstLevelString string              `json:"stringFirstLevel"`
	FirstLevelInt    int                 `json:"intFirstLevel"`
	FirstLevelFloat  float32             `json:"floatFirstLevel"`
	FirstLevelBool   bool                `json:"boolFirstLevel"`
	FirstLevelArray  []secondLevelStruct `json:"arrayFirstLevel"`
	FirstLevelStruct secondLevelStruct   `json:"structFirstLevel"`
	//TODO map
}
