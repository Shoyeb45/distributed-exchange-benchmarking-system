"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";
import { useAuth } from "@/hooks/use-auth";

const publicEntryPaths = new Set(["/", "/login"]);

export function SessionRedirect() {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (isLoading || !isAuthenticated || !publicEntryPaths.has(pathname)) {
      return;
    }

    router.replace("/app");
  }, [isAuthenticated, isLoading, pathname, router]);

  return null;
}
