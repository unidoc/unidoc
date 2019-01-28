/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import "errors"

var (
	errTypeCheck      = errors.New("type check error")
	errRangeCheck     = errors.New("range check error")
	errInvalidContext = errors.New("invaiid context")
	errRequiredField  = errors.New("required field missing")
	errNilReceiver    = errors.New("receiver pointer not initialized")
)
