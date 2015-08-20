package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (m *MessageHeaderSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], m)
	if err != nil {
		return err
	}
	m.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		m.Size = &element.DigitDataElement{}
		err = m.Size.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		m.HBCIVersion = &element.NumberDataElement{}
		err = m.HBCIVersion.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		m.DialogID = &element.IdentificationDataElement{}
		err = m.DialogID.UnmarshalHBCI(elements[3])
		if err != nil {
			return err
		}
	}
	if len(elements) > 4 && len(elements[4]) > 0 {
		m.Number = &element.NumberDataElement{}
		err = m.Number.UnmarshalHBCI(elements[4])
		if err != nil {
			return err
		}
	}
	if len(elements) > 5 && len(elements[5]) > 0 {
		m.Ref = &element.ReferencingMessageDataElement{}
		if len(elements)+1 > 5 {
			err = m.Ref.UnmarshalHBCI(bytes.Join(elements[5:], []byte("+")))
		} else {
			err = m.Ref.UnmarshalHBCI(elements[5])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
