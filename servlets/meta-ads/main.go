package main

import (
	"fmt"

	"github.com/extism/go-pdk"
)

// Schema related types for tool description
type SchemaProperty struct {
	Type        string  `json:"type"`
	Description string  `json:"description,omitempty"`
	Items       *schema `json:"items,omitempty"`
}

func prop(tpe, description string) SchemaProperty {
	return SchemaProperty{Type: tpe, Description: description}
}

func propWithItems(tpe, description string, itemsSchema schema) SchemaProperty {
	return SchemaProperty{
		Type:        tpe,
		Description: description,
		Items:       &itemsSchema,
	}
}

type schema = map[string]interface{}
type props = map[string]SchemaProperty

// Helper function to create a pointer
func some[T any](t T) *T {
	return &t
}

// Called when the tool is invoked
func Call(input CallToolRequest) (CallToolResult, error) {
	token, ok := pdk.GetConfig("OAUTH_TOKEN")
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("No OAUTH_TOKEN configured"),
			}},
		}, nil
	}

	args := input.Params.Arguments.(map[string]interface{})
	pdk.Log(pdk.LogDebug, fmt.Sprint("Args: ", args))

	switch input.Params.Name {
	// Account Management
	case GetBusinessesTool.Name:
		return getBusinesses(token, args)
	case GetAdAccountsTool.Name:
		return getAdAccounts(token, args)
	case GetAccountInfoTool.Name:
		return getAccountInfo(token, args)
	// Campaign Management
	case GetCampaignsTool.Name:
		return getCampaigns(token, args)
	case CreateCampaignTool.Name:
		return createCampaign(token, args)
	case UpdateCampaignTool.Name:
		return updateCampaign(token, args)
	case DeleteCampaignTool.Name:
		return deleteCampaign(token, args)
	case GetCampaignInfoTool.Name:
		return getCampaignInfo(token, args)
	// Ad Set Management
	case GetAdSetsTool.Name:
		return getAdSets(token, args)
	case CreateAdSetTool.Name:
		return createAdSet(token, args)
	case UpdateAdSetTool.Name:
		return updateAdSet(token, args)
	case DeleteAdSetTool.Name:
		return deleteAdSet(token, args)
	case GetAdSetInfoTool.Name:
		return getAdSetInfo(token, args)
	// Ad Management
	case GetAdsTool.Name:
		return getAds(token, args)
	case CreateAdTool.Name:
		return createAd(token, args)
	case UpdateAdTool.Name:
		return updateAd(token, args)
	case DeleteAdTool.Name:
		return deleteAd(token, args)
	case GetAdInfoTool.Name:
		return getAdInfo(token, args)
	// Insights
	case GetCampaignInsightsTool.Name:
		return getCampaignInsights(token, args)
	case GetAdSetInsightsTool.Name:
		return getAdSetInsights(token, args)
	case GetAdInsightsTool.Name:
		return getAdInsights(token, args)
	// Audiences
	case GetAudiencesTool.Name:
		return getAudiences(token, args)
	case CreateAudienceTool.Name:
		return createAudience(token, args)
	// Creatives
	case GetCreativesTool.Name:
		return getCreatives(token, args)
	case CreateCreativeTool.Name:
		return createCreative(token, args)
	// Bulk Operations
	case BulkUpdateTool.Name:
		return bulkUpdate(token, args)
	default:
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Unknown tool " + input.Params.Name),
			}},
		}, nil
	}
}

// Describe the tools provided by this servlet
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: MetaAdsReadOnlyTools,
	}, nil
}
