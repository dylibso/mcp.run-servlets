// Response interfaces
interface GoogleMapsResponse {
  status: string;
  error_message?: string;
}

interface GeocodeResponse extends GoogleMapsResponse {
  results: Array<{
    place_id: string;
    formatted_address: string;
    geometry: {
      location: {
        lat: number;
        lng: number;
      }
    };
    address_components: Array<{
      long_name: string;
      short_name: string;
      types: string[];
    }>;
  }>;
}

interface PlacesSearchResponse extends GoogleMapsResponse {
  results: Array<{
    name: string;
    place_id: string;
    formatted_address: string;
    geometry: {
      location: {
        lat: number;
        lng: number;
      }
    };
    rating?: number;
    types: string[];
  }>;
}

interface PlaceDetailsResponse extends GoogleMapsResponse {
  result: {
    name: string;
    place_id: string;
    formatted_address: string;
    formatted_phone_number?: string;
    website?: string;
    rating?: number;
    reviews?: Array<{
      author_name: string;
      rating: number;
      text: string;
      time: number;
    }>;
    opening_hours?: {
      weekday_text: string[];
      open_now: boolean;
    };
    geometry: {
      location: {
        lat: number;
        lng: number;
      }
    };
  };
}

interface DistanceMatrixResponse extends GoogleMapsResponse {
  origin_addresses: string[];
  destination_addresses: string[];
  rows: Array<{
    elements: Array<{
      status: string;
      duration: {
        text: string;
        value: number;
      };
      distance: {
        text: string;
        value: number;
      };
    }>;
  }>;
}

interface ElevationResponse extends GoogleMapsResponse {
  results: Array<{
    elevation: number;
    location: {
      lat: number;
      lng: number;
    };
    resolution: number;
  }>;
}

interface DirectionsResponse extends GoogleMapsResponse {
  routes: Array<{
    summary: string;
    legs: Array<{
      distance: {
        text: string;
        value: number;
      };
      duration: {
        text: string;
        value: number;
      };
      steps: Array<{
        html_instructions: string;
        distance: {
          text: string;
          value: number;
        };
        duration: {
          text: string;
          value: number;
        };
        travel_mode: string;
      }>;
    }>;
  }>;
}
