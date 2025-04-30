# Fetch Image Servlet

A servlet that fetches pictures and returns them as resources. Intended for interactive use (e.g. in Claude Desktop).

## What it does

Takes a URL, fetches the image and return it base64-encoded for inline-display.

## Usage

Call with:
```typescript
{
  `arguments`: {
    `url`: `https://example.com`,  // Required: URL to fetch
    `mime-type`: "image/png"       // The mime type to filter by
  }
}
```
