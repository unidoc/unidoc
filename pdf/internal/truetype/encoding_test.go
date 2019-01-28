/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMacEncoding(t *testing.T) {
	// Spot checks based on: https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6post.html
	assert.Equal(t, 258, len(macGlyphNames))
	assert.Equal(t, GlyphName(".notdef"), macGlyphNames[0])
	assert.Equal(t, GlyphName("space"), macGlyphNames[3])
	assert.Equal(t, GlyphName("comma"), macGlyphNames[15])
	assert.Equal(t, GlyphName("a"), macGlyphNames[68])
	assert.Equal(t, GlyphName("z"), macGlyphNames[93])
	assert.Equal(t, GlyphName("dcroat"), macGlyphNames[257])
}
