#!/usr/bin/env node

var fs = require('fs'),
readline = require('readline');
vm = require('vm');

var args = process.argv.slice(2);

if (args.length < 1) {
    console.log("Usage: eachline.js LDJ")
    process.exit(code=1)
}

var rd = readline.createInterface({
    input: fs.createReadStream(args[0]),
    output: process.stdout,
    terminal: false
});

var i = 0;
rd.on('line', function(line) {
    var obj = JSON.parse(line);
    i++;
    if (i % 1000000 == 0) {
        console.log(i);
    };
});

