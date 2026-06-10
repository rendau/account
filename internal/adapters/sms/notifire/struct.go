package notifire

type RouteSt struct {
	Channel string `json:"channel" yaml:"channel"`
	Ttl     int    `json:"ttl,omitempty" yaml:"ttl"`
}

type SendReqSt struct {
	To     string    `json:"to"`
	Text   string    `json:"text"`
	Sync   bool      `json:"sync"`
	Source string    `json:"source,omitempty"`
	Route  []RouteSt `json:"route,omitempty"`
}

type SendRepSt struct {
	ID string `json:"id"`
}

type ErrorRepSt struct {
	ErrorCode string `json:"error_code"`
}
