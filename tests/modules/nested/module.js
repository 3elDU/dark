import { message } from "../message.js";

console.log(`Message imported in a nested script: ${message}`);
console.log(
  `(from nested) import.meta.resolve('./test.js')=${import.meta.resolve("./test.js")}`,
);
