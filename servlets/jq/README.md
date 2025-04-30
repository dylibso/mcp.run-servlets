# jq

A simple `jq` servlet based on the popular [`jaq` library](https://docs.rs/jaq-core/2.1.0/jaq_core/). It applies a JQ expression to a JSON input and returns the result.


## Example


> You: can you extract the title of the window from this JSON payload? 

```json
    {"widget": {
        "debug": "on",
        "window": {
            "title": "Sample Konfabulator Widget",
            "name": "main_window",
            "width": 500,
            "height": 500
        },
        "image": {
            "src": "Images/Sun.png",
            "name": "sun1",
            "hOffset": 250,
            "vOffset": 250,
            "alignment": "center"
        },
        "text": {
            "data": "Click Here",
            "size": 36,
            "style": "bold",
            "name": "text1",
            "hOffset": 250,
            "vOffset": 100,
            "alignment": "center",
            "onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
        }
    }}
```

> Assistant: The title of the window from the JSON payload is `"Sample Konfabulator Widget"`.