import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import {
  useListWorkspaces,
  getListWorkspacesQueryKey,
  useCreateWorkspace,
} from "@/lib/api/generated/workspaces/workspaces";
import { Plus, Folder } from "lucide-react";

export const Route = createFileRoute("/")({
  component: IndexPage,
});

function IndexPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [newName, setNewName] = useState("");

  const { data: listResult, isLoading, isError } = useListWorkspaces();
  const workspaces = listResult?.data ?? [];

  const createMutation = useCreateWorkspace({
    mutation: {
      onSuccess: async (result) => {
        await queryClient.invalidateQueries({
          queryKey: getListWorkspacesQueryKey(),
        });
        setOpen(false);
        setNewName("");
        if (result.status === 201) {
          await navigate({
            to: "/workspaces/$wsId",
            params: { wsId: result.data.id },
          });
        }
      },
    },
  });

  const handleCreate = () => {
    if (!newName.trim() || createMutation.isPending) return;
    createMutation.mutate({ data: { name: newName.trim() } });
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-2xl mx-auto p-8">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-2xl font-bold">ワークスペース</h1>
          <Button onClick={() => setOpen(true)}>
            <Plus size={16} className="mr-2" />
            新規作成
          </Button>
        </div>

        {isLoading && <p className="text-muted-foreground text-sm">読み込み中...</p>}
        {isError && <p className="text-destructive text-sm">ワークスペースの取得に失敗しました</p>}

        {!isLoading && workspaces.length === 0 && (
          <Card className="text-center py-12">
            <CardContent>
              <p className="text-muted-foreground">ワークスペースがありません</p>
              <Button variant="outline" className="mt-4" onClick={() => setOpen(true)}>
                最初のワークスペースを作成する
              </Button>
            </CardContent>
          </Card>
        )}

        <div className="space-y-3">
          {workspaces.map((ws) => (
            <Card
              key={ws.id}
              className="cursor-pointer hover:shadow-md transition-shadow"
              onClick={() => navigate({ to: "/workspaces/$wsId", params: { wsId: ws.id } })}
            >
              <CardHeader className="py-4">
                <CardTitle className="flex items-center gap-2 text-base">
                  <Folder size={18} className="text-muted-foreground" />
                  {ws.name}
                </CardTitle>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新規ワークスペース作成</DialogTitle>
          </DialogHeader>
          <Input
            placeholder="ワークスペース名"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleCreate()}
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>
              キャンセル
            </Button>
            <Button onClick={handleCreate} disabled={createMutation.isPending || !newName.trim()}>
              作成
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
