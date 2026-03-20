import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { apiClient } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { format, startOfMonth, endOfMonth } from "date-fns";

type CategoryExpenseItem = components["schemas"]["CategoryExpenseItem"];
type AccountBalanceItem = components["schemas"]["AccountBalanceItem"];

export const Route = createFileRoute("/workspaces/$wsId/reports/")({
  component: ReportsPage,
});

function ReportsPage() {
  const { wsId } = Route.useParams();
  const [categoryExpenses, setCategoryExpenses] = useState<
    CategoryExpenseItem[]
  >([]);
  const [accountBalances, setAccountBalances] = useState<AccountBalanceItem[]>(
    [],
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const now = new Date();
    const from = format(startOfMonth(now), "yyyy-MM-dd'T'HH:mm:ss'Z'");
    const to = format(endOfMonth(now), "yyyy-MM-dd'T'HH:mm:ss'Z'");
    const asOf = format(now, "yyyy-MM-dd'T'HH:mm:ss'Z'");

    Promise.all([
      apiClient.GET("/workspaces/{wsId}/reports/category-expenses", {
        params: { path: { wsId }, query: { from, to } },
      }),
      apiClient.GET("/workspaces/{wsId}/reports/account-balances", {
        params: { path: { wsId }, query: { asOf } },
      }),
    ]).then(([catRes, balRes]) => {
      if (catRes.error) {
        setError("カテゴリ別支出の取得に失敗しました");
      } else {
        setCategoryExpenses(catRes.data ?? []);
      }
      if (balRes.error) {
        setError("口座残高の取得に失敗しました");
      } else {
        setAccountBalances(balRes.data ?? []);
      }
      setLoading(false);
    });
  }, [wsId]);

  const formatAmount = (val: string) =>
    new Intl.NumberFormat("ja-JP", {
      style: "currency",
      currency: "JPY",
    }).format(Number(val));

  const maxExpense = Math.max(
    ...categoryExpenses.map((c) => Number(c.totalExpense)),
    1,
  );

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-2">レポート</h1>
      <p className="text-sm text-muted-foreground mb-6">
        {format(new Date(), "yyyy年M月")}の集計
      </p>

      {loading && (
        <p className="text-muted-foreground text-sm">読み込み中...</p>
      )}
      {error && <p className="text-destructive text-sm mb-4">{error}</p>}

      {!loading && (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">カテゴリ別支出</CardTitle>
            </CardHeader>
            <CardContent>
              {categoryExpenses.length === 0 ? (
                <p className="text-muted-foreground text-sm">
                  データがありません
                </p>
              ) : (
                <div className="space-y-3">
                  {categoryExpenses
                    .sort(
                      (a, b) => Number(b.totalExpense) - Number(a.totalExpense),
                    )
                    .map((item) => (
                      <div key={item.categoryId} className="space-y-1">
                        <div className="flex justify-between text-sm">
                          <span>{item.categoryName}</span>
                          <span className="font-medium">
                            {formatAmount(item.totalExpense)}
                          </span>
                        </div>
                        <div className="h-2 bg-muted rounded-full overflow-hidden">
                          <div
                            className="h-full bg-primary rounded-full transition-all"
                            style={{
                              width: `${(Number(item.totalExpense) / maxExpense) * 100}%`,
                            }}
                          />
                        </div>
                      </div>
                    ))}
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-base">口座残高</CardTitle>
            </CardHeader>
            <CardContent>
              {accountBalances.length === 0 ? (
                <p className="text-muted-foreground text-sm">
                  データがありません
                </p>
              ) : (
                <div className="space-y-2">
                  {accountBalances.map((item) => (
                    <div
                      key={item.accountId}
                      className="flex justify-between py-2 border-b last:border-0"
                    >
                      <span className="text-sm">{item.accountName}</span>
                      <span
                        className={`font-medium text-sm ${Number(item.balance) >= 0 ? "text-green-600" : "text-red-600"}`}
                      >
                        {formatAmount(item.balance)}
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
