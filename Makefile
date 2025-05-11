CHOMPJS_PATH=internal/chompjs
FULL_C_PATH := $(PWD)/$(CHOMPJS_PATH)

DIR_OUT=lib

COMPILE_OPTS=-DNDEBUG -fwrapv -O2 -fstack-protector-strong -D_FORTIFY_SOURCE=2 -fPIC
LINKER_OPTS=-shared -Wl,-O1 -Wl,-Bsymbolic-functions -g -fwrapv -O2 -fstack-protector-strong -Wformat -Werror=format-security -Wdate-time -D_FORTIFY_SOURCE=2

shared: makedir
	gcc $(COMPILE_OPTS) -c $(FULL_C_PATH)/buffer.c -o $(DIR_OUT)/buffer.o -Wl,-Bsymbolic-functions
	gcc $(COMPILE_OPTS) -c $(FULL_C_PATH)/parser.c -o $(DIR_OUT)/parser.o -Wl,-Bsymbolic-functions
	gcc $(LINKER_OPTS) $(DIR_OUT)/*.o -o $(DIR_OUT)/libchompjs.so

makedir:
	@rm -rf $(DIR_OUT)
	@mkdir -p $(DIR_OUT)

build:
	go build -a ./...

test:
	go test -v ./...

coverage:
	@rm -f coverage.profile coverage.html
	@-go test -coverprofile=coverage.profile ./...
	@go tool cover -html=coverage.profile -o coverage.html

example:
	@rm -f main
	@go build ./examples/main.go
	./main

chompjs:
	./scripts/get_chompjs.sh
