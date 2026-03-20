import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  useListAccounts,
  getListAccountsQueryKey,
  useCreateAccount,
  useUpdateAccount,
  useDeleteAccount,
} from "@/lib/api/generated/accounts/accounts";
import type { Account, AccountType } from "@/lib/api/model";
import { Plus, Pencil, Trash2 } from "lucide-react";

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
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Account | null>(null);
  const [form, setForm] = useState(emptyForm);

  const { data: listResult, isLoading } = useListAccounts(wsId);
  const accounts = listResult?.data ?? [];

  const invalidate = () => queryClient.invalidateQueries({ queryKey: getListAccountsQueryKey(wsId) });

  const createMutation = useCreateAccount({
    mutation: {
      onSuccess: () => {
        invalidate();
        setOpen(false);
      },
      onError: () => setError("口座の作成に失敗しました"),
    },
  });
  const updateMutation = useUpdateAccount({
    mutation: {
      onSuccess: () => {
        invalidate();
        setOpen(false);
      },
      onError: () => setError("口座の更新に失敗しました"),
    },
  });
  const deleteMutation = useDeleteAccount({
    mutation: {
      onSuccess: () => {
        invalidate();
      },
      onError: () => setError("口座の削除に失敗しました"),
    },
  });

  const saving = createMutation.isPending || updateMutation.isPending;

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

  const handleSave = () => {
    if (!form.name.trim()) return;
    if (editTarget) {
      updateMutation.mutate({
        wsId,
        accountId: editTarget.id,
        data: { name: form.name, accountType: form.accountType },
      });
    } else {
      createMutation.mutate({
        wsId,
        data: {
          name: form.name,
          accountType: form.accountType,
          initialBalance: form.initialBalance,
        },
      });
    }
  };

  const handleDelete = (accountId: string) => {
    if (!confirm("この口座を削除しますか？")) return;
    deleteMutation.mutate({ wsId, accountId });
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

      {isLoading ? (
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
