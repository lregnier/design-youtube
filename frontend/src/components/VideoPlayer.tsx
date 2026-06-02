import { useEffect, useRef } from "react";
import Hls from "hls.js";

interface Props {
  manifestUrl: string;
  thumbnailUrl?: string;
  title: string;
}

export function VideoPlayer({ manifestUrl, thumbnailUrl, title }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    if (Hls.isSupported()) {
      const hls = new Hls();
      hls.loadSource(manifestUrl);
      hls.attachMedia(video);
      return () => hls.destroy();
    } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = manifestUrl;
    }
  }, [manifestUrl]);

  return (
    <video
      ref={videoRef}
      controls
      poster={thumbnailUrl}
      aria-label={title}
      style={{ width: "100%", maxHeight: "70vh", background: "#000", display: "block" }}
    />
  );
}
