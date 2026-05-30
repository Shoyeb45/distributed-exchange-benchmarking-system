"use client";

import { Toaster } from "sonner";
import { AuthProvider } from "@/context/auth-context";
import type { ReactNode } from "react";

export function GlobalProviders({ children }: { children: ReactNode }) {
  return (
    <AuthProvider>
      {children}
      <Toaster richColors position="top-right" closeButton />
    </AuthProvider>
  );
}
