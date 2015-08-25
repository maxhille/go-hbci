package segment

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (a *AccountTransactionResponseSegment) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment AccountTransactionResponse
	switch header.Version.Val() {
	case 5:
		segment = &AccountTransactionResponseSegmentV5{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 6:
		segment = &AccountTransactionResponseSegmentV6{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 7:
		segment = &AccountTransactionResponseSegmentV7{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	a.AccountTransactionResponse = segment
	return nil
}

func (a *AccountTransactionResponseSegmentV5) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], a)
	if err != nil {
		return err
	}
	a.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		a.BookedTransactions = &element.SwiftMT940DataElement{}
		err = a.BookedTransactions.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		a.UnbookedTransactions = &element.BinaryDataElement{}
		if len(elements)+1 > 2 {
			err = a.UnbookedTransactions.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = a.UnbookedTransactions.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountTransactionResponseSegmentV6) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], a)
	if err != nil {
		return err
	}
	a.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		a.BookedTransactions = &element.SwiftMT940DataElement{}
		err = a.BookedTransactions.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		a.UnbookedTransactions = &element.BinaryDataElement{}
		if len(elements)+1 > 2 {
			err = a.UnbookedTransactions.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = a.UnbookedTransactions.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountTransactionResponseSegmentV7) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], a)
	if err != nil {
		return err
	}
	a.Segment = seg
	if len(elements) > 1 && len(elements[1]) > 0 {
		a.BookedTransactions = &element.SwiftMT940DataElement{}
		err = a.BookedTransactions.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 && len(elements[2]) > 0 {
		a.UnbookedTransactions = &element.BinaryDataElement{}
		if len(elements)+1 > 2 {
			err = a.UnbookedTransactions.UnmarshalHBCI(bytes.Join(elements[2:], []byte("+")))
		} else {
			err = a.UnbookedTransactions.UnmarshalHBCI(elements[2])
		}
		if err != nil {
			return err
		}
	}
	return nil
}
