package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (s *SynchronisationResponseSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment SynchronisationResponse
	switch header.Version.Val() {
	case 3:
		segment = &SynchronisationResponseSegmentV3{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 4:
		segment = &SynchronisationResponseSegmentV4{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	s.SynchronisationResponse = segment
	return nil
}

func (s *SynchronisationResponseSegmentV3) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], s)
	if err != nil {
		return err
	}
	s.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		s.ClientSystemIDResponse = &element.IdentificationDataElement{}
		err = s.ClientSystemIDResponse.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		s.MessageNumberResponse = &element.NumberDataElement{}
		err = s.MessageNumberResponse.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		s.SignatureIDResponse = &element.NumberDataElement{}
		if len(elements)+1 > 3 {
			err = s.SignatureIDResponse.UnmarshalHBCI(bytes.Join(elements[3:], []byte("+")))
		} else {
			err = s.SignatureIDResponse.UnmarshalHBCI(elements[3])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SynchronisationResponseSegmentV4) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], s)
	if err != nil {
		return err
	}
	s.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		s.ClientSystemIDResponse = &element.IdentificationDataElement{}
		err = s.ClientSystemIDResponse.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		s.MessageNumberResponse = &element.NumberDataElement{}
		err = s.MessageNumberResponse.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	if len(elements) > 3 && len(elements[3]) > 0 {
		s.SignatureIDResponse = &element.NumberDataElement{}
		if len(elements)+1 > 3 {
			err = s.SignatureIDResponse.UnmarshalHBCI(bytes.Join(elements[3:], []byte("+")))
		} else {
			err = s.SignatureIDResponse.UnmarshalHBCI(elements[3])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
