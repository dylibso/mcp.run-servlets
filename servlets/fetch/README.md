# Fetch Servlet

A servlet that fetches web pages and converts them to markdown.

## What it does

Takes a URL, fetches the page content, strips out scripts and styles, and converts the HTML to markdown format.

## Usage

Call with:
```typescript
{
  arguments: {
    url: "https://example.com"  // Required: URL to fetch
  }
}
```

Returns the page content converted to markdown format.