export function getAccessTokenExpiryMs(accessToken: string): number | null {
  try {
    const parts = accessToken.split(".");
    if (parts.length !== 3) return null;

    const base64 = parts[1].replace(/-/g, "+").replace(/_/g, "/");
    const decoded = JSON.parse(atob(base64)) as { exp?: number };
    if (typeof decoded.exp !== "number") return null;

    return decoded.exp * 1000;
  } catch {
    return null;
  }
}

export function isAccessTokenExpired(accessToken: string): boolean {
  const expMs = getAccessTokenExpiryMs(accessToken);
  if (expMs === null) return false;
  return expMs <= Date.now();
}
