export type User = {
  id: number;
  name: string;
  email: string;
  avatarUrl: string;
  githubUsername: string;
};

export type MeResponse = User & {
  success: boolean;
};
