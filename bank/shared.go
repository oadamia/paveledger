package bank

import "time"

const (
	TASK_QUEUE    = "authorization"
	QUERY_HANDLER = "authorize"
)

type workflowStatus int

const (
	WORKFLOW_STATUS_NOT_FOUND workflowStatus = 0
	WORKFLOW_STATUS_RUNNING   workflowStatus = 1
	WORKFLOW_STATUS_CLOSED    workflowStatus = 2
)

const (
	AUTHORIZATION_CHANNEL = "AUTHORIZATION_CHANNEL"
	PRESENTMENT_CHANNEL   = "PRESENTMENT_CHANNEL"
)

const authorizationTimeout = 30 * time.Second
