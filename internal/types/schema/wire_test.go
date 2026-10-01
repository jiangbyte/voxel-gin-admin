package schema_test

import (
	"encoding/json"
	"testing"

	"voxel-gin-admin/internal/types/schema"
)

func TestWireIntMarshal(t *testing.T) {
	raw, err := json.Marshal(schema.WireIntValue(300))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `"300"` {
		t.Fatalf("got %s", raw)
	}
	raw, err = json.Marshal(schema.WireFlagValue(1))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `"1"` {
		t.Fatalf("flag got %s", raw)
	}
}
