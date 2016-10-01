package apns

import "encoding/json"

type BadgeNumber struct {
	Number uint
	IsSet  bool
}

func (b *BadgeNumber) Unset() {
	b.Number = 0
	b.IsSet = false
}

func (b *BadgeNumber) Set(number uint) {
	b.Number = number
	b.IsSet = true
}

func (b *BadgeNumber) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &b.Number)
	if err != nil {
		return err
	}

	b.IsSet = true
	return nil
}
