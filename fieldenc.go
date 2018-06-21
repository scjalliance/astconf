package astconf

import "reflect"

// PROPOSAL: Generate a unique ID for each field encoder then check
//           and set it in fieldEncoder.encode
// var nextID uint64
// id := atomic.AddUint64(&nextID, 1)

type fieldEncoder struct {
	index   []int
	elemEnc encoderFunc
	//omitEmpty bool
}

func newFieldEncoder(index []int, elemEnc encoderFunc) encoderFunc {
	enc := fieldEncoder{index: index, elemEnc: elemEnc}
	return enc.encode
}

func (fe fieldEncoder) encode(v reflect.Value, e *Encoder) error {
	fv := fieldByIndex(v, fe.index)
	//if !fv.IsValid() || fe.omitEmpty && isEmptyValue(fv) {
	if !fv.IsValid() {
		return nil
	}
	return fe.elemEnc(fv, e)
}
