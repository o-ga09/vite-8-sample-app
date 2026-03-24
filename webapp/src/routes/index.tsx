import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  useListWorkspaces,
  getListWorkspacesQueryKey,
  useCreateWorkspace,
} from "@/lib/api/generated/workspaces/workspaces";
import { Plus } from "lucide-react";

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
    <div className="min-h-screen flex flex-col bg-background font-body">
      {/* Top Band */}
      <header className="flex items-center h-16 px-14 shrink-0 bg-sidebar">
        {/* Logo */}
        <div className="flex items-center gap-[10px]">
          <div className="w-8 h-8 rounded flex items-center justify-center flex-shrink-0 bg-brand">
            <span className="text-brand-foreground font-bold text-xs leading-none">
              家
            </span>
          </div>
          <span className="text-sidebar-foreground font-bold text-xl font-display">
            家計簿
          </span>
        </div>
        <div className="flex-1" />
        {/* User Info */}
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded flex items-center justify-center flex-shrink-0 bg-brand">
            <span className="text-brand-foreground font-semibold text-xs">
              田
            </span>
          </div>
          <span className="text-sidebar-foreground text-sm">田中 太郎</span>
        </div>
      </header>

      {/* Body */}
      <main className="flex flex-col gap-10 px-14 py-[60px] flex-1">
        {/* Head Row */}
        <div className="flex items-center gap-4 w-full">
          {/* Title Block */}
          <div className="flex flex-col gap-1">
            <h1
              className="font-bold text-[32px] leading-tight font-display text-foreground"
              style={{ letterSpacing: "-1px" }}
            >
              ワークスペースを選択
            </h1>
            <p className="text-sm text-muted-foreground">
              参加中のワークスペースから選択するか、新しく作成してください
            </p>
          </div>
          <div className="flex-1" />
          {/* New Button */}
          <button
            type="button"
            onClick={() => setOpen(true)}
            className="flex items-center gap-2 px-5 py-3 rounded-lg bg-brand text-brand-foreground text-sm font-medium transition-all hover:opacity-90 active:scale-95"
          >
            <Plus size={16} />
            新規作成
          </button>
        </div>

        {/* Error */}
        {isError && (
          <p className="text-sm text-destructive">
            ワークスペースの取得に失敗しました
          </p>
        )}

        {/* Workspace Grid */}
        <div className="grid grid-cols-3 gap-6 w-full">
          {isLoading && (
            <p className="text-sm col-span-3 text-muted-foreground">
              読み込み中...
            </p>
          )}

          {workspaces.map((ws) => (
            <button
              key={ws.id}
              type="button"
              className="flex flex-col gap-5 p-8 rounded-xl text-left bg-primary border border-sidebar-accent transition-all hover:shadow-md hover:-translate-y-0.5 cursor-pointer"
              onClick={() =>
                navigate({ to: "/workspaces/$wsId", params: { wsId: ws.id } })
              }
            >
              {/* Icon */}
              <div className="w-12 h-12 rounded-lg flex-shrink-0 bg-brand" />
              {/* Card Body */}
              <div className="flex flex-col gap-1.5 w-full">
                <span className="font-bold text-lg text-primary-foreground font-display">
                  {ws.name}
                </span>
              </div>
            </button>
          ))}

          {/* New Workspace Card */}
          {!isLoading && (
            <button
              type="button"
              className="flex flex-col gap-5 p-8 rounded-xl text-left bg-card border border-border transition-all hover:bg-muted cursor-pointer"
              onClick={() => setOpen(true)}
            >
              {/* Plus Icon */}
              <div className="w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0 bg-muted">
                <Plus size={20} className="text-muted-foreground" />
              </div>
              {/* Label */}
              <div className="flex flex-col justify-center flex-1 w-full">
                <span className="font-bold text-lg text-muted-foreground font-display">
                  新しいワークスペースを作成
                </span>
              </div>
            </button>
          )}
        </div>
      </main>

      {/* Dialog */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="font-display">
              新規ワークスペース作成
            </DialogTitle>
          </DialogHeader>
          <Input
            placeholder="ワークスペース名"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleCreate()}
          />
          <DialogFooter>
            <button
              type="button"
              className="px-4 py-2 rounded-lg text-sm font-medium bg-secondary text-muted-foreground transition-colors hover:opacity-80"
              onClick={() => setOpen(false)}
            >
              キャンセル
            </button>
            <button
              type="button"
              className="px-4 py-2 rounded-lg text-sm font-medium bg-brand text-brand-foreground transition-all hover:opacity-90 active:scale-95 disabled:opacity-40"
              onClick={handleCreate}
              disabled={createMutation.isPending || !newName.trim()}
            >
              作成
            </button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
