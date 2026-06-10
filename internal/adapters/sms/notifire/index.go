package notifire

import (
	"github.com/rendau/dop/adapters/client/httpc"
)

type St struct {
	httpc  httpc.HttpC
	source string
	route  []RouteSt
}

// New creates the notifire sms-adapter. The httpc base-uri is expected to
// already contain the full send path with the account var (e.g. .../send/{account}).
// source and route are optional and are only sent when set.
func New(httpc httpc.HttpC, source string, route []RouteSt) *St {
	return &St{
		httpc:  httpc,
		source: source,
		route:  route,
	}
}

func (s *St) Send(phone string, msg string) bool {
	return s.send(&SendReqSt{
		To:   phone,
		Text: msg,
		Sync: true,
	})
}

func (s *St) SendAsync(phone string, msg string) bool {
	return s.send(&SendReqSt{
		To:   phone,
		Text: msg,
		Sync: false,
	})
}

func (s *St) send(req *SendReqSt) bool {
	req.Source = s.source
	req.Route = s.route

	_, err := s.httpc.Send(&httpc.OptionsSt{
		Method: "POST",
		Uri:    "send",

		ReqObj: req,
	})

	return err == nil
}
