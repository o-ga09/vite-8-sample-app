package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/o-ga09/vite-8-sample-app/internal/handler"
	inframd "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/hook"
	api "github.com/o-ga09/vite-8-sample-app/internal/oas"
	"github.com/o-ga09/vite-8-sample-app/internal/repository"
	"github.com/o-ga09/vite-8-sample-app/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// 環境変数から設定を読み込む
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return fmt.Errorf("DSN environment variable is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	corsOrigins := os.Getenv("CORS_ORIGINS")
	var allowOrigins []string
	if corsOrigins == "" {
		allowOrigins = []string{"*"}
	} else {
		for _, o := range strings.Split(corsOrigins, ",") {
			allowOrigins = append(allowOrigins, strings.TrimSpace(o))
		}
	}

	// DB接続
	db, err := inframd.NewDB(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// テナントスコープのテーブルにワークスペースフックを登録する
	hook.RegisterHooks(&dbgen.Members.SelectQueryHooks, &dbgen.Members.UpdateQueryHooks, &dbgen.Members.DeleteQueryHooks)
	hook.RegisterHooks(&dbgen.Accounts.SelectQueryHooks, &dbgen.Accounts.UpdateQueryHooks, &dbgen.Accounts.DeleteQueryHooks)
	hook.RegisterHooks(&dbgen.Categories.SelectQueryHooks, &dbgen.Categories.UpdateQueryHooks, &dbgen.Categories.DeleteQueryHooks)
	hook.RegisterHooks(&dbgen.Transactions.SelectQueryHooks, &dbgen.Transactions.UpdateQueryHooks, &dbgen.Transactions.DeleteQueryHooks)

	// リポジトリの生成
	workspaceRepo := repository.NewWorkspaceRepository()
	memberRepo := repository.NewMemberRepository()
	accountRepo := repository.NewAccountRepository()
	categoryRepo := repository.NewCategoryRepository()
	transactionRepo := repository.NewTransactionRepository()
	reportRepo := repository.NewReportRepository()

	// サービスの生成
	workspaceSvc := service.NewWorkspaceService(db, workspaceRepo)
	memberSvc := service.NewMemberService(db, memberRepo)
	accountSvc := service.NewAccountService(db, accountRepo)
	categorySvc := service.NewCategoryService(db, categoryRepo)
	transactionSvc := service.NewTransactionService(db, transactionRepo)
	reportSvc := service.NewReportService(db, reportRepo)

	// ハンドラの生成
	h := handler.New(handler.Options{
		WorkspaceService:   workspaceSvc,
		MemberService:      memberSvc,
		AccountService:     accountSvc,
		CategoryService:    categorySvc,
		TransactionService: transactionSvc,
		ReportService:      reportSvc,
	})

	// ogen サーバーの生成
	oasServer, err := api.NewServer(h)
	if err != nil {
		return fmt.Errorf("failed to create oas server: %w", err)
	}

	// Echo の設定
	e := echo.New()
	e.HideBanner = true
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: allowOrigins,
	}))
	e.GET("/metrics", echoprometheus.NewHandler())
	e.Any("/*", echo.WrapHandler(oasServer))

	// グレースフルシャットダウン
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		addr := ":" + port
		slog.Info("starting server", "addr", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return e.Shutdown(shutdownCtx)
}
