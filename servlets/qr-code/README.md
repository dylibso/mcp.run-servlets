# QR Code Servlet

A servlet that generates QR codes as PNG images.

## What it does

Takes input text (like a URL or message) and generates a QR code image.

## Usage

Call with:
```typescript
{
  arguments: {
    data: "https://example.com",  // Required: text to encode
    ecc: 4,                       // Optional: error correction (1-4, default 4)
  }
}
```

Returns a base64-encoded PNG image of the QR code.