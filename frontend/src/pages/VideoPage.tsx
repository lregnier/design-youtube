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
    <div style={{ maxWidth: 960, margin: "0 auto", padding: "24px 16px" }}>
      <Link to="/" style={{ fontSize: 13, color: "#555", textDecoration: "none" }}>
        ← Back
      </Link>

      {loading && <p>Loading…</p>}
      {error && <p style={{ color: "red" }}>Error: {error}</p>}

      {video && (
        <>
          {video.status === "ready" && video.manifestUrl ? (
            <VideoPlayer
              manifestUrl={video.manifestUrl}
              thumbnailUrl={video.thumbnailUrl ?? undefined}
              title={video.title}
            />
          ) : (
            <div style={{ background: "#111", aspectRatio: "16/9", display: "flex", alignItems: "center", justifyContent: "center" }}>
              <p style={{ color: "#aaa" }}>
                {video.status === "processing" ? "Processing…" : video.status === "failed" ? "Processing failed" : "Uploading…"}
              </p>
            </div>
          )}
          <h2 style={{ marginTop: 16 }}>{video.title}</h2>
          {video.description && <p style={{ color: "#444" }}>{video.description}</p>}
          <p style={{ fontSize: 12, color: "#888" }}>
            Uploaded {new Date(video.uploadedAt).toLocaleDateString()}
          </p>
        </>
      )}
    </div>
  );
}
