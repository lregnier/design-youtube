import type { components } from "./types.gen";

const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export type VideoSummary = components["schemas"]["VideoSummary"];
export type VideoDetail = components["schemas"]["VideoDetail"];
export type UploadInitRequest = components["schemas"]["UploadInitRequest"];
export type UploadInitResponse = components["schemas"]["UploadInitResponse"];
export type ConfirmChunkRequest = components["schemas"]["ConfirmChunkRequest"];
export type ConfirmChunkResponse = components["schemas"]["ConfirmChunkResponse"];
export type CompleteUploadRequest = components["schemas"]["CompleteUploadRequest"];

class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json", ...init?.headers },
    ...init,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new ApiError(res.status, body.error ?? res.statusText);
  }
  if (res.status === 200 && res.headers.get("content-length") === "0") return undefined as T;
  return res.json() as Promise<T>;
}

export function getVideos(): Promise<VideoSummary[]> {
  return request<VideoSummary[]>("/videos");
}

export function getVideo(videoId: string): Promise<VideoDetail> {
  return request<VideoDetail>(`/videos/${videoId}`);
}

export function initUpload(
  body: UploadInitRequest,
  secret: string
): Promise<UploadInitResponse> {
  return request<UploadInitResponse>("/videos/upload/init", {
    method: "POST",
    headers: { "X-Upload-Secret": secret },
    body: JSON.stringify(body),
  });
}

export function confirmChunk(
  videoId: string,
  body: ConfirmChunkRequest,
  secret: string
): Promise<ConfirmChunkResponse> {
  return request<ConfirmChunkResponse>(`/videos/${videoId}/upload/confirm-chunk`, {
    method: "POST",
    headers: { "X-Upload-Secret": secret },
    body: JSON.stringify(body),
  });
}

export function completeUpload(
  videoId: string,
  body: CompleteUploadRequest,
  secret: string
): Promise<void> {
  return request<void>(`/videos/${videoId}/upload/complete`, {
    method: "POST",
    headers: { "X-Upload-Secret": secret },
    body: JSON.stringify(body),
  });
}

export { ApiError };
