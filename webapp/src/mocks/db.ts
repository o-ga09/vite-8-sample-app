import type {
  Workspace,
  Member,
  Account,
  Category,
  Transaction,
} from "@/lib/api/model";

// --- Workspaces ---
export const workspaces: Workspace[] = [
  { id: "ws-1", name: "田中家計簿" },
  { id: "ws-2", name: "サンプル家計" },
];

// --- Members ---
export const members: Member[] = [
  {
    id: "m-1",
    workspaceId: "ws-1",
    email: "taro@example.com",
    displayName: "田中太郎",
    role: "owner",
  },
  {
    id: "m-2",
    workspaceId: "ws-1",
    email: "hanako@example.com",
    displayName: "田中花子",
    role: "editor",
  },
];

// --- Accounts ---
export const accounts: Account[] = [
  {
    id: "acc-1",
    workspaceId: "ws-1",
    name: "財布",
    accountType: "cash",
    initialBalance: "50000",
  },
  {
    id: "acc-2",
    workspaceId: "ws-1",
    name: "銀行口座",
    accountType: "bank",
    initialBalance: "1000000",
  },
  {
    id: "acc-3",
    workspaceId: "ws-1",
    name: "Suica",
    accountType: "e_money",
    initialBalance: "5000",
  },
];

// --- Categories ---
export const categories: Category[] = [
  {
    id: "cat-1",
    workspaceId: "ws-1",
    name: "食費",
    categoryType: "expense",
    parentId: null,
  },
  {
    id: "cat-2",
    workspaceId: "ws-1",
    name: "外食",
    categoryType: "expense",
    parentId: "cat-1",
  },
  {
    id: "cat-3",
    workspaceId: "ws-1",
    name: "給与",
    categoryType: "income",
    parentId: null,
  },
  {
    id: "cat-4",
    workspaceId: "ws-1",
    name: "交通費",
    categoryType: "expense",
    parentId: null,
  },
];

// --- Transactions ---
export const transactions: Transaction[] = [
  {
    id: "tx-1",
    workspaceId: "ws-1",
    transactionType: "income",
    accountId: "acc-2",
    counterpartyAccountId: null,
    categoryId: "cat-3",
    amount: "300000",
    occurredAt: "2026-03-25T09:00:00Z",
    description: "3月分給与",
  },
  {
    id: "tx-2",
    workspaceId: "ws-1",
    transactionType: "expense",
    accountId: "acc-1",
    counterpartyAccountId: null,
    categoryId: "cat-1",
    amount: "8500",
    occurredAt: "2026-03-18T12:30:00Z",
    description: "週末の買い物",
  },
  {
    id: "tx-3",
    workspaceId: "ws-1",
    transactionType: "expense",
    accountId: "acc-3",
    counterpartyAccountId: null,
    categoryId: "cat-4",
    amount: "2340",
    occurredAt: "2026-03-17T08:15:00Z",
    description: "電車代",
  },
  {
    id: "tx-4",
    workspaceId: "ws-1",
    transactionType: "expense",
    accountId: "acc-1",
    counterpartyAccountId: null,
    categoryId: "cat-2",
    amount: "3200",
    occurredAt: "2026-03-15T19:00:00Z",
    description: "外食",
  },
];

/** UUID v4 simplistic generator for mock */
export function newId(): string {
  return crypto.randomUUID();
}
