// +build linux
package flite

import (
	"unsafe"
	"fmt"
)

// #cgo CFLAGS: -I /usr/include/flite/
// #cgo LDFLAGS: -lflite -lflite_cmu_us_kal
// #include "flite.h"
// cst_voice* register_cmu_us_kal(const char *voxdir);
import "C"

var voice *C.cst_voice

func init() {
	C.flite_init()
	voice = C.register_cmu_us_kal(nil)
}

func TextToSpeech(path, text string) error {
	if voice == nil {
		return fmt.Errorf("Could not find default voice")
	}

	ctext := C.CString(text)
	cout := C.CString(path)

	C.flite_text_to_speech(ctext, voice, cout)
	C.free(unsafe.Pointer(ctext))
	C.free(unsafe.Pointer(cout))
	return nil
}
