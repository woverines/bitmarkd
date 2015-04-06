// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package block

// the starting block number
const (
	GenesisBlockNumber = 1
)

// this is block 1, the Genesis Block
// ----------------------------------

// LIVE Network
// ------------

var LiveGenesisBlock = Packed([]byte{
	0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x3d, 0xd5, 0x0d, 0x37,
	0xce, 0x54, 0xcb, 0x0e, 0x18, 0xd7, 0x94, 0x7d,
	0x02, 0xf7, 0x4f, 0x6d, 0xa6, 0xfc, 0xb6, 0x29,
	0xe5, 0xd1, 0xef, 0x5f, 0xa6, 0x1d, 0xa1, 0x0e,
	0x3e, 0xa7, 0x88, 0xaa, 0x4b, 0x42, 0x78, 0x54,
	0xff, 0xff, 0x00, 0x1d, 0x74, 0x5b, 0x29, 0x70,
	0x65, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xff, 0xff, 0x11, 0x02, 0x01, 0x00, 0x04,
	0x4b, 0x42, 0x78, 0x54, 0x08, 0x42, 0x4d, 0x52,
	0x4b, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x18, 0x6a, 0x00, 0x15, 0x42, 0x69,
	0x74, 0x6d, 0x61, 0x72, 0x6b, 0x20, 0x47, 0x65,
	0x6e, 0x65, 0x73, 0x69, 0x73, 0x20, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x00, 0x00, 0x00, 0x00, 0x01,
	0x00, 0x3d, 0xd5, 0x0d, 0x37, 0xce, 0x54, 0xcb,
	0x0e, 0x18, 0xd7, 0x94, 0x7d, 0x02, 0xf7, 0x4f,
	0x6d, 0xa6, 0xfc, 0xb6, 0x29, 0xe5, 0xd1, 0xef,
	0x5f, 0xa6, 0x1d, 0xa1, 0x0e, 0x3e, 0xa7, 0x88,
	0xaa,
})

// digest of the live genesis header
// 00000000122e6bcab6746c9b6328aaed08080dcc3970a617895ca1183c740a0a
var LiveGenesisDigest = Digest([...]byte{
	0x0a, 0x0a, 0x74, 0x3c, 0x18, 0xa1, 0x5c, 0x89,
	0x17, 0xa6, 0x70, 0x39, 0xcc, 0x0d, 0x08, 0x08,
	0xed, 0xaa, 0x28, 0x63, 0x9b, 0x6c, 0x74, 0xb6,
	0xca, 0x6b, 0x2e, 0x12, 0x00, 0x00, 0x00, 0x00,
})

// TEST Network
// ------------

var TestGenesisBlock = Packed([]byte{
	0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xaa, 0x5e, 0x56, 0xa6,
	0xd8, 0x3f, 0xd0, 0xfd, 0xc7, 0x58, 0x3f, 0x33,
	0xb1, 0x17, 0xb5, 0x4d, 0xfb, 0x30, 0x12, 0x1c,
	0xa9, 0x2d, 0xcc, 0xf3, 0x34, 0xc3, 0x5b, 0x67,
	0x8c, 0xdc, 0xe9, 0x3c, 0x4b, 0x42, 0x78, 0x54,
	0xff, 0xff, 0x00, 0x1d, 0x3b, 0xa8, 0xcc, 0xae,
	0x6d, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xff, 0xff, 0x11, 0x02, 0x01, 0x00, 0x04,
	0x4b, 0x42, 0x78, 0x54, 0x08, 0x42, 0x4d, 0x52,
	0x4b, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x20, 0x6a, 0x00, 0x1d, 0x42, 0x69,
	0x74, 0x6d, 0x61, 0x72, 0x6b, 0x20, 0x54, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x47, 0x65,
	0x6e, 0x65, 0x73, 0x69, 0x73, 0x20, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x00, 0x00, 0x00, 0x00, 0x01,
	0x00, 0xaa, 0x5e, 0x56, 0xa6, 0xd8, 0x3f, 0xd0,
	0xfd, 0xc7, 0x58, 0x3f, 0x33, 0xb1, 0x17, 0xb5,
	0x4d, 0xfb, 0x30, 0x12, 0x1c, 0xa9, 0x2d, 0xcc,
	0xf3, 0x34, 0xc3, 0x5b, 0x67, 0x8c, 0xdc, 0xe9,
	0x3c,
})

// digest of the test genesis header
// 00000000706b725a931946298216ae21453efd0b2ea3575bf97cee00035da020
var TestGenesisDigest = Digest([...]byte{
	0x20, 0xa0, 0x5d, 0x03, 0x00, 0xee, 0x7c, 0xf9,
	0x5b, 0x57, 0xa3, 0x2e, 0x0b, 0xfd, 0x3e, 0x45,
	0x21, 0xae, 0x16, 0x82, 0x29, 0x46, 0x19, 0x93,
	0x5a, 0x72, 0x6b, 0x70, 0x00, 0x00, 0x00, 0x00,
})
