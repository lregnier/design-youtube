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
        style={{
          width: "100%",
          aspectRatio: "16/9",
          objectFit: "cover",
          display: "block",
          borderRadius: 12,
        }}
      />
    ) : (
      <div
        style={{
          width: "100%",
          aspectRatio: "16/9",
          background: "#1a1a1a",
          borderRadius: 12,
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
            <div
              style={{
                width: 22,
                height: 22,
                border: "2.5px solid rgba(255,255,255,0.15)",
                borderTop: "2.5px solid rgba(255,255,255,0.7)",
                borderRadius: "50%",
                animation: "_vc-spin 0.8s linear infinite",
              }}
            />
          </>
        )}
        <span
          style={{
            fontSize: 13,
            fontWeight: 500,
            color: video.status === "failed" ? "#f44" : "rgba(255,255,255,0.55)",
          }}
        >
          {video.status === "failed" ? "Processing failed" : "Processing…"}
        </span>
      </div>
    );

  const info = (
    <div style={{ padding: "10px 2px 0" }}>
      <p
        style={{
          fontSize: 14,
          fontWeight: 500,
          lineHeight: "20px",
          color: "#0f0f0f",
          display: "-webkit-box",
          WebkitLineClamp: 2,
          WebkitBoxOrient: "vertical",
          overflow: "hidden",
        }}
      >
        {video.title}
      </p>
      <p style={{ fontSize: 13, color: "#606060", marginTop: 4 }}>{uploadedAt}</p>
    </div>
  );

  if (video.status === "ready") {
    return (
      <Link to={`/videos/${video.videoId}`} style={{ textDecoration: "none", display: "block" }}>
        {thumbnail}
        {info}
      </Link>
    );
  }

  return (
    <div style={{ cursor: "default" }}>
      {thumbnail}
      {info}
    </div>
  );
}
