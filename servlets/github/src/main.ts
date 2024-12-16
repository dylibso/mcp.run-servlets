import { CallToolRequest, CallToolResult, ContentType, ListToolsResult } from "./pdk";
import { GITHUB_TOOLS } from "./tools";
import { z } from 'zod';
import {
  CreateBranchSchema,
  CreateIssueSchema,
  CreateOrUpdateFileSchema,
  CreatePullRequestSchema,
  CreateRepositorySchema,
  ForkRepositorySchema,
  GetFileContentsSchema,
  GetIssueSchema,
  GitHubContentSchema,
  GitHubForkSchema,
  GitHubIssueSchema,
  GitHubPullRequestSchema,
  GitHubReferenceSchema,
  IssueCommentSchema,
  ListCommitsSchema,
  ListIssuesOptionsSchema,
  PushFilesSchema,
  SearchCodeSchema,
  SearchIssuesSchema,
  SearchRepositoriesSchema,
  SearchUsersSchema,
  UpdateIssueOptionsSchema
} from './types';

// Helper functions remain mostly unchanged, just moved to the top of the file
async function forkRepository(owner: string, repo: string, organization?: string) {
  // Implementation remains the same
}

// [Include all other helper functions from the original file]
// ... getDefaultBranchSHA, getFileContents, createIssue, etc.

/**
 * Called by mcpx to understand how and why to use this tool
 */
export function describeImpl(): ListToolsResult {
  return {
    tools: GITHUB_TOOLS
  };
}

/**
 * Called when the tool is invoked
 */
export function callImpl(input: CallToolRequest): CallToolResult {
  const apiKey = Config.get("github-token");
  if (!apiKey) {
    throw new Error("GitHub token not configured");
  }

  try {
    switch (input.params.name) {
      case "fork_repository": {
        const args = ForkRepositorySchema.parse(input.params.arguments);
        const fork = await forkRepository(args.owner, args.repo, args.organization);
        return {
          content: [{ type: ContentType.Text, text: JSON.stringify(fork, null, 2) }]
        };
      }

      case "create_branch": {
        const args = CreateBranchSchema.parse(input.params.arguments);
        let sha: string;
        if (args.from_branch) {
          // Implementation continues as in original
        } else {
          sha = await getDefaultBranchSHA(args.owner, args.repo);
        }
        const branch = await createBranch(args.owner, args.repo, {
          ref: args.branch,
          sha
        });
        return {
          content: [{ type: ContentType.Text, text: JSON.stringify(branch, null, 2) }]
        };
      }

      // [Include all other case handlers from the original switch statement]
      // The implementation of each case remains the same, just reformatted for MCPX

      default:
        throw new Error(`Unknown tool: ${input.params.name}`);
    }
  } catch (error) {
    if (error instanceof z.ZodError) {
      throw new Error(
        `Invalid arguments: ${error.errors
          .map((e) => `${e.path.join(".")}: ${e.message}`)
          .join(", ")}`
      );
    }
    throw error;
  }
}