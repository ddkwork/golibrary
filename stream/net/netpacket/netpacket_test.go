package netpacket

type (
	Doc struct {
		Api      string
		Function string
		Note     string
		Todo     string
		Chinese  string
	}
	doc struct{ infos []Doc }
)
