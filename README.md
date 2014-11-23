ottily
======

Ottily executes a javascript snippet on each line of an input file in parallel.

Installation
------------

    $ go get github.com/miku/ottily/cmd/ottily

Usage
-----

    $ ottily -h
    Usage of ottily:
      -cpuprofile="": write cpu profile to file
      -e="": execute argument on each line of input
      -s="": script to execute on each line of input
      -v=false: prints current program version
      -w=4: number of workers

Examples
--------

Noop -- is just an expensive `cat`.

    $ ottily datasets/simple.ldj
    {"name": "ottily", "language": "Golang"}

Inline script with `-e`:

    $ ottily -e 'output=input.length' datasets/simple.ldj
    40

    $ ottily -e 'o=JSON.parse(input); o["language"] = "Go"; output=JSON.stringify(o);' \
                datasets/simple.ldj

    {"language":"Go","name":"ottily"}

Pass a script file with `-s`:

    $ cat scripts/classified.js
    output = "CLASSIFIED"

    $ ottily -s scripts/classified.js datasets/simple.ldj
    CLASSIFIED

A string variable `input` is passed into the javascript snippet.
To produce output, set the `output` variable to the desired string value.
If `output` is set to null or is undefined, nothing is printed. That's it.

By default, `ottily` will run as many workers as there are cores.

Performance
-----------

Ottily is just a [156 LOC](https://github.com/miku/ottily/blob/6d81c71afe2a29fb5d3445b0813642285463ca6b/cmd/ottily/main.go) Go program. In many cases, there will be faster alternatives.

Given a file with 1 million lines, calculate the length of each line.

    $ time awk '{ print length }' datasets/1M.ldj  > /dev/null

    real    0m5.836s
    user    0m5.514s
    sys     0m0.314s

    $ time ottily -e 'output=input.length' datasets/1M.ldj > /dev/null

    real    0m11.758s
    user    0m38.204s
    sys     0m4.713s

Given a file with 1 million lines, one JSON document per line, add a new key to each JSON document.

    $ time jq -c '."about" = "jq"' datasets/1M.ldj > /dev/null

    real    2m31.063s
    user    2m29.517s
    sys     0m1.500s

    $ time ottily -e 'o=JSON.parse(input); o["about"] = "ot"; output=JSON.stringify(o);' \
                  datasets/1M.ldj > /dev/null

    real    5m55.630s
    user    20m25.530s
    sys     0m46.750s

Above tests were done on 4 cores.

Credits
-------

Ottily uses [otto](https://github.com/robertkrimen/otto) to run js.
