package segment

import (
	"testing"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func TestAccountInformationSegmentUnmarshalHBCI(t *testing.T) {
	test := "HIUPD:1:4:4+123456::280:10000000+12345+EUR+Muster+Max+Sichteinlagen++DKPAE:1'"

	account := &AccountInformationSegment{}

	err := account.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	v4 := &AccountInformationV4{
		AccountConnection:           element.NewAccountConnection(domain.AccountConnection{AccountID: "123456", CountryCode: 280, BankID: "10000000"}),
		UserID:                      element.NewIdentification("12345"),
		AccountCurrency:             element.NewCurrency("EUR"),
		Name1:                       element.NewAlphaNumeric("Muster", 27),
		Name2:                       element.NewAlphaNumeric("Max", 27),
		AccountProductID:            element.NewAlphaNumeric("Sichteinlagen", 30),
		AllowedBusinessTransactions: element.NewAllowedBusinessTransactions(domain.BusinessTransaction{ID: "DKPAE", NeededSignatures: 1}),
	}
	v4.Segment = NewReferencingBasicSegment(1, 4, v4)
	expected := &AccountInformationSegment{v4}

	expectedString := expected.String()
	actualString := account.String()

	if expectedString != actualString {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expectedString, actualString)
		t.Fail()
	}
}
