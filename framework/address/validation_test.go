package address_test

import (
	"testing"

	"github.com/jmcarbo/maddy/framework/address"
)

func TestValidMailboxName(t *testing.T) {
	if !address.ValidMailboxName("caddy.bug") {
		t.Error("caddy.bug should be valid mailbox name")
	}
}
