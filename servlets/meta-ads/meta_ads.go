package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
)

// Tool definitions
var (
	// Account Management Tools
	GetBusinessesTool = ToolDescription{
		Name:        "get-businesses",
		Description: "Get all businesses associated with the user account",
		InputSchema: schema{
			"type":       "object",
			"properties": props{},
		},
	}

	GetAdAccountsTool = ToolDescription{
		Name:        "get-ad-accounts",
		Description: "Get all ad accounts accessible to the user",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"business_id": prop("string", "Optional business ID to filter ad accounts"),
			},
		},
	}

	GetAccountInfoTool = ToolDescription{
		Name:        "get-account-info",
		Description: "Get detailed information about a specific ad account",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"account_id": prop("string", "The ad account ID (format: act_123456789)"),
			},
			"required": []string{"account_id"},
		},
	}

	// Campaign Management Tools
	GetCampaignsTool = ToolDescription{
		Name:        "get-campaigns",
		Description: "Get campaigns for a specific ad account",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"account_id": prop("string", "The ad account ID (format: act_123456789)"),
				"status":     prop("string", "Optional status filter (ACTIVE, PAUSED, ARCHIVED)"),
				"limit":      prop("integer", "Optional limit for number of campaigns to return (default: 25)"),
			},
			"required": []string{"account_id"},
		},
	}

	CreateCampaignTool = ToolDescription{
		Name:        "create-campaign",
		Description: "Create a new campaign",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"account_id": prop("string", "The ad account ID (format: act_123456789)"),
				"name":       prop("string", "Campaign name"),
				"objective":  prop("string", "Campaign objective (e.g., LINK_CLICKS, CONVERSIONS, BRAND_AWARENESS)"),
				"status":     prop("string", "Campaign status (ACTIVE, PAUSED) - defaults to PAUSED"),
				"special_ad_categories": propWithItems("array", "Special ad categories if applicable", schema{
					"type": "string",
					"enum": []string{"CREDIT", "EMPLOYMENT", "HOUSING", "ISSUES_ELECTIONS_POLITICS"},
				}),
			},
			"required": []string{"account_id", "name", "objective"},
		},
	}

	UpdateCampaignTool = ToolDescription{
		Name:        "update-campaign",
		Description: "Update an existing campaign",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"campaign_id": prop("string", "The campaign ID"),
				"name":        prop("string", "Optional new campaign name"),
				"status":      prop("string", "Optional new status (ACTIVE, PAUSED, ARCHIVED)"),
			},
			"required": []string{"campaign_id"},
		},
	}

	DeleteCampaignTool = ToolDescription{
		Name:        "delete-campaign",
		Description: "Delete a campaign (moves it to archived state)",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"campaign_id": prop("string", "The campaign ID"),
			},
			"required": []string{"campaign_id"},
		},
	}

	GetCampaignInfoTool = ToolDescription{
		Name:        "get-campaign-info",
		Description: "Get detailed information about a specific campaign",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"campaign_id": prop("string", "The campaign ID"),
			},
			"required": []string{"campaign_id"},
		},
	}

	// Ad Set Management Tools
	GetAdSetsTool = ToolDescription{
		Name:        "get-adsets",
		Description: "Get ad sets for a campaign or account",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"campaign_id": prop("string", "Optional campaign ID to filter ad sets"),
				"account_id":  prop("string", "Optional account ID to get all ad sets"),
				"status":      prop("string", "Optional status filter (ACTIVE, PAUSED, ARCHIVED)"),
				"limit":       prop("integer", "Optional limit for number of ad sets to return (default: 25)"),
			},
		},
	}

	CreateAdSetTool = ToolDescription{
		Name:        "create-adset",
		Description: "Create a new ad set",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"campaign_id":       prop("string", "The campaign ID"),
				"name":              prop("string", "Ad set name"),
				"daily_budget":      prop("integer", "Daily budget in cents (e.g., 1000 = $10.00)"),
				"lifetime_budget":   prop("integer", "Optional lifetime budget in cents"),
				"targeting":         prop("object", "Targeting specification as JSON object"),
				"billing_event":     prop("string", "Billing event (IMPRESSIONS, CLICKS, etc.) - defaults to IMPRESSIONS"),
				"optimization_goal": prop("string", "Optimization goal - defaults to LINK_CLICKS"),
				"status":            prop("string", "Ad set status (ACTIVE, PAUSED) - defaults to PAUSED"),
			},
			"required": []string{"campaign_id", "name", "targeting"},
		},
	}

	UpdateAdSetTool = ToolDescription{
		Name:        "update-adset",
		Description: "Update an existing ad set",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"adset_id":        prop("string", "The ad set ID"),
				"name":            prop("string", "Optional new ad set name"),
				"daily_budget":    prop("integer", "Optional new daily budget in cents"),
				"lifetime_budget": prop("integer", "Optional new lifetime budget in cents"),
				"targeting":       prop("object", "Optional new targeting specification"),
				"status":          prop("string", "Optional new status (ACTIVE, PAUSED, ARCHIVED)"),
			},
			"required": []string{"adset_id"},
		},
	}

	DeleteAdSetTool = ToolDescription{
		Name:        "delete-adset",
		Description: "Delete an ad set (moves it to archived state)",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"adset_id": prop("string", "The ad set ID"),
			},
			"required": []string{"adset_id"},
		},
	}

	GetAdSetInfoTool = ToolDescription{
		Name:        "get-adset-info",
		Description: "Get detailed information about a specific ad set",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"adset_id": prop("string", "The ad set ID"),
			},
			"required": []string{"adset_id"},
		},
	}

	// Ad Management Tools
	GetAdsTool = ToolDescription{
		Name:        "get-ads",
		Description: "Get ads for an ad set, campaign, or account",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"adset_id":    prop("string", "Optional ad set ID to filter ads"),
				"campaign_id": prop("string", "Optional campaign ID to filter ads"),
				"account_id":  prop("string", "Optional account ID to get all ads"),
				"status":      prop("string", "Optional status filter (ACTIVE, PAUSED, ARCHIVED)"),
				"limit":       prop("integer", "Optional limit for number of ads to return (default: 25)"),
			},
		},
	}

	CreateAdTool = ToolDescription{
		Name:        "create-ad",
		Description: "Create a new ad",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"adset_id":    prop("string", "The ad set ID"),
				"name":        prop("string", "Ad name"),
				"creative_id": prop("string", "Creative ID to use for this ad"),
				"status":      prop("string", "Ad status (ACTIVE, PAUSED) - defaults to PAUSED"),
			},
			"required": []string{"adset_id", "name", "creative_id"},
		},
	}

	UpdateAdTool = ToolDescription{
		Name:        "update-ad",
		Description: "Update an existing ad",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"ad_id":       prop("string", "The ad ID"),
				"name":        prop("string", "Optional new ad name"),
				"creative_id": prop("string", "Optional new creative ID"),
				"status":      prop("string", "Optional new status (ACTIVE, PAUSED, ARCHIVED)"),
			},
			"required": []string{"ad_id"},
		},
	}

	DeleteAdTool = ToolDescription{
		Name:        "delete-ad",
		Description: "Delete an ad (moves it to archived state)",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"ad_id": prop("string", "The ad ID"),
			},
			"required": []string{"ad_id"},
		},
	}

	GetAdInfoTool = ToolDescription{
		Name:        "get-ad-info",
		Description: "Get detailed information about a specific ad",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"ad_id": prop("string", "The ad ID"),
			},
			"required": []string{"ad_id"},
		},
	}

	// Insights Tools
	GetCampaignInsightsTool = ToolDescription{
		Name:        "get-campaign-insights",
		Description: "Get performance insights for campaigns",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"campaign_ids": propWithItems("array", "Array of campaign IDs", schema{"type": "string"}),
				"account_id":   prop("string", "Account ID to get insights for all campaigns"),
				"date_preset":  prop("string", "Date preset (today, yesterday, last_7d, last_30d, etc.)"),
				"time_range":   prop("object", "Custom time range with since and until dates"),
				"metrics": propWithItems("array", "Metrics to retrieve", schema{
					"type": "string",
					"enum": []string{"impressions", "clicks", "spend", "reach", "frequency", "ctr", "cpc", "cpm", "cpp", "cost_per_result", "conversions", "conversion_rate"},
				}),
			},
		},
	}

	GetAdSetInsightsTool = ToolDescription{
		Name:        "get-adset-insights",
		Description: "Get performance insights for ad sets",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"adset_ids":   propWithItems("array", "Array of ad set IDs", schema{"type": "string"}),
				"campaign_id": prop("string", "Campaign ID to get insights for all ad sets"),
				"account_id":  prop("string", "Account ID to get insights for all ad sets"),
				"date_preset": prop("string", "Date preset (today, yesterday, last_7d, last_30d, etc.)"),
				"time_range":  prop("object", "Custom time range with since and until dates"),
				"metrics": propWithItems("array", "Metrics to retrieve", schema{
					"type": "string",
					"enum": []string{"impressions", "clicks", "spend", "reach", "frequency", "ctr", "cpc", "cpm", "cpp", "cost_per_result", "conversions", "conversion_rate"},
				}),
			},
		},
	}

	GetAdInsightsTool = ToolDescription{
		Name:        "get-ad-insights",
		Description: "Get performance insights for ads",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"ad_ids":      propWithItems("array", "Array of ad IDs", schema{"type": "string"}),
				"adset_id":    prop("string", "Ad set ID to get insights for all ads"),
				"campaign_id": prop("string", "Campaign ID to get insights for all ads"),
				"account_id":  prop("string", "Account ID to get insights for all ads"),
				"date_preset": prop("string", "Date preset (today, yesterday, last_7d, last_30d, etc.)"),
				"time_range":  prop("object", "Custom time range with since and until dates"),
				"metrics": propWithItems("array", "Metrics to retrieve", schema{
					"type": "string",
					"enum": []string{"impressions", "clicks", "spend", "reach", "frequency", "ctr", "cpc", "cpm", "cpp", "cost_per_result", "conversions", "conversion_rate"},
				}),
			},
		},
	}

	// Audience Management Tools
	GetAudiencesTool = ToolDescription{
		Name:        "get-audiences",
		Description: "Get custom audiences for an ad account",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"account_id": prop("string", "The ad account ID (format: act_123456789)"),
				"limit":      prop("integer", "Optional limit for number of audiences to return (default: 25)"),
			},
			"required": []string{"account_id"},
		},
	}

	CreateAudienceTool = ToolDescription{
		Name:        "create-audience",
		Description: "Create a custom audience",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"account_id":  prop("string", "The ad account ID (format: act_123456789)"),
				"name":        prop("string", "Audience name"),
				"description": prop("string", "Optional audience description"),
				"subtype":     prop("string", "Audience subtype (CUSTOM, LOOKALIKE, WEBSITE, etc.)"),
			},
			"required": []string{"account_id", "name", "subtype"},
		},
	}

	// Creative Management Tools
	GetCreativesTool = ToolDescription{
		Name:        "get-creatives",
		Description: "Get creatives for an ad account",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"account_id": prop("string", "The ad account ID (format: act_123456789)"),
				"limit":      prop("integer", "Optional limit for number of creatives to return (default: 25)"),
			},
			"required": []string{"account_id"},
		},
	}

	CreateCreativeTool = ToolDescription{
		Name:        "create-creative",
		Description: "Create an ad creative",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"account_id":              prop("string", "The ad account ID (format: act_123456789)"),
				"name":                    prop("string", "Creative name"),
				"object_story_spec":       prop("object", "Object story specification for the creative"),
				"degrees_of_freedom_spec": prop("object", "Optional degrees of freedom specification"),
			},
			"required": []string{"account_id", "name", "object_story_spec"},
		},
	}

	// Bulk Operations
	BulkUpdateTool = ToolDescription{
		Name:        "bulk-update",
		Description: "Perform bulk updates on multiple campaigns, ad sets, or ads",
		InputSchema: schema{
			"type": "object",
			"properties": props{
				"operations": propWithItems("array", "Array of update operations", schema{
					"type": "object",
					"properties": props{
						"object_type": prop("string", "Type of object (campaign, adset, ad)"),
						"object_id":   prop("string", "ID of the object to update"),
						"updates":     prop("object", "Fields to update with their new values"),
					},
					"required": []string{"object_type", "object_id", "updates"},
				}),
			},
			"required": []string{"operations"},
		},
	}

	// All tools collection
	MetaAdsTools = []ToolDescription{
		GetBusinessesTool,
		GetAdAccountsTool,
		GetAccountInfoTool,
		GetCampaignsTool,
		CreateCampaignTool,
		UpdateCampaignTool,
		DeleteCampaignTool,
		GetCampaignInfoTool,
		GetAdSetsTool,
		CreateAdSetTool,
		UpdateAdSetTool,
		DeleteAdSetTool,
		GetAdSetInfoTool,
		GetAdsTool,
		CreateAdTool,
		UpdateAdTool,
		DeleteAdTool,
		GetAdInfoTool,
		GetCampaignInsightsTool,
		GetAdSetInsightsTool,
		GetAdInsightsTool,
		GetAudiencesTool,
		CreateAudienceTool,
		GetCreativesTool,
		CreateCreativeTool,
		BulkUpdateTool,
	}

	MetaAdsReadOnlyTools = []ToolDescription{
		GetBusinessesTool,
		GetAdAccountsTool,
		GetAccountInfoTool,
		GetCampaignsTool,

		GetCampaignInfoTool,
		GetAdSetsTool,

		GetAdSetInfoTool,
		GetAdsTool,

		GetAdInfoTool,
		GetCampaignInsightsTool,
		GetAdSetInsightsTool,
		GetAdInsightsTool,
		GetAudiencesTool,

		GetCreativesTool,

		BulkUpdateTool,
	}
)

// Account Management Functions
func getBusinesses(token string, args map[string]any) (CallToolResult, error) {
	client := NewMetaClient(token)

	resp, err := client.MakeRequest(pdk.MethodGet, "me/businesses", nil, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get businesses: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func getAdAccounts(token string, args map[string]any) (CallToolResult, error) {
	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,account_status,currency,timezone_name,business"},
	}

	if businessId, ok := args["business_id"].(string); ok && businessId != "" {
		resp, err := client.MakeRequest(pdk.MethodGet, fmt.Sprintf("%s/adaccounts", businessId), query, nil)
		if err != nil {
			return CallToolResult{
				IsError: some(true),
				Content: []Content{{
					Type: ContentTypeText,
					Text: some(fmt.Sprintf("Failed to get ad accounts: %s", err)),
				}},
			}, nil
		}
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(resp.Body)),
			}},
		}, nil
	}

	resp, err := client.MakeRequest(pdk.MethodGet, "me/adaccounts", query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ad accounts: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func getAccountInfo(token string, args map[string]any) (CallToolResult, error) {
	accountId, ok := args["account_id"].(string)
	if !ok || accountId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Account ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,account_status,currency,timezone_name,business,spend_cap,balance,amount_spent,created_time"},
	}

	resp, err := client.MakeRequest(pdk.MethodGet, accountId, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get account info: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

// Campaign Management Functions
func getCampaigns(token string, args map[string]any) (CallToolResult, error) {
	accountId, ok := args["account_id"].(string)
	if !ok || accountId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Account ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,objective,status,created_time,updated_time,start_time,stop_time,daily_budget,lifetime_budget"},
	}

	if status, ok := args["status"].(string); ok && status != "" {
		query["filtering"] = []string{fmt.Sprintf(`[{"field":"status","operator":"EQUAL","value":"%s"}]`, status)}
	}

	if limit, ok := args["limit"].(float64); ok && limit > 0 {
		query["limit"] = []string{fmt.Sprintf("%.0f", limit)}
	}

	resp, err := client.MakeRequest(pdk.MethodGet, fmt.Sprintf("%s/campaigns", accountId), query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get campaigns: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func createCampaign(token string, args map[string]any) (CallToolResult, error) {
	accountId, ok := args["account_id"].(string)
	if !ok || accountId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Account ID is required"),
			}},
		}, nil
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Campaign name is required"),
			}},
		}, nil
	}

	objective, ok := args["objective"].(string)
	if !ok || objective == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Campaign objective is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	formData := map[string]string{
		"name":      name,
		"objective": objective,
		"status":    "PAUSED", // Default to paused for safety
	}

	if status, ok := args["status"].(string); ok && status != "" {
		formData["status"] = status
	}

	if categories, ok := args["special_ad_categories"].([]any); ok && len(categories) > 0 {
		var categoryStrs []string
		for _, cat := range categories {
			if catStr, ok := cat.(string); ok {
				categoryStrs = append(categoryStrs, catStr)
			}
		}
		if len(categoryStrs) > 0 {
			categoriesJson, _ := json.Marshal(categoryStrs)
			formData["special_ad_categories"] = string(categoriesJson)
		}
	}

	resp, err := client.MakeFormRequest(fmt.Sprintf("%s/campaigns", accountId), formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create campaign: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func updateCampaign(token string, args map[string]any) (CallToolResult, error) {
	campaignId, ok := args["campaign_id"].(string)
	if !ok || campaignId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Campaign ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)
	formData := map[string]string{}

	if name, ok := args["name"].(string); ok && name != "" {
		formData["name"] = name
	}

	if status, ok := args["status"].(string); ok && status != "" {
		formData["status"] = status
	}

	if len(formData) == 0 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("At least one field to update is required"),
			}},
		}, nil
	}

	resp, err := client.MakeFormRequest(campaignId, formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to update campaign: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func deleteCampaign(token string, args map[string]any) (CallToolResult, error) {
	campaignId, ok := args["campaign_id"].(string)
	if !ok || campaignId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Campaign ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	resp, err := client.MakeRequest(pdk.MethodDelete, campaignId, nil, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to delete campaign: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func getCampaignInfo(token string, args map[string]any) (CallToolResult, error) {
	campaignId, ok := args["campaign_id"].(string)
	if !ok || campaignId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Campaign ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,objective,status,created_time,updated_time,start_time,stop_time,daily_budget,lifetime_budget,account_id,special_ad_categories"},
	}

	resp, err := client.MakeRequest(pdk.MethodGet, campaignId, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get campaign info: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

// Ad Set Management Functions
func getAdSets(token string, args map[string]any) (CallToolResult, error) {
	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,status,created_time,updated_time,daily_budget,lifetime_budget,billing_event,optimization_goal,targeting"},
	}

	var endpoint string
	if campaignId, ok := args["campaign_id"].(string); ok && campaignId != "" {
		endpoint = fmt.Sprintf("%s/adsets", campaignId)
	} else if accountId, ok := args["account_id"].(string); ok && accountId != "" {
		endpoint = fmt.Sprintf("%s/adsets", accountId)
	} else {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Either campaign_id or account_id is required"),
			}},
		}, nil
	}

	if status, ok := args["status"].(string); ok && status != "" {
		query["filtering"] = []string{fmt.Sprintf(`[{"field":"status","operator":"EQUAL","value":"%s"}]`, status)}
	}

	if limit, ok := args["limit"].(float64); ok && limit > 0 {
		query["limit"] = []string{fmt.Sprintf("%.0f", limit)}
	}

	resp, err := client.MakeRequest(pdk.MethodGet, endpoint, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ad sets: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func createAdSet(token string, args map[string]any) (CallToolResult, error) {
	campaignId, ok := args["campaign_id"].(string)
	if !ok || campaignId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Campaign ID is required"),
			}},
		}, nil
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad set name is required"),
			}},
		}, nil
	}

	targeting, ok := args["targeting"].(map[string]any)
	if !ok || len(targeting) == 0 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Targeting is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	formData := map[string]string{
		"name":              name,
		"campaign_id":       campaignId,
		"billing_event":     "IMPRESSIONS",
		"optimization_goal": "LINK_CLICKS",
		"status":            "PAUSED", // Default to paused
	}

	if dailyBudget, ok := args["daily_budget"].(float64); ok && dailyBudget > 0 {
		formData["daily_budget"] = fmt.Sprintf("%.0f", dailyBudget)
	}

	if lifetimeBudget, ok := args["lifetime_budget"].(float64); ok && lifetimeBudget > 0 {
		formData["lifetime_budget"] = fmt.Sprintf("%.0f", lifetimeBudget)
	}

	if billingEvent, ok := args["billing_event"].(string); ok && billingEvent != "" {
		formData["billing_event"] = billingEvent
	}

	if optimizationGoal, ok := args["optimization_goal"].(string); ok && optimizationGoal != "" {
		formData["optimization_goal"] = optimizationGoal
	}

	if status, ok := args["status"].(string); ok && status != "" {
		formData["status"] = status
	}

	// Convert targeting to JSON string
	targetingJson, err := json.Marshal(targeting)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal targeting: %s", err)),
			}},
		}, nil
	}
	formData["targeting"] = string(targetingJson)

	// Get account ID from campaign
	campaignInfo, err := client.MakeRequest(pdk.MethodGet, campaignId, map[string][]string{"fields": {"account_id"}}, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get campaign info: %s", err)),
			}},
		}, nil
	}

	var campaignData map[string]any
	if err := json.Unmarshal(campaignInfo.Body, &campaignData); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to parse campaign data: %s", err)),
			}},
		}, nil
	}

	accountId, ok := campaignData["account_id"].(string)
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Failed to get account ID from campaign"),
			}},
		}, nil
	}

	resp, err := client.MakeFormRequest(fmt.Sprintf("%s/adsets", accountId), formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create ad set: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func updateAdSet(token string, args map[string]any) (CallToolResult, error) {
	adsetId, ok := args["adset_id"].(string)
	if !ok || adsetId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad set ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)
	formData := map[string]string{}

	if name, ok := args["name"].(string); ok && name != "" {
		formData["name"] = name
	}

	if dailyBudget, ok := args["daily_budget"].(float64); ok && dailyBudget > 0 {
		formData["daily_budget"] = fmt.Sprintf("%.0f", dailyBudget)
	}

	if lifetimeBudget, ok := args["lifetime_budget"].(float64); ok && lifetimeBudget > 0 {
		formData["lifetime_budget"] = fmt.Sprintf("%.0f", lifetimeBudget)
	}

	if targeting, ok := args["targeting"].(map[string]any); ok && len(targeting) > 0 {
		targetingJson, err := json.Marshal(targeting)
		if err != nil {
			return CallToolResult{
				IsError: some(true),
				Content: []Content{{
					Type: ContentTypeText,
					Text: some(fmt.Sprintf("Failed to marshal targeting: %s", err)),
				}},
			}, nil
		}
		formData["targeting"] = string(targetingJson)
	}

	if status, ok := args["status"].(string); ok && status != "" {
		formData["status"] = status
	}

	if len(formData) == 0 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("At least one field to update is required"),
			}},
		}, nil
	}

	resp, err := client.MakeFormRequest(adsetId, formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to update ad set: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func deleteAdSet(token string, args map[string]any) (CallToolResult, error) {
	adsetId, ok := args["adset_id"].(string)
	if !ok || adsetId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad set ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	resp, err := client.MakeRequest(pdk.MethodDelete, adsetId, nil, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to delete ad set: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func getAdSetInfo(token string, args map[string]any) (CallToolResult, error) {
	adsetId, ok := args["adset_id"].(string)
	if !ok || adsetId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad set ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,status,created_time,updated_time,daily_budget,lifetime_budget,billing_event,optimization_goal,targeting,campaign_id,account_id"},
	}

	resp, err := client.MakeRequest(pdk.MethodGet, adsetId, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ad set info: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

// Ad Management Functions
func getAds(token string, args map[string]any) (CallToolResult, error) {
	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,status,created_time,updated_time,creative,adset_id,campaign_id,account_id"},
	}

	var endpoint string
	if adsetId, ok := args["adset_id"].(string); ok && adsetId != "" {
		endpoint = fmt.Sprintf("%s/ads", adsetId)
	} else if campaignId, ok := args["campaign_id"].(string); ok && campaignId != "" {
		endpoint = fmt.Sprintf("%s/ads", campaignId)
	} else if accountId, ok := args["account_id"].(string); ok && accountId != "" {
		endpoint = fmt.Sprintf("%s/ads", accountId)
	} else {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Either adset_id, campaign_id, or account_id is required"),
			}},
		}, nil
	}

	if status, ok := args["status"].(string); ok && status != "" {
		query["filtering"] = []string{fmt.Sprintf(`[{"field":"status","operator":"EQUAL","value":"%s"}]`, status)}
	}

	if limit, ok := args["limit"].(float64); ok && limit > 0 {
		query["limit"] = []string{fmt.Sprintf("%.0f", limit)}
	}

	resp, err := client.MakeRequest(pdk.MethodGet, endpoint, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ads: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func createAd(token string, args map[string]any) (CallToolResult, error) {
	adsetId, ok := args["adset_id"].(string)
	if !ok || adsetId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad set ID is required"),
			}},
		}, nil
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad name is required"),
			}},
		}, nil
	}

	creativeId, ok := args["creative_id"].(string)
	if !ok || creativeId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Creative ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	formData := map[string]string{
		"name":     name,
		"adset_id": adsetId,
		"creative": fmt.Sprintf(`{"creative_id":"%s"}`, creativeId),
		"status":   "PAUSED", // Default to paused
	}

	if status, ok := args["status"].(string); ok && status != "" {
		formData["status"] = status
	}

	// Get account ID from ad set
	adsetInfo, err := client.MakeRequest(pdk.MethodGet, adsetId, map[string][]string{"fields": {"account_id"}}, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ad set info: %s", err)),
			}},
		}, nil
	}

	var adsetData map[string]any
	if err := json.Unmarshal(adsetInfo.Body, &adsetData); err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to parse ad set data: %s", err)),
			}},
		}, nil
	}

	accountId, ok := adsetData["account_id"].(string)
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Failed to get account ID from ad set"),
			}},
		}, nil
	}

	resp, err := client.MakeFormRequest(fmt.Sprintf("%s/ads", accountId), formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create ad: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func updateAd(token string, args map[string]any) (CallToolResult, error) {
	adId, ok := args["ad_id"].(string)
	if !ok || adId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)
	formData := map[string]string{}

	if name, ok := args["name"].(string); ok && name != "" {
		formData["name"] = name
	}

	if creativeId, ok := args["creative_id"].(string); ok && creativeId != "" {
		formData["creative"] = fmt.Sprintf(`{"creative_id":"%s"}`, creativeId)
	}

	if status, ok := args["status"].(string); ok && status != "" {
		formData["status"] = status
	}

	if len(formData) == 0 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("At least one field to update is required"),
			}},
		}, nil
	}

	resp, err := client.MakeFormRequest(adId, formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to update ad: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func deleteAd(token string, args map[string]any) (CallToolResult, error) {
	adId, ok := args["ad_id"].(string)
	if !ok || adId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	resp, err := client.MakeRequest(pdk.MethodDelete, adId, nil, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to delete ad: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func getAdInfo(token string, args map[string]any) (CallToolResult, error) {
	adId, ok := args["ad_id"].(string)
	if !ok || adId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Ad ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,status,created_time,updated_time,creative,adset_id,campaign_id,account_id"},
	}

	resp, err := client.MakeRequest(pdk.MethodGet, adId, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ad info: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

// Insights Functions
func getCampaignInsights(token string, args map[string]any) (CallToolResult, error) {
	client := NewMetaClient(token)

	query := map[string][]string{}

	// Set default metrics if none provided
	metrics := []string{"impressions", "clicks", "spend", "reach", "ctr", "cpc", "cpm"}
	if providedMetrics, ok := args["metrics"].([]any); ok && len(providedMetrics) > 0 {
		metrics = []string{}
		for _, m := range providedMetrics {
			if metricStr, ok := m.(string); ok {
				metrics = append(metrics, metricStr)
			}
		}
	}
	query["fields"] = []string{strings.Join(metrics, ",")}

	// Handle date range
	if datePreset, ok := args["date_preset"].(string); ok && datePreset != "" {
		query["date_preset"] = []string{datePreset}
	} else if timeRange, ok := args["time_range"].(map[string]any); ok {
		if since, ok := timeRange["since"].(string); ok {
			if until, ok := timeRange["until"].(string); ok {
				query["time_range"] = []string{fmt.Sprintf(`{"since":"%s","until":"%s"}`, since, until)}
			}
		}
	} else {
		query["date_preset"] = []string{"last_7d"} // Default
	}

	var endpoint string
	if campaignIds, ok := args["campaign_ids"].([]any); ok && len(campaignIds) > 0 {
		// Handle multiple campaigns - need to make multiple requests
		var allInsights []map[string]any
		for _, id := range campaignIds {
			if campaignId, ok := id.(string); ok {
				resp, err := client.MakeRequest(pdk.MethodGet, fmt.Sprintf("%s/insights", campaignId), query, nil)
				if err != nil {
					continue // Skip failed requests
				}
				var insights map[string]any
				if err := json.Unmarshal(resp.Body, &insights); err == nil {
					if data, ok := insights["data"].([]any); ok {
						for _, insight := range data {
							if insightMap, ok := insight.(map[string]any); ok {
								allInsights = append(allInsights, insightMap)
							}
						}
					}
				}
			}
		}
		result := map[string]any{
			"data": allInsights,
		}
		jsonResult, _ := json.Marshal(result)
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(jsonResult)),
			}},
		}, nil
	} else if accountId, ok := args["account_id"].(string); ok && accountId != "" {
		endpoint = fmt.Sprintf("%s/insights", accountId)
		query["level"] = []string{"campaign"}
	} else {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Either campaign_ids or account_id is required"),
			}},
		}, nil
	}

	resp, err := client.MakeRequest(pdk.MethodGet, endpoint, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get campaign insights: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func getAdSetInsights(token string, args map[string]any) (CallToolResult, error) {
	client := NewMetaClient(token)

	query := map[string][]string{}

	// Set default metrics if none provided
	metrics := []string{"impressions", "clicks", "spend", "reach", "ctr", "cpc", "cpm"}
	if providedMetrics, ok := args["metrics"].([]any); ok && len(providedMetrics) > 0 {
		metrics = []string{}
		for _, m := range providedMetrics {
			if metricStr, ok := m.(string); ok {
				metrics = append(metrics, metricStr)
			}
		}
	}
	query["fields"] = []string{strings.Join(metrics, ",")}

	// Handle date range
	if datePreset, ok := args["date_preset"].(string); ok && datePreset != "" {
		query["date_preset"] = []string{datePreset}
	} else if timeRange, ok := args["time_range"].(map[string]any); ok {
		if since, ok := timeRange["since"].(string); ok {
			if until, ok := timeRange["until"].(string); ok {
				query["time_range"] = []string{fmt.Sprintf(`{"since":"%s","until":"%s"}`, since, until)}
			}
		}
	} else {
		query["date_preset"] = []string{"last_7d"} // Default
	}

	var endpoint string
	if adsetIds, ok := args["adset_ids"].([]any); ok && len(adsetIds) > 0 {
		// Handle multiple ad sets
		var allInsights []map[string]any
		for _, id := range adsetIds {
			if adsetId, ok := id.(string); ok {
				resp, err := client.MakeRequest(pdk.MethodGet, fmt.Sprintf("%s/insights", adsetId), query, nil)
				if err != nil {
					continue // Skip failed requests
				}
				var insights map[string]any
				if err := json.Unmarshal(resp.Body, &insights); err == nil {
					if data, ok := insights["data"].([]any); ok {
						for _, insight := range data {
							if insightMap, ok := insight.(map[string]any); ok {
								allInsights = append(allInsights, insightMap)
							}
						}
					}
				}
			}
		}
		result := map[string]any{
			"data": allInsights,
		}
		jsonResult, _ := json.Marshal(result)
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(jsonResult)),
			}},
		}, nil
	} else if campaignId, ok := args["campaign_id"].(string); ok && campaignId != "" {
		endpoint = fmt.Sprintf("%s/insights", campaignId)
		query["level"] = []string{"adset"}
	} else if accountId, ok := args["account_id"].(string); ok && accountId != "" {
		endpoint = fmt.Sprintf("%s/insights", accountId)
		query["level"] = []string{"adset"}
	} else {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Either adset_ids, campaign_id, or account_id is required"),
			}},
		}, nil
	}

	resp, err := client.MakeRequest(pdk.MethodGet, endpoint, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ad set insights: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func getAdInsights(token string, args map[string]any) (CallToolResult, error) {
	client := NewMetaClient(token)

	query := map[string][]string{}

	// Set default metrics if none provided
	metrics := []string{"impressions", "clicks", "spend", "reach", "ctr", "cpc", "cpm"}
	if providedMetrics, ok := args["metrics"].([]any); ok && len(providedMetrics) > 0 {
		metrics = []string{}
		for _, m := range providedMetrics {
			if metricStr, ok := m.(string); ok {
				metrics = append(metrics, metricStr)
			}
		}
	}
	query["fields"] = []string{strings.Join(metrics, ",")}

	// Handle date range
	if datePreset, ok := args["date_preset"].(string); ok && datePreset != "" {
		query["date_preset"] = []string{datePreset}
	} else if timeRange, ok := args["time_range"].(map[string]any); ok {
		if since, ok := timeRange["since"].(string); ok {
			if until, ok := timeRange["until"].(string); ok {
				query["time_range"] = []string{fmt.Sprintf(`{"since":"%s","until":"%s"}`, since, until)}
			}
		}
	} else {
		query["date_preset"] = []string{"last_7d"} // Default
	}

	var endpoint string
	if adIds, ok := args["ad_ids"].([]any); ok && len(adIds) > 0 {
		// Handle multiple ads
		var allInsights []map[string]any
		for _, id := range adIds {
			if adId, ok := id.(string); ok {
				resp, err := client.MakeRequest(pdk.MethodGet, fmt.Sprintf("%s/insights", adId), query, nil)
				if err != nil {
					continue // Skip failed requests
				}
				var insights map[string]any
				if err := json.Unmarshal(resp.Body, &insights); err == nil {
					if data, ok := insights["data"].([]any); ok {
						for _, insight := range data {
							if insightMap, ok := insight.(map[string]any); ok {
								allInsights = append(allInsights, insightMap)
							}
						}
					}
				}
			}
		}
		result := map[string]any{
			"data": allInsights,
		}
		jsonResult, _ := json.Marshal(result)
		return CallToolResult{
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(string(jsonResult)),
			}},
		}, nil
	} else if adsetId, ok := args["adset_id"].(string); ok && adsetId != "" {
		endpoint = fmt.Sprintf("%s/insights", adsetId)
		query["level"] = []string{"ad"}
	} else if campaignId, ok := args["campaign_id"].(string); ok && campaignId != "" {
		endpoint = fmt.Sprintf("%s/insights", campaignId)
		query["level"] = []string{"ad"}
	} else if accountId, ok := args["account_id"].(string); ok && accountId != "" {
		endpoint = fmt.Sprintf("%s/insights", accountId)
		query["level"] = []string{"ad"}
	} else {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Either ad_ids, adset_id, campaign_id, or account_id is required"),
			}},
		}, nil
	}

	resp, err := client.MakeRequest(pdk.MethodGet, endpoint, query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get ad insights: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

// Audience Management Functions
func getAudiences(token string, args map[string]any) (CallToolResult, error) {
	accountId, ok := args["account_id"].(string)
	if !ok || accountId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Account ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,description,audience_type,subtype,time_created,time_updated,approximate_count"},
	}

	if limit, ok := args["limit"].(float64); ok && limit > 0 {
		query["limit"] = []string{fmt.Sprintf("%.0f", limit)}
	}

	resp, err := client.MakeRequest(pdk.MethodGet, fmt.Sprintf("%s/customaudiences", accountId), query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get audiences: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func createAudience(token string, args map[string]any) (CallToolResult, error) {
	accountId, ok := args["account_id"].(string)
	if !ok || accountId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Account ID is required"),
			}},
		}, nil
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Audience name is required"),
			}},
		}, nil
	}

	subtype, ok := args["subtype"].(string)
	if !ok || subtype == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Audience subtype is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	formData := map[string]string{
		"name":    name,
		"subtype": subtype,
	}

	if description, ok := args["description"].(string); ok && description != "" {
		formData["description"] = description
	}

	resp, err := client.MakeFormRequest(fmt.Sprintf("%s/customaudiences", accountId), formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create audience: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

// Creative Management Functions
func getCreatives(token string, args map[string]any) (CallToolResult, error) {
	accountId, ok := args["account_id"].(string)
	if !ok || accountId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Account ID is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	query := map[string][]string{
		"fields": {"id,name,status,object_story_spec,degrees_of_freedom_spec,body,title,image_url,video_id"},
	}

	if limit, ok := args["limit"].(float64); ok && limit > 0 {
		query["limit"] = []string{fmt.Sprintf("%.0f", limit)}
	}

	resp, err := client.MakeRequest(pdk.MethodGet, fmt.Sprintf("%s/adcreatives", accountId), query, nil)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to get creatives: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

func createCreative(token string, args map[string]any) (CallToolResult, error) {
	accountId, ok := args["account_id"].(string)
	if !ok || accountId == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Account ID is required"),
			}},
		}, nil
	}

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Creative name is required"),
			}},
		}, nil
	}

	objectStorySpec, ok := args["object_story_spec"].(map[string]any)
	if !ok || len(objectStorySpec) == 0 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Object story spec is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)

	formData := map[string]string{
		"name": name,
	}

	// Convert object story spec to JSON
	objectStorySpecJson, err := json.Marshal(objectStorySpec)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to marshal object story spec: %s", err)),
			}},
		}, nil
	}
	formData["object_story_spec"] = string(objectStorySpecJson)

	if degreesOfFreedomSpec, ok := args["degrees_of_freedom_spec"].(map[string]any); ok && len(degreesOfFreedomSpec) > 0 {
		degreesOfFreedomSpecJson, err := json.Marshal(degreesOfFreedomSpec)
		if err == nil {
			formData["degrees_of_freedom_spec"] = string(degreesOfFreedomSpecJson)
		}
	}

	resp, err := client.MakeFormRequest(fmt.Sprintf("%s/adcreatives", accountId), formData)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to create creative: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(resp.Body)),
		}},
	}, nil
}

// Bulk Operations
func bulkUpdate(token string, args map[string]any) (CallToolResult, error) {
	operations, ok := args["operations"].([]any)
	if !ok || len(operations) == 0 {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Operations array is required"),
			}},
		}, nil
	}

	client := NewMetaClient(token)
	var results []map[string]any

	for i, op := range operations {
		opMap, ok := op.(map[string]any)
		if !ok {
			results = append(results, map[string]any{
				"index":   i,
				"success": false,
				"error":   "Invalid operation format",
			})
			continue
		}

		objectType, ok := opMap["object_type"].(string)
		if !ok {
			results = append(results, map[string]any{
				"index":   i,
				"success": false,
				"error":   "Missing object_type",
			})
			continue
		}

		objectId, ok := opMap["object_id"].(string)
		if !ok {
			results = append(results, map[string]any{
				"index":   i,
				"success": false,
				"error":   "Missing object_id",
			})
			continue
		}

		updates, ok := opMap["updates"].(map[string]any)
		if !ok {
			results = append(results, map[string]any{
				"index":   i,
				"success": false,
				"error":   "Missing updates",
			})
			continue
		}

		// Convert updates to form data
		formData := map[string]string{}
		for key, value := range updates {
			switch v := value.(type) {
			case string:
				formData[key] = v
			case float64:
				formData[key] = fmt.Sprintf("%.0f", v)
			case bool:
				if v {
					formData[key] = "true"
				} else {
					formData[key] = "false"
				}
			case map[string]any:
				jsonValue, err := json.Marshal(v)
				if err == nil {
					formData[key] = string(jsonValue)
				}
			}
		}

		// Execute the update
		resp, err := client.MakeFormRequest(objectId, formData)
		if err != nil {
			results = append(results, map[string]any{
				"index":       i,
				"object_id":   objectId,
				"object_type": objectType,
				"success":     false,
				"error":       err.Error(),
			})
		} else {
			results = append(results, map[string]any{
				"index":       i,
				"object_id":   objectId,
				"object_type": objectType,
				"success":     true,
				"response":    string(resp.Body),
			})
		}
	}

	// Return summary of results
	summary := map[string]any{
		"total_operations": len(operations),
		"successful":       0,
		"failed":           0,
		"results":          results,
	}

	for _, result := range results {
		if success, ok := result["success"].(bool); ok && success {
			summary["successful"] = summary["successful"].(int) + 1
		} else {
			summary["failed"] = summary["failed"].(int) + 1
		}
	}

	jsonResult, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(fmt.Sprintf("Failed to format results: %s", err)),
			}},
		}, nil
	}

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(string(jsonResult)),
		}},
	}, nil
}
