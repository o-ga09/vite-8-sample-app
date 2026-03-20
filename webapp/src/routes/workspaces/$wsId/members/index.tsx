import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { apiClient } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { Plus, Pencil, Trash2 } from "lucide-react";

type Member = components["schemas"]["Member"];
type MemberRole = components["schemas"]["MemberRole"];

export const Route = createFileRoute("/workspaces/$wsId/members/")({
  component: MembersPage,
});

const ROLE_LABELS: Record<MemberRole, string> = {
  owner: "オーナー",
  editor: "編集者",
  viewer: "閲覧者",
};

const ROLE_VARIANT: Record<MemberRole, "default" | "secondary" | "outline"> = {
  owner: "default",
  editor: "secondary",
  viewer: "outline",
};

const emptyForm = {
  email: "",
  displayName: "",
  role: "viewer" as MemberRole,
};

function MembersPage() {
  const { wsId } = Route.useParams();
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Member | null>(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);

  const fetchMembers = useCallback(async () => {
    setLoading(true);
    const { data, error: err } = await apiClient.GET("/workspaces/{wsId}/members", {
      params: { path: { wsId } },
    });
    if (err) {
      setError("メンバーの取得に失敗しました");
    } else {
      setMembers(data ?? []);
    }
    setLoading(false);
  }, [wsId]);

  useEffect(() => {
    fetchMembers();
  }, [fetchMembers]);

  const openCreate = () => {
    setEditTarget(null);
    setForm(emptyForm);
    setOpen(true);
  };

  const openEdit = (member: Member) => {
    setEditTarget(member);
    setForm({
      email: member.email,
      displayName: member.displayName,
      role: member.role,
    });
    setOpen(true);
  };

  const handleSave = async () => {
    if (!form.displayName.trim()) return;
    setSaving(true);

    if (editTarget) {
      const { error: err } = await apiClient.PUT("/workspaces/{wsId}/members/{memberId}", {
        params: { path: { wsId, memberId: editTarget.id } },
        body: { displayName: form.displayName, role: form.role },
      });
      if (err) setError("メンバーの更新に失敗しました");
    } else {
      const { error: err } = await apiClient.POST("/workspaces/{wsId}/members", {
        params: { path: { wsId } },
        body: { email: form.email, displayName: form.displayName, role: form.role },
      });
      if (err) setError("メンバーの追加に失敗しました");
    }

    setSaving(false);
    setOpen(false);
    fetchMembers();
  };

  const handleDelete = async (memberId: string) => {
    if (!confirm("このメンバーを削除しますか？")) return;
    const { error: err } = await apiClient.DELETE("/workspaces/{wsId}/members/{memberId}", {
      params: { path: { wsId, memberId } },
    });
    if (err) setError("メンバーの削除に失敗しました");
    else fetchMembers();
  };

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">メンバー</h1>
        <Button onClick={openCreate}>
          <Plus size={16} className="mr-2" />
          メンバーを追加
        </Button>
      </div>

      {error && <p className="text-destructive text-sm mb-4">{error}</p>}

      {loading ? (
        <p className="text-muted-foreground text-sm">読み込み中...</p>
      ) : members.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">メンバーがいません</p>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>表示名</TableHead>
                <TableHead>メール</TableHead>
                <TableHead>役割</TableHead>
                <TableHead className="w-24"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {members.map((member) => (
                <TableRow key={member.id}>
                  <TableCell className="font-medium">{member.displayName}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">{member.email}</TableCell>
                  <TableCell>
                    <Badge variant={ROLE_VARIANT[member.role]}>{ROLE_LABELS[member.role]}</Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button size="icon" variant="ghost" onClick={() => openEdit(member)}>
                        <Pencil size={14} />
                      </Button>
                      <Button size="icon" variant="ghost" onClick={() => handleDelete(member.id)}>
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
            <DialogTitle>{editTarget ? "メンバーを編集" : "メンバーを追加"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            {!editTarget && (
              <Input
                placeholder="メールアドレス"
                type="email"
                value={form.email}
                onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
              />
            )}
            <Input
              placeholder="表示名"
              value={form.displayName}
              onChange={(e) => setForm((f) => ({ ...f, displayName: e.target.value }))}
            />
            <Select
              value={form.role}
              onValueChange={(v) => setForm((f) => ({ ...f, role: v as MemberRole }))}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="owner">オーナー</SelectItem>
                <SelectItem value="editor">編集者</SelectItem>
                <SelectItem value="viewer">閲覧者</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>
              キャンセル
            </Button>
            <Button onClick={handleSave} disabled={saving || !form.displayName.trim()}>
              {editTarget ? "更新" : "追加"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
