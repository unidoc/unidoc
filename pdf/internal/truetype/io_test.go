/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"encoding/binary"
	"io"
	"os"
	"testing"
)

/*
Comparing use of binary.Read vs binary.BigEndian.... direct use.
Benchmark results indicate that the performance in all cases is pretty comparable,
within 5% difference. Thus choosing the simplest option (first one).

BenchmarkBinaryRead-8               3000            509553 ns/op
--- BENCH: BenchmarkBinaryRead-8
    io_test.go:26: Result: 1633005 (N: 1)
    io_test.go:26: Result: 163300500 (N: 100)
    io_test.go:26: Result: 4899015000 (N: 3000)
BenchmarkBinaryRead2-8              3000            492651 ns/op
--- BENCH: BenchmarkBinaryRead2-8
    io_test.go:49: Result: 1633005 (N: 1)
    io_test.go:49: Result: 163300500 (N: 100)
    io_test.go:49: Result: 4899015000 (N: 3000)
BenchmarkBinaryRead3-8              3000            508578 ns/op
--- BENCH: BenchmarkBinaryRead3-8
    io_test.go:76: Result: 1633005 (N: 1)
    io_test.go:76: Result: 163300500 (N: 100)
    io_test.go:76: Result: 4899015000 (N: 3000)
PASS
ok      github.com/unidoc/unidoc/pdf/internal/truetype  10.357s
*/

func BenchmarkBinaryRead(b *testing.B) {
	f, err := os.Open("../../creator/testdata/FreeSans.ttf")
	if err != nil {
		b.Fatalf("Error: %v", err)
	}
	defer f.Close()

	var sum int64
	for i := 0; i < b.N; i++ {
		f.Seek(0, io.SeekStart)
		for j := 0; j < 100; j++ {
			var val offset16
			binary.Read(f, binary.BigEndian, &val)
			sum += int64(val)
		}
	}
	b.Logf("Result: %d (N: %d)", sum, b.N)
}

func BenchmarkBinaryRead2(b *testing.B) {
	f, err := os.Open("../../creator/testdata/FreeSans.ttf")
	if err != nil {
		b.Fatalf("Error: %v", err)
	}
	defer f.Close()

	var sum int64
	for i := 0; i < b.N; i++ {
		f.Seek(0, io.SeekStart)
		for j := 0; j < 100; j++ {
			data := make([]byte, 2)
			_, err = io.ReadFull(f, data)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			val := offset16(binary.BigEndian.Uint16(data))
			sum += int64(val)
		}
	}
	b.Logf("Result: %d (N: %d)", sum, b.N)
}

func BenchmarkBinaryRead3(b *testing.B) {
	f, err := os.Open("../../creator/testdata/FreeSans.ttf")
	if err != nil {
		b.Fatalf("Error: %v", err)
	}
	defer f.Close()

	readOffset16 := func(r io.Reader) (offset16, error) {
		var val offset16
		err := binary.Read(f, binary.BigEndian, &val)
		return val, err
	}

	var sum int64
	for i := 0; i < b.N; i++ {
		f.Seek(0, io.SeekStart)
		for j := 0; j < 100; j++ {
			val, err := readOffset16(f)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			sum += int64(val)
		}
	}
	b.Logf("Result: %d (N: %d)", sum, b.N)
}
