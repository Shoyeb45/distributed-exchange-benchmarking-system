"use client";

import {
  createContext,
  useCallback,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api-client";
import { apiBaseUrl } from "@/lib/config";
import type { MeResponse, User } from "@/types/user";

type AuthContextValue = {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (accessToken: string, refreshToken: string) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
  restoreSession: () => Promise<boolean>;
  loginWithGitHub: () => void;
};

export const AuthContext = createContext<AuthContextValue | null>(null);

function toUser(data: MeResponse): User {
  return {
    id: data.id,
    name: data.name,
    email: data.email,
    avatarUrl: data.avatarUrl,
    githubUsername: data.githubUsername,
  };
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [hasToken, setHasToken] = useState(false);

  const syncTokenState = useCallback(() => {
    setHasToken(api.hasSession());
  }, []);

  const refreshUser = useCallback(async () => {
    if (!api.getAccessToken()) {
      setUser(null);
      return;
    }

    try {
      const profile = await api.get<MeResponse>("/auth/me");
      setUser(toUser(profile));
    } catch {
      setUser(null);
    }
  }, []);

  const restoreSession = useCallback(async () => {
    if (!api.hasSession()) return false;

    const valid = await api.ensureValidAccessToken();
    if (!valid) return false;

    await refreshUser();
    return api.hasSession() && Boolean(api.getAccessToken());
  }, [refreshUser]);

  const login = useCallback(
    async (accessToken: string, refreshToken: string) => {
      api.setTokens({ accessToken, refreshToken });
      syncTokenState();
      await refreshUser();
    },
    [refreshUser, syncTokenState],
  );

  const logout = useCallback(() => {
    api.clearAuth();
    setUser(null);
    syncTokenState();
    router.replace("/login");
  }, [router, syncTokenState]);

  const loginWithGitHub = useCallback(() => {
    window.location.href = `${apiBaseUrl}/auth/github`;
  }, []);

  useEffect(() => {
    const unsubscribe = api.onAuthChange(syncTokenState);
    return unsubscribe;
  }, [syncTokenState]);

  useEffect(() => {
    let active = true;

    const boot = async () => {
      syncTokenState();

      if (!api.hasSession()) {
        if (active) setIsLoading(false);
        return;
      }

      await restoreSession();
      if (active) setIsLoading(false);
    };

    void boot();

    return () => {
      active = false;
    };
  }, [restoreSession, syncTokenState]);

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: hasToken,
        isLoading,
        login,
        logout,
        refreshUser,
        restoreSession,
        loginWithGitHub,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
