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

  return (
    <Link to={`/videos/${video.videoId}`} style={{ textDecoration: "none", color: "inherit" }}>
      <div style={{ border: "1px solid #ddd", borderRadius: 8, overflow: "hidden", cursor: "pointer" }}>
        <img
          src={video.thumbnailUrl}
          alt={video.title}
          style={{ width: "100%", aspectRatio: "16/9", objectFit: "cover", display: "block" }}
        />
        <div style={{ padding: "8px 12px" }}>
          <p style={{ margin: 0, fontWeight: 600, fontSize: 14 }}>{video.title}</p>
          <p style={{ margin: "4px 0 0", fontSize: 12, color: "#666" }}>{uploadedAt}</p>
        </div>
      </div>
    </Link>
  );
}
