package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	inframd "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
	models "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/shopspring/decimal"
)

func main() {
	if err := run(); err != nil {
		slog.Error("seed failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	dsn := os.Getenv("PSQL_DSN")
	if dsn == "" {
		return fmt.Errorf("PSQL_DSN environment variable is required")
	}

	db, err := inframd.NewDB(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// ワークスペースを作成する
	ws, err := models.Workspaces.Insert(&models.WorkspaceSetter{
		ID:   omit.From(uuid.New()),
		Name: omit.From("マイワークスペース"),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	slog.Info("created workspace", "id", ws.ID, "name", ws.Name)

	// メンバーを作成する
	_, err = models.Members.Insert(&models.MemberSetter{
		ID:          omit.From(uuid.New()),
		WorkspaceID: omit.From(ws.ID),
		Email:       omit.From("owner@example.com"),
		DisplayName: omit.From("オーナー"),
		Role:        omit.From(enums.MemberRoleOwner),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create owner: %w", err)
	}

	_, err = models.Members.Insert(&models.MemberSetter{
		ID:          omit.From(uuid.New()),
		WorkspaceID: omit.From(ws.ID),
		Email:       omit.From("member@example.com"),
		DisplayName: omit.From("メンバー"),
		Role:        omit.From(enums.MemberRoleMember),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}
	slog.Info("created 2 members")

	// 口座を作成する
	cash, err := models.Accounts.Insert(&models.AccountSetter{
		ID:             omit.From(uuid.New()),
		WorkspaceID:    omit.From(ws.ID),
		Name:           omit.From("財布"),
		AccountType:    omit.From(enums.AccountTypeCash),
		InitialBalance: omit.From(decimal.NewFromFloat(10000)),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create cash account: %w", err)
	}

	bank, err := models.Accounts.Insert(&models.AccountSetter{
		ID:             omit.From(uuid.New()),
		WorkspaceID:    omit.From(ws.ID),
		Name:           omit.From("銀行口座"),
		AccountType:    omit.From(enums.AccountTypeBank),
		InitialBalance: omit.From(decimal.NewFromFloat(500000)),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create bank account: %w", err)
	}

	_, err = models.Accounts.Insert(&models.AccountSetter{
		ID:             omit.From(uuid.New()),
		WorkspaceID:    omit.From(ws.ID),
		Name:           omit.From("クレジットカード"),
		AccountType:    omit.From(enums.AccountTypeCreditCard),
		InitialBalance: omit.From(decimal.NewFromFloat(0)),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create credit account: %w", err)
	}
	slog.Info("created 3 accounts")

	// カテゴリを作成する（収入）
	salary, err := models.Categories.Insert(&models.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("給与"),
		CategoryType: omit.From(enums.CategoryTypeIncome),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create salary category: %w", err)
	}

	_, err = models.Categories.Insert(&models.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("副業"),
		CategoryType: omit.From(enums.CategoryTypeIncome),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create side income category: %w", err)
	}

	// カテゴリを作成する（支出）
	food, err := models.Categories.Insert(&models.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("食費"),
		CategoryType: omit.From(enums.CategoryTypeExpense),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create food category: %w", err)
	}

	transport, err := models.Categories.Insert(&models.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("交通費"),
		CategoryType: omit.From(enums.CategoryTypeExpense),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create transport category: %w", err)
	}

	utility, err := models.Categories.Insert(&models.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("光熱費"),
		CategoryType: omit.From(enums.CategoryTypeExpense),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create utility category: %w", err)
	}

	_, err = models.Categories.Insert(&models.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("娯楽費"),
		CategoryType: omit.From(enums.CategoryTypeExpense),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create entertainment category: %w", err)
	}

	_, err = models.Categories.Insert(&models.CategorySetter{
		ID:           omit.From(uuid.New()),
		WorkspaceID:  omit.From(ws.ID),
		Name:         omit.From("医療費"),
		CategoryType: omit.From(enums.CategoryTypeExpense),
	}).One(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create medical category: %w", err)
	}
	slog.Info("created 7 categories")

	// 取引データを作成する（当月1か月分）
	now := time.Now()
	baseMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

	type txSeed struct {
		day    int
		amt    float64
		txType enums.TransactionType
		acctID uuid.UUID
		catID  uuid.UUID
		desc   string
	}

	seeds := []txSeed{
		{1, 300000, enums.TransactionTypeIncome, bank.ID, salary.ID, "月給"},
		{2, 800, enums.TransactionTypeExpense, cash.ID, food.ID, "昼食"},
		{3, 1200, enums.TransactionTypeExpense, cash.ID, food.ID, "夕食"},
		{4, 5000, enums.TransactionTypeExpense, bank.ID, transport.ID, "定期券"},
		{5, 3000, enums.TransactionTypeExpense, cash.ID, food.ID, "スーパー"},
		{7, 8000, enums.TransactionTypeExpense, bank.ID, utility.ID, "電気代"},
		{8, 4000, enums.TransactionTypeExpense, bank.ID, utility.ID, "ガス代"},
		{10, 1500, enums.TransactionTypeExpense, cash.ID, food.ID, "ランチ"},
		{11, 2500, enums.TransactionTypeExpense, cash.ID, food.ID, "夕食"},
		{12, 3000, enums.TransactionTypeExpense, cash.ID, food.ID, "外食"},
		{13, 500, enums.TransactionTypeExpense, cash.ID, transport.ID, "バス"},
		{14, 15000, enums.TransactionTypeExpense, bank.ID, food.ID, "食材まとめ買い"},
		{15, 50000, enums.TransactionTypeIncome, bank.ID, salary.ID, "賞与"},
		{16, 2000, enums.TransactionTypeExpense, cash.ID, food.ID, "昼食"},
		{17, 6000, enums.TransactionTypeExpense, bank.ID, utility.ID, "水道代"},
		{18, 1000, enums.TransactionTypeExpense, cash.ID, transport.ID, "電車"},
		{20, 12000, enums.TransactionTypeExpense, cash.ID, food.ID, "食料品"},
		{22, 4500, enums.TransactionTypeExpense, cash.ID, food.ID, "外食"},
		{25, 2000, enums.TransactionTypeExpense, cash.ID, transport.ID, "タクシー"},
		{28, 9000, enums.TransactionTypeExpense, bank.ID, food.ID, "月末まとめ買い"},
	}

	for i, s := range seeds {
		d := baseMonth.AddDate(0, 0, s.day-1)
		_, err = models.Transactions.Insert(&models.TransactionSetter{
			ID:              omit.From(uuid.New()),
			WorkspaceID:     omit.From(ws.ID),
			TransactionType: omit.From(s.txType),
			AccountID:       omitnull.From(s.acctID),
			CategoryID:      omitnull.From(s.catID),
			Amount:          omit.From(decimal.NewFromFloat(s.amt)),
			OccurredAt:      omit.From(d),
			Description:     omitnull.From(s.desc),
		}).One(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to create transaction %d: %w", i, err)
		}
	}

	slog.Info("seed completed", "workspace_id", ws.ID, "transactions", len(seeds))
	return nil
}
