/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package contentstream

import "errors"

var (
	errInvalidOperand = errors.New("invalid operand")
	errTypeCheck      = errors.New("type check error")
	errRangeCheck     = errors.New("range check error")
	errNotFound       = errors.New("not found")
)
