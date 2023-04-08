package assets

import (
	"fmt"
	"strings"
)

func GroupErrors(fns ...func() error) error {
	var errs []string
	for _, fn := range fns {
		err := fn()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}
