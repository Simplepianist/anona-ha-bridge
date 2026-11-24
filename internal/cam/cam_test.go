package cam

import (
	"reflect"
	"testing"
)

func TestParseEventTagsJSON(t *testing.T) {
	// Shape of a captured detection push; ids and timestamp are made up.
	detection := `{"businessType":3000001,"deviceChannel":96001001,"deviceId":"00000000000000000000000000000000","deviceModule":96001,"deviceName":"Küche","deviceType":96,"eventId":"AAAAAAAAAAAAAAAA","eventTs":1790000000000}`
	tests := []struct {
		name   string
		packet []byte
		want   []string
	}{
		{"json detection", []byte(detection), []string{"motion"}},
		{"json without eventId", []byte(`{"businessType":3000001}`), nil},
		{"broken json", []byte(`{"eventId":`), nil},
		{"protobuf status frame", []byte{0x08, 0x32, 0x10, 0x03}, nil},
		{"empty", nil, nil},
	}
	for _, tt := range tests {
		if got := ParseEventTags(tt.packet); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}
