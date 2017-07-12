#!/usr/bin/env bash

TEST_COMMAND="./pattern-getter.rb"

OUTPUT_CODE=0;

OUTPUT_MD_FILE=examples/README.md
rm -f $OUTPUT_MD_FILE
cat examples/header.md >> examples/README.md
for EXAMPLE_DIR in $(find examples -type d -depth 1); do
	echo $EXAMPLE_DIR
	CONF_FILE=$EXAMPLE_DIR/conf.jsonw
	PASSING_FILE=$EXAMPLE_DIR/passing-data.csv
	INPUT_MD_FILE=$EXAMPLE_DIR/input.md
	echo "## $(basename $EXAMPLE_DIR | sed -e 's/[-_]/ /g')" >> $OUTPUT_MD_FILE
	cat $INPUT_MD_FILE >> $OUTPUT_MD_FILE
	echo -e "\n" >> $OUTPUT_MD_FILE
	sed -e 's/^/    /' $CONF_FILE >> $OUTPUT_MD_FILE
	echo -e "\nThis will produce the following output\n" >> $OUTPUT_MD_FILE
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
		echo "    $ $TEST_COMMAND $INPUT" >> $OUTPUT_MD_FILE
		echo -e "    $ACTUAL_OUTPUT\n     " >> $OUTPUT_MD_FILE
	done < <(sed -e 's/#.*//' $PASSING_FILE | grep '.')
done
cat examples/footer.md >> examples/README.md

echo "There were $OUTPUT_CODE failing tests"
if [ "$OUTPUT_CODE" != "0" ]; then
	exit 1
fi
