const file = process.argv[2];
const field = process.argv[3];

const content = JSON.parse(require('fs').readFileSync(file).toString());
const fn = new Function("json", "return json." + field + ";");
console.log(fn(content));
