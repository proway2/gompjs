package gompjs

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

var defaultLoader UnmarshalFunc = json.Unmarshal

type args struct {
	inputStr      string
	unicodeEscape bool
	loader        UnmarshalFunc
}

type tests []struct {
	name    string
	args    args
	want    any
	wantErr bool
	skip    bool
}

var objectTests = tests{
	{
		name:    "Single keyword",
		args:    args{inputStr: "{'hello': 'world'}"},
		want:    map[string]any{"hello": "world"},
		wantErr: false,
	},
	{
		name:    "Two keywords",
		args:    args{inputStr: "{'hello': 'world', 'my': 'master'}"},
		want:    map[string]any{"hello": "world", "my": "master"},
		wantErr: false,
	},
	{
		name:    "Nested JSON",
		args:    args{inputStr: "{'hello': 'world', 'my': {'master': 'of Orion'}, 'test': 'xx'}"},
		want:    map[string]any{"hello": "world", "my": map[string]any{"master": "of Orion"}, "test": "xx"},
		wantErr: false,
	},
	{
		name:    "Empty JSON",
		args:    args{inputStr: "{}"},
		want:    map[string]any{},
		wantErr: false,
	},
	{
		name:    "JSON with a number",
		args:    args{inputStr: "{\"num\": 126}"},
		want:    map[string]any{"num": float64(126)},
		wantErr: false,
	},
}

var listTests = tests{
	{
		name:    "Empty list",
		args:    args{inputStr: "[]"},
		want:    []any{},
		wantErr: false,
	},
	{
		name:    "Empty nested list - 3rd level",
		args:    args{inputStr: "[[[]]]"},
		want:    []any{[]any{[]any{}}},
		wantErr: false,
	},
	{
		name:    "Nested list - 3rd level with a number",
		args:    args{inputStr: "[[[1]]]"},
		want:    []any{[]any{[]any{float64(1)}}},
		wantErr: false,
	},
	{
		name:    "Just a list with a number",
		args:    args{inputStr: "[1]"},
		want:    []any{float64(1)},
		wantErr: false,
	},
	{
		name:    "List with several numbers",
		args:    args{inputStr: "[1, 2, 3, 4]"},
		want:    []any{float64(1), float64(2), float64(3), float64(4)},
		wantErr: false,
	},
	{
		name:    "List with several symbols",
		args:    args{inputStr: "['h', 'e', 'l', 'l', 'o']"},
		want:    []any{"h", "e", "l", "l", "o"},
		wantErr: false,
	},
	{
		name:    "Empty extremely nested list",
		args:    args{inputStr: "[[[[[[[[[[[[[[[1]]]]]]]]]]]]]]]"},
		want:    []any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{[]any{float64(1)}}}}}}}}}}}}}}},
		wantErr: false,
	},
}

var mixedTests = tests{
	{
		name:    "Empty list + single element list",
		args:    args{inputStr: "{'hello': [], 'world': [0]}"},
		want:    map[string]any{"hello": []any{}, "world": []any{float64(0)}},
		wantErr: false,
	},
	{
		name:    "List with numbers",
		args:    args{inputStr: "{'hello': [1, 2, 3, 4]}"},
		want:    map[string]any{"hello": []any{float64(1), float64(2), float64(3), float64(4)}},
		wantErr: false,
	},
	{
		name:    "Lists with nested objects",
		args:    args{inputStr: "[{'a':12}, {'b':33}]"},
		want:    []any{map[string]any{"a": float64(12)}, map[string]any{"b": float64(33)}},
		wantErr: false,
	},
	{
		name:    "Boolean as keys and values",
		args:    args{inputStr: "[false, {'true': true, `pies`: \"kot\"}, false,]"},
		want:    []any{false, map[string]any{"true": true, "pies": "kot"}, false},
		wantErr: false,
	},
	{
		name: "Object with number",
		args: args{inputStr: "{a:1,b:1,c:1,d:1,e:1,f:1,g:1,h:1,i:1,j:1}"},
		want: map[string]any{
			"a": float64(1),
			"b": float64(1),
			"c": float64(1),
			"d": float64(1),
			"e": float64(1),
			"f": float64(1),
			"g": float64(1),
			"h": float64(1),
			"i": float64(1),
			"j": float64(1),
		},
		wantErr: false,
	},
	{
		name: "Mixed and nested object",
		args: args{inputStr: "{'a':[{'b':1},{'c':[{'d':{'f':{'g':[1,2]}}},{'e':1}]}]}"},
		want: map[string]any{
			"a": []any{
				map[string]any{"b": float64(1)},
				map[string]any{"c": []any{
					map[string]any{"d": map[string]any{"f": map[string]any{
						"g": []any{float64(1), float64(2)}}}},
					map[string]any{"e": float64(1)},
				},
				},
			},
		},
		wantErr: false,
	},
}

var standardValues = tests{
	{
		name:    "Basic - numbers",
		args:    args{inputStr: "{'hello': 12, 'world': 10002.21}"},
		want:    map[string]any{"hello": float64(12), "world": 10002.21},
		wantErr: false,
	},
	{
		name:    "List with numbers",
		args:    args{inputStr: "[12, -323, 0.32, -32.22, .2, - 4]"},
		want:    []any{float64(12), float64(-323), 0.32, -32.22, 0.2, float64(-4)},
		wantErr: false,
	},
	{
		name:    "Object with negative numbers",
		args:    args{inputStr: "{\"a\": -12, \"b\": - 5}"},
		want:    map[string]any{"a": float64(-12), "b": float64(-5)},
		wantErr: false,
	},
	{
		name:    "Object with booleans and null",
		args:    args{inputStr: "{'a': true, 'b': false, 'c': null}"},
		want:    map[string]any{"a": true, "b": false, "c": nil},
		wantErr: false,
	},
	{
		name:    "List with a single Unicode code",
		args:    args{inputStr: "[\"\\uD834\\uDD1E\"]"},
		want:    []any{"𝄞"},
		wantErr: false,
	},
	{
		name: "Object with strange string",
		args: args{inputStr: "{'a': '123\\'456\\n'}"},
		want: map[string]any{"a": "123'456\n"},
	},
	{
		name:    "List with a single Unicode code 2",
		args:    args{inputStr: "['\u00E9']"},
		want:    []any{"é"},
		wantErr: false,
	},
	{
		name:    "Object with nested object with a key as Unicode code",
		args:    args{inputStr: "{\"cache\":{\"\u002Ftest\u002F\": 0}}"},
		want:    map[string]any{"cache": map[string]any{"/test/": float64(0)}},
		wantErr: false,
	},
	{
		name:    "Floating-point exponential value",
		args:    args{inputStr: "{\"a\": 3.125e7}"},
		want:    map[string]any{"a": 3.125e7},
		wantErr: false,
	},
	{
		name:    "Object with many quotes",
		args:    args{inputStr: `{"a": "b'"}`},
		want:    map[string]any{"a": "b'"},
		wantErr: false,
	},
	{
		name:    "Object with negative numbers as values",
		args:    args{inputStr: `{"a": .99, "b": -.1}`},
		want:    map[string]any{"a": 0.99, "b": -0.1},
		wantErr: false,
	},
	{
		name:    "List with ellipsises",
		args:    args{inputStr: `["/* ... */", "// ..."]`},
		want:    []any{"/* ... */", "// ..."},
		wantErr: false,
	},
	{
		name:    "Object with a list",
		args:    args{inputStr: `{"inclusions":["/*","/"]}`},
		want:    map[string]any{"inclusions": []any{"/*", "/"}},
		wantErr: false,
	},
}

var strangeValues = tests{
	{
		name:    "Object with numbers",
		args:    args{inputStr: "{abc: 100, dev: 200}"},
		want:    map[string]any{"abc": float64(100), "dev": float64(200)},
		wantErr: false,
	},
	{
		name:    "Long key",
		args:    args{inputStr: "{abcdefghijklmnopqrstuvwxyz: 12}"},
		want:    map[string]any{"abcdefghijklmnopqrstuvwxyz": float64(12)},
		wantErr: false,
	},
	{
		name:    "Function as a value",
		args:    args{inputStr: "{age: function(yearBorn,thisYear) {return thisYear - yearBorn;}}"},
		want:    map[string]any{"age": "function(yearBorn,thisYear) {return thisYear - yearBorn;}"},
		wantErr: false,
	},
	{
		name:    "Function as a value with many brackets",
		args:    args{inputStr: "{\"abc\": function() {return '])))))))))))))))';}}"},
		want:    map[string]any{"abc": "function() {return '])))))))))))))))';}"},
		wantErr: false,
	},
	{
		name:    "Undefined value",
		args:    args{inputStr: `{"a": undefined}`},
		want:    map[string]any{"a": "undefined"},
		wantErr: false,
	},
	{
		name:    "List with undefined values",
		args:    args{inputStr: "[undefined, undefined]"},
		want:    []any{"undefined", "undefined"},
		wantErr: false,
	},
	{
		name:    "Underscore and dollar sign in keys",
		args:    args{inputStr: "{_a: 1, $b: 2}"},
		want:    map[string]any{"_a": float64(1), "$b": float64(2)},
		wantErr: false,
	},
	{
		name:    "Regex as a value",
		args:    args{inputStr: "{regex: /a[^d]{1,12}/i}"},
		want:    map[string]any{"regex": "/a[^d]{1,12}/i"},
		wantErr: false,
	},
	{
		name:    "Function as a value with an escape symbol",
		args:    args{inputStr: "{'a': function(){return '\"'}}"},
		want:    map[string]any{"a": `function(){return '"'}`},
		wantErr: false,
	},
	{
		name: "Numbers as keys",
		args: args{inputStr: "{1: 1, 2: 2, 3: 3, 4: 4}"},
		want: map[string]any{"1": float64(1), "2": float64(2), "3": float64(3), "4": float64(4)},
	},
	{
		name:    "Incomplete floating point numbers as values",
		args:    args{inputStr: "{'a': 121.}"},
		want:    map[string]any{"a": 121.0},
		wantErr: false,
	},
	{
		name:    "No quotes around key",
		args:    args{inputStr: "{abc : 100}"},
		want:    map[string]any{"abc": float64(100)},
		wantErr: false,
	},
	{
		name:    "Spaces between a key and a value",
		args:    args{inputStr: "{abc     :       100}"},
		want:    map[string]any{"abc": float64(100)},
		wantErr: false,
	},
	{
		name:    "Value is unquotted string with space",
		args:    args{inputStr: "{abc: name }"},
		want:    map[string]any{"abc": "name"},
		wantErr: false,
	},
	{
		name:    "Value with \\t",
		args:    args{inputStr: "{abc: name\t}"},
		want:    map[string]any{"abc": "name"},
		wantErr: false,
	},
	{
		name:    "Value with \\n",
		args:    args{inputStr: "{abc: value\n}"},
		want:    map[string]any{"abc": "value"},
		wantErr: false,
	},
	{
		name:    "Value is unquotted string",
		args:    args{inputStr: "{abc:  name}"},
		want:    map[string]any{"abc": "name"},
		wantErr: false,
	},
	{
		name:    "Value with \\t 2",
		args:    args{inputStr: "{abc: \tname}"},
		want:    map[string]any{"abc": "name"},
		wantErr: false,
	},
	{
		name:    "Value with \\n 2",
		args:    args{inputStr: "{abc: \nvalue}"},
		want:    map[string]any{"abc": "value"},
		wantErr: false,
	},
}

var strangeInput = tests{
	{
		name: "Some additional not relevant data added",
		args: args{inputStr: `{"a": {"b": [12, 13, 14]}}text text`},
		want: map[string]any{"a": map[string]any{"b": []any{float64(12), float64(13), float64(14)}}},
	},
	{
		name: "JS variable declared",
		args: args{inputStr: `var test = {"a": {"b": [12, 13, 14]}}`},
		want: map[string]any{"a": map[string]any{"b": []any{float64(12), float64(13), float64(14)}}},
	},
	{
		name: "Form-feed symbol before value",
		args: args{inputStr: "{\"a\":\r\n10}"},
		want: map[string]any{"a": float64(10)},
	},
	{
		name: "Form-feed and line-feed added at the end",
		args: args{inputStr: "{'foo': 0,\r\n}"},
		want: map[string]any{"foo": float64(0)},
	},
	{
		name: "Just strange keys",
		args: args{inputStr: "{truefalse: 0, falsefalse: 1, nullnull: 2}"},
		want: map[string]any{"truefalse": float64(0), "falsefalse": float64(1), "nullnull": float64(2)},
	},
}

var integetNumericValuesTests = tests{
	{
		args: args{inputStr: "[0]"},
		want: []any{float64(0)},
	},
	{
		args: args{inputStr: "[1]"},
		want: []any{float64(1)},
	},
	{
		args: args{inputStr: "[12]"},
		want: []any{float64(12)},
	},
	{
		args: args{inputStr: "[12_12]"},
		want: []any{float64(1212)},
	},
	{
		args: args{inputStr: "[0x12]"},
		want: []any{float64(18)},
	},
	{
		args: args{inputStr: "[0xab]"},
		want: []any{float64(171)},
	},
	{
		args: args{inputStr: "[0xAB]"},
		want: []any{float64(171)},
	},
	{
		args: args{inputStr: "[0X12]"},
		want: []any{float64(18)},
	},
	{
		args: args{inputStr: "[0Xab]"},
		want: []any{float64(171)},
	},
	{
		args: args{inputStr: "[0XAB]"},
		want: []any{float64(171)},
	},
	{
		args: args{inputStr: "[01234]"},
		want: []any{float64(668)},
	},
	{
		args: args{inputStr: "[0o1234]"},
		want: []any{float64(668)},
	},
	{
		args: args{inputStr: "[0O1234]"},
		want: []any{float64(668)},
	},
	{
		args: args{inputStr: "[0b1111]"},
		want: []any{float64(15)},
	},
	{
		args: args{inputStr: "[0B1111]"},
		want: []any{float64(15)},
	},
	{
		name: "Negative zero - THIS TEST's OUTPUT DIFFERS FROM PYTHON!!!",
		args: args{inputStr: "[-0]"},
		want: []any{float64(0)}, // originally in chompjs: [-0]
	},
	{
		args: args{inputStr: "[-1]"},
		want: []any{float64(-1)},
	},
	{
		args: args{inputStr: "[-12]"},
		want: []any{float64(-12)},
	},
	{
		args: args{inputStr: "[-12_12]"},
		want: []any{float64(-1212)},
	},
	{
		args: args{inputStr: "[-0x12]"},
		want: []any{float64(-18)},
	},
	{
		args: args{inputStr: "[-0xab]"},
		want: []any{float64(-171)},
	},
	{
		args: args{inputStr: "[-0xAB]"},
		want: []any{float64(-171)},
	},
	{
		args: args{inputStr: "[-0x12]"},
		want: []any{float64(-18)},
	},
	{
		args: args{inputStr: "[-0Xab]"},
		want: []any{float64(-171)},
	},
	{
		args: args{inputStr: "[-0XAB]"},
		want: []any{float64(-171)},
	},
	{
		args: args{inputStr: "[-01234]"},
		want: []any{float64(-668)},
	},
	{
		args: args{inputStr: "[-0o1234]"},
		want: []any{float64(-668)},
	},
	{
		args: args{inputStr: "[-0O1234]"},
		want: []any{float64(-668)},
	},
	{
		args: args{inputStr: "[-0b1111]"},
		want: []any{float64(-15)},
	},
	{
		args: args{inputStr: "[-0B1111]"},
		want: []any{float64(-15)},
	},
}

var floatNumericValuesTests = tests{
	{
		args: args{inputStr: "[0.32]"},
		want: []any{0.32},
	},
	{
		args: args{inputStr: "[-0.32]"},
		want: []any{-0.32},
	},
	{
		args: args{inputStr: "[.32]"},
		want: []any{0.32},
	},
	{
		args: args{inputStr: "[-.32]"},
		want: []any{-0.32},
	},
	{
		args: args{inputStr: "[12.]"},
		want: []any{12.0},
	},
	{
		args: args{inputStr: "[-12.]"},
		want: []any{-12.0},
	},
	{
		args: args{inputStr: "[12.32]"},
		want: []any{12.32},
	},
	{
		args: args{inputStr: "[-12.12]"},
		want: []any{-12.12},
	},
	{
		args: args{inputStr: "[3.1415926]"},
		want: []any{3.1415926},
	},
	{
		args: args{inputStr: "[.123456789]"},
		want: []any{0.123456789},
	},
	{
		args: args{inputStr: "[.0123]"},
		want: []any{0.0123},
	},
	{
		args: args{inputStr: "[0.0123]"},
		want: []any{0.0123},
	},
	{
		args: args{inputStr: "[-.0123]"},
		want: []any{-.0123},
	},
	{
		args: args{inputStr: "[-0.0123]"},
		want: []any{-0.0123},
	},
	{
		args: args{inputStr: "[3.1E+12]"},
		want: []any{3.1e+12},
	},
	{
		args: args{inputStr: "[3.1e+12]"},
		want: []any{3.1e+12},
	},
	{
		args: args{inputStr: "[.1E+12]"},
		want: []any{.1e+12},
	},
	{
		args: args{inputStr: "[.1e+12]"},
		want: []any{.1e+12},
	},
}

var commentsTests = tests{
	{
		args: args{inputStr: `
			var obj = {
				// Comment
				x: "X", // Comment
			};
		`},
		want: map[string]any{"x": "X"},
	},
	{
		args: args{inputStr: `
			var /* Comment */ obj = /* Comment */ {
				/* Comment */
				x: /* Comment */ "X", /* Comment */
			};
		`},
		want: map[string]any{"x": "X"},
	},
	{
		args: args{inputStr: `[/*...*/1,2,3,/*...*/4,5,6]`},
		want: []any{float64(1), float64(2), float64(3), float64(4), float64(5), float64(6)},
	},
}

var exceptionsTests = tests{
	{
		args:    args{inputStr: "}{"},
		wantErr: true,
	},
	{
		args:    args{inputStr: ""},
		wantErr: true,
	},
}

var malformedInputTests = tests{
	{
		args:    args{inputStr: "{whose: 's's', category_name: '>'}"},
		wantErr: true,
	},
}

var errorMessagesTests = tests{
	{
		args:    args{inputStr: `{"test": """}`},
		wantErr: true,
	},
}

var unicodeEscapeTests = tests{
	{
		args: args{inputStr: "{\\\"a\\\": 12}", unicodeEscape: true},
		want: map[string]any{"a": float64(12)},
	},
}

var jsonNonStrictTests = tests{
	{
		args: args{inputStr: `["\n"]`},
		want: []any{"\n"},
	},
	{
		args: args{inputStr: `{'a': '\"\"', 'b': '\\\\', 'c': '\t\n'}`},
		want: map[string]any{"a": `""`, "b": "\\\\", "c": "\t\n"},
	},
	{
		name: "DOESN'T WORK!!!",
		args: args{inputStr: `
		var myObj = {
            myMethod: function(params) {
                // ...
            },
            myValue: 100
		}`},
		want: map[string]any{"myMethod": "function(params) {\n                    // ...\n                }", "myValue": float64(100)},
		skip: true,
	},
}

var jsonLoaderTests = tests{
	// some tests are dupes and present in other variables
	{
		args: args{inputStr: "[]"},
		want: []any{},
	},
	{
		args: args{inputStr: "[1, 2, 3]"},
		want: []any{float64(1), float64(2), float64(3)},
	},
	{
		args: args{inputStr: "var x = [1, 2, 3, 4, 5,]"},
		want: []any{float64(1), float64(2), float64(3), float64(4), float64(5)},
	},
	{
		args: args{inputStr: "{}"},
		want: map[string]any{},
	},
	{
		args: args{inputStr: "{'a': 12, 'b': 13, 'c': 14}"},
		want: map[string]any{"a": float64(12), "b": float64(13), "c": float64(14)},
	},
	{
		args: args{inputStr: "var x = {'a': 12, 'b': 13, 'c': 14}"},
		want: map[string]any{"a": float64(12), "b": float64(13), "c": float64(14)},
	},
}

func runner(t *testing.T, ut *tests) {
	for _, tt := range *ut {
		if tt.skip {
			t.Logf("The `%v/%v` test doesn't work, supposedly. Skipping it.", t.Name(), tt.name)
			continue
		}
		if tt.args.loader == nil {
			tt.args.loader = defaultLoader
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJsObject(&(tt.args.inputStr), tt.args.unicodeEscape, tt.args.loader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJsObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJsObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseObject(t *testing.T) {
	runner(t, &objectTests)
}

func TestParseList(t *testing.T) {
	runner(t, &listTests)
}

func TestParseMixed(t *testing.T) {
	runner(t, &mixedTests)
}

func TestStandardValues(t *testing.T) {
	runner(t, &standardValues)
}

func TestNaN(t *testing.T) {
	// this doesn't work as expected
	t.Skip("NaN values aren't supported by Go")
	var parseNaN = tests{
		{
			name:    "NaN",
			args:    args{inputStr: `{"A": NaN}`},
			want:    map[string]any{"A": math.NaN()},
			wantErr: false,
		},
	}
	runner(t, &parseNaN)
}

func TestStrangeValues(t *testing.T) {
	runner(t, &strangeValues)
}

func TestStrangeInput(t *testing.T) {
	runner(t, &strangeInput)
}

func TestIntegerNumericValues(t *testing.T) {
	runner(t, &integetNumericValuesTests)
}

func TestFloatNumericValues(t *testing.T) {
	runner(t, &floatNumericValuesTests)
}

func TestComments(t *testing.T) {
	runner(t, &commentsTests)
}

func TestException(t *testing.T) {
	runner(t, &exceptionsTests)
}

func TestMalformedInput(t *testing.T) {
	runner(t, &malformedInputTests)
}

func TestErrorMessages(t *testing.T) {
	runner(t, &errorMessagesTests)
}

func TestUnicodeEscape(t *testing.T) {
	runner(t, &unicodeEscapeTests)
}

func TestJsonNonStrict(t *testing.T) {
	runner(t, &jsonNonStrictTests)
}

func TestJsonLoader(t *testing.T) {
	runner(t, &jsonLoaderTests)
}

func TestParseJsObjects2(t *testing.T) {
	type args struct {
		inputStr      string
		unicodeEscape bool
		omitEmpty     bool
		loader        UnmarshalFunc
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		{
			args: args{inputStr: "[\"Test\\nDrive\"]\n{\"Test\": \"Drive\"}"},
			want: []any{[]any{"Test\nDrive"}, map[string]any{"Test": "Drive"}},
		},
		{
			args: args{inputStr: ""},
			want: []any{},
		},
		{
			args: args{inputStr: "aaaaaaaaaaaaaaaa"},
			want: []any{},
		},
		{
			args: args{inputStr: "         "},
			want: []any{},
		},
		{
			args: args{inputStr: "      {'a': 12}"},
			want: []any{map[string]any{"a": float64(12)}},
		},
		{
			args: args{inputStr: "[1, 2, 3, 4]xxxxxxxxxxxxxxxxxxxxxxxx"},
			want: []any{[]any{float64(1), float64(2), float64(3), float64(4)}},
		},
		{
			args: args{inputStr: "[12] [13] [14]"},
			want: []any{[]any{float64(12)}, []any{float64(13)}, []any{float64(14)}},
		},
		{
			args: args{inputStr: "[10] {'a': [1, 1, 1,]}"},
			want: []any{[]any{float64(10)}, map[string]any{"a": []any{float64(1), float64(1), float64(1)}}},
		},
		{
			args: args{inputStr: "[1][1][1]"},
			want: []any{[]any{float64(1)}, []any{float64(1)}, []any{float64(1)}},
		},
		{
			args: args{inputStr: "[1] [2] {'a': "},
			want: []any{[]any{float64(1)}, []any{float64(2)}},
		},
		{
			args: args{inputStr: "[]"},
			want: []any{[]any{}},
		},
		{
			args: args{inputStr: "[][][][]"},
			want: []any{[]any{}, []any{}, []any{}, []any{}},
		},
		{
			args: args{inputStr: "{}"},
			want: []any{map[string]any{}},
		},
		{
			args: args{inputStr: "{}{}{}{}"},
			want: []any{map[string]any{}, map[string]any{}, map[string]any{}, map[string]any{}},
		},
		{
			args: args{inputStr: "{{}}{{}}"},
			want: []any{},
		},
		{
			args: args{inputStr: "[[]][[]]"},
			want: []any{[]any{[]any{}}, []any{[]any{}}},
		},
		{
			args: args{inputStr: "{am: 'ab'}\n{'ab': 'xx'}"},
			want: []any{map[string]any{"am": "ab"}, map[string]any{"ab": "xx"}},
		},
		{
			args: args{inputStr: "function(a, b, c){ /* ... */ }({\"a\": 12}, Null, [1, 2, 3])"},
			want: []any{map[string]any{}, map[string]any{"a": float64(12)}, []any{float64(1), float64(2), float64(3)}},
		},
		{
			args: args{inputStr: "{\"a\": 12, broken}{\"c\": 100}"},
			want: []any{map[string]any{"c": float64(100)}},
		},
		{
			args: args{inputStr: "[12,,,,21][211,,,][12,12][12,,,21]"},
			want: []any{[]any{float64(12), float64(12)}},
		},
		{
			args: args{inputStr: "[1][][2]", omitEmpty: true},
			want: []any{[]any{float64(1)}, []any{float64(2)}},
		},
		{
			args: args{inputStr: "{'a': 12}{}{'b': 13}", omitEmpty: true},
			want: []any{map[string]any{"a": float64(12)}, map[string]any{"b": float64(13)}},
		},
		{
			args: args{inputStr: "[][][][][][][][][]", omitEmpty: true},
			want: []any{},
		},
		{
			args: args{inputStr: "{}{}{}{}{}{}{}{}{}", omitEmpty: true},
			want: []any{},
		},
	}
	for _, tt := range tests {
		if tt.args.loader == nil {
			tt.args.loader = defaultLoader
		}
		t.Run(tt.name, func(t *testing.T) {
			dataChannel, errChannel := ParseJsObjects(&(tt.args.inputStr), tt.args.unicodeEscape, tt.args.omitEmpty, tt.args.loader)
			var got []any
			var parseErr error
		OuterLabel:
			for {
				select {
				case data, ok := <-dataChannel:
					if !ok {
						// channel closed
						break OuterLabel
					}
					got = append(got, data)
				case err, ok := <-errChannel:
					if ok { // if error returned
						parseErr = err
					}
					break OuterLabel
				}
			}
			if (parseErr != nil) != tt.wantErr {
				t.Errorf("ParseJsObjects() error = %v, wantErr %v", parseErr, tt.wantErr)
				return
			}
			if len(got) == len(tt.want) && len(got) == 0 {
				// empty slices, we can't use reflect.DeepEqual on empty slices
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJsObjects() = %v, want %v", got, tt.want)
			}
		})
	}
}
