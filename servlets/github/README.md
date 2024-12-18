# GitHub Servlet

Currently supported: 

Issues: 

- `gh-create-issue` Create issue
- `gh-get-issue` Get issue
- `gh-update-issue` Add comments to an issue
- `gh-add-issue-comment` Read the contents of an issue

Files:

- `gh-get-file-contents` Get contents and metadata about a file on a branch
- `gh-create-or-update-file` Create or a update a file on a branch
- `gh-push-files` Bulk push files to a branch

Branches: 

- `gh-create-branch` Create a new branch
- `gh-create-pull-request` Create a PR from a branch



## Config

Requires the following config keys:

- **api-key** with a [GitHub Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

## Permissions

Requires access to the following **domains**:

- `api.github.com`

## Example

— can you add a funny comment on issue #1 from repo evacchi/mcpx-target-playground ?

<blockquote>
I'll add a playful comment to keep with the "don't panic" theme of the issue.

I added a comment with a reference to The Hitchhiker's Guide to the Galaxy, since it fits well with the "don't panic" theme of the issue! The comment has been successfully posted.
</blockquote>

