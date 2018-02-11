package gzip_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stairlin/kargo/context"
	"github.com/stairlin/kargo/pkg/testutil"
	"github.com/stairlin/kargo/pkg/unit"
	"github.com/stairlin/kargo/plugin/process/gzip"
)

func TestVerbatim(t *testing.T) {
	proc := gzip.Processor{}
	if err := proc.Init(); err != nil {
		t.Fatal(err)
	}

	expect := testutil.GenRandBytes(t, int(64*unit.MB))
	plain := ioutil.NopCloser(bytes.NewReader(expect))
	ctx := context.Background()

	// Encode/Decode
	encoded, err := proc.Encode(ctx, plain)
	if err != nil {
		t.Fatal("Error encoding", err)
	}
	decoded, err := proc.Decode(ctx, encoded)
	if err != nil {
		t.Fatal("Error decoding", err)
	}

	// Tests
	got, err := ioutil.ReadAll(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if len(expect) != len(got) {
		t.Errorf("expect length of %d, but got %d", len(expect), len(got))
	}
	if string(expect) != string(got) {
		t.Errorf("expect text %s, but got %s",
			testutil.Truncate(expect, 140), testutil.Truncate(got, 140),
		)
	}

	// Close readers
	if err := encoded.Close(); err != nil {
		t.Error("encoded close err", err)
	}
	if err := decoded.Close(); err != nil {
		t.Error("decoded close err", err)
	}
}

func TestDiff(t *testing.T) {
	proc := gzip.Processor{}
	if err := proc.Init(); err != nil {
		t.Fatal(err)
	}

	expect := testutil.GenRandBytes(t, int(128*unit.KB))
	plain := ioutil.NopCloser(bytes.NewReader(expect))
	ctx := context.Background()

	// Encode/Decode
	encoded, err := proc.Encode(ctx, plain)
	if err != nil {
		t.Fatal("Error encoding", err)
	}

	// Tests
	got, err := ioutil.ReadAll(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if string(expect) == string(got) {
		t.Error("expect encoded and decoded data to be different")
	}

	// Close readers
	if err := encoded.Close(); err != nil {
		t.Error("encoded close err", err)
	}
}
