#!/usr/bin/env bash

#TEST_COMMAND="./config-expression"
TEST_COMMAND="go/confexpr-go"

OUTPUT_CODE=0;

OUTPUT_MD_FILE=examples/README.md
rm -f $OUTPUT_MD_FILE
cat examples/header.md >> examples/README.md
for EXAMPLE_DIR in $(find examples -type d -depth 1); do
	echo $EXAMPLE_DIR
	CONF_FILE=$EXAMPLE_DIR/conf.jsonw
	INPUT_MD_FILE=$EXAMPLE_DIR/input.md
	echo "## $(basename $EXAMPLE_DIR | sed -e 's/[-_]/ /g')" >> $OUTPUT_MD_FILE
	cat $INPUT_MD_FILE >> $OUTPUT_MD_FILE
	echo -e "\n" >> $OUTPUT_MD_FILE
	sed -e 's/^/    /' $CONF_FILE >> $OUTPUT_MD_FILE


	PASSING_FILE=$EXAMPLE_DIR/passing-data.csv
	if [ -f $PASSING_FILE ]; then
		echo -e "\nThis will produce the following output\n" >> $OUTPUT_MD_FILE
		while read -r LINE; do
			echo $LINE | grep '#' > /dev/null
			CODE=$?
			if [ "$CODE" == "0" ]; then
				echo -e "$(echo -e "$LINE" | sed -e 's/#//')\n" >> $OUTPUT_MD_FILE
			else
				read -r INPUT EXPECTED_OUTPUT <<< "$LINE"
				ACTUAL_OUTPUT=$($TEST_COMMAND lookup --file $CONF_FILE $INPUT)
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
				echo "    $ $TEST_COMMAND lookup $INPUT" >> $OUTPUT_MD_FILE
				echo -e "    '$ACTUAL_OUTPUT'\n     " >> $OUTPUT_MD_FILE
			fi
		done < <(cat $PASSING_FILE | grep '.')
	fi

	FAILING_FILE=$EXAMPLE_DIR/failing-data.csv
	if [ -f $FAILING_FILE ]; then
		echo -e "\nThe following examples will fail:\n" >> $OUTPUT_MD_FILE
		while read -r LINE; do
			echo $LINE | grep '#' > /dev/null
			CODE=$?
			if [ "$CODE" == "0" ]; then
				echo -e "$(echo -e "$LINE" | sed -e 's/#//')\n" >> $OUTPUT_MD_FILE
			else
				read -r INPUT <<< "$LINE"
				ACTUAL_OUTPUT=$($TEST_COMMAND lookup --file $CONF_FILE $INPUT)
				ACTUAL_CODE=$?
				RESULT="PASS"
				if [ "$ACTUAL_CODE" == "0" ]; then
					OUTPUT_CODE=$((OUTPUT_CODE+1));
					RESULT="FAIL";
				fi
				if [ "$RESULT" == "PASS" ]; then
					echo -e "\x1B[1;32m$RESULT\x1B[0m - Input: '$INPUT' Actual: '$ACTUAL_CODE' Expected: '1'"
				else
					echo -e "\x1B[1;31m$RESULT\x1B[0m - Input: '$INPUT' Actual: '$ACTUAL_CODE' Expected: '1'"
				fi
				echo "    $ $TEST_COMMAND lookup $INPUT" >> $OUTPUT_MD_FILE
				#echo -e "    '$ACTUAL_OUTPUT'\n     " >> $OUTPUT_MD_FILE
				$TEST_COMMAND lookup --file $CONF_FILE $INPUT | sed -e 's/^/    /' >> $OUTPUT_MD_FILE
			fi
		done < <(cat $FAILING_FILE | grep '.')
	fi

	EXPLAIN_FILE=$EXAMPLE_DIR/explain-data.csv
	if [ -f $EXPLAIN_FILE ]; then
		echo -e "\nThe explain command can be used to gain insight into why the tool returned a certain value for an expression\n" >> $OUTPUT_MD_FILE
		while read -r LINE; do
			echo $LINE | grep '#' > /dev/null
			CODE=$?
			if [ "$CODE" == "0" ]; then
				echo -e "$(echo -e "$LINE" | sed -e 's/#//')\n" >> $OUTPUT_MD_FILE
			else
				read -r INPUT <<< "$LINE"
				echo "    $ $TEST_COMMAND explain $INPUT" >> $OUTPUT_MD_FILE
				$TEST_COMMAND explain --no-color --file $CONF_FILE $INPUT | sed -e 's/^/    /' >> $OUTPUT_MD_FILE
			fi
		done < <(cat $EXPLAIN_FILE | grep '.')
	fi
done
cat examples/footer.md >> examples/README.md

echo "There were $OUTPUT_CODE failing tests"
if [ "$OUTPUT_CODE" != "0" ]; then
	exit 1
fi
