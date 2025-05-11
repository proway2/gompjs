package chompjs

/*
#cgo LDFLAGS: -lm

#include <stdlib.h>
#include "parser.h"
struct Lexer lexer;

*/
import "C"
import (
	"fmt"
	"unsafe"
)

func FixString(input *string) (*string, error) {
	inputStr := C.CString(*input)
	defer C.free(unsafe.Pointer(inputStr))
	C.init_lexer(&C.lexer, inputStr)
	for C.lexer.lexer_status == C.CAN_ADVANCE {
		C.advance(&C.lexer)
	}
	parsedString := C.GoString(C.lexer.output.data)
	C.release_lexer(&C.lexer)
	if C.lexer.lexer_status == C.ERROR {
		err := fmt.Errorf("error parsing input near character %d", uint64(C.lexer.input_position))
		return nil, err
	}
	return &parsedString, nil
}

func FixStrings(input *string) (<-chan *string, <-chan error) {
	dataChannel := make(chan *string)
	// this channel is created but not actually being used, as the original code doesn't raise on errors
	errChannel := make(chan error, 1)

	inputStr := C.CString(*input)

	go func() {
		defer close(dataChannel)
		defer close(errChannel)
		defer C.free(unsafe.Pointer(inputStr))

		// json_iter_new (parser.h)
		C.init_lexer(&C.lexer, inputStr)

		// json_iter_dealloc (parser.h)
		defer C.release_lexer(&C.lexer)

		// json_iter_next (parser.h)
		for {
			for C.lexer.lexer_status == C.CAN_ADVANCE {
				C.advance(&C.lexer)
			}
			if C.lexer.output.index == 1 {
				return
			}
			// THIS CODE, IF ENABLED, MAY CAUSE OR MAY NOT ERRORS!!!
			// this is not from the original code ->
			// if C.lexer.lexer_status == C.ERROR {
			// 	err := fmt.Errorf("error parsing input near character %d", uint64(C.lexer.input_position))
			// 	errChannel <- err
			// 	return
			// }
			// <-
			// writing correct data into the channel
			parsedString := C.GoString(C.lexer.output.data)
			dataChannel <- &parsedString
			// from json_iter_next (parser.h)
			C.reset_lexer_output(&C.lexer)
		}
	}()
	return dataChannel, errChannel
}
