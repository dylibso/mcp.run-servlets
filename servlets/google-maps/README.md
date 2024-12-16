# Google Maps

A port of the [Google Maps server][mcp-gmaps] to the mcp.run Servlet API.

[mcp-gmaps]: https://github.com/modelcontextprotocol/servers/tree/0d9cee028030246279b180c3e9197715e2073b38/src/google-maps

## Config

Requires the following config keys:

- **api-key** with a [Google Maps API Key](https://developers.google.com/maps/documentation/javascript/get-api-key#create-api-keys).

## Permissions

Requires access to the following **domains**:

- `maps.googleapis.com`

## Example

— hey claude can you give me directions from the Tour Eiffel to the Champs Elisees on foot?

<blockquote>
I'll help you get walking directions from the Eiffel Tower to the Champs-Élysées.

Here are the walking directions from the Eiffel Tower to the Champs-Élysées:

The walk is about 2.1 km and should take around 31 minutes. Here's the route:

1. From the Eiffel Tower, head southeast toward Avenue Gustave Eiffel
2. Turn left onto Avenue Silvestre de Sacy
3. Turn right onto Avenue de la Bourdonnais
4. Turn left onto Rue de Monttessuy
5. Turn left onto Avenue Rapp and follow it
6. Cross the Seine using the Pont de l'Alma bridge
7. Turn left onto Avenue du Président Wilson
8. Turn right onto Avenue Marceau and follow it
9. Turn right onto Rue Galilée
10. Finally, turn left and you'll arrive at the Champs-Élysées

This is a scenic walk that takes you across the Seine and through some beautiful Parisian neighborhoods. You'll get nice views of the city along the way. Would you like me to suggest any points of interest you might pass during this walk?
</blockquote>

— yes anything interesting?

<blockquote>
Let me check for notable places along this route.

Let me highlight some interesting spots you'll encounter on your walk:

1. Near your starting point, you can see Gustave Eiffel's original office (Bureau de Gustave Eiffel) at the top of the Eiffel Tower - it's been preserved as it was in his time.

2. As you walk along Avenue Rapp, you'll pass the beautiful Fontaine du Gros Caillou, a historic fountain from the 19th century.

3. When you cross the Seine on Pont de l'Alma, you'll pass by the Flame of Liberty (near Place Diana), which has become an unofficial memorial to Princess Diana, as it's located above the tunnel where her fatal accident occurred.

4. Along the Seine, you'll see the beautiful Quais de la Seine, which are listed as a UNESCO World Heritage site. The riverside walkways offer great views of the city.

5. As you approach the Champs-Élysées, you'll be entering the "Golden Triangle" (Triangle d'Or), one of Paris's most luxurious neighborhoods, known for its high-end fashion boutiques and beautiful Haussmann architecture.

This walk takes you through some of Paris's most elegant neighborhoods, from the more residential 7th arrondissement to the luxurious 8th arrondissement. Would you like more specific details about any of these locations?
</blockquote>


## Tools

1. `geocode`
   - Convert address to coordinates
   - Input: `address` (string)
   - Returns: location, formatted_address, place_id

2. `reverse_geocode`
   - Convert coordinates to address
   - Inputs:
     - `latitude` (number)
     - `longitude` (number)
   - Returns: formatted_address, place_id, address_components

3. `search_places`
   - Search for places using text query
   - Inputs:
     - `query` (string)
     - `location` (optional): { latitude: number, longitude: number }
     - `radius` (optional): number (meters, max 50000)
   - Returns: array of places with names, addresses, locations

4. `get_place_details`
   - Get detailed information about a place
   - Input: `place_id` (string)
   - Returns: name, address, contact info, ratings, reviews, opening hours

5. `get_distance_matrix`
   - Calculate distances and times between points
   - Inputs:
     - `origins` (string[])
     - `destinations` (string[])
     - `mode` (optional): "driving" | "walking" | "bicycling" | "transit"
   - Returns: distances and durations matrix

6. `get_elevation`
   - Get elevation data for locations
   - Input: `locations` (array of {latitude, longitude})
   - Returns: elevation data for each point

7. `get_directions`
   - Get directions between points
   - Inputs:
     - `origin` (string)
     - `destination` (string)
     - `mode` (optional): "driving" | "walking" | "bicycling" | "transit"
   - Returns: route details with steps, distance, duration

