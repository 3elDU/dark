import { message } from "./message.js";

console.log(`Hello from a module script! The message is: ${message}`);

console.log(`import.meta.url=${import.meta.url}`);
console.log(
  `import.meta.resolve('./test.js')=${import.meta.resolve("./test.js")}`,
);
