package util

import (
	"testing"
	"time"

	"warehousestock/internal/constants"
)

func TestGenerateTokenAndParse(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateToken(secret, time.Hour, 3, "13800000001", constants.RoleAdmin)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != 3 || claims.Phone != "13800000001" || claims.Role != constants.RoleAdmin {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	for _, tc := range []string{"", "x", "a.b.c"} {
		if _, err := ParseToken("s", tc); err == nil {
			t.Errorf("expected error for %q", tc)
		}
	}
}

func TestFormatters(t *testing.T) {
	if OrderTypeText(constants.OrderTypeInbound) != "入库" {
		t.Error("order type text mismatch")
	}
	if OrderStatusText(constants.OrderStatusCompleted) != "已完成" {
		t.Error("order status text mismatch")
	}
	if WarehouseStatusText("active") != "启用" {
		t.Error("warehouse status text mismatch")
	}
	if got := FormatStockNumber(5, "件"); got != "5 件" {
		t.Errorf("FormatStockNumber = %q", got)
	}
}

func TestOrderValidators(t *testing.T) {
	if !constants.IsValidOrderType(constants.OrderTypeTransfer) {
		t.Error("transfer should be valid")
	}
	if constants.IsValidOrderType("bogus") {
		t.Error("bogus should be invalid")
	}
	if !constants.IsValidOrderStatus(constants.OrderStatusProcessing) {
		t.Error("processing should be valid")
	}
}
