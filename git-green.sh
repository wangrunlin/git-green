#!/bin/bash
for i in {365..1000};do
	# write text
	datetime=`date -v -"$i"d "+%Y-%m-%dT%H:%M:%S"`

	git checkout -B git-green

	echo "$(date "+%Y-%m-%d") add $datetime\n" >> git-green.md

	git add git-green.md
	git commit -m "Update git-green.md"

	# modify date
	git commit --amend --date="$datetime" -m "Update git-green.md"

done;

# git push


