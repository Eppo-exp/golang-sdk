package eppoclient

import "errors"

var (
	ErrSubjectAllocation           = errors.New("subject is not part of any allocation")
	ErrFlagNotEnabled              = errors.New("the experiment or flag is not enabled")
	ErrFlagConfigurationNotFound   = errors.New("flag configuration not found")
	ErrBanditConfigurationNotFound = errors.New("bandit configuration not found")
)
