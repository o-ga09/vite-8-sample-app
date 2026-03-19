package handler

import (
	oas "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
)

// Options holds the dependencies for the Handler.
type Options struct {
	WorkspaceService  service.WorkspaceService
	MemberService     service.MemberService
	AccountService    service.AccountService
	CategoryService   service.CategoryService
	TransactionService service.TransactionService
	ReportService     service.ReportService
}

// Handler implements oas.Handler.
type Handler struct {
	workspaceSvc  service.WorkspaceService
	memberSvc     service.MemberService
	accountSvc    service.AccountService
	categorySvc   service.CategoryService
	transactionSvc service.TransactionService
	reportSvc     service.ReportService
}

var _ oas.Handler = (*Handler)(nil)

// New creates a new Handler with the given options.
func New(opts Options) *Handler {
	return &Handler{
		workspaceSvc:  opts.WorkspaceService,
		memberSvc:     opts.MemberService,
		accountSvc:    opts.AccountService,
		categorySvc:   opts.CategoryService,
		transactionSvc: opts.TransactionService,
		reportSvc:     opts.ReportService,
	}
}

// notFound returns an ErrorResponse for resource not found.
func notFound(msg string) *oas.ErrorResponse {
	return &oas.ErrorResponse{Message: msg}
}

// badRequest returns an ErrorResponse for bad request.
func badRequest(msg string) *oas.ErrorResponse {
	return &oas.ErrorResponse{Message: msg}
}
