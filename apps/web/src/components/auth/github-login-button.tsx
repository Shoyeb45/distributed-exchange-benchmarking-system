"use client";

import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";

export function GithubLoginButton() {
  const { loginWithGitHub } = useAuth();

  return (
    <Button variant="primary" onClick={loginWithGitHub}>
      Continue with GitHub
    </Button>
  );
}
