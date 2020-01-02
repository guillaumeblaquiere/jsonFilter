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
					"RootString": {"value1"},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
					RootInt:    0,
					RootFloat:  0,
					RootBool:   false,
					RootArray:  nil,
					RootStruct: SubStruct{},
				},
			}},
			want: []testStruct{
				{
					RootString: "value1",
					RootInt:    0,
					RootFloat:  0,
					RootBool:   false,
					RootArray:  nil,
					RootStruct: SubStruct{},
				},
			},
			wantErr: false,
		},
		{
			name: "Composite filter",
			fields: fields{
				options: defaultOption,
				filter: map[string][]string{
					"RootArray.SubString": {"string1"},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
					RootInt:    0,
					RootFloat:  0,
					RootBool:   false,
					RootArray: []SubStruct{
						{SubString: "string1"},
					},
					RootStruct: SubStruct{},
				},
				{
					RootString: "value2",
					RootInt:    0,
					RootFloat:  0,
					RootBool:   false,
					RootArray: []SubStruct{
						{SubString: "string2"},
					},
					RootStruct: SubStruct{},
				},
				{
					RootString: "value3",
					RootInt:    0,
					RootFloat:  0,
					RootBool:   false,
					RootArray: []SubStruct{
						{SubString: "string1"},
						{SubString: "string2"},
					},
					RootStruct: SubStruct{},
				},
			}},
			want: []testStruct{
				{
					RootString: "value1",
					RootInt:    0,
					RootFloat:  0,
					RootBool:   false,
					RootArray: []SubStruct{
						{SubString: "string1"},
					},
					RootStruct: SubStruct{},
				},
				{
					RootString: "value3",
					RootInt:    0,
					RootFloat:  0,
					RootBool:   false,
					RootArray: []SubStruct{
						{SubString: "string1"},
						{SubString: "string2"},
					},
					RootStruct: SubStruct{},
				},
			},
			wantErr: false,
		},
		{
			name: "Minimal filter",
			fields: fields{
				options: defaultOption,
				filter: map[string][]string{
					"RootString": {"value1"},
				},
			},
			args: args{entries: testStruct{
				RootString: "value1",
				RootInt:    0,
				RootFloat:  0,
				RootBool:   false,
				RootArray:  nil,
				RootStruct: SubStruct{},
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
				filterValue: "RootString=val2",
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
				filterValue: "RootString=val2",
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
				filterValue: "RootString=val2:RootString=val1",
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
					"stringRoot": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootString": {"val1"},
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
					"stringRoot": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootString": {"val1"},
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
					"RootStruct.SubString": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootStruct.SubString": {"val1"},
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
					"structRoot.stringSub": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootStruct.SubString": {"val1"},
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
					"RootStruct.stringSub": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootStruct.SubString": {"val1"},
			},
			wantErr: false,
		},
		{
			name: "filter on Map with mixed tag/name",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: map[string][]string{
					"RootMap.entry1.stringSub": {"valMap"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootMap.entry1.SubString": {"valMap"},
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
					"RootArray": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootArray": {"val1"},
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
					"arrayRoot": {"val1"},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: map[string][]string{
				"RootArray": {"val1"},
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
				filterKey: "RootString",
				entryValues: reflect.ValueOf(testStruct{
					RootString: "value1",
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
				filterKey: "RootInt",
				entryValues: reflect.ValueOf(testStruct{
					RootInt: 2,
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
				filterKey: "RootFloat",
				entryValues: reflect.ValueOf(testStruct{
					RootFloat: -1.2,
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
				filterKey: "RootBool",
				entryValues: reflect.ValueOf(testStruct{
					RootBool: true,
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
				filterKey: "RootArray",
				entryValues: reflect.ValueOf(testStruct{
					RootArray: []SubStruct{
						{SubString: "string1"},
						{SubString: "string2"},
					},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(SubStruct{"string1"}),
				reflect.ValueOf(SubStruct{"string2"}),
			},
		},
		{
			name: "Ok Struct",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootStruct",
				entryValues: reflect.ValueOf(testStruct{
					RootStruct: SubStruct{"string3"},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(SubStruct{"string3"}),
			},
		},
		{
			name: "Ok Ptr Struct",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootPtrStruct",
				entryValues: reflect.ValueOf(testStruct{
					RootPtrStruct: &SubStruct{"string3"},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(&SubStruct{"string3"}),
			},
		},
		{
			name: "Ok nil Ptr",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootPtrStruct.SubString",
				entryValues: reflect.ValueOf(testStruct{
					RootPtrStruct: nil,
				}),
			},
			want: []reflect.Value{reflect.ValueOf(testStruct{}.RootPtrStruct)},
		},
		{
			name: "Ok Ptr Struct evaluated",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootPtrStruct.SubString",
				entryValues: reflect.ValueOf(testStruct{
					RootPtrStruct: &SubStruct{"string3"},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf("string3"),
			},
		},
		{
			name: "Ok Map",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootMap",
				entryValues: reflect.ValueOf(testStruct{
					RootMap: map[string]SubStruct{"entry1": {"string3"}},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(map[string]SubStruct{"entry1": {"string3"}}),
			},
		},
		{
			name: "Ok Map with entry",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootMap.entry1",
				entryValues: reflect.ValueOf(testStruct{
					RootMap: map[string]SubStruct{"entry1": {"string3"}},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(SubStruct{"string3"}),
			},
		},
		{
			name: "noKey", // This case should never occur
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey:   "noKey",
				entryValues: reflect.ValueOf(testStruct{}),
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
			name: "RootString by name",
			args: args{
				filterKey: "RootString",
				t:         reflect.TypeOf(testStruct{}),
			},
			wantFieldName: getField(reflect.TypeOf(testStruct{}), 0),
		},
		{
			name: "RootString by tag",
			args: args{
				filterKey: "stringRoot",
				t:         reflect.TypeOf(testStruct{}),
			},
			wantFieldName: getField(reflect.TypeOf(testStruct{}), 0),
		},
		{
			name: "Ptr get elem",
			args: args{
				filterKey: "ptrStructRoot",
				t:         reflect.TypeOf(&testStruct{}),
			},
			wantFieldName: getField(reflect.TypeOf(testStruct{}), 6),
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

type SubStruct struct {
	SubString string `json:"stringSub,omitempty"`
}

type testStruct struct {
	RootString    string               `json:"stringRoot,omitempty"`
	RootInt       int                  `json:"intRoot,omitempty"`
	RootFloat     float32              `json:"floatRoot,omitempty"`
	RootBool      bool                 `json:"boolRoot,omitempty"`
	RootArray     []SubStruct          `json:"arrayRoot,omitempty"`
	RootStruct    SubStruct            `json:"structRoot,omitempty"`
	RootPtrStruct *SubStruct           `json:"ptrStructRoot,omitempty"`
	RootMap       map[string]SubStruct `json:"mapRoot,omitempty"`
}
