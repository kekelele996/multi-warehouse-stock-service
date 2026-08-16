package service

import (
	"testing"

	"warehousestock/internal/constants"
)

func TestU64AndItoa(t *testing.T) {
	if u64(0) != "0" || u64(99) != "99" {
		t.Error("u64 mismatch")
	}
	if itoa(42) != "42" {
		t.Error("itoa mismatch")
	}
}

func TestContains(t *testing.T) {
	if !contains(constants.UserRoleValues, constants.RoleOperator) {
		t.Error("operator should be in roles")
	}
	if contains(constants.UserRoleValues, "bogus") {
		t.Error("bogus should not be in roles")
	}
}

func TestGenOrderNo(t *testing.T) {
	for _, tc := range []struct{ t, prefix string }{
		{constants.OrderTypeInbound, "IN"},
		{constants.OrderTypeOutbound, "OUT"},
		{constants.OrderTypeTransfer, "TR"},
		{constants.OrderTypeInventoryCheck, "CK"},
	} {
		no := genOrderNo(tc.t)
		if len(no) < 4 || no[:len(tc.prefix)] != tc.prefix {
			t.Errorf("genOrderNo(%s) = %s, want prefix %s", tc.t, no, tc.prefix)
		}
	}
}
