import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { apiClient } from "@/lib/api/client";
import type { components } from "@/lib/api/schema";
import { Plus, Folder } from "lucide-react";

type Workspace = components["schemas"]["Workspace"];

export const Route = createFileRoute("/")({
  component: IndexPage,
});

function IndexPage() {
  const navigate = useNavigate();
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [newName, setNewName] = useState("");
  const [creating, setCreating] = useState(false);

  const fetchWorkspaces = async () => {
    setLoading(true);
    const { data, error: err } = await apiClient.GET("/workspaces");
    if (err) {
      setError("ワークスペースの取得に失敗しました");
    } else {
      setWorkspaces(data ?? []);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchWorkspaces();
  }, []);

  const handleCreate = async () => {
    if (!newName.trim()) return;
    setCreating(true);
    const { data, error: err } = await apiClient.POST("/workspaces", {
      body: { name: newName.trim() },
    });
    setCreating(false);
    if (err) {
      setError("ワークスペースの作成に失敗しました");
      return;
    }
    setOpen(false);
    setNewName("");
    if (data) {
      await navigate({ to: "/workspaces/$wsId", params: { wsId: data.id } });
    } else {
      fetchWorkspaces();
    }
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

        {loading && <p className="text-muted-foreground text-sm">読み込み中...</p>}
        {error && <p className="text-destructive text-sm">{error}</p>}

        {!loading && workspaces.length === 0 && (
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
            <Button onClick={handleCreate} disabled={creating || !newName.trim()}>
              作成
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
