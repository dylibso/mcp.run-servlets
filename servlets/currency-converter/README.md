# Currency Converter Servlet

A servlet that converts amounts between different currencies using current exchange rates.

## What it does

Takes an amount and currency codes, returns the converted amount using current rates from fxratesapi.com.

## Usage

Call with:
```typescript
{
  arguments: {
    amount: 100,           // Required: amount to convert
    from: "USD",          // Required: source currency code
    to: "EUR"            // Required: target currency code
  }
}
```

Returns the converted amount as a string.