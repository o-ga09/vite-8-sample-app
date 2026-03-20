import { createFileRoute, Link, Outlet, useParams } from "@tanstack/react-router";
import { BarChart3, CreditCard, FolderTree, LayoutDashboard, Users, Wallet } from "lucide-react";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/workspaces/$wsId")({
  component: WorkspaceLayout,
});

const navItems = [
  { to: "/workspaces/$wsId", label: "ダッシュボード", icon: LayoutDashboard },
  { to: "/workspaces/$wsId/transactions", label: "取引", icon: CreditCard },
  { to: "/workspaces/$wsId/accounts", label: "口座", icon: Wallet },
  { to: "/workspaces/$wsId/categories", label: "カテゴリ", icon: FolderTree },
  { to: "/workspaces/$wsId/members", label: "メンバー", icon: Users },
  { to: "/workspaces/$wsId/reports", label: "レポート", icon: BarChart3 },
] as const;

function WorkspaceLayout() {
  const { wsId } = useParams({ from: "/workspaces/$wsId" });

  return (
    <div className="flex min-h-screen bg-background">
      <aside className="w-60 border-r bg-card flex flex-col">
        <div className="p-4 border-b">
          <Link to="/" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
            ← ワークスペース一覧
          </Link>
        </div>
        <nav className="flex-1 p-2 space-y-1">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              params={{ wsId }}
              activeOptions={{ exact: item.to.endsWith("/$wsId/") }}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors",
                "hover:bg-accent hover:text-accent-foreground",
                "[&.active]:bg-accent [&.active]:text-accent-foreground [&.active]:font-medium",
              )}
            >
              <item.icon size={16} />
              {item.label}
            </Link>
          ))}
        </nav>
      </aside>
      <main className="flex-1 overflow-auto">
        <Outlet />
      </main>
    </div>
  );
}
