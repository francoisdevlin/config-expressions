#!/usr/bin/env bash

TEST_COMMAND="./pattern-getter.rb"
TEST_DIR=test-data/hello
CONF_FILE=$TEST_DIR/conf.jsonw
PASSING_FILE=$TEST_DIR/passing-data.csv

OUTPUT_CODE=0;
while read -r INPUT EXPECTED_OUTPUT; do
	ACTUAL_OUTPUT=$($TEST_COMMAND --file $CONF_FILE $INPUT)
	ACTUAL_CODE=$?
	RESULT="PASS"
	if [ "$ACTUAL_CODE" != "0" ]; then
		OUTPUT_CODE=$((OUTPUT_CODE+1));
		RESULT="FAIL";
	elif [ "$ACTUAL_OUTPUT" != "$EXPECTED_OUTPUT" ]; then
		OUTPUT_CODE=$((OUTPUT_CODE+1));
		RESULT="FAIL";
	fi
	if [ "$RESULT" == "PASS" ]; then
		echo -e "\x1B[1;32m$RESULT\x1B[0m - Input: '$INPUT' Actual: '$ACTUAL_OUTPUT' Expected: '$EXPECTED_OUTPUT'"
	else
		echo -e "\x1B[1;31m$RESULT\x1B[0m - Input: '$INPUT' Actual: '$ACTUAL_OUTPUT' Expected: '$EXPECTED_OUTPUT'"
	fi

done < <(sed -e 's/#.*//' $PASSING_FILE | grep '.')

echo "There were $OUTPUT_CODE failing tests"
if [ "$OUTPUT_CODE" != "0" ]; then
	exit 1
fi
