package articles

import (
	"shanraq.org/pkg/site"
	"testing"
)

func TestServiceLinkOffAndMsg(t *testing.T) {
	svc := map[string]ServiceView{
		"listing_submission": {On: true, Msg: ""},
		"ad_orders":          {On: false, Msg: "Временно недоступно"},
	}

	// An available service is not disabled.
	if site.ServiceOff(svc, "listing_submission") {
		t.Error("available service should not be off")
	}
	// A turned-off service is disabled and carries its notice.
	if !site.ServiceOff(svc, "ad_orders") {
		t.Error("off service should be reported off")
	}
	if got := site.ServiceMsg(svc, "ad_orders"); got != "Временно недоступно" {
		t.Errorf("msg = %q, want the notice", got)
	}
	// An unknown code is treated as available (never hide a link on a missing flag).
	if site.ServiceOff(svc, "nonexistent") {
		t.Error("unknown code must default to available")
	}
	if site.ServiceMsg(svc, "nonexistent") != "" {
		t.Error("unknown code has no message")
	}
	// A nil map is safe.
	if site.ServiceOff(nil, "ad_orders") {
		t.Error("nil svc map should default to available")
	}
}
