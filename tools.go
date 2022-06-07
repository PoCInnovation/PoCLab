package main

func formatAuthor(author string) string {
	if author[0] != '@' {
		return "@randomUser"
	}
	return author
}
