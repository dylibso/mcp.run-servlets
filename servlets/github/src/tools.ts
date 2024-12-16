import { Tool } from "@modelcontextprotocol/sdk/types.js";
import { zodToJsonSchema } from 'zod-to-json-schema';
import {
  CreateOrUpdateFileSchema,
  SearchRepositoriesSchema,
  CreateRepositorySchema,
  GetFileContentsSchema,
  PushFilesSchema,
  CreateIssueSchema,
  CreatePullRequestSchema,
  ForkRepositorySchema,
  CreateBranchSchema,
  ListCommitsSchema,
  ListIssuesOptionsSchema,
  UpdateIssueOptionsSchema,
  IssueCommentSchema,
  SearchCodeSchema,
  SearchIssuesSchema,
  SearchUsersSchema,
  GetIssueSchema
} from './types.js';

export const CREATE_UPDATE_FILE_TOOL: Tool = {
  name: "create_or_update_file",
  description: "Create or update a single file in a GitHub repository",
  inputSchema: zodToJsonSchema(CreateOrUpdateFileSchema)
};

export const SEARCH_REPOSITORIES_TOOL: Tool = {
  name: "search_repositories",
  description: "Search for GitHub repositories",
  inputSchema: zodToJsonSchema(SearchRepositoriesSchema)
};

export const CREATE_REPOSITORY_TOOL: Tool = {
  name: "create_repository",
  description: "Create a new GitHub repository in your account",
  inputSchema: zodToJsonSchema(CreateRepositorySchema)
};

export const GET_FILE_CONTENTS_TOOL: Tool = {
  name: "get_file_contents",
  description: "Get the contents of a file or directory from a GitHub repository",
  inputSchema: zodToJsonSchema(GetFileContentsSchema)
};

export const PUSH_FILES_TOOL: Tool = {
  name: "push_files",
  description: "Push multiple files to a GitHub repository in a single commit",
  inputSchema: zodToJsonSchema(PushFilesSchema)
};

export const CREATE_ISSUE_TOOL: Tool = {
  name: "create_issue",
  description: "Create a new issue in a GitHub repository",
  inputSchema: zodToJsonSchema(CreateIssueSchema)
};

export const CREATE_PULL_REQUEST_TOOL: Tool = {
  name: "create_pull_request",
  description: "Create a new pull request in a GitHub repository",
  inputSchema: zodToJsonSchema(CreatePullRequestSchema)
};

export const FORK_REPOSITORY_TOOL: Tool = {
  name: "fork_repository",
  description: "Fork a GitHub repository to your account or specified organization",
  inputSchema: zodToJsonSchema(ForkRepositorySchema)
};

export const CREATE_BRANCH_TOOL: Tool = {
  name: "create_branch",
  description: "Create a new branch in a GitHub repository",
  inputSchema: zodToJsonSchema(CreateBranchSchema)
};

export const LIST_COMMITS_TOOL: Tool = {
  name: "list_commits",
  description: "Get list of commits of a branch in a GitHub repository",
  inputSchema: zodToJsonSchema(ListCommitsSchema)
};

export const LIST_ISSUES_TOOL: Tool = {
  name: "list_issues",
  description: "List issues in a GitHub repository with filtering options",
  inputSchema: zodToJsonSchema(ListIssuesOptionsSchema)
};

export const UPDATE_ISSUE_TOOL: Tool = {
  name: "update_issue",
  description: "Update an existing issue in a GitHub repository",
  inputSchema: zodToJsonSchema(UpdateIssueOptionsSchema)
};

export const ADD_ISSUE_COMMENT_TOOL: Tool = {
  name: "add_issue_comment",
  description: "Add a comment to an existing issue",
  inputSchema: zodToJsonSchema(IssueCommentSchema)
};

export const SEARCH_CODE_TOOL: Tool = {
  name: "search_code",
  description: "Search for code across GitHub repositories",
  inputSchema: zodToJsonSchema(SearchCodeSchema)
};

export const SEARCH_ISSUES_TOOL: Tool = {
  name: "search_issues",
  description: "Search for issues and pull requests across GitHub repositories",
  inputSchema: zodToJsonSchema(SearchIssuesSchema)
};

export const SEARCH_USERS_TOOL: Tool = {
  name: "search_users",
  description: "Search for users on GitHub",
  inputSchema: zodToJsonSchema(SearchUsersSchema)
};

export const GET_ISSUE_TOOL: Tool = {
  name: "get_issue",
  description: "Get details of a specific issue in a GitHub repository",
  inputSchema: zodToJsonSchema(GetIssueSchema)
};

export const GITHUB_TOOLS = [
  CREATE_UPDATE_FILE_TOOL,
  SEARCH_REPOSITORIES_TOOL,
  CREATE_REPOSITORY_TOOL,
  GET_FILE_CONTENTS_TOOL,
  PUSH_FILES_TOOL,
  CREATE_ISSUE_TOOL,
  CREATE_PULL_REQUEST_TOOL,
  FORK_REPOSITORY_TOOL,
  CREATE_BRANCH_TOOL,
  LIST_COMMITS_TOOL,
  LIST_ISSUES_TOOL,
  UPDATE_ISSUE_TOOL,
  ADD_ISSUE_COMMENT_TOOL,
  SEARCH_CODE_TOOL,
  SEARCH_ISSUES_TOOL,
  SEARCH_USERS_TOOL,
  GET_ISSUE_TOOL
] as const;