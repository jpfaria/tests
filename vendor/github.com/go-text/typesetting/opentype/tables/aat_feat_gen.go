// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from aat_feat_src.go. DO NOT EDIT

func (item *FeatureSettingName) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.Setting = binary.BigEndian.Uint16(src[0:])
	item.NameIndex = binary.BigEndian.Uint16(src[2:])
}

func ParseFeat(src []byte) (Feat, int, error) {
	var item Feat
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading Feat: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.version = binary.BigEndian.Uint32(src[0:])
	item.featureNameCount = binary.BigEndian.Uint16(src[4:])
	item.none1 = binary.BigEndian.Uint16(src[6:])
	item.none2 = binary.BigEndian.Uint32(src[8:])
	n += 12

	{
		arrayLength := int(item.featureNameCount)

		offset := 12
		for i := 0; i < arrayLength; i++ {
			elem, read, err := ParseFeatureName(src[offset:], src)
			if err != nil {
				return item, 0, fmt.Errorf("reading Feat: %s", err)
			}
			item.Names = append(item.Names, elem)
			offset += read
		}
		n = offset
	}
	return item, n, nil
}

func ParseFeatureName(src []byte, parentSrc []byte) (FeatureName, int, error) {
	var item FeatureName
	n := 0
	if L := len(src); L < 12 {
		return item, 0, fmt.Errorf("reading FeatureName: "+"EOF: expected length: 12, got %d", L)
	}
	_ = src[11] // early bound checking
	item.Feature = binary.BigEndian.Uint16(src[0:])
	item.nSettings = binary.BigEndian.Uint16(src[2:])
	offsetSettingTable := int(binary.BigEndian.Uint32(src[4:]))
	item.FeatureFlags = binary.BigEndian.Uint16(src[8:])
	item.NameIndex = binary.BigEndian.Uint16(src[10:])
	n += 12

	{

		if offsetSettingTable != 0 { // ignore null offset
			if L := len(parentSrc); L < offsetSettingTable {
				return item, 0, fmt.Errorf("reading FeatureName: "+"EOF: expected length: %d, got %d", offsetSettingTable, L)
			}

			arrayLength := int(item.nSettings)

			if L := len(parentSrc); L < offsetSettingTable+arrayLength*4 {
				return item, 0, fmt.Errorf("reading FeatureName: "+"EOF: expected length: %d, got %d", offsetSettingTable+arrayLength*4, L)
			}

			item.SettingTable = make([]FeatureSettingName, arrayLength) // allocation guarded by the previous check
			for i := range item.SettingTable {
				item.SettingTable[i].mustParse(parentSrc[offsetSettingTable+i*4:])
			}
			offsetSettingTable += arrayLength * 4
		}
	}
	return item, n, nil
}
