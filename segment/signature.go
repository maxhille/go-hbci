package segment

import (
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func NewPinTanSignatureHeaderSegment(controlReference string, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	s := &SignatureHeaderSegment{
		SecurityFunction:         element.NewAlphaNumeric("999", 3),
		SecurityControlRef:       element.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: element.NewAlphaNumeric("1", 3),
		SecuritySupplierRole:     element.NewAlphaNumeric("1", 3),
		SecurityID:               element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        element.NewNumber(0, 16),
		SecurityDate:             element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		HashAlgorithm:            element.NewDefaultHashAlgorithm(),
		SignatureAlgorithm:       element.NewRDHSignatureAlgorithm(),
		KeyName:                  element.NewKeyName(keyName),
	}
	s.Segment = NewBasicSegment(2, s)
	return s
}

func NewRDHSignatureHeaderSegment(controlReference string, signatureId int, clientSystemId string, keyName domain.KeyName) *SignatureHeaderSegment {
	s := &SignatureHeaderSegment{
		SecurityFunction:         element.NewAlphaNumeric("1", 3),
		SecurityControlRef:       element.NewAlphaNumeric(controlReference, 14),
		SecurityApplicationRange: element.NewAlphaNumeric("1", 3),
		SecuritySupplierRole:     element.NewAlphaNumeric("1", 3),
		SecurityID:               element.NewRDHSecurityIdentification(element.SecurityHolderMessageSender, clientSystemId),
		SecurityRefNumber:        element.NewNumber(signatureId, 16),
		SecurityDate:             element.NewSecurityDate(element.SecurityTimestamp, time.Now()),
		HashAlgorithm:            element.NewDefaultHashAlgorithm(),
		SignatureAlgorithm:       element.NewRDHSignatureAlgorithm(),
		KeyName:                  element.NewKeyName(keyName),
	}
	s.Segment = NewBasicSegment(2, s)
	return s
}

type SignatureHeaderSegment struct {
	Segment
	// "1" for NRO, Non-Repudiation of Origin (RDH)
	// "2" for AUT, Message Origin Authentication (DDV)
	// "999" for PIN/TAN
	SecurityFunction   *element.AlphaNumericDataElement
	SecurityControlRef *element.AlphaNumericDataElement
	// "1" for SHM (SignatureHeader and HBCI-Data)
	// "2" for SHT (SignatureHeader to SignatureEnd)
	SecurityApplicationRange *element.AlphaNumericDataElement
	// "1" for ISS, Herausgeber der signierten Nachricht (z.B. Erfasser oder Erstsignatur)
	// "3" for CON, der Unterzeichnete unterstützt den Inhalt der Nachricht (z.B. bei Zweitsignatur)
	// "4" for WIT, der Unterzeichnete ist Zeuge (z.B. Übermittler), aber für den Inhalt der Nachricht nicht verantwortlich)
	SecuritySupplierRole *element.AlphaNumericDataElement
	SecurityID           *element.SecurityIdentificationDataElement
	SecurityRefNumber    *element.NumberDataElement
	SecurityDate         *element.SecurityDateDataElement
	HashAlgorithm        *element.HashAlgorithmDataElement
	SignatureAlgorithm   *element.SignatureAlgorithmDataElement
	KeyName              *element.KeyNameDataElement
	Certificate          *element.CertificateDataElement
}

func (s *SignatureHeaderSegment) init() {
	*s.SecurityFunction = *new(element.AlphaNumericDataElement)
	*s.SecurityControlRef = *new(element.AlphaNumericDataElement)
	*s.SecurityApplicationRange = *new(element.AlphaNumericDataElement)
	*s.SecuritySupplierRole = *new(element.AlphaNumericDataElement)
	*s.SecurityID = *new(element.SecurityIdentificationDataElement)
	*s.SecurityRefNumber = *new(element.NumberDataElement)
	*s.SecurityDate = *new(element.SecurityDateDataElement)
	*s.HashAlgorithm = *new(element.HashAlgorithmDataElement)
	*s.SignatureAlgorithm = *new(element.SignatureAlgorithmDataElement)
	*s.KeyName = *new(element.KeyNameDataElement)
	*s.Certificate = *new(element.CertificateDataElement)
}
func (s *SignatureHeaderSegment) version() int         { return 3 }
func (s *SignatureHeaderSegment) id() string           { return "HNSHK" }
func (s *SignatureHeaderSegment) referencedId() string { return "" }
func (s *SignatureHeaderSegment) sender() string       { return senderBoth }

func (s *SignatureHeaderSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.SecurityFunction,
		s.SecurityControlRef,
		s.SecurityApplicationRange,
		s.SecuritySupplierRole,
		s.SecurityID,
		s.SecurityRefNumber,
		s.SecurityDate,
		s.HashAlgorithm,
		s.SignatureAlgorithm,
		s.KeyName,
		s.Certificate,
	}
}

func NewSignatureEndSegment(number int, controlReference string) *SignatureEndSegment {
	s := &SignatureEndSegment{
		SecurityControlRef: element.NewAlphaNumeric(controlReference, 14),
	}
	s.Segment = NewBasicSegment(number, s)
	return s
}

type SignatureEndSegment struct {
	Segment
	SecurityControlRef *element.AlphaNumericDataElement
	Signature          *element.BinaryDataElement
	PinTan             *element.PinTanDataElement
}

func (s *SignatureEndSegment) init() {
	*s.SecurityControlRef = *new(element.AlphaNumericDataElement)
	*s.Signature = *new(element.BinaryDataElement)
	*s.PinTan = *new(element.PinTanDataElement)
}
func (s *SignatureEndSegment) version() int         { return 1 }
func (s *SignatureEndSegment) id() string           { return "HNSHA" }
func (s *SignatureEndSegment) referencedId() string { return "" }
func (s *SignatureEndSegment) sender() string       { return senderBoth }

func (s *SignatureEndSegment) elements() []element.DataElement {
	return []element.DataElement{
		s.SecurityControlRef,
		s.Signature,
		s.PinTan,
	}
}

func (s *SignatureEndSegment) SetSignature(signature []byte) {
	s.Signature = element.NewBinary(signature, 512)
}

func (s *SignatureEndSegment) SetPinTan(pin, tan string) {
	s.PinTan = element.NewPinTan(pin, tan)
}
