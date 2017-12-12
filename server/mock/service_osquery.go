// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import (
	"context"
	"encoding/json"

	"github.com/kolide/fleet/server/kolide"
)

var _ kolide.OsqueryService = (*TLSService)(nil)

type EnrollAgentFunc func(ctx context.Context, enrollSecret string, hostIdentifier string) (nodeKey string, err error)

type AuthenticateHostFuncI func(ctx context.Context, nodeKey string) (host *kolide.Host, err error)

type GetClientConfigFunc func(ctx context.Context) (config *kolide.OsqueryConfig, err error)

type GetDistributedQueriesFunc func(ctx context.Context) (queries map[string]string, accelerate uint, err error)

type SubmitDistributedQueryResultsFunc func(ctx context.Context, results kolide.OsqueryDistributedQueryResults, statuses map[string]string) (err error)

type SubmitStatusLogsFunc func(ctx context.Context, logs []json.RawMessage) (err error)

type SubmitResultLogsFunc func(ctx context.Context, logs []json.RawMessage) (err error)

type TLSService struct {
	EnrollAgentFunc        EnrollAgentFunc
	EnrollAgentFuncInvoked bool

	AuthenticateHostFunc        AuthenticateHostFuncI
	AuthenticateHostFuncInvoked bool

	GetClientConfigFunc        GetClientConfigFunc
	GetClientConfigFuncInvoked bool

	GetDistributedQueriesFunc        GetDistributedQueriesFunc
	GetDistributedQueriesFuncInvoked bool

	SubmitDistributedQueryResultsFunc        SubmitDistributedQueryResultsFunc
	SubmitDistributedQueryResultsFuncInvoked bool

	SubmitStatusLogsFunc        SubmitStatusLogsFunc
	SubmitStatusLogsFuncInvoked bool

	SubmitResultLogsFunc        SubmitResultLogsFunc
	SubmitResultLogsFuncInvoked bool
}

func (s *TLSService) EnrollAgent(ctx context.Context, enrollSecret string, hostIdentifier string) (nodeKey string, err error) {
	s.EnrollAgentFuncInvoked = true
	return s.EnrollAgentFunc(ctx, enrollSecret, hostIdentifier)
}

func (s *TLSService) AuthenticateHost(ctx context.Context, nodeKey string) (host *kolide.Host, err error) {
	s.AuthenticateHostFuncInvoked = true
	return s.AuthenticateHostFunc(ctx, nodeKey)
}

func (s *TLSService) GetClientConfig(ctx context.Context) (config *kolide.OsqueryConfig, err error) {
	s.GetClientConfigFuncInvoked = true
	return s.GetClientConfigFunc(ctx)
}

func (s *TLSService) GetDistributedQueries(ctx context.Context) (queries map[string]string, accelerate uint, err error) {
	s.GetDistributedQueriesFuncInvoked = true
	return s.GetDistributedQueriesFunc(ctx)
}

func (s *TLSService) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults, statuses map[string]string) (err error) {
	s.SubmitDistributedQueryResultsFuncInvoked = true
	return s.SubmitDistributedQueryResultsFunc(ctx, results, statuses)
}

func (s *TLSService) SubmitStatusLogs(ctx context.Context, logs []json.RawMessage) (err error) {
	s.SubmitStatusLogsFuncInvoked = true
	return s.SubmitStatusLogsFunc(ctx, logs)
}

func (s *TLSService) SubmitResultLogs(ctx context.Context, logs []json.RawMessage) (err error) {
	s.SubmitResultLogsFuncInvoked = true
	return s.SubmitResultLogsFunc(ctx, logs)
}
