package seeds

import (
	"base/app/stories"
)

func DefaultStories() []stories.Story {
	return []stories.Story{
		{
			Title:       "The Hobbit",
			Description: "The Hobbit, or There and Back Again is a children's fantasy novel by English author J. R. R. Tolkien. It was published on 21 September 1937 to wide critical acclaim, being nominated for the Carnegie Medal and awarded a prize from the New York Herald Tribune for best juvenile fiction.",
			Cover:       "https://upload.wikimedia.org/wikipedia/en/3/30/Hobbit_cover.JPG",
			UserID:      1,
			Category:    "Fantasy",
		},
		{
			Title:       "The Lord of the Rings",
			Description: "The Lord of the Rings is an epic high-fantasy novel by English author and scholar J. R. R. Tolkien. Set in Middle-earth, the world at some distant time in the past, the story began as a sequel to Tolkien's 1937 children's book The Hobbit, but eventually developed into a much larger work.",
			Cover:       "https://upload.wikimedia.org/wikipedia/en/8/8e/Lord_of_the_Rings_-_The_Fellowship_of_the_Ring.jpg",
			UserID:      1,
			Category:    "Fantasy",
		},
	}
}
