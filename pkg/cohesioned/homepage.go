package cohesioned

type Homepage struct {
	*Auditable
	HeaderTagline         string            `datastore:"page_header_tagline" json:"page_header_tagline"`
	HeaderSubtext         string            `datastore:"page_header_subtext" json:"page_header_subtext"`
	FeaturesHeaderTagline string            `datastore:"features_header_tagline" json:"features_header_tagline"`
	FeaturesHeaderSubtext string            `datastore:"features_header_subtext" json:"features_header_subtext"`
	Highlights            []Highlight       `datastore:"higlights" json:"highlights"`
	Features              []Feature         `datastore:"features" json:"features"`
	Reviews               []Review          `datastore:"reviews" json:"reviews"`
	Pricing               []PricingDetails  `datastore:"pricing" json:"pricing"`
	SocialMediaLinks      []SocialMediaLink `datastore:"social_media_links" json:"social_media_links"`
}

type Highlight struct {
	*Auditable
	Description string `datastore:"description" json:"description"`
	Subtext     string `datastore:"subtext" json:"subtext"`
}

type Feature struct {
	*Auditable
	Description string `datastore:"description" json:"description"`
	Subtext     string `datastore:"subtext" json:"subtext"`
	//TODO - Image       []byte
}

type Review struct {
	*Auditable
	Blurb     string `datastore:"blurb" json:"blurb"`
	FullName  string `datastore:"fullname" json:"fullname"`
	AvatarURL string `datastore:"avatar" json:"avatar"`
}

type PricingDetails struct {
	*Auditable
	Description string   `datastore:"description" json:"description"`
	Cost        string   `datastore:"cost" json:"cost"`
	Frequency   string   `datastore:"frequency" json:"frequency"`
	Details     []string `datastore:"details" json:"details"`
}

type SocialMediaLink struct {
	*Auditable
	//TODO - is there an enum type?
	Type    string `datastore:"type" json:"type"`
	LinkURL string `datastore:"url" json:"url"`
}
