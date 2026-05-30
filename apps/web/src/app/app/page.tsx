"use client";

import Image from "next/image";
import { ProtectedRoute } from "@/components/protected-route";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";

function AppShell() {
  const { user, logout } = useAuth();

  return (
    <main className="mx-auto flex min-h-full max-w-3xl flex-col gap-8 px-6 py-16">
      <div className="space-y-2">
        <h1 className="text-2xl font-semibold text-foreground">Workspace</h1>
        <p className="text-muted">Signed in via GitHub OAuth.</p>
      </div>

      {user ? (
        <section className="flex items-center gap-4 rounded-lg border border-border bg-background p-4">
          <Image
            src={user.avatarUrl}
            alt={user.name}
            width={56}
            height={56}
            className="rounded-full"
          />
          <div className="space-y-1">
            <p className="font-medium text-foreground">{user.name}</p>
            <p className="text-sm text-muted">@{user.githubUsername}</p>
            <p className="text-sm text-muted">{user.email}</p>
          </div>
        </section>
      ) : (
        <p className="text-muted">Loading profile...</p>
      )}

      <Button variant="ghost" onClick={logout}>
        Sign out
      </Button>
    </main>
  );
}

export default function AppPage() {
  return (
    <ProtectedRoute>
      <AppShell />
    </ProtectedRoute>
  );
}
