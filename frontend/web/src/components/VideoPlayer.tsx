import { useEffect, useRef } from "react";
import Hls from "hls.js";
import Plyr from "plyr";
import "plyr/dist/plyr.css";

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

    let hls: Hls | null = null;
    let plyr: Plyr | null = null;

    if (Hls.isSupported()) {
      hls = new Hls();
      hls.loadSource(manifestUrl);
      hls.attachMedia(video);

      hls.on(Hls.Events.MANIFEST_PARSED, (_, data) => {
        const qualities = [0, ...data.levels.map((l) => l.height)];

        plyr = new Plyr(video, {
          controls: ["play", "progress", "current-time", "mute", "volume", "settings", "fullscreen"],
          settings: ["quality"],
          i18n: { qualityLabel: { 0: "Auto" } },
          quality: {
            default: 0,
            options: qualities,
            forced: true,
            onChange: (quality: number) => {
              if (!hls) return;
              if (quality === 0) {
                hls.currentLevel = -1;
              } else {
                const idx = data.levels.findIndex((l) => l.height === quality);
                if (idx !== -1) hls!.currentLevel = idx;
              }
            },
          },
        });
      });
    } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = manifestUrl;
      plyr = new Plyr(video, {
        controls: ["play", "progress", "current-time", "mute", "volume", "fullscreen"],
      });
    }

    return () => {
      plyr?.destroy();
      hls?.destroy();
    };
  }, [manifestUrl]);

  return (
    <div style={{ width: "100%", background: "#000" }}>
      <video
        ref={videoRef}
        poster={thumbnailUrl}
        aria-label={title}
        style={{ width: "100%", display: "block" }}
      />
    </div>
  );
}
