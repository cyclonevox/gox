package gox

import (
	"bytes"
	"reflect"
	"testing"
)

// General Test funcs
func TestNewSnowflake(t *testing.T) {
	_, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	_, err = NewSnowflake(5000)
	if err == nil {
		t.Fatalf("no error creating NewSnowflake, %s", err)
	}

}

// lazy check if Next will create duplicate Ids would be good to later enhance this with more smarts
func TestNextDuplicateId(t *testing.T) {
	node, _ := NewSnowflake(1)

	var x, y Id
	for i := 0; i < 1000000; i++ {
		y = node.Next()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

// I feel like there's probably a better way
func TestRace(t *testing.T) {
	node, _ := NewSnowflake(1)

	go func() {
		for i := 0; i < 1000000000; i++ {

			NewSnowflake(1)
		}
	}()

	for i := 0; i < 4000; i++ {

		node.Next()
	}

}

// Converters/Parsers Test funcs
// We should have funcs here to test conversion both ways for everything
func TestPrintAll(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	id := node.Next()

	t.Logf("Int64    : %#v", id.Int64())
	t.Logf("String   : %#v", id.String())
	t.Logf("Base2    : %#v", id.Base2())
	t.Logf("Base32   : %#v", id.Base32())
	t.Logf("Base36   : %#v", id.Base36())
	t.Logf("Base58   : %#v", id.Base58())
	t.Logf("Base64   : %#v", id.Base64())
	t.Logf("Bytes    : %#v", id.Bytes())
	t.Logf("IntBytes : %#v", id.IntBytes())
}

func TestInt64(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	oId := node.Next()
	i := oId.Int64()

	pId := ParseInt64(i)
	if pId != oId {
		t.Fatalf("pId %v != oId %v", pId, oId)
	}

	mi := int64(1116766490855473152)
	pId = ParseInt64(mi)
	if pId.Int64() != mi {
		t.Fatalf("pId %v != mi %v", pId.Int64(), mi)
	}

}

func TestString(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	oId := node.Next()
	si := oId.String()

	pId, err := ParseString(si)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	if pId != oId {
		t.Fatalf("pId %v != oId %v", pId, oId)
	}

	ms := `1116766490855473152`
	_, err = ParseString(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseString(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

func TestBase2(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	oId := node.Next()
	i := oId.Base2()

	pId, err := ParseBase2(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pId != oId {
		t.Fatalf("pId %v != oId %v", pId, oId)
	}

	ms := `111101111111101110110101100101001000000000000000000000000000`
	_, err = ParseBase2(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `1112316766490855473152`
	_, err = ParseBase2(ms)
	if err == nil {
		t.Fatalf("no error parsing %s", ms)
	}
}

func TestBase32(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	for i := 0; i < 100; i++ {

		sf := node.Next()
		b32i := sf.Base32()
		psf, err := ParseBase32([]byte(b32i))
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func TestBase36(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	oId := node.Next()
	i := oId.Base36()

	pId, err := ParseBase36(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pId != oId {
		t.Fatalf("pId %v != oId %v", pId, oId)
	}

	ms := `8hgmw4blvlkw`
	_, err = ParseBase36(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `68h5gmw443blv2lk1w`
	_, err = ParseBase36(ms)
	if err == nil {
		t.Fatalf("no error parsing, %s", err)
	}
}

func TestBase58(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	for i := 0; i < 10; i++ {

		sf := node.Next()
		b58 := sf.Base58()
		psf, err := ParseBase58([]byte(b58))
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func TestBase64(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	oId := node.Next()
	i := oId.Base64()

	pId, err := ParseBase64(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pId != oId {
		t.Fatalf("pId %v != oId %v", pId, oId)
	}

	ms := `MTExNjgxOTQ5NDY2MDk5NzEyMA==`
	_, err = ParseBase64(ms)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}

	ms = `MTExNjgxOTQ5NDY2MDk5NzEyMA`
	_, err = ParseBase64(ms)
	if err == nil {
		t.Fatalf("no error parsing, %s", err)
	}
}

func TestBytes(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	oId := node.Next()
	i := oId.Bytes()

	pId, err := ParseBytes(i)
	if err != nil {
		t.Fatalf("error parsing, %s", err)
	}
	if pId != oId {
		t.Fatalf("pId %v != oId %v", pId, oId)
	}

	ms := []byte{0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err != nil {
		t.Fatalf("error parsing, %#v", err)
	}

	ms = []byte{0xFF, 0xFF, 0xFF, 0x31, 0x31, 0x31, 0x36, 0x38, 0x32, 0x31, 0x36, 0x37, 0x39, 0x35, 0x37, 0x30, 0x34, 0x31, 0x39, 0x37, 0x31, 0x32}
	_, err = ParseBytes(ms)
	if err == nil {
		t.Fatalf("no error parsing, %#v", err)
	}
}

func TestIntBytes(t *testing.T) {
	node, err := NewSnowflake(0)
	if err != nil {
		t.Fatalf("error creating NewSnowflake, %s", err)
	}

	oId := node.Next()
	i := oId.IntBytes()

	pId := ParseIntBytes(i)
	if pId != oId {
		t.Fatalf("pId %v != oId %v", pId, oId)
	}

	ms := [8]uint8{0xf, 0x7f, 0xc0, 0xfc, 0x2f, 0x80, 0x0, 0x0}
	mi := int64(1116823421972381696)
	pId = ParseIntBytes(ms)
	if pId.Int64() != mi {
		t.Fatalf("pId %v != mi %v", pId.Int64(), mi)
	}

}

// Marshall Test Methods
func TestMarshalJSON(t *testing.T) {
	id := Id(13587)
	expected := "\"13587\""

	bytes, err := id.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error during MarshalJSON")
	}

	if string(bytes) != expected {
		t.Fatalf("Got %s, expected %s", string(bytes), expected)
	}
}

func TestMarshalsIntBytes(t *testing.T) {
	id := Id(13587).IntBytes()
	expected := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x35, 0x13}
	if !bytes.Equal(id[:], expected) {
		t.Fatalf("Expected Id to be encoded as %v, got %v", expected, id)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tt := []struct {
		json        string
		expectedId  Id
		expectedErr error
	}{
		{`"13587"`, 13587, nil},
		{`1`, 0, JSONSyntaxError{[]byte(`1`)}},
		{`"invalid`, 0, JSONSyntaxError{[]byte(`"invalid`)}},
	}

	for _, tc := range tt {
		var id Id
		err := id.UnmarshalJSON([]byte(tc.json))
		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Fatalf("Expected to get error '%s' decoding JSON, but got '%s'", tc.expectedErr, err)
		}

		if id != tc.expectedId {
			t.Fatalf("Expected to get Id '%s' decoding JSON, but got '%s'", tc.expectedId, id)
		}
	}
}

// Benchmark Methods
func BenchmarkParseBase32(b *testing.B) {
	node, _ := NewSnowflake(1)
	sf := node.Next()
	b32i := sf.Base32()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ParseBase32([]byte(b32i))
	}
}

func BenchmarkBase32(b *testing.B) {
	node, _ := NewSnowflake(1)
	sf := node.Next()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base32()
	}
}

func BenchmarkParseBase58(b *testing.B) {
	node, _ := NewSnowflake(1)
	sf := node.Next()
	b58 := sf.Base58()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ParseBase58([]byte(b58))
	}
}

func BenchmarkBase58(b *testing.B) {
	node, _ := NewSnowflake(1)
	sf := node.Next()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sf.Base58()
	}
}

func BenchmarkNext(b *testing.B) {
	node, _ := NewSnowflake(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Next()
	}
}

func BenchmarkNextMaxSequence(b *testing.B) {
	NodeBits = 1
	StepBits = 21
	node, _ := NewSnowflake(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Next()
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	node, _ := NewSnowflake(1)
	id := node.Next()
	bytes, _ := id.MarshalJSON()

	var id2 Id

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = id2.UnmarshalJSON(bytes)
	}
}

func BenchmarkMarshal(b *testing.B) {
	node, _ := NewSnowflake(1)
	id := node.Next()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = id.MarshalJSON()
	}
}
