import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { apiClient } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { Plus, Pencil, Trash2 } from "lucide-react";
import { format } from "date-fns";

type Transaction = components["schemas"]["Transaction"];
type TransactionType = components["schemas"]["TransactionType"];

export const Route = createFileRoute("/workspaces/$wsId/transactions/")({
  component: TransactionsPage,
});

const TX_TYPE_LABELS: Record<TransactionType, string> = {
  income: "収入",
  expense: "支出",
  transfer: "振替",
};

const TX_TYPE_VARIANT: Record<TransactionType, "default" | "destructive" | "secondary"> = {
  income: "default",
  expense: "destructive",
  transfer: "secondary",
};

const emptyForm = {
  transactionType: "expense" as TransactionType,
  amount: "",
  occurredAt: format(new Date(), "yyyy-MM-dd'T'HH:mm"),
  description: "",
  accountId: "",
  categoryId: "",
};

function TransactionsPage() {
  const { wsId } = Route.useParams();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Transaction | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);

  const fetchTransactions = useCallback(async () => {
    setLoading(true);
    const { data, error: err } = await apiClient.GET("/workspaces/{wsId}/transactions", {
      params: { path: { wsId } },
    });
    if (err) {
      setError("取引の取得に失敗しました");
    } else {
      setTransactions(data ?? []);
    }
    setLoading(false);
  }, [wsId]);

  useEffect(() => {
    fetchTransactions();
  }, [fetchTransactions]);

  const openCreate = () => {
    setEditTarget(null);
    setForm(emptyForm);
    setOpen(true);
  };

  const openEdit = (tx: Transaction) => {
    setEditTarget(tx);
    setForm({
      transactionType: tx.transactionType,
      amount: tx.amount,
      occurredAt: format(new Date(tx.occurredAt), "yyyy-MM-dd'T'HH:mm"),
      description: tx.description ?? "",
      accountId: tx.accountId ?? "",
      categoryId: tx.categoryId ?? "",
    });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!form.amount || !form.occurredAt) return;
    setSaving(true);

    const body = {
      transactionType: form.transactionType,
      amount: form.amount,
      occurredAt: new Date(form.occurredAt).toISOString(),
      description: form.description || null,
      accountId: form.accountId || null,
      categoryId: form.categoryId || null,
      counterpartyAccountId: null,
    };

    if (editTarget) {
      const { error: err } = await apiClient.PUT("/workspaces/{wsId}/transactions/{txId}", {
        params: { path: { wsId, txId: editTarget.id } },
        body,
      });
      if (err) setError("取引の更新に失敗しました");
    } else {
      const { error: err } = await apiClient.POST("/workspaces/{wsId}/transactions", {
        params: { path: { wsId } },
        body,
      });
      if (err) setError("取引の作成に失敗しました");
    }

    setSaving(false);
    setOpen(false);
    fetchTransactions();
  };

  const handleDelete = async (txId: string) => {
    if (!confirm("この取引を削除しますか？")) return;
    const { error: err } = await apiClient.DELETE("/workspaces/{wsId}/transactions/{txId}", {
      params: { path: { wsId, txId } },
    });
    if (err) setError("取引の削除に失敗しました");
    else fetchTransactions();
  };

  const formatAmount = (val: string) =>
    new Intl.NumberFormat("ja-JP", {
      style: "currency",
      currency: "JPY",
    }).format(Number(val));

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">取引</h1>
        <Button onClick={openCreate}>
          <Plus size={16} className="mr-2" />
          取引を追加
        </Button>
      </div>

      {error && <p className="text-destructive text-sm mb-4">{error}</p>}

      {loading ? (
        <p className="text-muted-foreground text-sm">読み込み中...</p>
      ) : transactions.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">取引がありません</p>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>日付</TableHead>
                <TableHead>種別</TableHead>
                <TableHead>金額</TableHead>
                <TableHead>説明</TableHead>
                <TableHead className="w-24"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.map((tx) => (
                <TableRow key={tx.id}>
                  <TableCell className="text-sm">{format(new Date(tx.occurredAt), "yyyy/MM/dd")}</TableCell>
                  <TableCell>
                    <Badge variant={TX_TYPE_VARIANT[tx.transactionType]}>{TX_TYPE_LABELS[tx.transactionType]}</Badge>
                  </TableCell>
                  <TableCell className="font-medium">{formatAmount(tx.amount)}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">{tx.description ?? "-"}</TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button size="icon" variant="ghost" onClick={() => openEdit(tx)}>
                        <Pencil size={14} />
                      </Button>
                      <Button size="icon" variant="ghost" onClick={() => handleDelete(tx.id)}>
                        <Trash2 size={14} />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editTarget ? "取引を編集" : "取引を追加"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Select
              value={form.transactionType}
              onValueChange={(v) =>
                setForm((f) => ({
                  ...f,
                  transactionType: v as TransactionType,
                }))
              }
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="income">収入</SelectItem>
                <SelectItem value="expense">支出</SelectItem>
                <SelectItem value="transfer">振替</SelectItem>
              </SelectContent>
            </Select>
            <Input
              placeholder="金額"
              type="number"
              value={form.amount}
              onChange={(e) => setForm((f) => ({ ...f, amount: e.target.value }))}
            />
            <Input
              type="datetime-local"
              value={form.occurredAt}
              onChange={(e) => setForm((f) => ({ ...f, occurredAt: e.target.value }))}
            />
            <Input
              placeholder="説明（任意）"
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>
              キャンセル
            </Button>
            <Button onClick={handleSave} disabled={saving}>
              {editTarget ? "更新" : "作成"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
