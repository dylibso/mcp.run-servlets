# Meta Ads


This project provides a comprehensive suite of tools designed for managing and analyzing advertising campaigns on the Meta advertising platform. By leveraging these tools, users can efficiently create, monitor, and optimize their advertising efforts across various campaigns, ad sets, and ads.

## Features

- **Business Management:** Retrieve all businesses associated with a user account for better organization and management.
- **Ad Account Access:** Access all ad accounts available to the user, allowing for streamlined management of advertising resources.
- **Detailed Account Insights:** Obtain comprehensive information about specific ad accounts, including performance metrics and settings.
- **Campaign Management:** Retrieve and manage campaigns linked to ad accounts, including filtering options based on status and limits.
- **Ad Set Management:** Access ad sets associated with specific campaigns or ad accounts, facilitating detailed analysis and organization.
- **Ad Management:** Retrieve ads linked to particular ad sets, campaigns, or accounts, enabling thorough oversight of advertising materials.
- **Performance Insights:** Gain insights into campaign, ad set, and ad performance, including metrics such as impressions, clicks, and conversions.
- **Custom Audience Management:** Manage custom audiences within ad accounts to target advertising more effectively.
- **Creative Management:** Retrieve and manage creative assets used in advertising efforts.
- **Bulk Operations:** Perform bulk updates on multiple campaigns, ad sets, or ads for efficient management.

## Configuration

To utilize these tools, authentication via OAuth is required. The tools are configured to interact with the Meta advertising platform, specifically through the Facebook Graph API. A network connection to `graph.facebook.com` is mandatory to access and manage advertising resources effectively.

## Tool Descriptions

### Business Management Tool

- **Description:** Retrieve all businesses associated with the user account.
- **Input Requirements:** None.

### Ad Account Tool

- **Description:** Access all ad accounts available to the user. Optionally, filter by a specific business ID.
- **Input Requirements:** 
  - `business_id`: (optional) The ID of the business to filter ad accounts.

### Account Info Tool

- **Description:** Obtain detailed information about a specific ad account.
- **Input Requirements:**
  - `account_id`: The ad account ID (format: act_123456789).

### Campaign Management Tool

- **Description:** Retrieve campaigns for a specific ad account, with options to filter by status and limit the number of campaigns returned.
- **Input Requirements:**
  - `account_id`: The ad account ID (format: act_123456789).
  - `limit`: (optional) The maximum number of campaigns to return.
  - `status`: (optional) Filter for campaign statuses (ACTIVE, PAUSED, ARCHIVED).

### Ad Set Management Tool

- **Description:** Retrieve ad sets for a specific campaign or ad account, with options to filter by status and limit the number returned.
- **Input Requirements:**
  - `account_id`: (optional) The ad account ID to filter ad sets.
  - `campaign_id`: (optional) The campaign ID to filter ad sets.
  - `limit`: (optional) The maximum number of ad sets to return.
  - `status`: (optional) Filter for ad set statuses (ACTIVE, PAUSED, ARCHIVED).

### Ad Management Tool

- **Description:** Retrieve ads for a specific ad set, campaign, or ad account, with options to filter by status and limit the number returned.
- **Input Requirements:**
  - `account_id`: (optional) The ad account ID to filter ads.
  - `adset_id`: (optional) The ad set ID to filter ads.
  - `campaign_id`: (optional) The campaign ID to filter ads.
  - `limit`: (optional) The maximum number of ads to return.
  - `status`: (optional) Filter for ad statuses (ACTIVE, PAUSED, ARCHIVED).

### Performance Insights Tool

- **Description:** Get performance insights for campaigns, ad sets, and ads, with customizable metrics and time ranges.
- **Input Requirements:**
  - **For campaigns:** 
    - `account_id`: The account ID to get insights for.
    - `campaign_ids`: An array of campaign IDs for insights.
    - `date_preset`: Preset value for date ranges.
    - `metrics`: Metrics to retrieve.
    - `time_range`: Optional custom time range.
  - **For ad sets:** 
    - Similar to above but focusing on ad sets.
  - **For ads:** 
    - Similar to above but focusing on ads.

### Audience Management Tool

- **Description:** Retrieve custom audiences within a specific ad account.
- **Input Requirements:**
  - `account_id`: The ad account ID (format: act_123456789).
  - `limit`: (optional) The maximum number of audiences to return.

### Creative Management Tool

- **Description:** Retrieve creatives associated with an ad account.
- **Input Requirements:**
  - `account_id`: The ad account ID (format: act_123456789).
  - `limit`: (optional) The maximum number of creatives to return.

### Bulk Update Tool

- **Description:** Perform bulk updates on multiple campaigns, ad sets, or ads for streamlined management.
- **Input Requirements:**
  - `operations`: An array of update operations specifying the type of object to update, object ID, and fields to update.

This suite of tools is designed to enhance the management and analysis of advertising activities on the Meta platform, ensuring users can efficiently optimize their advertising strategies.

## How to regsiter an oauth client

### Creating an app
1. Go to https://developers.facebook.com/apps
2. Click on "Create App"
3. Fill in form, click on "Next"
4. In the use cases, click on "Other"
5. Select "Business activity"
6. In Business Portfolio pick the one for mcp.run/dylibso. Note: we will need to verify the business before being able to publish the app for all users, then click on "Create app"

### Basic settings
1. Click on "App settings" -> "Basic", take note of App ID and App secret
2. Click on "Add platform" -> Website, and then in Site URL write "https://www.mcp.run/"
3. Note that you can start business verification from this page

![image](https://github.com/user-attachments/assets/17268b36-ff48-4fe5-af93-315b75156b12)

### Advanced Settings
1. Click on "App settings" -> "Advanced"
2. in "Share Redirect Domain Allow List" add "www.mcp.run" and "mcp.run"

### Setting up Marketing API
1. Go to `https://developers.facebook.com/apps/:appId/add/` or click on "Add Product" on the sidebar
![image](https://github.com/user-attachments/assets/4cb69641-a223-4260-80b9-51647324d32b)
3. Click on "setup up" on Marketing API
4. you don't need to do anything else here

### Setting up Facebook Login for Businesses
1. Go back to "Add product" and then set up "Facebook Login for Business"
2. Go to `https://developers.facebook.com/apps/:appId/business-login/settings/`
3. In `Valid OAuth Redirect URIs` add `https://www.mcp.run/api/login/servlet/dylibso/meta-ads/callback`

## OAuth configuration

- Scopes: `ads_management`, `ads_read`, `read_insights`
- Allowed domains: `graph.facebook.com`
- Authorize URL: `https://www.facebook.com/dialog/oauth`
- Token URL: `https://graph.facebook.com/oauth/access_token`