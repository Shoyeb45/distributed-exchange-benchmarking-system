import type { Metadata } from "next";
import "./globals.css";
import { GlobalProviders } from "@/components/providers/global-providers";
import { SessionRedirect } from "@/components/auth/session-redirect";

export const metadata: Metadata = {
  title: "Distributed Benchmarking Platform",
  description: "",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full antialiased">
      <body className="flex min-h-full flex-col">
        <GlobalProviders>
          <SessionRedirect />
          {children}
        </GlobalProviders>
      </body>
    </html>
  );
}
