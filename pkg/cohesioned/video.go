package cohesioned

import "time"

type Video struct {
	id       int64
	Created  time.Time `datastore:"created" json:"created"`
	Title    string
	Category Taxonomy
	//TODO - Teacher, Tags, Related Videos, FAQs
}
