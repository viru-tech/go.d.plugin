// SPDX-License-Identifier: GPL-3.0-or-later

package iprange

import (
	"bytes"
	"fmt"
	"math/big"
	"net"
)

// Family represents IP Range address-family.
type Family uint8

const (
	// V4Family is IPv4 address-family.
	V4Family Family = iota
	// V6Family is IPv6 address-family.
	V6Family
)

// Range represents an IP range.
type Range interface {
	Family() Family
	Contains(ip net.IP) bool
	Size() *big.Int
	fmt.Stringer
}

// New returns new IP Range.
// If it is not a valid range (start and end IPs have different address-families, or start > end),
// New returns nil.
func New(start, end net.IP) Range {
	if isV4RangeValid(start, end) {
		return v4Range{Start: start, End: end}
	}
	if isV6RangeValid(start, end) {
		return v6Range{Start: start, End: end}
	}
	return nil
}

type v4Range struct {
	Start net.IP
	End   net.IP
}

// String returns the string form of the range.
func (r v4Range) String() string {
	return fmt.Sprintf("%s-%s", r.Start, r.End)
}

// Family returns the range address family.
func (r v4Range) Family() Family {
	return V4Family
}

// Contains reports whether the range includes IP.
func (r v4Range) Contains(ip net.IP) bool {
	return bytes.Compare(ip, r.Start) >= 0 && bytes.Compare(ip, r.End) <= 0
}

// Size reports the number of IP addresses in the range.
func (r v4Range) Size() *big.Int {
	return big.NewInt(v4ToInt(r.End) - v4ToInt(r.Start) + 1)
}

type v6Range struct {
	Start net.IP
	End   net.IP
}

// String returns the string form of the range.
func (r v6Range) String() string {
	return fmt.Sprintf("%s-%s", r.Start, r.End)
}

// Family returns the range address family.
func (r v6Range) Family() Family {
	return V6Family
}

// Contains reports whether the range includes IP.
func (r v6Range) Contains(ip net.IP) bool {
	return bytes.Compare(ip, r.Start) >= 0 && bytes.Compare(ip, r.End) <= 0
}

// Size reports the number of IP addresses in the range.
func (r v6Range) Size() *big.Int {
	size := big.NewInt(0)
	size.Add(size, big.NewInt(0).SetBytes(r.End))
	size.Sub(size, big.NewInt(0).SetBytes(r.Start))
	size.Add(size, big.NewInt(1))
	return size
}

func v4ToInt(ip net.IP) int64 {
	ip = ip.To4()
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}
