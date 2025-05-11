package gompjs

import (
	"strconv"

	"github.com/proway2/gompjs/internal/chompjs"
)

type UnmarshalFunc func([]byte, any) error

func ParseJsObject(inputStr *string, unicodeEscape bool, loader UnmarshalFunc) (any, error) {
	var err error
	if unicodeEscape {
		if inputStr, err = decodeUnicodeEscape(inputStr); err != nil {
			return nil, err
		}
	}
	var parsedString *string
	if parsedString, err = chompjs.FixString(inputStr); err != nil {
		return nil, err
	}
	var res any
	byteParsedString := []byte(*parsedString)
	if err = parseString(loader, &byteParsedString, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func ParseJsObjects(inputStr *string, unicodeEscape, omitEmpty bool, loader UnmarshalFunc) (<-chan any, <-chan error) {
	dataChannel := make(chan any)
	errChannel := make(chan error, 1)
	var err error
	if unicodeEscape {
		if inputStr, err = decodeUnicodeEscape(inputStr); err != nil {
			errChannel <- err
			return dataChannel, errChannel
		}
	}
	go func() {
		defer close(dataChannel)
		defer close(errChannel)
		chompjsResCh, chompjsErrCh := chompjs.FixStrings(inputStr)
		for {
			select {
			case parsedString, ok := <-chompjsResCh:
				if !ok {
					return
				}
				var element any
				byteParsedString := []byte(*parsedString)
				if err := parseString(loader, &byteParsedString, &element); err != nil {
					// Original Python code skips on loader error
					// try:
					// 	data = loader(raw_data, *loader_args, **loader_kwargs)
					// except ValueError:
					// 	continue
					continue
				}
				if omitEmpty {
					switch v := element.(type) {
					case []any:
						if len(v) == 0 {
							continue
						}
					case map[string]any:
						if len(v) == 0 {
							continue
						}
					}
				}
				dataChannel <- element
			case err := <-chompjsErrCh:
				if err != nil {
					errChannel <- err
					return
				}
			}
		}
	}()
	return dataChannel, errChannel
}

func parseString(loader UnmarshalFunc, data *[]byte, v any) error {
	if err := loader(*data, &v); err != nil {
		return err
	}
	return nil
}

func decodeUnicodeEscape(s *string) (*string, error) {
	// quotes must be addef for strconv.Unquote to recogniz a string
	quoted := `"` + *s + `"`
	unquoted, err := strconv.Unquote(quoted)
	if err != nil {
		return nil, err
	}
	return &unquoted, nil
}
