import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { apiClient } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { TrendingUp, TrendingDown, ArrowRightLeft, Hash } from "lucide-react";
import { format, startOfMonth, endOfMonth } from "date-fns";

type DashboardReport = components["schemas"]["DashboardReport"];

export const Route = createFileRoute("/workspaces/$wsId/")({
  component: DashboardPage,
});

function DashboardPage() {
  const { wsId } = Route.useParams();
  const [report, setReport] = useState<DashboardReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const now = new Date();
    const from = format(startOfMonth(now), "yyyy-MM-dd'T'HH:mm:ss'Z'");
    const to = format(endOfMonth(now), "yyyy-MM-dd'T'HH:mm:ss'Z'");

    apiClient
      .GET("/workspaces/{wsId}/reports/dashboard", {
        params: { path: { wsId }, query: { from, to } },
      })
      .then(({ data, error: err }) => {
        if (err) {
          setError("ダッシュボードデータの取得に失敗しました");
        } else {
          setReport(data ?? null);
        }
        setLoading(false);
      });
  }, [wsId]);

  const formatAmount = (val: string) =>
    new Intl.NumberFormat("ja-JP", { style: "currency", currency: "JPY" }).format(Number(val));

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6">ダッシュボード</h1>
      <p className="text-sm text-muted-foreground mb-6">{format(new Date(), "yyyy年M月")}の集計</p>

      {loading && <p className="text-muted-foreground text-sm">読み込み中...</p>}
      {error && <p className="text-destructive text-sm">{error}</p>}

      {report && (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <TrendingUp size={16} className="text-green-500" />
                収入合計
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold text-green-600">
                {formatAmount(report.totalIncome)}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <TrendingDown size={16} className="text-red-500" />
                支出合計
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold text-red-600">{formatAmount(report.totalExpense)}</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <ArrowRightLeft size={16} className="text-blue-500" />
                収支
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p
                className={`text-2xl font-bold ${Number(report.netFlow) >= 0 ? "text-green-600" : "text-red-600"}`}
              >
                {formatAmount(report.netFlow)}
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                <Hash size={16} />
                取引件数
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{report.transactionCount}</p>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
