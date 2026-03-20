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

type Account = components["schemas"]["Account"];
type AccountType = components["schemas"]["AccountType"];

export const Route = createFileRoute("/workspaces/$wsId/accounts/")({
  component: AccountsPage,
});

const ACCOUNT_TYPE_LABELS: Record<AccountType, string> = {
  cash: "現金",
  bank: "銀行",
  credit_card: "クレジット",
  e_money: "電子マネー",
  investment: "投資",
};

const emptyForm = {
  name: "",
  accountType: "bank" as AccountType,
  initialBalance: "0",
};

function AccountsPage() {
  const { wsId } = Route.useParams();
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Account | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);

  const fetchAccounts = useCallback(async () => {
    setLoading(true);
    const { data, error: err } = await apiClient.GET("/workspaces/{wsId}/accounts", {
      params: { path: { wsId } },
    });
    if (err) {
      setError("口座の取得に失敗しました");
    } else {
      setAccounts(data ?? []);
    }
    setLoading(false);
  }, [wsId]);

  useEffect(() => {
    fetchAccounts();
  }, [fetchAccounts]);

  const openCreate = () => {
    setEditTarget(null);
    setForm(emptyForm);
    setOpen(true);
  };

  const openEdit = (account: Account) => {
    setEditTarget(account);
    setForm({
      name: account.name,
      accountType: account.accountType,
      initialBalance: account.initialBalance,
    });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!form.name.trim()) return;
    setSaving(true);

    if (editTarget) {
      const { error: err } = await apiClient.PUT("/workspaces/{wsId}/accounts/{accountId}", {
        params: { path: { wsId, accountId: editTarget.id } },
        body: { name: form.name, accountType: form.accountType },
      });
      if (err) setError("口座の更新に失敗しました");
    } else {
      const { error: err } = await apiClient.POST("/workspaces/{wsId}/accounts", {
        params: { path: { wsId } },
        body: {
          name: form.name,
          accountType: form.accountType,
          initialBalance: form.initialBalance,
        },
      });
      if (err) setError("口座の作成に失敗しました");
    }

    setSaving(false);
    setOpen(false);
    fetchAccounts();
  };

  const handleDelete = async (accountId: string) => {
    if (!confirm("この口座を削除しますか？")) return;
    const { error: err } = await apiClient.DELETE("/workspaces/{wsId}/accounts/{accountId}", {
      params: { path: { wsId, accountId } },
    });
    if (err) setError("口座の削除に失敗しました");
    else fetchAccounts();
  };

  const formatAmount = (val: string) =>
    new Intl.NumberFormat("ja-JP", {
      style: "currency",
      currency: "JPY",
    }).format(Number(val));

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">口座</h1>
        <Button onClick={openCreate}>
          <Plus size={16} className="mr-2" />
          口座を追加
        </Button>
      </div>

      {error && <p className="text-destructive text-sm mb-4">{error}</p>}

      {loading ? (
        <p className="text-muted-foreground text-sm">読み込み中...</p>
      ) : accounts.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">口座がありません</p>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>口座名</TableHead>
                <TableHead>種別</TableHead>
                <TableHead>初期残高</TableHead>
                <TableHead className="w-24"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {accounts.map((account) => (
                <TableRow key={account.id}>
                  <TableCell className="font-medium">{account.name}</TableCell>
                  <TableCell>
                    <Badge variant="secondary">{ACCOUNT_TYPE_LABELS[account.accountType]}</Badge>
                  </TableCell>
                  <TableCell>{formatAmount(account.initialBalance)}</TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button size="icon" variant="ghost" onClick={() => openEdit(account)}>
                        <Pencil size={14} />
                      </Button>
                      <Button size="icon" variant="ghost" onClick={() => handleDelete(account.id)}>
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
            <DialogTitle>{editTarget ? "口座を編集" : "口座を追加"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Input
              placeholder="口座名"
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            />
            <Select
              value={form.accountType}
              onValueChange={(v) => setForm((f) => ({ ...f, accountType: v as AccountType }))}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {(Object.keys(ACCOUNT_TYPE_LABELS) as AccountType[]).map((type) => (
                  <SelectItem key={type} value={type}>
                    {ACCOUNT_TYPE_LABELS[type]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {!editTarget && (
              <Input
                placeholder="初期残高"
                type="number"
                value={form.initialBalance}
                onChange={(e) => setForm((f) => ({ ...f, initialBalance: e.target.value }))}
              />
            )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>
              キャンセル
            </Button>
            <Button onClick={handleSave} disabled={saving || !form.name.trim()}>
              {editTarget ? "更新" : "作成"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
