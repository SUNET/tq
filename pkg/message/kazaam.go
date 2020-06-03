package message

import (
	"github.com/qntfy/kazaam"
)

func NewKazaam(spec string) *kazaam.Kazaam {
	kz, err := kazaam.NewKazaam(spec)
	if err != nil {
		Log.Fatalf("Unable to create kazaam from %s: %s", spec, err.Error())
	}
	return kz
}

func KazaamHandler(kz *kazaam.Kazaam, o Message) (Message, error) {
	j, err := FromJson(o)
	if err != nil {
		return nil, err
	}

	j, err = kz.Transform(j)
	if err != nil {
		return nil, err
	}

	o, err = ToJson(j)
	if err != nil {
		return nil, err
	}

	return o, nil
}
