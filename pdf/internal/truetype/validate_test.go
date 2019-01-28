/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"fmt"
	"testing"
	"time"

	"github.com/unidoc/unidoc/common"
)

func init() {
	//common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	common.SetLogger(common.NewConsoleLogger(common.LogLevelInfo))
}

func TestFontValidation(t *testing.T) {
	testcases := []struct {
		fontPath string
	}{
		{
			"../../creator/testdata/FreeSans.ttf",
		},
		{
			"../../creator/testdata/wts11.ttf",
		},
		{
			"../../creator/testdata/roboto/Roboto-BoldItalic.ttf",
		},
	}

	for _, tcase := range testcases {
		t.Logf("%s", tcase.fontPath)
		fmt.Printf("==== %s\n", tcase.fontPath)
		common.Log.Debug("==== %s", tcase.fontPath)
		start := time.Now()
		err := ValidateFile(tcase.fontPath)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		diff := time.Now().Sub(start)
		t.Logf("- took: %s", diff.String())
	}

}
