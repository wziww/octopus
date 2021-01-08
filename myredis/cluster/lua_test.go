package cluster

import (
	"testing"
)

func TestLua(t *testing.T) {
	if AddSlotsLua(0, 999) != "for i=0,999, 1 do redis.call(\"CLUSTER\",\"ADDSLOTS\",i)  end return 'success'" {
		t.Fatal("lua error")
	}
}
