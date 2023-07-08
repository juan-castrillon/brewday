package main

import (
	"brewday/internal/notifications"
	"log"
)

type Click struct {
	Url string `json:"url"`
}

type ClientNot struct {
	Click       Click  `json:"click"`
	BigImageURL string `json:"bigImageUrl"`
}

func main() {
	app_token := "Aopj09R0IeCqD6r"
	gotify_url := "http://localhost:8080"
	n := notifications.NewNotifier(app_token, gotify_url)
	err := n.Send("**Click me!**\n- List1\n- List 2", "Test")
	if err != nil {
		log.Fatal(err)
	}
	// a, err := json.MarshalIndent(extras, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(a))
}
