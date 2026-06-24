import { useEffect, useRef, useState } from "react";
import { getVideos, type VideoSummary } from "../api/client";
import { VideoCard } from "../components/VideoCard";

export function HomePage() {
  const [videos, setVideos] = useState<VideoSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  function loadVideos() {
    getVideos()
      .then((data) => {
        setVideos(data);
        const anyProcessing = data.some((v) => v.status === "processing");
        if (anyProcessing && !intervalRef.current) {
          intervalRef.current = setInterval(loadVideos, 5000);
        } else if (!anyProcessing && intervalRef.current) {
          clearInterval(intervalRef.current);
          intervalRef.current = null;
        }
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    loadVideos();
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, []);

  return (
    <div style={{ maxWidth: 1280, margin: "0 auto", padding: "24px 24px 60px" }}>
      {loading && (
        <p style={{ color: "#606060", fontSize: 14, padding: "40px 0" }}>Loading…</p>
      )}
      {error && (
        <p style={{ color: "#c00", fontSize: 14, padding: "40px 0" }}>Error: {error}</p>
      )}

      {!loading && !error && videos.length === 0 && (
        <div style={{ textAlign: "center", padding: "100px 0", color: "#606060" }}>
          <svg width="64" height="64" viewBox="0 0 64 64" fill="none" aria-hidden="true" style={{ marginBottom: 20, opacity: 0.3 }}>
            <rect width="64" height="64" rx="12" fill="#0f0f0f" />
            <polygon points="26,18 50,32 26,46" fill="#fff" />
          </svg>
          <p style={{ fontSize: 18, fontWeight: 500, color: "#0f0f0f", marginBottom: 8 }}>No videos yet</p>
          <p style={{ fontSize: 14 }}>Upload your first video to get started</p>
        </div>
      )}

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fill, minmax(300px, 1fr))",
          gap: "40px 16px",
        }}
      >
        {videos.map((v) => (
          <VideoCard key={v.videoId} video={v} />
        ))}
      </div>
    </div>
  );
}
