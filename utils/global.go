package utils

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	GLOBAL_WORKFLOW_MAP  map[string]string
	GLOBAL_LOGGER        *zap.Logger
	GLOBAL_CRON_WORKER   *cron.Cron
	GLOBAL_USER_IDENTITY bool
)
