import Link from "next/link";
import { SessionExpiredBanner } from "@/components/auth/session-expired-banner";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <main className="mx-auto flex min-h-full max-w-3xl flex-col justify-center gap-6 px-6 py-16">
      <SessionExpiredBanner />

      <div className="space-y-3">
        <h1 className="text-3xl font-semibold text-foreground">
          Distributed Exchange Benchmarking
        </h1>
        <p className="text-muted">
          Landing page placeholder. Sign in to access the workspace.
        </p>
      </div>
      <div>
        
        <Button variant="secondary" asChild>
          <Link href="/login">Sign in</Link>
        </Button>
      </div>
    </main>
  );
}
