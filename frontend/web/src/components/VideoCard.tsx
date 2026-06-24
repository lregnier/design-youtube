import { Link } from "react-router-dom";
import type { VideoSummary } from "../api/client";

interface Props {
  video: VideoSummary;
}

export function VideoCard({ video }: Props) {
  const uploadedAt = new Date(video.uploadedAt).toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });

  const thumbnail =
    video.status === "ready" ? (
      <img
        src={video.thumbnailUrl}
        alt={video.title}
        style={{ width: "100%", aspectRatio: "16/9", objectFit: "cover", display: "block" }}
      />
    ) : (
      <div
        style={{
          width: "100%",
          aspectRatio: "16/9",
          background: "#111",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          gap: 10,
        }}
      >
        {video.status === "processing" && (
          <>
            <style>{`@keyframes _vc-spin { to { transform: rotate(360deg); } }`}</style>
            <div style={{
              width: 20,
              height: 20,
              border: "2px solid rgba(170,170,170,0.25)",
              borderTop: "2px solid #aaa",
              borderRadius: "50%",
              animation: "_vc-spin 0.8s linear infinite",
            }} />
          </>
        )}
        <span style={{ fontSize: 13, color: video.status === "failed" ? "#c00" : "#aaa" }}>
          {video.status === "failed" ? "Failed" : "Processing…"}
        </span>
      </div>
    );

  const cardContent = (
    <div style={{ border: "1px solid #ddd", borderRadius: 8, overflow: "hidden" }}>
      {thumbnail}
      <div style={{ padding: "8px 12px" }}>
        <p style={{ margin: 0, fontWeight: 600, fontSize: 14 }}>{video.title}</p>
        <p style={{ margin: "4px 0 0", fontSize: 12, color: "#666" }}>{uploadedAt}</p>
      </div>
    </div>
  );

  if (video.status === "ready") {
    return (
      <Link to={`/videos/${video.videoId}`} style={{ textDecoration: "none", color: "inherit" }}>
        {cardContent}
      </Link>
    );
  }

  return <div style={{ cursor: "default" }}>{cardContent}</div>;
}
