import { ToolDescription } from "./pdk";

// Tool definitions
export const GEOCODE_TOOL: ToolDescription = {
  name: "maps_geocode",
  description: "Convert an address into geographic coordinates",
  inputSchema: {
    type: "object",
    properties: {
      address: {
        type: "string",
        description: "The address to geocode"
      }
    },
    required: ["address"]
  }
};

export const REVERSE_GEOCODE_TOOL: ToolDescription = {
  name: "maps_reverse_geocode",
  description: "Convert coordinates into an address",
  inputSchema: {
    type: "object",
    properties: {
      latitude: {
        type: "number",
        description: "Latitude coordinate"
      },
      longitude: {
        type: "number",
        description: "Longitude coordinate"
      }
    },
    required: ["latitude", "longitude"]
  }
};

export const SEARCH_PLACES_TOOL: ToolDescription = {
  name: "maps_search_places",
  description: "Search for places using Google Places API",
  inputSchema: {
    type: "object",
    properties: {
      query: {
        type: "string",
        description: "Search query"
      },
      location: {
        type: "object",
        properties: {
          latitude: { 
            type: "number",
            description: "Latitude coordinate for the center point" 
          },
          longitude: { 
            type: "number",
            description: "Longitude coordinate for the center point"
          }
        },
        description: "Optional center point for the search"
      },
      radius: {
        type: "number",
        description: "Search radius in meters (max 50000)"
      }
    },
    required: ["query"]
  }
};

export const PLACE_DETAILS_TOOL: ToolDescription = {
  name: "maps_place_details",
  description: "Get detailed information about a specific place",
  inputSchema: {
    type: "object",
    properties: {
      place_id: {
        type: "string",
        description: "The place ID to get details for"
      }
    },
    required: ["place_id"]
  }
};

export const DISTANCE_MATRIX_TOOL: ToolDescription = {
  name: "maps_distance_matrix",
  description: "Calculate travel distance and time for multiple origins and destinations",
  inputSchema: {
    type: "object",
    properties: {
      origins: {
        type: "array",
        items: { type: "string" },
        description: "Array of origin addresses or coordinates"
      },
      destinations: {
        type: "array",
        items: { type: "string" },
        description: "Array of destination addresses or coordinates"
      },
      mode: {
        type: "string",
        description: "Travel mode (driving, walking, bicycling, transit)",
        enum: ["driving", "walking", "bicycling", "transit"]
      }
    },
    required: ["origins", "destinations"]
  }
};

export const ELEVATION_TOOL: ToolDescription = {
  name: "maps_elevation",
  description: "Get elevation data for locations on the earth",
  inputSchema: {
    type: "object",
    properties: {
      locations: {
        type: "array",
        items: {
          type: "object",
          properties: {
            latitude: { 
              type: "number",
              description: "Latitude coordinate of the location"
            },
            longitude: { 
              type: "number",
              description: "Longitude coordinate of the location"
            }
          },
          required: ["latitude", "longitude"]
        },
        description: "Array of locations to get elevation for"
      }
    },
    required: ["locations"]
  }
};

export const DIRECTIONS_TOOL: ToolDescription = {
  name: "maps_directions",
  description: "Get directions between two points",
  inputSchema: {
    type: "object",
    properties: {
      origin: {
        type: "string",
        description: "Starting point address or coordinates"
      },
      destination: {
        type: "string",
        description: "Ending point address or coordinates"
      },
      mode: {
        type: "string",
        description: "Travel mode (driving, walking, bicycling, transit)",
        enum: ["driving", "walking", "bicycling", "transit"]
      }
    },
    required: ["origin", "destination"]
  }
};

export const ALL_TOOLS = [
  GEOCODE_TOOL,
  REVERSE_GEOCODE_TOOL,
  SEARCH_PLACES_TOOL,
  PLACE_DETAILS_TOOL,
  DISTANCE_MATRIX_TOOL,
  ELEVATION_TOOL,
  DIRECTIONS_TOOL,
];
