# Pyro

Pyro is a simple interpreted programming language built in Go.

## Features

Pyro supports:

- Standard Output (`print`)
- Function Declarations
- Global Variables
- Arithmetic Expressions (`+`, `-`, `*`, `/`, `%`)
- Comparison Operators (`<`, `<=`, `>`, `>=`, `==`, `!=`)
- Logical Operators (`and`, `or`, `!`)
- Control Flow  
  - `if` / `else`  
  - `while` loops  
  - `for` loops
- Blocks & Scoping
- Closures and Lexical Scoping
- Tree-Walk Interpreter Architecture


## ⚙️ Dependencies

- [Golang](https://golang.org/) (Go 1.18 or later)

## Running Pyro

To run a `.pyro` file:

```bash
./pyro <filename>.pyro
```
## Sample Code

Here’s a sample Pyro program that prints the FizzBuzz sequence:

```pyro
fun FizzBuzz(n) {
  for (var i = 1; i <= n; i = i + 1) {
    if (i % 3 == 0) {
      if (i % 5 == 0) {
        print "FizzBuzz";
      } else {
        print "Fizz";
      }
    } else {
      if (i % 5 == 0) {
        print "Buzz";
      } else {
        print i;
      }
    }
  }
}

FizzBuzz(15);
```

## Credits
Pyro is based on the book [Crafting Interpreters](https://craftinginterpreters.com/) by Bob Nystrom.
