package main

import (
	r "github.com/christopherhesse/rethinkgo"
)

// Tag for JSON
type Tag struct {
	Name  string
	Count int
}

// GetTags fetches tags from rethinkdb
func GetTags(dbSession *r.Session, userID string) []Tag {
	var response []interface{}
	tagMap := make(map[string]int)
	var tags []Tag

	err := r.Db("magnet").
		Table("bookmarks").
		Filter(r.Row.Attr("User").
		Eq(userID)).
		WithFields("Tags").
		Run(dbSession).
		All(&response)

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
