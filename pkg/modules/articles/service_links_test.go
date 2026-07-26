package articles

import "testing"

func TestServiceLinkOffAndMsg(t *testing.T) {
	svc := map[string]ServiceView{
		"listing_submission": {On: true, Msg: ""},
		"ad_orders":          {On: false, Msg: "Временно недоступно"},
	}

	// An available service is not disabled.
	if serviceLinkOff(svc, "listing_submission") {
		t.Error("available service should not be off")
	}
	// A turned-off service is disabled and carries its notice.
	if !serviceLinkOff(svc, "ad_orders") {
		t.Error("off service should be reported off")
	}
	if got := serviceLinkMsg(svc, "ad_orders"); got != "Временно недоступно" {
		t.Errorf("msg = %q, want the notice", got)
	}
	// An unknown code is treated as available (never hide a link on a missing flag).
	if serviceLinkOff(svc, "nonexistent") {
		t.Error("unknown code must default to available")
	}
	if serviceLinkMsg(svc, "nonexistent") != "" {
		t.Error("unknown code has no message")
	}
	// A nil map is safe.
	if serviceLinkOff(nil, "ad_orders") {
		t.Error("nil svc map should default to available")
	}
}
