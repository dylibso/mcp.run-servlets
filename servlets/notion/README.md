# Notion MCP Tool

A Model Context Protocol plugin that provides comprehensive access to the Notion API, allowing seamless integration with Notion's blocks, pages, databases, users, and comments functionality.

## Features

- **Block Operations**
  - Append child blocks
  - Retrieve blocks and their children
  - Delete blocks
- **Page Operations**
  - Retrieve pages
  - Update page properties
- **Database Operations**
  - Create and query databases
  - Retrieve and update database properties
  - Create database items
- **User Management**
  - List workspace users
  - Retrieve user information
  - Access bot user details
- **Comment System**
  - Create comments
  - Retrieve comments for blocks
- **Search Functionality**
  - Search across pages and databases
  - Filter and sort results

## Setup

  * Create a Notion integration at [https://www.notion.so/my-integrations](https://www.notion.so/my-integrations)
  * Get your integration token
  * Set the `NOTION_TOKEN` when installing

## Usage

### Tool Names and Descriptions

Each Notion operation is exposed as a separate tool with the prefix `notion_`. Here are the available tools:

#### Block Operations
- `notion_append_block_children`: Append new blocks to a parent block
- `notion_retrieve_block`: Get a specific block's information
- `notion_retrieve_block_children`: List a block's child blocks
- `notion_delete_block`: Remove a block

#### Page Operations
- `notion_retrieve_page`: Get page information
- `notion_update_page_properties`: Modify page properties

#### Database Operations
- `notion_create_database`: Create a new database
- `notion_query_database`: Search and filter database items
- `notion_retrieve_database`: Get database information
- `notion_update_database`: Modify database properties
- `notion_create_database_item`: Add a new item to a database

#### User Operations
- `notion_list_users`: List workspace users
- `notion_retrieve_user`: Get user information
- `notion_retrieve_bot_user`: Get current bot user details

#### Comment Operations
- `notion_create_comment`: Create a new comment
- `notion_retrieve_comments`: Get comments for a block

#### Search Operations
- `notion_search`: Search across pages and databases

### Example Usage

Here's an example of how to use the plugin to create a new database item:

```json
{
  "name": "notion_create_database_item",
  "arguments": {
    "database_id": "your-database-id",
    "properties": {
      "Name": {
        "title": [
          {
            "text": {
              "content": "New Item"
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
}
```

