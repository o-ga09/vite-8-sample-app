import { createFileRoute } from "@tanstack/react-router";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useGetDashboard } from "@/lib/api/generated/reports/reports";
import { TrendingUp, TrendingDown, ArrowRightLeft, Hash } from "lucide-react";
import { format, startOfMonth, endOfMonth } from "date-fns";

export const Route = createFileRoute("/workspaces/$wsId/")({
  component: DashboardPage,
});

function DashboardPage() {
  const { wsId } = Route.useParams();

  const now = new Date();
  const from = format(startOfMonth(now), "yyyy-MM-dd'T'HH:mm:ss'Z'");
  const to = format(endOfMonth(now), "yyyy-MM-dd'T'HH:mm:ss'Z'");

  const { data: dashboardResult, isLoading, isError } = useGetDashboard(wsId, { from, to });
  const report = dashboardResult?.data;

  const formatAmount = (val: string) =>
    new Intl.NumberFormat("ja-JP", {
      style: "currency",
      currency: "JPY",
    }).format(Number(val));

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6">ダッシュボード</h1>
      <p className="text-sm text-muted-foreground mb-6">{format(new Date(), "yyyy年M月")}の集計</p>

      {isLoading && <p className="text-muted-foreground text-sm">読み込み中...</p>}
      {isError && <p className="text-destructive text-sm">ダッシュボードデータの取得に失敗しました</p>}

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
              <p className="text-2xl font-bold text-green-600">{formatAmount(report.totalIncome)}</p>
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
              <p className={`text-2xl font-bold ${Number(report.netFlow) >= 0 ? "text-green-600" : "text-red-600"}`}>
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
