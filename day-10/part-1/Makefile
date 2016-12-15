.SHELL=bash

input.png: input.dot
	dot -Tpng input.dot > input.png

input.dot: input.txt
	(echo 'digraph G {'; gsed -r -e 's/^value ([0-9]+) goes to (\w+) ([0-9]+)/value_\1 -> \2_\3/' -e 's/^bot ([0-9]+) gives low to (\w+) ([0-9]+) and high to (\w+) ([0-9]+)/bot_\1 -> \2_\3[label="L"]; bot_\1 -> \4_\5[label="H"]/'; echo '}') < input.txt > input.dot 

clean:
	rm input.{dot,png}
