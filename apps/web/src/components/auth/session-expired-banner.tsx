"use client";

import { useSearchParams } from "next/navigation";
import { Suspense } from "react";

function Banner() {
  const searchParams = useSearchParams();
  const expired = searchParams.get("reason") === "session_expired";

  if (!expired) return null;

  return (
    <div
      role="alert"
      className="rounded-lg border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
    >
      Your session expired. Please sign in again to continue.
    </div>
  );
}

export function SessionExpiredBanner() {
  return (
    <Suspense fallback={null}>
      <Banner />
    </Suspense>
  );
}
