package bookmark

const (
	TypeBookmark  = 1
	TypeFolder    = 2
	TypeSeparator = 3

	DefaultIndex = -1

	MaxTagLength = 100

	GUIDRoot    = "root________"
	GUIDMenu    = "menu________"
	GUIDToolbar = "toolbar_____"
	GUIDUnfiled = "unfiled_____"
	GUIDMobile  = "mobile______"
	GUIDTag     = "tags________"

	GUIDVirtMenu    = "menu_______v"
	GUIDVirtToolbar = "toolbar____v"
	GUIDVirtUnfiled = "unfiled___v"
	GUIDVirtMobile  = "mobile____v"
)

type Item struct {
	GUID string `json:"guid"`

	ParentGUID string `json:"parentGuid"`

	Title string `json:"title"`

	Index int `json:"index"`

	DateAdded int64 `json:"dateAdded"`

	LastModified int64 `json:"lastModified"`

	ID int `json:"id"`

	TypeCode int `json:"typeCode"`

	Type string `json:"type"`

	Root string `json:"root"`

	Children []*Item `json:"children"`

	Annos   []Anno `json:"annos"`
	URI     string `json:"uri"`
	IconURI string `json:"iconuri"`
	Keyword string `json:"keyword"`
	Charset string `json:"charset"`
	Tags    string `json:"tags"`
}

type Anno struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Expires int    `json:"expires"`
	Flags   int    `json:"flags"`
}
