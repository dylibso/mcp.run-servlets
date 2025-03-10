# Unsplash Servlet

[Unsplash](https://unsplash.com) is a collection of high-quality photos.
This servlet allows you to search and retrieve pictures from the archive.
Use in conjuction with the `fetch` servlet to fetch and display the URLs.

The servlet serves also as a tutorial for Zig support.

## Config

`API_KEY`: Your "Access Key" on https://unsplash.com/oauth/applications

## Hosts

- You should allow access to `api.unsplash.com`.

## Example:


fetch a photo of Milan, Piazza Duomo

> I'll help you fetch a photo of Milan's Piazza del Duomo using the Unsplash API
>
>> (uses `unsplash`)
>
> Let me fetch this specific photo to show you:
>
>> (uses `fetch`)
>
> This is a stunning aerial view of Milan's Piazza del Duomo, taken from one of the upper levels of the Cathedral (Duomo di Milano). The image captures the intricate Gothic architecture of the Duomo on the left, with its elaborate spires and statues, while showing the expansive square below filled with people. On the right side, you can see the elegant Galleria Vittorio Emanuele II, one of the world's oldest shopping galleries, with its distinctive neoclassical architecture. The photo provides a great perspective on how this historic square serves as the heart of Milan, where Gothic architecture meets classical design and modern city life.

The retrieved image: https://unsplash.com/photos/people-on-street-near-concrete-buildings-BBigUtobEYo
