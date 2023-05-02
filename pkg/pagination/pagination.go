package pagination

// why did you add pagination module in pkg?

// Page defines the page parameters
type Page struct {
	// Cursor describes the position in the database to start from
	Cursor int `json:"cursor"`

	// Limit describes the number of records per request
	Limit int `json:"limit"`

	//Filter describe filter params - by DeadlineAt or CreateAt
	Filter string `json:"filter"`

	//FilterOption describe filterOption params - DESC or ASC
	FilterOption string `json:"filterOption"`
}

// PageInfo describes the results page information
type PageInfo struct {
	// Page describes original request
	Page Page `json:"page"`

	// NextCursor describes the position of the next page
	NextCursor int `json:"nextCursor"`
}
