package resolvers

import (
	"errors"
	"fmt"

	"github.com/bdreece/hopper/pkg/utils"
	"gorm.io/gorm"
)

func handleError(err error, logger utils.Logger) error {
	logger.Errorf("An error has occurred: %v\n", err)
	return err
}
func handleQueryError(err error, logger utils.Logger, field string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(fmt.Errorf("%s not found", field), err)
	}
	return handleError(err, logger)
}
