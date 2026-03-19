package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	enums "github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/dbenums"
	"github.com/o-ga09/vite-8-sample-app/internal/infra/dbgen/factory"
	inframd "github.com/o-ga09/vite-8-sample-app/internal/infra/db"
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

	dsn := os.Getenv("DSN")
	if dsn == "" {
		return fmt.Errorf("DSN environment variable is required")
	}

	db, err := inframd.NewDB(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	f := factory.New()

	// ワークスペースを作成する
	ws, err := f.NewWorkspace(
		factory.WorkspaceMods.Name("マイワークスペース"),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	slog.Info("created workspace", "id", ws.ID, "name", ws.Name)

	// メンバーを作成する
	_, err = f.NewMember(
		factory.MemberMods.WorkspaceID(ws.ID),
		factory.MemberMods.Email("owner@example.com"),
		factory.MemberMods.DisplayName("オーナー"),
		factory.MemberMods.Role(enums.MemberRoleOwner),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create owner: %w", err)
	}

	_, err = f.NewMember(
		factory.MemberMods.WorkspaceID(ws.ID),
		factory.MemberMods.Email("member@example.com"),
		factory.MemberMods.DisplayName("メンバー"),
		factory.MemberMods.Role(enums.MemberRoleMember),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}
	slog.Info("created 2 members")

	// 口座を作成する
	cash, err := f.NewAccount(
		factory.AccountMods.WorkspaceID(ws.ID),
		factory.AccountMods.Name("財布"),
		factory.AccountMods.AccountType(enums.AccountTypeCash),
		factory.AccountMods.InitialBalance(decimal.NewFromFloat(10000)),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create cash account: %w", err)
	}

	bank, err := f.NewAccount(
		factory.AccountMods.WorkspaceID(ws.ID),
		factory.AccountMods.Name("銀行口座"),
		factory.AccountMods.AccountType(enums.AccountTypeBank),
		factory.AccountMods.InitialBalance(decimal.NewFromFloat(500000)),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create bank account: %w", err)
	}

	_, err = f.NewAccount(
		factory.AccountMods.WorkspaceID(ws.ID),
		factory.AccountMods.Name("クレジットカード"),
		factory.AccountMods.AccountType(enums.AccountTypeCreditCard),
		factory.AccountMods.InitialBalance(decimal.NewFromFloat(0)),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create credit account: %w", err)
	}
	slog.Info("created 3 accounts")

	// カテゴリを作成する（収入）
	salary, err := f.NewCategory(
		factory.CategoryMods.WorkspaceID(ws.ID),
		factory.CategoryMods.Name("給与"),
		factory.CategoryMods.CategoryType(enums.CategoryTypeIncome),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create salary category: %w", err)
	}

	_, err = f.NewCategory(
		factory.CategoryMods.WorkspaceID(ws.ID),
		factory.CategoryMods.Name("副業"),
		factory.CategoryMods.CategoryType(enums.CategoryTypeIncome),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create side income category: %w", err)
	}

	// カテゴリを作成する（支出）
	food, err := f.NewCategory(
		factory.CategoryMods.WorkspaceID(ws.ID),
		factory.CategoryMods.Name("食費"),
		factory.CategoryMods.CategoryType(enums.CategoryTypeExpense),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create food category: %w", err)
	}

	transport, err := f.NewCategory(
		factory.CategoryMods.WorkspaceID(ws.ID),
		factory.CategoryMods.Name("交通費"),
		factory.CategoryMods.CategoryType(enums.CategoryTypeExpense),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create transport category: %w", err)
	}

	utility, err := f.NewCategory(
		factory.CategoryMods.WorkspaceID(ws.ID),
		factory.CategoryMods.Name("光熱費"),
		factory.CategoryMods.CategoryType(enums.CategoryTypeExpense),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create utility category: %w", err)
	}

	_, err = f.NewCategory(
		factory.CategoryMods.WorkspaceID(ws.ID),
		factory.CategoryMods.Name("娯楽費"),
		factory.CategoryMods.CategoryType(enums.CategoryTypeExpense),
	).Create(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create entertainment category: %w", err)
	}

	_, err = f.NewCategory(
		factory.CategoryMods.WorkspaceID(ws.ID),
		factory.CategoryMods.Name("医療費"),
		factory.CategoryMods.CategoryType(enums.CategoryTypeExpense),
	).Create(ctx, db)
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
		_, err = f.NewTransaction(
			factory.TransactionMods.WorkspaceID(ws.ID),
			factory.TransactionMods.AccountID(null.From(s.acctID)),
			factory.TransactionMods.CategoryID(null.From(s.catID)),
			factory.TransactionMods.TransactionType(s.txType),
			factory.TransactionMods.Amount(decimal.NewFromFloat(s.amt)),
			factory.TransactionMods.OccurredAt(d),
			factory.TransactionMods.Description(null.From(s.desc)),
		).Create(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to create transaction %d: %w", i, err)
		}
	}

	slog.Info("seed completed", "workspace_id", ws.ID, "transactions", len(seeds))
	return nil
}
