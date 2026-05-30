"use client";

import { Suspense, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useAuth } from "@/hooks/use-auth";

function CallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);

  const accessToken = searchParams.get("accessToken");
  const refreshToken = searchParams.get("refreshToken");
  const missingTokens = !accessToken || !refreshToken;

  useEffect(() => {
    if (missingTokens) return;

    const completeLogin = async () => {
      try {
        await login(accessToken!, refreshToken!);
        router.replace("/app");
      } catch {
        setError("Could not complete sign in. Please try again.");
      }
    };

    void completeLogin();
  }, [accessToken, login, missingTokens, refreshToken, router]);

  if (missingTokens) {
    return (
      <main className="mx-auto max-w-md px-6 py-16 text-destructive">
        Missing tokens in callback URL.
      </main>
    );
  }

  if (error) {
    return (
      <main className="mx-auto max-w-md px-6 py-16 text-destructive">{error}</main>
    );
  }

  return (
    <main className="mx-auto max-w-md px-6 py-16 text-muted">
      Completing sign in...
    </main>
  );
}

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <main className="mx-auto max-w-md px-6 py-16 text-muted">
          Completing sign in...
        </main>
      }
    >
      <CallbackHandler />
    </Suspense>
  );
}
