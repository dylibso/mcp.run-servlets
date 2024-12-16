# Time Servlet

Time operations plugin. It provides the following operations:
        
- `get_time_utc`: Returns the current time in the UTC timezone. Takes no parameters.
- `parse_time`: Takes a `time` string in RFC2822 format and returns the timestamp in UTC timezone.
- `time_offset`: Takes integer `time` and `offset` parameters. Adds a time offset to a given timestamp and returns the new timestamp in UTC timezone.
        
## Examples

Get current time:

```typescript
{
  `name`: `get_time_utc`
}
// returns:
{
    "utc_time" : "1734085548",
    "utc_time_rfc2822" : "Fri, 13 Dec 2024 10:25:48 +0000"
}
```

What would be last Sunday?

```typescript
{
  `name`: `time_offset`,
  `time`: 1734085548,
  `offset`: -432000
}
// returns:
{
    "utc_time" : "1733653548",
    "utc_time_rfc2822" : "Sun, 8 Dec 2024 10:25:48 +0000"
}
```
