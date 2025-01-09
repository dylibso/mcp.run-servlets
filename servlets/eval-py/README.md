# Python Eval Servlet

A simple servlet that evaluates Python code in a CPython Wasm sandbox and returns the result.

## What it does

Takes Python code as input, evaluates it using `exec()`, and returns the stdout output as a string.

## Usage

Call with:

```json
{
  "arguments": {
    "code": "print(2 + 2)"  // Required: Python code to evaluate
  }
}
```

Returns:

```
4
```
