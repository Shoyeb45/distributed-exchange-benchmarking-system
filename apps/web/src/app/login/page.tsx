import { GithubLoginButton } from "@/components/auth/github-login-button";

export default function LoginPage() {
  return (
    <main className="mx-auto flex min-h-full max-w-md flex-col justify-center gap-6 px-6 py-16">
      <div className="space-y-2">
        <h1 className="text-2xl font-semibold text-foreground">Sign in</h1>
        <p className="text-muted">Use your GitHub account to continue.</p>
      </div>
      <GithubLoginButton />
    </main>
  );
}
