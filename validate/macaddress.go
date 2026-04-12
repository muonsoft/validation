package validate

import (
	"bytes"
	"errors"
	"net"
)

// ErrInvalidMAC is returned when the value is not a valid 48-bit MAC address for the
// configured checks, or when an unknown Symfony type name was set via [WithMacAddressType].
var ErrInvalidMAC = errors.New("invalid MAC address")

// MacAddressType represents Symfony\Component\Validator\Constraints\MacAddress "type".
type MacAddressType string

// MAC address type values match Symfony\Component\Validator\Constraints\MacAddress.
const (
	MacAddressTypeAll                       MacAddressType = "all"
	MacAddressTypeAllNoBroadcast            MacAddressType = "all_no_broadcast"
	MacAddressTypeLocalAll                  MacAddressType = "local_all"
	MacAddressTypeLocalNoBroadcast          MacAddressType = "local_no_broadcast"
	MacAddressTypeLocalUnicast              MacAddressType = "local_unicast"
	MacAddressTypeLocalMulticast            MacAddressType = "local_multicast"
	MacAddressTypeLocalMulticastNoBroadcast MacAddressType = "local_multicast_no_broadcast"
	MacAddressTypeUniversalAll              MacAddressType = "universal_all"
	MacAddressTypeUniversalUnicast          MacAddressType = "universal_unicast"
	MacAddressTypeUniversalMulticast        MacAddressType = "universal_multicast"
	MacAddressTypeUnicastAll                MacAddressType = "unicast_all"
	MacAddressTypeMulticastAll              MacAddressType = "multicast_all"
	MacAddressTypeMulticastNoBroadcast      MacAddressType = "multicast_no_broadcast"
	MacAddressTypeBroadcast                 MacAddressType = "broadcast"
)

// MacAddressOptions configures [MacAddress] validation (Symfony MacAddress "type").
type MacAddressOptions struct {
	Type MacAddressType
}

func newMacAddressOptions() *MacAddressOptions {
	return &MacAddressOptions{Type: MacAddressTypeAll}
}

// WithMacAddressType sets the Symfony-compatible MAC class filter (default [MacAddressTypeAll]).
func WithMacAddressType(t MacAddressType) func(*MacAddressOptions) {
	return func(o *MacAddressOptions) {
		o.Type = t
	}
}

// MacAddress validates that value is a 48-bit IEEE 802 MAC address accepted by [net.ParseMAC]
// (colon, hyphen, or dot-separated forms), then applies the Symfony [MacAddress] "type" rules
// on the first octet’s I/G and U/L bits and on the broadcast address ff:ff:ff:ff:ff:ff.
//
// Returns [ErrInvalidMAC] when the string is not a 48-bit MAC, when it is not parseable, or
// when an unknown type name was set via [WithMacAddressType].
func MacAddress(value string, options ...func(*MacAddressOptions)) error {
	opts := newMacAddressOptions()
	for _, opt := range options {
		opt(opts)
	}
	if _, ok := macAddressTypePredicates[opts.Type]; !ok {
		return ErrInvalidMAC
	}
	hw, err := net.ParseMAC(value)
	if err != nil || len(hw) != 6 {
		return ErrInvalidMAC
	}
	if !macAddressMatchesType(hw, opts.Type) {
		return ErrInvalidMAC
	}
	return nil
}

var macBroadcast6 = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

// macAddressClass holds I/G (unicast vs multicast), U/L (universal vs local), and broadcast
// flags derived from a 48-bit MAC, matching Symfony MacAddressValidator.
type macAddressClass struct {
	unicast   bool
	local     bool
	broadcast bool
}

func macAddressClassOf(hw net.HardwareAddr) macAddressClass {
	first := hw[0]
	return macAddressClass{
		unicast:   first&1 == 0,
		local:     first&2 != 0,
		broadcast: bytes.Equal(hw, macBroadcast6),
	}
}

// macAddressTypePredicates maps Symfony MacAddress "type" to a predicate on [macAddressClass].
var macAddressTypePredicates = map[MacAddressType]func(macAddressClass) bool{
	MacAddressTypeAll: func(c macAddressClass) bool {
		return true
	},
	MacAddressTypeAllNoBroadcast: func(c macAddressClass) bool {
		return !c.broadcast
	},
	MacAddressTypeLocalAll: func(c macAddressClass) bool {
		return c.local
	},
	MacAddressTypeLocalNoBroadcast: func(c macAddressClass) bool {
		return c.local && !c.broadcast
	},
	MacAddressTypeLocalUnicast: func(c macAddressClass) bool {
		return c.local && c.unicast
	},
	MacAddressTypeLocalMulticast: func(c macAddressClass) bool {
		return c.local && !c.unicast
	},
	MacAddressTypeLocalMulticastNoBroadcast: func(c macAddressClass) bool {
		return c.local && !c.unicast && !c.broadcast
	},
	MacAddressTypeUniversalAll: func(c macAddressClass) bool {
		return !c.local
	},
	MacAddressTypeUniversalUnicast: func(c macAddressClass) bool {
		return !c.local && c.unicast
	},
	MacAddressTypeUniversalMulticast: func(c macAddressClass) bool {
		return !c.local && !c.unicast
	},
	MacAddressTypeUnicastAll: func(c macAddressClass) bool {
		return c.unicast
	},
	MacAddressTypeMulticastAll: func(c macAddressClass) bool {
		return !c.unicast
	},
	MacAddressTypeMulticastNoBroadcast: func(c macAddressClass) bool {
		return !c.unicast && !c.broadcast
	},
	MacAddressTypeBroadcast: func(c macAddressClass) bool {
		return c.broadcast
	},
}

func macAddressMatchesType(hw net.HardwareAddr, typ MacAddressType) bool {
	pred, ok := macAddressTypePredicates[typ]
	if !ok {
		return false
	}
	return pred(macAddressClassOf(hw))
}
