#!/bin/sh

CURL=`which curl`

COMMIT_HASH="97dea8436c7ca680a11ce558f70597fc7621e17f"
BASE_URL="https://raw.githubusercontent.com/Nykakin/chompjs/$COMMIT_HASH/_chompjs"

BUFFER_C="$BASE_URL/buffer.c"
BUFFER_H="$BASE_URL/buffer.h"
PARSER_C="$BASE_URL/parser.c"
PARSER_H="$BASE_URL/parser.h"

PATH_TO_CHOMPJS="./internal/chompjs"
$CURL -s -o $PATH_TO_CHOMPJS/buffer.c $BUFFER_C
$CURL -s -o $PATH_TO_CHOMPJS/buffer.h $BUFFER_H
$CURL -s -o $PATH_TO_CHOMPJS/parser.c $PARSER_C
$CURL -s -o $PATH_TO_CHOMPJS/parser.h $PARSER_H
