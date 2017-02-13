package domain

import (
	"bytes"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encoderType"
	"github.com/distributed-vision/go-resources/version/versionType"
)

type idTest struct {
	scopeId     []byte
	idRoot      []byte
	incarnation *uint32
	crcLength   uint
	versionType versionType.VersionType
	encodedId   string
}

var incarnation0 = uint32(0)
var incarnation1 = uint32(1)
var incarnation1011 = uint32(1011)
var incarnation201011 = uint32(201011)

var idTests = []idTest{
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versionType.UNVERSIONED, "H7C"},
	{base62.MustDecode("1"), base62.MustDecode("12"), nil, 0, versionType.UNVERSIONED, "H8C"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), nil, 0, versionType.UNVERSIONED, "4keRnR"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versionType.UNVERSIONED, "5tbciO"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 0, versionType.UNVERSIONED, "5tbciP"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1011, 0, versionType.UNVERSIONED, "OQJ9VFj"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation201011, 0, versionType.UNVERSIONED, "6kzTmoi1Vj"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation0, 0, versionType.UNVERSIONED, "1d9wvIxu4"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation1, 0, versionType.UNVERSIONED, "1d9wvIxu5"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation1011, 0, versionType.UNVERSIONED, "6kzKBJzG9T"},
	{base62.MustDecode("1"), base62.MustDecode("a12b"), &incarnation201011, 0, versionType.UNVERSIONED, "1rvLheMVxt1pz"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 8, versionType.UNVERSIONED, "1QGxu"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 16, versionType.UNVERSIONED, "1QHEQ"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 32, versionType.UNVERSIONED, "1QHUw"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versionType.NUMERIC, "1QHlS"},
	{base62.MustDecode("1"), base62.MustDecode("2"), nil, 0, versionType.SEMANTIC, "1QIpW"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 8, versionType.UNVERSIONED, "5tciuX"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 16, versionType.UNVERSIONED, "5tdp6f"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 32, versionType.UNVERSIONED, "5tevIn"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versionType.NUMERIC, "5tg1Uu"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation0, 0, versionType.SEMANTIC, "5tkQHQ"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 32, versionType.NUMERIC, "5tjK5J"},
	{base62.MustDecode("1"), base62.MustDecode("2"), &incarnation1, 8, versionType.SEMANTIC, "5tlWTZ"},
}

func TestDomainIdFormatting(t *testing.T) {

	for _, test := range idTests {
		id, err := ToId(test.scopeId, test.idRoot, test.incarnation, test.crcLength, test.versionType)

		if err != nil {
			t.Errorf("TestDomainIdFormatting: ToId failed with err %s:", err)
		}

		if base62.Encode(id) != test.encodedId {
			t.Errorf("TestDomainIdFormatting: ToId failed expected %s: got %s", test.encodedId, base62.Encode(id))
		}

		if !bytes.Equal(id[:ScopeLength(id)], test.scopeId) {
			t.Errorf("TestDomainIdFormatting: ToId failed expected scopeId %v: got %v", test.scopeId, id[:ScopeLength(id)])
		}

		if !bytes.Equal(test.idRoot, IdRootValue(id)) {
			t.Errorf("TestDomainIdFormatting: ToId failed expected idRoot %v: got %v", test.idRoot, IdRootValue(id))
		}

		if test.incarnation == nil {
			if IncarnationValue(id) != nil {
				t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", test.incarnation, *IncarnationValue(id))
			}
		} else {
			if IncarnationValue(id) == nil {
				t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", test.incarnation, IncarnationValue(id))
			} else {
				if *test.incarnation != *IncarnationValue(id) {
					t.Errorf("TestDomainIdFormatting: ToId failed expected incarnation %v: got %v", *test.incarnation, *IncarnationValue(id))
				}
			}
		}

		crcLengthValue, err := CrcLengthValue(id)

		if err != nil {
			t.Errorf("TestDomainIdFormatting: ToId failed with err %s:", err)
		}

		if test.crcLength != crcLengthValue {
			t.Errorf("TestDomainIdFormatting: ToId failed expected crcLength %v: got %v", test.crcLength, crcLengthValue)
		}

		versionTypeValue, err := VersionTypeValue(id)

		if err != nil {
			t.Errorf("TestDomainIdFormatting: ToId failed with err %s:", err)
		}

		if test.versionType != versionTypeValue {
			t.Errorf("TestDomainIdFormatting: ToId failed expected versionTypeValue %v: got %v", test.versionType, versionTypeValue)
		}
	}
}

func TestGoTypeIdDomainDefinition(t *testing.T) {
	id, err := DecodeId(encoderType.BASE62, "2", "0")

	if err != nil {
		t.Errorf("TestDomainIdFormatting: ToId failed with err %s:", err)
	}

	if base62.Encode(id) != "YAC" {
		t.Errorf("TestDomainIdFormatting: ToId failed expected %s: got %s", "YAC", base62.Encode(id))
	}

	id = MustDecodeId(encoderType.BASE62, "2", "0")

	if base62.Encode(id) != "YAC" {
		t.Errorf("TestDomainIdFormatting: ToId failed expected %s: got %s", "YAC", base62.Encode(id))
	}

	incarnation := uint32(0)

	id, err = DecodeId(encoderType.BASE62, "T", "0", &incarnation)

	if base62.Encode(id) != "2DAEU1w" {
		t.Errorf("TestDomainIdFormatting: ToId failed expected %s: got %s", "2DAEU1w", base62.Encode(id))
	}

}
