import axios, {
  type AxiosError,
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from "axios";
import { toast } from "sonner";
import { apiBaseUrl } from "@/lib/config";
import { authStorage } from "@/lib/auth-storage";
import { getAccessTokenExpiryMs } from "@/lib/jwt";

export type Tokens = {
  accessToken: string;
  refreshToken: string;
};

type TokenResponse = {
  success: boolean;
  accessToken: string;
  refreshToken: string;
};

type ApiErrorBody = {
  message?: string;
};

type RetryConfig = InternalAxiosRequestConfig & { _retry?: boolean };

const REFRESH_BEFORE_MS = 5 * 1000;

function getErrorMessage(error: AxiosError<ApiErrorBody>): string {
  const msg = error.response?.data?.message ?? error.message;
  if (typeof msg === "string" && msg.trim()) return msg;
  if (error.response?.status === 401) {
    return "Session expired. Please sign in again.";
  }
  return "Something went wrong. Please try again.";
}

class ApiClient {
  private client: AxiosInstance;
  private refreshPromise: Promise<string> | null = null;
  private refreshTimerId: ReturnType<typeof setTimeout> | null = null;
  private authListeners = new Set<() => void>();

  constructor(baseURL: string) {
    this.client = axios.create({
      baseURL,
      headers: { "Content-Type": "application/json" },
    });
    this.attachInterceptors();
  }

  onAuthChange(listener: () => void): () => void {
    this.authListeners.add(listener);
    return () => this.authListeners.delete(listener);
  }

  getAccessToken(): string | null {
    return authStorage.getAccessToken();
  }

  getRefreshToken(): string | null {
    return authStorage.getRefreshToken();
  }

  hasSession(): boolean {
    return authStorage.hasSession();
  }

  setTokens(tokens: Tokens): void {
    authStorage.setTokens(tokens.accessToken, tokens.refreshToken);
    this.notifyAuthChange();
    this.startTokenRefreshTimer();
  }

  clearAuth(): void {
    this.stopTokenRefreshTimer();
    authStorage.clear();
    this.notifyAuthChange();
  }

  startTokenRefreshTimer(): void {
    this.stopTokenRefreshTimer();

    const accessToken = this.getAccessToken();
    if (!accessToken) return;

    const expMs = getAccessTokenExpiryMs(accessToken);
    if (expMs === null) return;

    const refreshAt = expMs - REFRESH_BEFORE_MS;
    const delay = refreshAt - Date.now();

    if (delay <= 0) {
      void this.refreshAccessToken({ redirectOnFailure: true }).catch(() => {});
      return;
    }

    this.refreshTimerId = setTimeout(() => {
      this.refreshTimerId = null;
      void this.refreshAccessToken({ redirectOnFailure: true }).catch(() => {});
    }, delay);
  }

  stopTokenRefreshTimer(): void {
    if (this.refreshTimerId !== null) {
      clearTimeout(this.refreshTimerId);
      this.refreshTimerId = null;
    }
  }

  async refreshAccessToken(options?: {
    redirectOnFailure?: boolean;
  }): Promise<string> {
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    this.refreshPromise = this.performRefresh(options);
    try {
      return await this.refreshPromise;
    } finally {
      this.refreshPromise = null;
    }
  }

  private async performRefresh(options?: {
    redirectOnFailure?: boolean;
  }): Promise<string> {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) {
      throw new Error("missing refresh token");
    }

    try {
      const { data } = await axios.post<TokenResponse>(
        `${apiBaseUrl}/auth/refresh`,
        { refreshToken },
        { headers: { "Content-Type": "application/json" } },
      );

      this.setTokens({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
      });

      return data.accessToken;
    } catch (error) {
      this.clearAuth();

      const msg = axios.isAxiosError(error)
        ? getErrorMessage(error)
        : "Session expired. Please sign in again.";

      toast.error(msg);

      if (options?.redirectOnFailure && typeof window !== "undefined") {
        window.location.href = "/?reason=session_expired";
      }

      throw error;
    }
  }

  async ensureValidAccessToken(): Promise<boolean> {
    const accessToken = this.getAccessToken();
    const refreshToken = this.getRefreshToken();

    if (!accessToken || !refreshToken) return false;

    this.startTokenRefreshTimer();

    const expMs = getAccessTokenExpiryMs(accessToken);
    const needsRefresh = expMs !== null && expMs - REFRESH_BEFORE_MS <= Date.now();

    if (!needsRefresh) return true;

    try {
      await this.refreshAccessToken({ redirectOnFailure: true });
      return true;
    } catch {
      return false;
    }
  }

  private notifyAuthChange(): void {
    this.authListeners.forEach((listener) => listener());
  }

  private attachInterceptors(): void {
    this.client.interceptors.request.use((config) => {
      const token = this.getAccessToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const config = error.config as RetryConfig | undefined;
        if (!config || config._retry || error.response?.status !== 401) {
          return Promise.reject(error);
        }

        if (this.isRefreshRequest(config.url)) {
          return Promise.reject(error);
        }

        config._retry = true;

        try {
          const newToken = await this.refreshAccessToken({
            redirectOnFailure: true,
          });
          config.headers.Authorization = `Bearer ${newToken}`;
          return this.client(config);
        } catch {
          return Promise.reject(error);
        }
      },
    );
  }

  private isRefreshRequest(url?: string): boolean {
    return Boolean(url?.includes("/auth/refresh"));
  }

  get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "GET", url });
  }

  post<T>(url: string, body?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "POST", url, data: body });
  }

  put<T>(url: string, body?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "PUT", url, data: body });
  }

  patch<T>(
    url: string,
    body?: unknown,
    config?: AxiosRequestConfig,
  ): Promise<T> {
    return this.request<T>({ ...config, method: "PATCH", url, data: body });
  }

  delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "DELETE", url });
  }

  private async request<T>(config: AxiosRequestConfig): Promise<T> {
    const response = await this.client.request<T>(config);
    return response.data;
  }
}

export const api = new ApiClient(apiBaseUrl);
