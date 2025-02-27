# BlueSky Servlet

Implements a reduced subset of the BlueSky OpenAPI spec.
Currently it supports [posting and replying to a post][post] and [getting threads][thread].

Additionally, it provides the `latest_mentions` tool, returning a list
a list of mentions to the given handle within a specified amount of time 
(e.g. "5 minutes ago", `5m`); this tool uses the [search][search] feature
under the hood.

[post]: https://docs.bsky.app/docs/api/com-atproto-repo-put-record
[search]: https://docs.bsky.app/docs/api/app-bsky-feed-search-posts
[thread]: https://docs.bsky.app/docs/api/app-bsky-feed-get-post-thread

## Config

- `BASE_URL` defaults to `https://bsky.social`
- `HANDLE` your handle
- `PASSWORD` use an [app password](https://bsky.app/settings/app-passwords)

## Domains

You should grant access to `bsky.social`.

## Example

> You are a helpful bot called mcprunbot.bsky.social
> You always use your tools if possible.
>
> Find all your recent mentions, within 5 minutes from now. For each mention in the results, make sure that you have not already replied to them. If you have not replied to a mention, then make sure you reply.
>
> Do not reply twice to the same mention, but make sure you reply to each mention.


> I have checked the recent mentions of `@mcprunbot.bsky.social` and verified the replies:
>
> 1. **Mention 1**: Request to search for a tutorial to build a bot using mcp.run.
>    - **Status**: Already replied with a link to a tutorial on creating a Telegram bot using mcp.run.
>
> 2. **Mention 2**: Request to tell a joke to `@nilslice.bsky.social`.
>    - **Status**: Already replied with a joke: "Why don't scientists trust atoms? Because they make up everything! 😄"
>
> 3. **Mention 3**: Request to tell a little about mcp.run.
>    - **Status**: Already replied with a brief explanation of mcp.run and a link to the site.
>
> All recent mentions have been replied to, and no further action is needed.
