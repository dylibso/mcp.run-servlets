import { CallToolRequest, CallToolResult, ContentType, ListToolsResult, ToolDescription } from "./pdk";
import * as tools from "./tools";

/**
 * Called when the tool is invoked.
 * If you support multiple tools, you must switch on the input.params.name to detect which tool is being called.
 * The name will match one of the tool names returned from "describe".
 *
 * @param {CallToolRequest} input - The incoming tool request from the LLM
 * @returns {CallToolResult} The servlet's response to the given tool call
 */
export function callImpl(input: CallToolRequest): CallToolResult {
  const args = input.params.arguments;
  const apiKey = Config.get("api-key")!;
  switch (input.params.name) {
    case tools.GEOCODE_TOOL.name:
      return handleGeocode(apiKey, args.address);
    case tools.REVERSE_GEOCODE_TOOL.name:
      return handleReverseGeocode(apiKey, args.latitude, args.longitude);
    case tools.SEARCH_PLACES_TOOL.name:
      return handlePlaceSearch(apiKey, args.query, args.location, args.radius);
    case tools.PLACE_DETAILS_TOOL.name:
      return handlePlaceDetails(apiKey, args.place_id);
    case tools.DISTANCE_MATRIX_TOOL.name:
      return handleDistanceMatrix(apiKey, args.origins, args.destinations, args.mode);
    case tools.ELEVATION_TOOL.name:
      return handleElevation(apiKey, args.locations);
    case tools.DIRECTIONS_TOOL.name:
      return handleDirections(apiKey, args.origin, args.destination, args.mode);
    default:
      throw new Error(`Unknown tool name: ${input.params.name}`);
  }
}

/**
 * Called by mcpx to understand how and why to use this tool.
 * Note: Your servlet configs will not be set when this function is called,
 * so do not rely on config in this function
 *
 * @returns {ListToolsResult} The tools' descriptions, supporting multiple tools from a single servlet.
 */
export function describeImpl(): ListToolsResult {
  return {
    tools: tools.ALL_TOOLS
  };
}

function handleGeocode(apiKey: string, address: string): CallToolResult {
  const url = new URL("https://maps.googleapis.com/maps/api/geocode/json");
  url.searchParams.append("address", address);
  url.searchParams.append("key", apiKey);

  const response = Http.request({ method: "GET", url: url.toString() });
  const data = JSON.parse(response.body) as GeocodeResponse;

  if (data.status !== "OK") {
    return {
      content: [{
        type: ContentType.Text,
        text: `Geocoding failed: ${data.error_message || data.status}`
      }],
      isError: true
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: JSON.stringify({
        location: data.results[0].geometry.location,
        formatted_address: data.results[0].formatted_address,
        place_id: data.results[0].place_id
      }, null, 2)
    }],
    isError: false
  };
}

function handleReverseGeocode(apiKey: string, latitude: number, longitude: number): CallToolResult {
  const url = new URL("https://maps.googleapis.com/maps/api/geocode/json");
  url.searchParams.append("latlng", `${latitude},${longitude}`);
  url.searchParams.append("key", apiKey);

  const response = Http.request({ method: "GET", url: url.toString() });
  const data = JSON.parse(response.body) as GeocodeResponse;

  if (data.status !== "OK") {
    return {
      content: [{
        type: ContentType.Text,
        text: `Reverse geocoding failed: ${data.error_message || data.status}`
      }],
      isError: true
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: JSON.stringify({
        formatted_address: data.results[0].formatted_address,
        place_id: data.results[0].place_id,
        address_components: data.results[0].address_components
      }, null, 2)
    }],
    isError: false
  };
}


function handlePlaceSearch(
  apiKey: string,
  query: string,
  location?: { latitude: number; longitude: number },
  radius?: number
): CallToolResult {
  const url = new URL("https://maps.googleapis.com/maps/api/place/textsearch/json");
  url.searchParams.append("query", query);
  url.searchParams.append("key", apiKey);

  if (location) {
    url.searchParams.append("location", `${location.latitude},${location.longitude}`);
  }
  if (radius) {
    url.searchParams.append("radius", radius.toString());
  }

  const response = Http.request({ method: "GET", url: url.toString() });
  const data = JSON.parse(response.body) as PlacesSearchResponse;

  if (data.status !== "OK") {
    return {
      content: [{
        type: ContentType.Text,
        text: `Place search failed: ${data.error_message || data.status}`
      }],
      isError: true
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: JSON.stringify({
        places: data.results.map((place) => ({
          name: place.name,
          formatted_address: place.formatted_address,
          location: place.geometry.location,
          place_id: place.place_id,
          rating: place.rating,
          types: place.types
        }))
      }, null, 2)
    }],
    isError: false
  };
}

function handlePlaceDetails(apiKey: string, placeId: string): CallToolResult {
  const url = new URL("https://maps.googleapis.com/maps/api/place/details/json");
  url.searchParams.append("place_id", placeId);
  url.searchParams.append("key", apiKey);

  const response = Http.request({
    method: "GET",
    url: url.toString()
  });
  const data = JSON.parse(response.body) as PlaceDetailsResponse;

  if (data.status !== "OK") {
    return {
      content: [{
        type: ContentType.Text,
        text: `Place details request failed: ${data.error_message || data.status}`
      }],
      isError: true
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: JSON.stringify({
        name: data.result.name,
        formatted_address: data.result.formatted_address,
        location: data.result.geometry.location,
        formatted_phone_number: data.result.formatted_phone_number,
        website: data.result.website,
        rating: data.result.rating,
        reviews: data.result.reviews,
        opening_hours: data.result.opening_hours
      }, null, 2)
    }],
    isError: false
  };
}

function handleDistanceMatrix(
  apiKey: string,
  origins: string[],
  destinations: string[],
  mode: "driving" | "walking" | "bicycling" | "transit" = "driving"
): CallToolResult {
  const url = new URL("https://maps.googleapis.com/maps/api/distancematrix/json");
  url.searchParams.append("origins", origins.join("|"));
  url.searchParams.append("destinations", destinations.join("|"));
  url.searchParams.append("mode", mode);
  url.searchParams.append("key", apiKey);

  const response = Http.request({ method: "GET", url: url.toString() });
  const data = JSON.parse(response.body) as DistanceMatrixResponse;

  if (data.status !== "OK") {
    return {
      content: [{
        type: ContentType.Text,
        text: `Distance matrix request failed: ${data.error_message || data.status}`
      }],
      isError: true
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: JSON.stringify({
        origin_addresses: data.origin_addresses,
        destination_addresses: data.destination_addresses,
        results: data.rows.map((row) => ({
          elements: row.elements.map((element) => ({
            status: element.status,
            duration: element.duration,
            distance: element.distance
          }))
        }))
      }, null, 2)
    }],
    isError: false
  };
}

function handleElevation(apiKey: string, locations: Array<{ latitude: number; longitude: number }>): CallToolResult {
  const url = new URL("https://maps.googleapis.com/maps/api/elevation/json");
  const locationString = locations
    .map((loc) => `${loc.latitude},${loc.longitude}`)
    .join("|");
  url.searchParams.append("locations", locationString);
  url.searchParams.append("key", apiKey);

  const response = Http.request({ method: "GET", url: url.toString() });
  const data = JSON.parse(response.body) as ElevationResponse;

  if (data.status !== "OK") {
    return {
      content: [{
        type: ContentType.Text,
        text: `Elevation request failed: ${data.error_message || data.status}`
      }],
      isError: true
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: JSON.stringify({
        results: data.results.map((result) => ({
          elevation: result.elevation,
          location: result.location,
          resolution: result.resolution
        }))
      }, null, 2)
    }],
    isError: false
  };
}


function handleDirections(
  apiKey: string,
  origin: string,
  destination: string,
  mode: "driving" | "walking" | "bicycling" | "transit" = "driving"
): CallToolResult {
  const url = new URL("https://maps.googleapis.com/maps/api/directions/json");
  url.searchParams.append("origin", origin);
  url.searchParams.append("destination", destination);
  url.searchParams.append("mode", mode);
  url.searchParams.append("key", apiKey);

  const response = Http.request({ method: "GET", url: url.toString() });
  const data = JSON.parse(response.body) as DirectionsResponse;

  if (data.status !== "OK") {
    return {
      content: [{
        type: ContentType.Text,
        text: `Directions request failed: ${data.error_message || data.status}`
      }],
      isError: true
    };
  }

  return {
    content: [{
      type: ContentType.Text,
      text: JSON.stringify({
        routes: data.routes.map((route) => ({
          summary: route.summary,
          distance: route.legs[0].distance,
          duration: route.legs[0].duration,
          steps: route.legs[0].steps.map((step) => ({
            instructions: step.html_instructions,
            distance: step.distance,
            duration: step.duration,
            travel_mode: step.travel_mode
          }))
        }))
      }, null, 2)
    }],
    isError: false
  };
}
