package main

type Wrap struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image"`
}

type Wraps []Wrap
