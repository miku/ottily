ottily
======

Ottily executes a javascript snippet on each line of an input file in parallel.

Noop.

    $ ottily -i datasets/simple.ldj
    {"name": "ottily", "language": "Golang"}

Inline script with -e:

    $ ottily -i datasets/simple.ldj -e 'output=input.length'
    40

    $ ottily -i datasets/simple.ldj -e 'o=JSON.parse(input); o["language"] = "Go"; output=JSON.stringify(o);'
    {"language":"Go","name":"ottily"}

Pass a script file:

    $ ottily -i datasets/simple.ldj -s scripts/classified.js
    CLASSIFIED
