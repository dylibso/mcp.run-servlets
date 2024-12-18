# Amadeus Flight API

Amadeus is a multinational tech company that provides software for the global travel and tourism industry. 

The company's core business involves connecting travel providers (like airlines, hotels, and car rental companies) with travel sellers (such as travel agencies and websites). Their systems help manage reservations, inventory, and departure control for airlines, while also providing booking capabilities for travel agents and online travel platforms.

One of Amadeus's key strengths is its flight booking system, which processes hundreds of millions of bookings annually. They also provide important technological infrastructure that powers many airline websites and reservation systems worldwide.

This servlet makes available the [Flight Offer Search endpoint](https://developers.amadeus.com/self-service/category/flights/api-doc/flight-offers-search) as a tool.

NOTE:  Depending on the credentials and host you will use, you might get real-time or simulated data.

You can get free development access to the test endpoints by creating an account at https://developers.amadeus.com/register

## Config

- `base-url`: the base URL of the API. For the test system use `https://test.api.amadeus.com`
- `api-key`: your self-service or production API key 
- `api-secret`: your self-service or production API Secret

## Hosts

- You should allow access to `*.amadeus.com` or at least `test.api.amadeus.com`

## Example:

— hi Claude I'd like to spend my holidays in Paris, can you find me a flight from Milan Malpensa on Xmas eve, and return on January new year's day?


<blockquote>
I'll help you search for flights from Milan Malpensa (MXP) to Paris for the holiday season.

I've found several flight options for your Christmas holiday trip to Paris. All flights are operated by Lufthansa and include a stopover in Munich. Here are the best options:

1. **Afternoon Departure (13:15)**
   - Outbound: Dec 24, MXP 13:15 → CDG 17:25 (4h 10m total)
   - Return: Jan 1, CDG 09:25 → MXP 12:35 (3h 10m total)
   - Price: €241.32

2. **Early Morning Departure (06:00)**
   - Outbound: Dec 24, MXP 06:00 → CDG 11:30 (5h 30m total)
   - Return: Jan 1, CDG 09:25 → MXP 12:35 (3h 10m total)
   - Price: €241.32

All flights include:
- 1 checked bag
- Standard seat reservation
- Option to purchase meals
- Possibility to change flights (with fee)

Would you like to know more details about any of these options or would you like me to search for different times?
</blockquote>

— is there any non-stop option?

<blockquote>
Let me search for direct flights between Milan Malpensa and Paris.

Yes, I found several direct flight options! You have choices between Air France (flying to Charles de Gaulle - CDG) and Vueling (flying to Orly - ORY). Here are the best options:

1. **Air France to CDG**
   - Outbound: Dec 24, MXP 18:40 → CDG 20:15 (1h 35m)
   - Return: Jan 1, CDG 09:05 → MXP 10:35 (1h 30m)
   - Price: €247.16
   - Includes: 1 checked bag, snack, beverage, and standard seat selection

2. **Vueling to ORY**
   - Outbound: Dec 24, MXP 20:45 → ORY 22:15 (1h 30m)
   - Return: Jan 1, ORY 06:40 → MXP 08:10 (1h 30m)
   - Price: €250.62
   - Includes: 25kg checked baggage allowance

These direct flights are only slightly more expensive than the connecting flights (around €5-10 more) but save you 2-4 hours of travel time each way. Would you like me to provide more details about any of these options?

Note that the Air France flight arrives at Charles de Gaulle (CDG), which is connected to central Paris by RER B train, while the Vueling flight uses Orly (ORY), which is typically closer to southern Paris and connected by the Orlyval shuttle to RER B.
</blockquote>