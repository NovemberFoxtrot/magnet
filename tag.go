package main

// Tag for JSON
type Tag struct {
	Name  string
	Count int
}

// GetTags fetches tags from rethinkdb
func GetTags(connection *Connection, userID string) []Tag {
	tagMap := make(map[string]int)
	var tags []Tag

	response, err := connection.GetTags(userID)

	if err == nil {
		// Search por repeated tags and count them
		for _, tagsMap := range response {
			for _, tag := range tagsMap.(map[string]interface{})["Tags"].([]interface{}) {
				if _, ok := tagMap[tag.(string)]; ok {
					tagMap[tag.(string)]++
				} else {
					tagMap[tag.(string)] = 1
				}
			}
		}

		// Then put them in a tag map
		tags = make([]Tag, len(tagMap))
		i := 0
		for tag, count := range tagMap {
			tags[i] = Tag{Name: tag, Count: count}
			i++
		}
	}

	return tags
}
