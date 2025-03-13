# Pantry

A simple JSON storage servlet based on [Pantry](https://getpantry.cloud). Use it to give long-term memory to your servlets or your tasks.

Register a Pantry with just an email address on [Pantry](https://getpantry.cloud), you will be given a `PANTRY_ID` you can use to store and retrieve data.

## Config

- `PANTRY_ID`: the ID of the pantry you want to use

## Domains

- getpantry.cloud

## Tools
`pantry` takes an `action` (`get` or `post`) and a `basket` to identify the data. If `post` then it also takes a `body` parameter in JSON format.
