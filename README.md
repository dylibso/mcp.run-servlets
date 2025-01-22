# Official mcp.run Servlets

These are the official servlets for the @dylibso account on [mcp.run](https://mcp.run).
Merging changes here automatically deploys them to all users.

## Generating Servlets with LLMs 

Using Cursor, Cline, or other AI-enabled IDEs have their own ways of integrating source context - use this repo as a starting point to give them examples. All the code in the `servlets` directory is great fine-tuning material. 

If you are using a Chat interface, you can copy/paste context generated from here: 
https://gitingest.com/dylibso/mcp.run-servlets

**Note:** to reduce context used, add these patterns to the "Exclude" dropdown so they are not tokenized:

```
**/package*.json, **/go.mod, **/go.sum, **/go.work*, **/LICENSE, LICENSE, .gitignore, **/.gitignore
```

