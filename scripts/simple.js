// Replace a nested value
//
//     $ ottily -i dataset/sample.js -s scripts/simple.js
//
obj = JSON.parse(input);
obj["003"] = "DE-15";
for (i = 0; i < obj["245"].length; i++) {
	obj["245"][i]["a"] = "THE Book";
}
output = JSON.stringify(obj);
