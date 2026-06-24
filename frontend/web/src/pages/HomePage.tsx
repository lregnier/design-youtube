import { useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";
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
    <div style={{ maxWidth: 1200, margin: "0 auto", padding: "24px 16px" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 24 }}>Videos</h1>
        <Link to="/upload" style={{ padding: "8px 16px", background: "#e00", color: "#fff", borderRadius: 6, textDecoration: "none", fontSize: 14, fontWeight: 600 }}>
          Upload
        </Link>
      </div>

      {loading && <p>Loading…</p>}
      {error && <p style={{ color: "red" }}>Error: {error}</p>}

      {!loading && !error && videos.length === 0 && (
        <p style={{ color: "#888", textAlign: "center", marginTop: 60 }}>No videos yet</p>
      )}

      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: 20 }}>
        {videos.map((v) => (
          <VideoCard key={v.videoId} video={v} />
        ))}
      </div>
    </div>
  );
}
