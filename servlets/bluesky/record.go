package main

type Record struct {
	Type      string  `json:"$type"`
	Text      string  `json:"text"`
	Facets    []Facet `json:"facets,omitempty"`
	CreatedAt string  `json:"createdAt"`
	Reply     *Reply  `json:"reply,omitempty"`
}

type Reply struct {
	Root   Post `json:"root"`
	Parent Post `json:"parent"`
}

type Post struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}
