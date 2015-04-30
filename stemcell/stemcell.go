package stemcell

import (
	"errors"
)

func Create(args []interface{}) (result interface{}, err error) {
	imagePath, ok := args[0].(string)
	if !ok {
		return nil, errors.New("Unexpected argument where image_path should be")
	}
	return "fake-stemcell-id" + imagePath, nil
}

func Delete(args []interface{}) (result interface{}, err error) {
	return nil, nil
}
