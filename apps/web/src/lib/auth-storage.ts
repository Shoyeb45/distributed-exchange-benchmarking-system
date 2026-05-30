const ACCESS_KEY = "access_token";
const REFRESH_KEY = "refresh_token";

function read(key: string): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(key);
}

function write(key: string, value: string): void {
  localStorage.setItem(key, value);
}

export const authStorage = {
  getAccessToken(): string | null {
    return read(ACCESS_KEY);
  },

  getRefreshToken(): string | null {
    return read(REFRESH_KEY);
  },

  setTokens(accessToken: string, refreshToken: string): void {
    write(ACCESS_KEY, accessToken);
    write(REFRESH_KEY, refreshToken);
  },

  clear(): void {
    localStorage.removeItem(ACCESS_KEY);
    localStorage.removeItem(REFRESH_KEY);
  },

  hasSession(): boolean {
    return Boolean(read(ACCESS_KEY) && read(REFRESH_KEY));
  },
};
