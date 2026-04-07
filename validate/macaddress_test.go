package validate

import (
	"errors"
	"testing"
)

func TestMacAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		opts    []func(*MacAddressOptions)
		wantErr bool
	}{
		{name: "colon", value: "00:1a:2b:3c:4d:5e"},
		{name: "hyphen", value: "00-1a-2b-3c-4d-5e"},
		{name: "dot", value: "001a.2b3c.4d5e"},
		{name: "uppercase", value: "00:1A:2B:3C:4D:5E"},
		{name: "invalid char", value: "00:1g:2b:3c:4d:5e", wantErr: true},
		{name: "too short", value: "00:11:22:33:44", wantErr: true},
		{name: "eui64 rejected", value: "02:00:5e:10:00:00:00:01", wantErr: true},
		{name: "empty", value: "", wantErr: true},
		{name: "garbage", value: "not-a-mac", wantErr: true},
		{
			name:    "broadcast passes all",
			value:   "ff:ff:ff:ff:ff:ff",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "broadcast fails all_no_broadcast",
			value:   "ff:ff:ff:ff:ff:ff",
			opts:    []func(*MacAddressOptions){MacAddressType(MacAddressTypeAllNoBroadcast)},
			wantErr: true,
		},
		{
			name:  "broadcast type",
			value: "ff:ff:ff:ff:ff:ff",
			opts:  []func(*MacAddressOptions){MacAddressType(MacAddressTypeBroadcast)},
		},
		{
			name:    "unicast fails broadcast type",
			value:   "00:1a:2b:3c:4d:5e",
			opts:    []func(*MacAddressOptions){MacAddressType(MacAddressTypeBroadcast)},
			wantErr: true,
		},
		{
			name:  "local unicast",
			value: "02:00:00:00:00:01",
			opts:  []func(*MacAddressOptions){MacAddressType(MacAddressTypeLocalUnicast)},
		},
		{
			name:    "universal unicast rejects local",
			value:   "02:00:00:00:00:01",
			opts:    []func(*MacAddressOptions){MacAddressType(MacAddressTypeUniversalUnicast)},
			wantErr: true,
		},
		{
			name:  "universal unicast",
			value: "00:1a:2b:3c:4d:5e",
			opts:  []func(*MacAddressOptions){MacAddressType(MacAddressTypeUniversalUnicast)},
		},
		{
			name:    "local all rejects universal",
			value:   "00:1a:2b:3c:4d:5e",
			opts:    []func(*MacAddressOptions){MacAddressType(MacAddressTypeLocalAll)},
			wantErr: true,
		},
		{
			name:  "multicast all",
			value: "01:00:5e:00:00:01",
			opts:  []func(*MacAddressOptions){MacAddressType(MacAddressTypeMulticastAll)},
		},
		{
			name:    "unicast all rejects multicast",
			value:   "01:00:5e:00:00:01",
			opts:    []func(*MacAddressOptions){MacAddressType(MacAddressTypeUnicastAll)},
			wantErr: true,
		},
		{
			name:    "multicast no broadcast rejects broadcast",
			value:   "ff:ff:ff:ff:ff:ff",
			opts:    []func(*MacAddressOptions){MacAddressType(MacAddressTypeMulticastNoBroadcast)},
			wantErr: true,
		},
		{
			name:    "unknown type",
			value:   "00:1a:2b:3c:4d:5e",
			opts:    []func(*MacAddressOptions){MacAddressType("bogus")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := MacAddress(tt.value, tt.opts...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("MacAddress(%q): want error", tt.value)
				}
				if !errors.Is(err, ErrInvalidMAC) {
					t.Fatalf("MacAddress(%q): got %v, want %v", tt.value, err, ErrInvalidMAC)
				}
			} else if err != nil {
				t.Fatalf("MacAddress(%q): %v", tt.value, err)
			}
		})
	}
}
