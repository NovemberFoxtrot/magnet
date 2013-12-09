package main

import (
	r "github.com/christopherhesse/rethinkgo"
	m "github.com/codegangsta/martini"
	s "github.com/gorilla/sessions"
	h "net/http"
	"strconv"
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

// GetTagHandler fetches books for a given tag
func GetTagHandler(params m.Params, req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	_, userID := GetUserData(cs, req)
	var response []interface{}
	page, _ := strconv.ParseInt(params["page"], 10, 16)

	err := r.Db("magnet").
		Table("bookmarks").
		Filter(r.Row.Attr("User").
		Eq(userID).
		And(r.Row.Attr("Tags").
		Contains(params["tag"]))).
		OrderBy(r.Desc("Created")).
		Skip(50 * page).
		Limit(50).
		Run(dbSession).
		All(&response)

	if err != nil {
		WriteJSONResponse(200, true, "Error getting bookmarks for tag "+params["tag"], req, w)
	} else {
		JSONDataResponse(200, false, response, req, w)
	}
}
