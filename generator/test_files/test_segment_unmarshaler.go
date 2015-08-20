package test_files

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (t *TestSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], t)
	if err != nil {
		return err
	}
	t.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		t.Abc = &element.AlphaNumericDataElement{}
		err = t.Abc.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		t.Xyz = &element.NumberDataElement{}
		if len(elements)+1 > 2 {
			err = t.Xyz.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = t.Xyz.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
