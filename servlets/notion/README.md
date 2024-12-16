# Notion MCP Tool

A Model Context Protocol (MCP) tool for interacting with the Notion API. This tool provides a comprehensive interface to Notion's API endpoints, allowing you to manipulate blocks, pages, databases, users, and more.

## Features

- **Block Operations**
  - Append block children
  - Retrieve blocks
  - List block children
  - Delete blocks

- **Page Operations**
  - Retrieve pages
  - Update page properties

- **Database Operations**
  - Create databases
  - Query databases
  - Retrieve database details
  - Update database properties
  - Create database items

- **Comment Operations**
  - Create comments
  - Retrieve comments

- **User Operations**
  - List workspace users
  - Retrieve user details
  - Get bot user information

- **Search Functionality**
  - Search across pages and databases

## Setup

1. Create a Notion integration at https://www.notion.so/my-integrations
2. Get your integration token
3. Configure the token in MCP:
```bash
mcpx config set NOTION_TOKEN "your-integration-token"
```

## Usage

The tool exposes a single `notion` command with various operations. Here are some example usages:

### Query a Database

```json
{
  "operation": "query_database",
  "database_id": "your-database-id",
  "filter": {
    "property": "Status",
    "select": {
      "equals": "Done"
    }
  },
  "sorts": [
    {
      "property": "Created",
      "direction": "descending"
    }
  ]
}
```

### Create a Page in a Database

```json
{
  "operation": "create_database_item",
  "database_id": "your-database-id",
  "properties": {
    "Name": {
      "title": [
        {
          "text": {
            "content": "New task"
          }
        }
      ]
    },
    "Status": {
      "select": {
        "name": "In Progress"
      }
    }
  }
}
```

### Retrieve Block Children

```json
{
  "operation": "retrieve_block_children",
  "block_id": "your-block-id",
  "page_size": 50
}
```

## Error Handling

The tool provides detailed error messages when:
- Required parameters are missing
- The Notion API returns an error
- The integration token is invalid or missing
- The request format is incorrect

## Integration Permissions

Make sure your Notion integration has the necessary capabilities enabled for the operations you want to perform:
- Read content
- Update content
- Insert content
- Read comments
- Create comments
- Read user information (requires Enterprise plan for some operations)


