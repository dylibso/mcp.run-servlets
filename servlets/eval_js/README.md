# JavaScript Eval Servlet

A simple servlet that evaluates JavaScript code in a QuickJS Wasm sandbox and returns the result.

## What it does

Takes JavaScript code as input, evaluates it using `eval()`, and returns the result as a string.

## Usage

Call with:
```typescript
{
  arguments: {
    code: "2 + 2"  // Required: JavaScript code to evaluate
  }
}
```

Returns:
```typescript
"4"
```

Note: Don't use `console.log()` in the code - the result needs to be returned directly.
