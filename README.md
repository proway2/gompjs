# Gompjs

## Purpose

This Go package is inspired by the excellent Python library [chompjs](https://github.com/Nykakin/chompjs) created by @nykakin.  
Its primary purpose is web scraping. The creation rationale is well documented in the original library.  
This package serves as a wrapper around C code taken from the original library.

## Tests

Unit tests were adapted from the original Python package for Go compatibility.

### Limitations

When using `encoding/json`, the following features are unavailable compared to the Python library:

* `NaN` values return an error (test skipped). Consider patching the original `parser.c` to return `NaN` as string `"NaN"` by removing relevant lines.  
  Alternative library: https://github.com/xhhuango/json (supports `NaN`, `+Inf`, `-Inf` parsing)
* Control characters in JSON aren't supported. This test fails:

```javascript
var myObj = {
    myMethod: function(params) {
        // ...
    },
    myValue: 100
}
```

In Python, this only works with: `parse_js_object(in_data, loader_kwargs={'strict': False})`

## Available Functions

Package gompjs provides:

```go
// Equivalent to chompjs.parse_js_object
func ParseJsObject(inputStr *string, unicodeEscape bool, loader UnmarshalFunc) (any, error)

// Equivalent to chompjs.parse_js_objects
func ParseJsObjects(inputStr *string, unicodeEscape, omitEmpty bool, loader UnmarshalFunc) (<-chan any, <-chan error)
```

The `UnmarshalFunc` type mirrors `encoding/json`'s Unmarshal signature, enabling compatibility with third-party JSON libraries:

```go
type UnmarshalFunc func([]byte, any) error
```

## Usage

Import the package:

```go
import "github.com/proway2/gompjs/pkg/gompjs"
```

See [./examples](./examples/) for implementation samples.

## ToDo

1. Comprehensive compatibility testing with other Go JSON parsers
2. Memory leak verification
3. Performance benchmarking (including large real-world JSON datasets)
4. Proper documentation

## Examples

Run examples via:

```bash
make example
```
