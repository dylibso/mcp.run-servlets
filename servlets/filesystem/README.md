# Filesystem Servlet

A servlet that performs file system operations.

## What it does

Provides four file system operations:
- Read file contents
- Write content to a file
- List directory contents
- Search for files by name pattern

## Tools

### read_file

Reads the contents of a file at the specified path.

```json
{
  "path": "path/to/file.txt"
}
```

### write_file

Writes content to a file at the specified path.

```json
{
  "path": "path/to/file.txt",
  "content": "Hello, world!"
}
```

### list_directory

Lists the contents of a directory, marking items as either [DIR] or [FILE].

```json
{
  "path": "path/to/directory"
}
```

### search_files

Recursively searches for files matching a pattern in the given directory.

```json
{
  "path": "path/to/search",
  "pattern": "*.txt"
}
```
