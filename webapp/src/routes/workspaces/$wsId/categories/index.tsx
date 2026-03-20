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
  useListCategories,
  getListCategoriesQueryKey,
  useCreateCategory,
  useUpdateCategory,
  useDeleteCategory,
} from "@/lib/api/generated/categories/categories";
import type { Category, CategoryType } from "@/lib/api/model";
import { Plus, Pencil, Trash2 } from "lucide-react";

export const Route = createFileRoute("/workspaces/$wsId/categories/")({
  component: CategoriesPage,
});

const CATEGORY_TYPE_LABELS: Record<CategoryType, string> = {
  income: "収入",
  expense: "支出",
};

const emptyForm = {
  name: "",
  categoryType: "expense" as CategoryType,
  parentId: "",
};

function CategoriesPage() {
  const { wsId } = Route.useParams();
  const queryClient = useQueryClient();
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [editTarget, setEditTarget] = useState<Category | null>(null);
  const [form, setForm] = useState(emptyForm);

  const { data: listResult, isLoading } = useListCategories(wsId);
  const categories = listResult?.data ?? [];

  const invalidate = () =>
    queryClient.invalidateQueries({
      queryKey: getListCategoriesQueryKey(wsId),
    });

  const createMutation = useCreateCategory({
    mutation: {
      onSuccess: () => {
        invalidate();
        setOpen(false);
      },
      onError: () => setError("カテゴリの作成に失敗しました"),
    },
  });
  const updateMutation = useUpdateCategory({
    mutation: {
      onSuccess: () => {
        invalidate();
        setOpen(false);
      },
      onError: () => setError("カテゴリの更新に失敗しました"),
    },
  });
  const deleteMutation = useDeleteCategory({
    mutation: {
      onSuccess: () => {
        invalidate();
      },
      onError: () => setError("カテゴリの削除に失敗しました"),
    },
  });

  const saving = createMutation.isPending || updateMutation.isPending;

  const openCreate = () => {
    setEditTarget(null);
    setForm(emptyForm);
    setOpen(true);
  };

  const openEdit = (category: Category) => {
    setEditTarget(category);
    setForm({
      name: category.name,
      categoryType: category.categoryType,
      parentId: category.parentId ?? "",
    });
    setOpen(true);
  };

  const handleSave = () => {
    if (!form.name.trim()) return;
    const body = {
      name: form.name,
      categoryType: form.categoryType,
      parentId: form.parentId || null,
    };
    if (editTarget) {
      updateMutation.mutate({ wsId, categoryId: editTarget.id, data: body });
    } else {
      createMutation.mutate({ wsId, data: body });
    }
  };

  const handleDelete = (categoryId: string) => {
    if (!confirm("このカテゴリを削除しますか？")) return;
    deleteMutation.mutate({ wsId, categoryId });
  };

  const getParentName = (parentId: string | null | undefined) => {
    if (!parentId) return "-";
    return categories.find((c) => c.id === parentId)?.name ?? "-";
  };

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">カテゴリ</h1>
        <Button onClick={openCreate}>
          <Plus size={16} className="mr-2" />
          カテゴリを追加
        </Button>
      </div>

      {error && <p className="text-destructive text-sm mb-4">{error}</p>}

      {isLoading ? (
        <p className="text-muted-foreground text-sm">読み込み中...</p>
      ) : categories.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">カテゴリがありません</p>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>カテゴリ名</TableHead>
                <TableHead>種別</TableHead>
                <TableHead>親カテゴリ</TableHead>
                <TableHead className="w-24"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {categories.map((category) => (
                <TableRow key={category.id}>
                  <TableCell className="font-medium">
                    {category.parentId && <span className="text-muted-foreground mr-2">└</span>}
                    {category.name}
                  </TableCell>
                  <TableCell>
                    <Badge variant={category.categoryType === "income" ? "default" : "destructive"}>
                      {CATEGORY_TYPE_LABELS[category.categoryType]}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">{getParentName(category.parentId)}</TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button size="icon" variant="ghost" onClick={() => openEdit(category)}>
                        <Pencil size={14} />
                      </Button>
                      <Button size="icon" variant="ghost" onClick={() => handleDelete(category.id)}>
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
            <DialogTitle>{editTarget ? "カテゴリを編集" : "カテゴリを追加"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <Input
              placeholder="カテゴリ名"
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            />
            <Select
              value={form.categoryType}
              onValueChange={(v) => setForm((f) => ({ ...f, categoryType: v as CategoryType }))}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="income">収入</SelectItem>
                <SelectItem value="expense">支出</SelectItem>
              </SelectContent>
            </Select>
            <Select
              value={form.parentId || "none"}
              onValueChange={(v) =>
                setForm((f) => ({
                  ...f,
                  parentId: v === "none" || !v ? "" : v,
                }))
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="親カテゴリ（任意）" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">なし</SelectItem>
                {categories
                  .filter((c) => !c.parentId && c.id !== editTarget?.id)
                  .map((c) => (
                    <SelectItem key={c.id} value={c.id}>
                      {c.name}
                    </SelectItem>
                  ))}
              </SelectContent>
            </Select>
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
