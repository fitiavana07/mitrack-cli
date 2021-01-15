package account

import "github.com/fitiavana07/mitrack/pkg/account"

type accountTypeValue account.Type

func newAccountTypeValue(p *account.Type) *accountTypeValue {
	*p = account.TypeAsset
	return (*accountTypeValue)(p)
}

func (t *accountTypeValue) Set(val string) error {
	got, err := account.TypeFromString(val)
	if err != nil {
		return err
	}
	*t = accountTypeValue(got)
	return nil
}

const accountTypeValueType = "account type"

func (t *accountTypeValue) Type() string {
	return accountTypeValueType
}

func (t *accountTypeValue) String() string {
	return account.Type(*t).String()
}
