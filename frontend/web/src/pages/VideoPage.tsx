import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { getVideo, type VideoDetail } from "../api/client";
import { VideoPlayer } from "../components/VideoPlayer";

export function VideoPage() {
  const { videoId } = useParams<{ videoId: string }>();
  const [video, setVideo] = useState<VideoDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!videoId) return;
    getVideo(videoId)
      .then(setVideo)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [videoId]);

  return (
    <div style={{ maxWidth: 1000, margin: "0 auto", padding: "20px 24px 60px" }}>
      <Link
        to="/"
        style={{
          display: "inline-flex",
          alignItems: "center",
          gap: 6,
          fontSize: 13,
          color: "#606060",
          textDecoration: "none",
          marginBottom: 20,
        }}
      >
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
          <path d="M10 3L5 8l5 5" stroke="#606060" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
        Back
      </Link>

      {loading && <p style={{ color: "#606060", fontSize: 14 }}>Loading…</p>}
      {error && <p style={{ color: "#c00", fontSize: 14 }}>Error: {error}</p>}

      {video && (
        <>
          {video.status === "ready" && video.manifestUrl ? (
            <VideoPlayer
              manifestUrl={video.manifestUrl}
              thumbnailUrl={video.thumbnailUrl ?? undefined}
              title={video.title}
            />
          ) : (
            <div
              style={{
                aspectRatio: "16/9",
                background: "#1a1a1a",
                borderRadius: 12,
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                gap: 12,
              }}
            >
              {video.status === "processing" && (
                <>
                  <style>{`@keyframes _vp-spin { to { transform: rotate(360deg); } }`}</style>
                  <div
                    style={{
                      width: 32,
                      height: 32,
                      border: "3px solid rgba(255,255,255,0.12)",
                      borderTop: "3px solid rgba(255,255,255,0.6)",
                      borderRadius: "50%",
                      animation: "_vp-spin 0.9s linear infinite",
                    }}
                  />
                </>
              )}
              <p style={{ color: "rgba(255,255,255,0.5)", fontSize: 14, fontWeight: 500 }}>
                {video.status === "processing"
                  ? "Processing…"
                  : video.status === "failed"
                  ? "Processing failed"
                  : "Uploading…"}
              </p>
            </div>
          )}

          <div style={{ marginTop: 20 }}>
            <h1 style={{ fontSize: 20, fontWeight: 600, lineHeight: "28px", color: "#0f0f0f" }}>
              {video.title}
            </h1>
            <p style={{ fontSize: 14, color: "#606060", marginTop: 8 }}>
              {new Date(video.uploadedAt).toLocaleDateString(undefined, {
                year: "numeric",
                month: "long",
                day: "numeric",
              })}
            </p>

            {video.description && (
              <div
                style={{
                  marginTop: 16,
                  padding: "14px 16px",
                  background: "#f2f2f2",
                  borderRadius: 12,
                  fontSize: 14,
                  lineHeight: "22px",
                  color: "#0f0f0f",
                  whiteSpace: "pre-wrap",
                }}
              >
                {video.description}
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
