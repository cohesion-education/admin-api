package cohesioned

import "time"

//Homepage the model for my homepage
type Homepage struct {
	Auditable
	Header       *Header       `json:"header" datastore:"header"`
	Features     *Features     `json:"features" datastore:"features"`
	Testimonials *Testimonials `json:"testimonials" datastore:"testimonials"`
	Pricing      *Pricing      `json:"pricing" datastore:"pricing"`
	// SocialMediaLinks      []SocialMediaLink `datastore:"social_media_links" json:"social_media_links"`
}

func NewHomepage(id int64) *Homepage {
	h := &Homepage{}
	h.Header = &Header{}
	h.Features = &Features{Highlights: []*Highlight{}}
	h.Testimonials = &Testimonials{List: []*Testimonial{}}
	h.Pricing = &Pricing{List: []*PricingDetail{}}

	h.GCPPersisted.id = id
	h.Auditable.Created = time.Now()

	return h
}

type Header struct {
	Title    string `datastore:"title" json:"title"`
	Subtitle string `datastore:"subtitle" json:"subtitle"`
}

type Features struct {
	Title      string       `datastore:"title" json:"title"`
	Subtitle   string       `datastore:"subtitle" json:"subtitle"`
	Highlights []*Highlight `datastore:"higlights" json:"highlights"`
}

type Highlight struct {
	Auditable
	Title           string `datastore:"title" json:"title"`
	Description     string `datastore:"description" json:"description"`
	FaIconClassName string `datastore:"fa_icon_classname" json:"iconClassName"`
}

type Testimonials struct {
	List []*Testimonial `json:"list" datastore:"list"`
}

type Testimonial struct {
	Auditable
	Blurb     string `datastore:"blurb" json:"text"`
	FullName  string `datastore:"fullname" json:"name"`
	AvatarURL string `datastore:"avatar" json:"avatar"`
}

type Pricing struct {
	Title    string           `datastore:"title" json:"title"`
	Subtitle string           `datastore:"subtitle" json:"subtitle"`
	List     []*PricingDetail `json:"list" datastore:"list"`
}

type PricingDetail struct {
	Auditable
	Title    string `json:"title"`
	Price    string `datastore:"price" json:"price"`
	Duration string `datastore:"duration" json:"duration"`
}

type SocialMediaLink struct {
	Auditable
	//TODO - is there an enum type?
	Type    string `datastore:"type" json:"type"`
	LinkURL string `datastore:"url" json:"url"`
}
