# Filesystem Servlet

A servlet that performs file system operations.

## What it does

Provides four file system operations:
- Read file contents
- Write content to a file
- List directory contents
- Search for files by name pattern

## Usage

Call with one of these operations:

```typescript
// Read a file
{
  arguments: {
    name: "read_file",
    path: "/path/to/file"
  }
}

// Write to a file
{
  arguments: {
    name: "write_file",
    path: "/path/to/file",
    content: "content to write"
  }
}

// List directory contents
{
  arguments: {
    name: "list_directory",
    path: "/path/to/dir"
  }
}

// Search for files
{
  arguments: {
    name: "search_files",
    path: "/path/to/search",
    pattern: "search-term"
  }
}
```