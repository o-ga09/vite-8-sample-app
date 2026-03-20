import { http, HttpResponse } from "msw";
import { workspaces, members, accounts, categories, transactions, newId } from "./db";
import type {
  Workspace,
  Member,
  Account,
  Category,
  Transaction,
  CreateWorkspaceRequest,
  UpdateWorkspaceRequest,
  CreateMemberRequest,
  UpdateMemberRequest,
  CreateAccountRequest,
  UpdateAccountRequest,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CreateTransactionRequest,
  UpdateTransactionRequest,
} from "@/lib/api/model";

const BASE = import.meta.env.VITE_API_URL ?? "http://localhost:18080";

export const handlers = [
  // ============================================================
  // Workspaces
  // ============================================================
  http.get(`${BASE}/workspaces`, () => HttpResponse.json(workspaces)),

  http.post(`${BASE}/workspaces`, async ({ request }) => {
    const body = (await request.json()) as CreateWorkspaceRequest;
    const ws: Workspace = { id: newId(), name: body.name };
    workspaces.push(ws);
    return HttpResponse.json(ws, { status: 201 });
  }),

  http.get(`${BASE}/workspaces/:wsId`, ({ params }) => {
    const ws = workspaces.find((w) => w.id === params["wsId"]);
    if (!ws) return HttpResponse.json({ message: "Not Found" }, { status: 404 });
    return HttpResponse.json(ws);
  }),

  http.put(`${BASE}/workspaces/:wsId`, async ({ params, request }) => {
    const ws = workspaces.find((w) => w.id === params["wsId"]);
    if (!ws) return HttpResponse.json({ message: "Not Found" }, { status: 404 });
    const body = (await request.json()) as UpdateWorkspaceRequest;
    if (body.name) ws.name = body.name;
    return HttpResponse.json(ws);
  }),

  http.delete(`${BASE}/workspaces/:wsId`, ({ params }) => {
    const idx = workspaces.findIndex((w) => w.id === params["wsId"]);
    if (idx !== -1) workspaces.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  // ============================================================
  // Members
  // ============================================================
  http.get(`${BASE}/workspaces/:wsId/members`, ({ params }) =>
    HttpResponse.json(members.filter((m) => m.workspaceId === params["wsId"])),
  ),

  http.post(`${BASE}/workspaces/:wsId/members`, async ({ params, request }) => {
    const body = (await request.json()) as CreateMemberRequest;
    const member: Member = {
      id: newId(),
      workspaceId: params["wsId"] as string,
      email: body.email,
      displayName: body.displayName,
      role: body.role,
    };
    members.push(member);
    return HttpResponse.json(member, { status: 201 });
  }),

  http.put(`${BASE}/workspaces/:wsId/members/:memberId`, async ({ params, request }) => {
    const member = members.find((m) => m.id === params["memberId"]);
    if (!member) return HttpResponse.json({ message: "Not Found" }, { status: 404 });
    const body = (await request.json()) as UpdateMemberRequest;
    if (body.displayName) member.displayName = body.displayName;
    if (body.role) member.role = body.role;
    return HttpResponse.json(member);
  }),

  http.delete(`${BASE}/workspaces/:wsId/members/:memberId`, ({ params }) => {
    const idx = members.findIndex((m) => m.id === params["memberId"]);
    if (idx !== -1) members.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  // ============================================================
  // Accounts
  // ============================================================
  http.get(`${BASE}/workspaces/:wsId/accounts`, ({ params }) =>
    HttpResponse.json(accounts.filter((a) => a.workspaceId === params["wsId"])),
  ),

  http.post(`${BASE}/workspaces/:wsId/accounts`, async ({ params, request }) => {
    const body = (await request.json()) as CreateAccountRequest;
    const account: Account = {
      id: newId(),
      workspaceId: params["wsId"] as string,
      name: body.name,
      accountType: body.accountType,
      initialBalance: body.initialBalance ?? "0",
    };
    accounts.push(account);
    return HttpResponse.json(account, { status: 201 });
  }),

  http.put(`${BASE}/workspaces/:wsId/accounts/:accountId`, async ({ params, request }) => {
    const account = accounts.find((a) => a.id === params["accountId"]);
    if (!account) return HttpResponse.json({ message: "Not Found" }, { status: 404 });
    const body = (await request.json()) as UpdateAccountRequest;
    if (body.name) account.name = body.name;
    if (body.accountType) account.accountType = body.accountType;
    return HttpResponse.json(account);
  }),

  http.delete(`${BASE}/workspaces/:wsId/accounts/:accountId`, ({ params }) => {
    const idx = accounts.findIndex((a) => a.id === params["accountId"]);
    if (idx !== -1) accounts.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  // ============================================================
  // Categories
  // ============================================================
  http.get(`${BASE}/workspaces/:wsId/categories`, ({ params }) =>
    HttpResponse.json(categories.filter((c) => c.workspaceId === params["wsId"])),
  ),

  http.post(`${BASE}/workspaces/:wsId/categories`, async ({ params, request }) => {
    const body = (await request.json()) as CreateCategoryRequest;
    const category: Category = {
      id: newId(),
      workspaceId: params["wsId"] as string,
      name: body.name,
      categoryType: body.categoryType,
      parentId: body.parentId ?? null,
    };
    categories.push(category);
    return HttpResponse.json(category, { status: 201 });
  }),

  http.put(`${BASE}/workspaces/:wsId/categories/:categoryId`, async ({ params, request }) => {
    const category = categories.find((c) => c.id === params["categoryId"]);
    if (!category) return HttpResponse.json({ message: "Not Found" }, { status: 404 });
    const body = (await request.json()) as UpdateCategoryRequest;
    if (body.name) category.name = body.name;
    if (body.categoryType) category.categoryType = body.categoryType;
    if ("parentId" in body) category.parentId = body.parentId ?? null;
    return HttpResponse.json(category);
  }),

  http.delete(`${BASE}/workspaces/:wsId/categories/:categoryId`, ({ params }) => {
    const idx = categories.findIndex((c) => c.id === params["categoryId"]);
    if (idx !== -1) categories.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  // ============================================================
  // Transactions
  // ============================================================
  http.get(`${BASE}/workspaces/:wsId/transactions`, ({ params }) =>
    HttpResponse.json(transactions.filter((t) => t.workspaceId === params["wsId"])),
  ),

  http.post(`${BASE}/workspaces/:wsId/transactions`, async ({ params, request }) => {
    const body = (await request.json()) as CreateTransactionRequest;
    const transaction: Transaction = {
      id: newId(),
      workspaceId: params["wsId"] as string,
      transactionType: body.transactionType,
      accountId: body.accountId ?? null,
      counterpartyAccountId: body.counterpartyAccountId ?? null,
      categoryId: body.categoryId ?? null,
      amount: body.amount,
      occurredAt: body.occurredAt,
      description: body.description ?? null,
    };
    transactions.push(transaction);
    return HttpResponse.json(transaction, { status: 201 });
  }),

  http.put(`${BASE}/workspaces/:wsId/transactions/:txId`, async ({ params, request }) => {
    const transaction = transactions.find((t) => t.id === params["txId"]);
    if (!transaction) return HttpResponse.json({ message: "Not Found" }, { status: 404 });
    const body = (await request.json()) as UpdateTransactionRequest;
    if (body.transactionType) transaction.transactionType = body.transactionType;
    if (body.amount) transaction.amount = body.amount;
    if (body.occurredAt) transaction.occurredAt = body.occurredAt;
    if ("description" in body) transaction.description = body.description ?? null;
    if ("accountId" in body) transaction.accountId = body.accountId ?? null;
    if ("categoryId" in body) transaction.categoryId = body.categoryId ?? null;
    return HttpResponse.json(transaction);
  }),

  http.delete(`${BASE}/workspaces/:wsId/transactions/:txId`, ({ params }) => {
    const idx = transactions.findIndex((t) => t.id === params["txId"]);
    if (idx !== -1) transactions.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  // ============================================================
  // Reports
  // ============================================================
  http.get(`${BASE}/workspaces/:wsId/reports/dashboard`, ({ params }) => {
    const wsTransactions = transactions.filter((t) => t.workspaceId === params["wsId"]);
    const totalIncome = wsTransactions
      .filter((t) => t.transactionType === "income")
      .reduce((sum, t) => sum + Number(t.amount), 0);
    const totalExpense = wsTransactions
      .filter((t) => t.transactionType === "expense")
      .reduce((sum, t) => sum + Number(t.amount), 0);
    return HttpResponse.json({
      workspaceId: params["wsId"],
      totalIncome: String(totalIncome),
      totalExpense: String(totalExpense),
      netFlow: String(totalIncome - totalExpense),
      transactionCount: wsTransactions.length,
    });
  }),

  http.get(`${BASE}/workspaces/:wsId/reports/category-expenses`, ({ params }) => {
    const wsTransactions = transactions.filter(
      (t) => t.workspaceId === params["wsId"] && t.transactionType === "expense",
    );
    const expenseMap = new Map<string, number>();
    for (const tx of wsTransactions) {
      if (!tx.categoryId) continue;
      expenseMap.set(tx.categoryId, (expenseMap.get(tx.categoryId) ?? 0) + Number(tx.amount));
    }
    const wsCategories = categories.filter((c) => c.workspaceId === params["wsId"]);
    const result = Array.from(expenseMap.entries()).map(([categoryId, total]) => ({
      categoryId,
      categoryName: wsCategories.find((c) => c.id === categoryId)?.name ?? categoryId,
      totalExpense: String(total),
    }));
    return HttpResponse.json(result);
  }),

  http.get(`${BASE}/workspaces/:wsId/reports/account-balances`, ({ params }) => {
    const wsAccounts = accounts.filter((a) => a.workspaceId === params["wsId"]);
    const result = wsAccounts.map((a) => {
      const txSum = transactions
        .filter((t) => t.accountId === a.id)
        .reduce((sum, t) => {
          if (t.transactionType === "income") return sum + Number(t.amount);
          if (t.transactionType === "expense") return sum - Number(t.amount);
          return sum;
        }, 0);
      return {
        accountId: a.id,
        accountName: a.name,
        accountType: a.accountType,
        currentBalance: String(Number(a.initialBalance) + txSum),
      };
    });
    return HttpResponse.json(result);
  }),
];
