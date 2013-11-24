package magnet

import r "github.com/christopherhesse/rethinkgo"

type Tag struct {
	Title   string
	Count   int
	User    string
	Created int32
}

func GetTags(dbSession *r.Session, userId string) []Tag {
	var tags []Tag

	r.Db("magnet").
		Table("tags").
		Filter(r.Row.Attr("User").
		Eq(userId)).
		OrderBy("id").
		Run(dbSession).
		All(&tags)

	return tags
}
