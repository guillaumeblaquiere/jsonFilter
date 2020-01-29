package jsonFilter

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFilter_ApplyFilter(t *testing.T) {
	type fields struct {
		options *Options
		filter  []kov
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
			name: "Minimal filter equals",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootString",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"value1"},
					},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
				},
			}},
			want: []testStruct{
				{
					RootString: "value1",
				},
			},
			wantErr: false,
		},
		{
			name: "Minimal filter not equals",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootString",
						Operator: defaultOption.NotEqualKeyValueSeparator,
						Values:   []string{"value2"},
					},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
				},
				{
					RootString: "value2",
				},
			}},
			want: []testStruct{
				{
					RootString: "value1",
				},
			},
			wantErr: false,
		},
		{
			name: "Filter not equals with 2 values",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootString",
						Operator: defaultOption.NotEqualKeyValueSeparator,
						Values:   []string{"value2", "value3"},
					},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
				},
				{
					RootString: "value2",
				},
			}},
			want: []testStruct{
				{
					RootString: "value1",
				},
			},
			wantErr: false,
		},
		{
			name: "Minimal filter lower than",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootInt",
						Operator: defaultOption.LowerThanKeyValueSeparator,
						Values:   []string{"11"},
					},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
					RootInt:    10,
				},
				{
					RootString: "value2",
					RootInt:    11,
				},
			}},
			want: []testStruct{
				{
					RootString: "value1",
					RootInt:    10,
				},
			},
			wantErr: false,
		},
		{
			name: "Minimal filter greater than",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootInt",
						Operator: defaultOption.GreaterThanKeyValueSeparator,
						Values:   []string{"10"},
					},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
					RootInt:    10,
				},
				{
					RootString: "value2",
					RootInt:    11,
				},
			}},
			want: []testStruct{
				{
					RootString: "value2",
					RootInt:    11,
				},
			},
			wantErr: false,
		},
		{
			name: "Minimal filter greater than float",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootFloat",
						Operator: defaultOption.GreaterThanKeyValueSeparator,
						Values:   []string{"10.5"},
					},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
					RootFloat:  10.4,
				},
				{
					RootString: "value2",
					RootFloat:  10.6,
				},
			}},
			want: []testStruct{
				{
					RootString: "value2",
					RootFloat:  10.6,
				},
			},
			wantErr: false,
		},
		{
			name: "Minimal filter greater than negative",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootInt",
						Operator: defaultOption.GreaterThanKeyValueSeparator,
						Values:   []string{"-11"},
					},
				},
			},
			args: args{entries: []testStruct{
				{
					RootString: "value1",
					RootInt:    -10,
				},
				{
					RootString: "value2",
					RootInt:    -11,
				},
			}},
			want: []testStruct{
				{
					RootString: "value1",
					RootInt:    -10,
				},
			},
			wantErr: false,
		},
		{
			name: "Composite filter",
			fields: fields{
				options: defaultOption,
				filter: []kov{
					{
						Key:      "RootArray.SubString",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"string1"},
					},
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
				filter: []kov{
					{
						Key:      "RootString",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"value1"},
					},
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
		filter  []kov
	}
	type args struct {
		o *Options
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOption *Options
	}{
		{
			name:       "Correct Option",
			fields:     fields{},
			args:       args{o: defaultOption},
			wantOption: defaultOption,
		},
		{
			name:       "nil option",
			fields:     fields{},
			args:       args{o: nil},
			wantOption: defaultOption,
		},
		{
			name:   "negative max depth",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     -1,
				EqualKeyValueSeparator:       "",
				NotEqualKeyValueSeparator:    "",
				LowerThanKeyValueSeparator:   "",
				GreaterThanKeyValueSeparator: "",
				ValueSeparator:               "",
				KeysSeparator:                "",
				ComposedKeySeparator:         "",
			}},
			wantOption: defaultOption,
		},
		{
			name:   "empty EqualKeyValueSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     0,
				EqualKeyValueSeparator:       "",
				NotEqualKeyValueSeparator:    "!=",
				LowerThanKeyValueSeparator:   "<",
				GreaterThanKeyValueSeparator: ">",
				ValueSeparator:               ",",
				KeysSeparator:                ":",
				ComposedKeySeparator:         ".",
			}},
			wantOption: defaultOption,
		},
		{
			name:   "empty NotEqualKeyValueSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     0,
				EqualKeyValueSeparator:       "=",
				NotEqualKeyValueSeparator:    "",
				LowerThanKeyValueSeparator:   "<",
				GreaterThanKeyValueSeparator: ">",
				ValueSeparator:               ",",
				KeysSeparator:                ":",
				ComposedKeySeparator:         ".",
			}},
			wantOption: defaultOption,
		},
		{
			name:   "empty LowerThanKeyValueSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     0,
				EqualKeyValueSeparator:       "=",
				NotEqualKeyValueSeparator:    "!=",
				LowerThanKeyValueSeparator:   "",
				GreaterThanKeyValueSeparator: ">",
				ValueSeparator:               ",",
				KeysSeparator:                ":",
				ComposedKeySeparator:         ".",
			}},
			wantOption: defaultOption,
		},
		{
			name:   "empty GreaterThanKeyValueSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     0,
				EqualKeyValueSeparator:       "=",
				NotEqualKeyValueSeparator:    "!=",
				LowerThanKeyValueSeparator:   "<",
				GreaterThanKeyValueSeparator: "",
				ValueSeparator:               ",",
				KeysSeparator:                ":",
				ComposedKeySeparator:         ".",
			}},
			wantOption: defaultOption,
		},
		{
			name:   "empty ValueSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     0,
				EqualKeyValueSeparator:       "=",
				NotEqualKeyValueSeparator:    "!=",
				LowerThanKeyValueSeparator:   "<",
				GreaterThanKeyValueSeparator: ">",
				ValueSeparator:               "",
				KeysSeparator:                ":",
				ComposedKeySeparator:         ".",
			}},
			wantOption: defaultOption,
		},
		{
			name:   "empty KeysSeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     0,
				EqualKeyValueSeparator:       "=",
				NotEqualKeyValueSeparator:    "!=",
				LowerThanKeyValueSeparator:   "<",
				GreaterThanKeyValueSeparator: ">",
				ValueSeparator:               ",",
				KeysSeparator:                "",
				ComposedKeySeparator:         ".",
			}},
			wantOption: defaultOption,
		},
		{
			name:   "empty ComposedKeySeparator",
			fields: fields{},
			args: args{o: &Options{
				MaxDepth:                     0,
				EqualKeyValueSeparator:       "=",
				NotEqualKeyValueSeparator:    "!=",
				LowerThanKeyValueSeparator:   "<",
				GreaterThanKeyValueSeparator: ">",
				ValueSeparator:               ",",
				KeysSeparator:                ":",
				ComposedKeySeparator:         "",
			}},
			wantOption: defaultOption,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			f.SetOptions(tt.args.o)
			if !reflect.DeepEqual(f.options, tt.wantOption) {
				t.Errorf("SetOptions() got = %v, wanted %v", f.options, tt.wantOption)
			}
		})
	}
}

func TestFilter_Init(t *testing.T) {
	type fields struct {
		options *Options
		filter  []kov
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
		filter  []kov
	}
	type args struct {
		filterMap []kov
		t         reflect.Type
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFilter []kov
		wantErr    bool
	}{
		{
			name: "minimal test on Name",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: []kov{
					{
						Key:      "stringRoot",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
				filterMap: []kov{
					{
						Key:      "stringRoot",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
				filterMap: []kov{
					{
						Key:      "noKey",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{},
			wantErr:    true,
		},
		{
			name: "filter in struct on Name",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: []kov{
					{
						Key:      "RootStruct.SubString",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootStruct.SubString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
				filterMap: []kov{
					{
						Key:      "structRoot.stringSub",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootStruct.SubString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
				filterMap: []kov{
					{
						Key:      "RootStruct.stringSub",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootStruct.SubString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
				filterMap: []kov{
					{
						Key:      "RootMap.entry1.stringSub",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootMap.entry1.SubString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
			},
			wantErr: false,
		},
		{
			name: "Filter in array of ptr",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: []kov{
					{
						Key:      "RootArrayPtr.RootString",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootArrayPtr.RootString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
			},
			wantErr: false,
		},
		{
			name: "Filter in map of ptr",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: []kov{
					{
						Key:      "RootMapPtr.RootString",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootMapPtr.RootString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
				filterMap: []kov{
					{
						Key:      "RootArray",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootArray",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
			},
			wantErr: false,
		},
		{
			name: "Filter in array values on Name, deeper",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: []kov{
					{
						Key:      "RootArray.SubString",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootArray.SubString",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
			},
			wantErr: false,
		},
		{
			name: "Filter on Matrix",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterMap: []kov{
					{
						Key:      "matrix",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "Matrix",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
				filterMap: []kov{
					{
						Key:      "arrayRoot",
						Operator: defaultOption.EqualKeyValueSeparator,
						Values:   []string{"val1"},
					},
				},
				t: reflect.TypeOf(testStruct{}),
			},
			wantFilter: []kov{
				{
					Key:      "RootArray",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
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
		filter  []kov
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
			name: "Ok Array of struct",
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
			name: "Ok Array of Ptr struct",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootArrayPtr",
				entryValues: reflect.ValueOf(testStruct{
					RootArrayPtr: []*testStruct{
						{RootString: "string1"},
						{RootString: "string2"},
					},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(testStruct{RootString: "string1"}),
				reflect.ValueOf(testStruct{RootString: "string2"}),
			},
		},
		{
			name: "Ok Array of Ptr struct, deeper",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootArrayPtr.RootString",
				entryValues: reflect.ValueOf(testStruct{
					RootArrayPtr: []*testStruct{
						{RootString: "string1"},
						{RootString: "string2"},
						nil,
					},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf("string1"),
				reflect.ValueOf("string2"),
			},
		},
		{
			name: "Ok Array of string",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootArray",
				entryValues: reflect.ValueOf(testStruct{
					RootArraySimple: []string{
						"string1",
						"string2",
					},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf([]string{
					"string1",
					"string2",
				}),
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
					RootPtrStruct: &testStruct{RootString: "string3"},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(testStruct{RootString: "string3"}),
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
				filterKey: "RootPtrStruct.RootString",
				entryValues: reflect.ValueOf(testStruct{
					RootPtrStruct: &testStruct{RootString: "string3"},
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
			name: "Ok Map ptr",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootPtrStruct.RootPtrStruct",
				entryValues: reflect.ValueOf(testStruct{
					RootPtrStruct: &testStruct{RootPtrStruct: &testStruct{RootString: "val1"}},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(testStruct{RootString: "val1"}),
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
			name: "Ok Map with incorrect entry",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootMap.entry",
				entryValues: reflect.ValueOf(testStruct{
					RootMap: map[string]SubStruct{"entry1": {"string3"}},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(SubStruct{"string3"}),
			},
		},
		{
			name: "Ok Map with Array of String entry",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootMap.entry1",
				entryValues: reflect.ValueOf(testStruct{
					RootMapArrayOfSimple: map[string][]string{"entry1": {"string3"}},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf([]string{"string3"}),
			},
		},
		{
			name: "Ok Map simple string with entry",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootMap.entry1",
				entryValues: reflect.ValueOf(testStruct{
					RootMapSimple: map[string]string{"entry1": "string3"},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf("string3"),
			},
		},
		{
			name: "Ok Map with Array of Struct entry",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootMap.entry1",
				entryValues: reflect.ValueOf(testStruct{
					RootMapArrayOfStruct: map[string][]SubStruct{
						"entry1": {
							SubStruct{"string3"},
							SubStruct{"string4"},
						},
						"entry2": {
							SubStruct{"string3"},
							SubStruct{"string4"},
						},
					},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf([]SubStruct{{"string3"}, {"string4"}}),
			},
		},
		{
			name: "Ok Map with Array of Struct entry, deeper",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "RootMap.entry1.SubString",
				entryValues: reflect.ValueOf(testStruct{
					RootMapArrayOfStruct: map[string][]SubStruct{
						"entry1": {
							SubStruct{"string3"},
							SubStruct{"string4"},
						},
						"entry2": {
							SubStruct{"string3"},
							SubStruct{"string4"},
						},
					},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf(SubStruct{"string3"}),
				reflect.ValueOf(SubStruct{"string4"}),
			},
		},
		{
			name: "Ok Matrix",
			fields: fields{
				options: defaultOption,
				filter:  nil,
			},
			args: args{
				filterKey: "Matrix",
				entryValues: reflect.ValueOf(testStruct{
					Matrix: [][]string{
						{"AA", "AB", "AC"},
						{"BA", "BB", "BC"},
						{"CA", "CB", "CC"},
					},
				}),
			},
			want: []reflect.Value{
				reflect.ValueOf("AA"),
				reflect.ValueOf("AB"),
				reflect.ValueOf("AC"),
				reflect.ValueOf("BA"),
				reflect.ValueOf("BB"),
				reflect.ValueOf("BC"),
				reflect.ValueOf("CA"),
				reflect.ValueOf("CB"),
				reflect.ValueOf("CC"),
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
		filter  []kov
	}
	type args struct {
		filterValue string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantFilterMap []kov
		wantErr       bool
	}{
		{
			name: "minimal filter",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1" + defaultOption.EqualKeyValueSeparator + "val1"},
			wantFilterMap: []kov{
				{
					Key:      "key1",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
			},
			wantErr: false,
		},
		{
			name: "filter multi values",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1" + defaultOption.EqualKeyValueSeparator + "val1"},
			wantFilterMap: []kov{
				{
					Key:      "key1",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
			},
			wantErr: false,
		},
		{
			name: "multi filter single value",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1" + defaultOption.EqualKeyValueSeparator + "val1"},
			wantFilterMap: []kov{
				{
					Key:      "key1",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
			},
			wantErr: false,
		},
		{
			name: "multi filter same key",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1" + defaultOption.EqualKeyValueSeparator + "val1" +
				defaultOption.KeysSeparator + "key1" + defaultOption.NotEqualKeyValueSeparator + "val2"},
			wantFilterMap: []kov{
				{
					Key:      "key1",
					Operator: defaultOption.EqualKeyValueSeparator,
					Values:   []string{"val1"},
				},
				{
					Key:      "key1",
					Operator: defaultOption.NotEqualKeyValueSeparator,
					Values:   []string{"val2"},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong multi filter same key same op",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1" + defaultOption.EqualKeyValueSeparator + "val1" +
				defaultOption.KeysSeparator + "key1" + defaultOption.EqualKeyValueSeparator + "val2"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "complex filter",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args: args{filterValue: "key1" + defaultOption.NotEqualKeyValueSeparator + "val1,val2,val3:" +
				"key2" + defaultOption.LowerThanKeyValueSeparator + "4.5:" +
				"key3" + defaultOption.GreaterThanKeyValueSeparator + "-5"},
			wantFilterMap: []kov{
				{
					Key:      "key1",
					Operator: defaultOption.NotEqualKeyValueSeparator,
					Values:   []string{"val1", "val2", "val3"},
				},
				{
					Key:      "key2",
					Operator: defaultOption.LowerThanKeyValueSeparator,
					Values:   []string{"4.5"},
				},
				{
					Key:      "key3",
					Operator: defaultOption.GreaterThanKeyValueSeparator,
					Values:   []string{"-5"},
				},
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
			name: "Wrong filter: no key",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "=val1"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "Wrong filter: no key, no value",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "="},
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
					MaxDepth:               1,
					EqualKeyValueSeparator: "=",
					ValueSeparator:         ",",
					KeysSeparator:          ":",
					ComposedKeySeparator:   ".",
				},
				filter: nil, //always null at parsing time
			},
			args:          args{filterValue: "key1.tooDeep=val1"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "Wrong filter: GT multiple values",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "key1" + defaultOption.GreaterThanKeyValueSeparator + "1,2"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "Wrong filter: LT multiple values",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "key1" + defaultOption.LowerThanKeyValueSeparator + "1,2"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "Wrong filter: GT no numeric values",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "k1" + defaultOption.GreaterThanKeyValueSeparator + "v1"},
			wantFilterMap: nil,
			wantErr:       true,
		},
		{
			name: "Wrong filter: LT no numeric values",
			fields: fields{
				options: defaultOption,
				filter:  nil, //always null at parsing time
			},
			args:          args{filterValue: "k1" + defaultOption.LowerThanKeyValueSeparator + "v2"},
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
			wantFieldName: getField(reflect.TypeOf(testStruct{}), 7),
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
	RootString           string                 `json:"stringRoot,omitempty"`
	RootInt              int                    `json:"intRoot,omitempty"`
	RootFloat            float32                `json:"floatRoot,omitempty"`
	RootBool             bool                   `json:"boolRoot,omitempty"`
	RootArray            []SubStruct            `json:"arrayRoot,omitempty"`
	RootArraySimple      []string               `json:"arrayRootSimple,omitempty"`
	RootStruct           SubStruct              `json:"structRoot,omitempty"`
	RootPtrStruct        *testStruct            `json:"ptrStructRoot,omitempty"`
	RootMap              map[string]SubStruct   `json:"mapRoot,omitempty"`
	RootMapSimple        map[string]string      `json:"mapRootString,omitempty"`
	RootMapArrayOfSimple map[string][]string    `json:"mapRootArrayOfString,omitempty"`
	RootMapArrayOfStruct map[string][]SubStruct `json:"mapRootArrayOfString,omitempty"`
	RootArrayPtr         []*testStruct          `json:"arrayRootPtr,omitempty"`
	RootMapPtr           map[string]*testStruct `json:"mapRootPtr,omitempty"`
	Matrix               [][]string             `json:"matrix,omitempty"`
}

func TestFilter_getFilterAndValue(t *testing.T) {
	type fields struct {
		options *Options
		filter  []kov
	}
	type args struct {
		filter string
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantFilterKeyValues []string
		wantOpFound         string
	}{
		{
			name: "EqualFound",
			fields: fields{
				options: defaultOption,
			},
			args:                args{"k1" + defaultOption.EqualKeyValueSeparator + "v1,v2"},
			wantFilterKeyValues: []string{"k1", "v1,v2"},
			wantOpFound:         defaultOption.EqualKeyValueSeparator,
		},
		{
			name: "NotEqualFound",
			fields: fields{
				options: defaultOption,
			},
			args:                args{"k1" + defaultOption.NotEqualKeyValueSeparator + "v1,v2"},
			wantFilterKeyValues: []string{"k1", "v1,v2"},
			wantOpFound:         defaultOption.NotEqualKeyValueSeparator,
		},
		{
			name: "GreaterThanFound",
			fields: fields{
				options: defaultOption,
			},
			args:                args{"k1" + defaultOption.GreaterThanKeyValueSeparator + "v1"},
			wantFilterKeyValues: []string{"k1", "v1"},
			wantOpFound:         defaultOption.GreaterThanKeyValueSeparator,
		},
		{
			name: "LowerThanFound",
			fields: fields{
				options: defaultOption,
			},
			args:                args{"k1" + defaultOption.LowerThanKeyValueSeparator + "v1"},
			wantFilterKeyValues: []string{"k1", "v1"},
			wantOpFound:         defaultOption.LowerThanKeyValueSeparator,
		},
		{
			name: "NotFound",
			fields: fields{
				options: defaultOption,
			},
			args:                args{"k1,v1,v2"},
			wantFilterKeyValues: []string{},
			wantOpFound:         "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
				options: tt.fields.options,
				filter:  tt.fields.filter,
			}
			gotFilterKeyValues, gotOpFound := f.getFilterAndValue(tt.args.filter)
			if fmt.Sprint(gotFilterKeyValues) != fmt.Sprint(tt.wantFilterKeyValues) { //issue with deepequal on NotFound
				t.Errorf("getFilterAndValue() gotFilterKeyValues = %v, want %v", gotFilterKeyValues, tt.wantFilterKeyValues)
			}
			if gotOpFound != tt.wantOpFound {
				t.Errorf("getFilterAndValue() gotOpFound = %v, want %v", gotOpFound, tt.wantOpFound)
			}
		})
	}
}
