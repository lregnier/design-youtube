import { useEffect, useRef, useState } from "react";
import Hls, { Level } from "hls.js";

interface Props {
  manifestUrl: string;
  thumbnailUrl?: string;
  title: string;
}

export function VideoPlayer({ manifestUrl, thumbnailUrl, title }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [levels, setLevels] = useState<Level[]>([]);
  const [selectedLevel, setSelectedLevel] = useState<number>(-1);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    if (Hls.isSupported()) {
      const hls = new Hls();
      hlsRef.current = hls;
      hls.loadSource(manifestUrl);
      hls.attachMedia(video);
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        setLevels(hls.levels);
      });
      return () => {
        hls.destroy();
        hlsRef.current = null;
        setLevels([]);
        setSelectedLevel(-1);
      };
    } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = manifestUrl;
    }
  }, [manifestUrl]);

  function selectLevel(index: number) {
    const hls = hlsRef.current;
    if (hls) hls.currentLevel = index;
    setSelectedLevel(index);
    setIsOpen(false);
  }

  const activeLabel = selectedLevel === -1 ? "Auto" : `${levels[selectedLevel]?.height}p`;

  return (
    <div style={{ position: "relative", width: "100%", background: "#000" }}>
      <video
        ref={videoRef}
        controls
        poster={thumbnailUrl}
        aria-label={title}
        style={{ width: "100%", maxHeight: "70vh", display: "block" }}
      />

      {levels.length > 0 && (
        <div style={{ position: "absolute", bottom: 48, right: 12, zIndex: 10 }}>
          {isOpen && (
            <ul style={{
              listStyle: "none",
              margin: 0,
              padding: "4px 0",
              background: "rgba(0,0,0,0.85)",
              borderRadius: 4,
              marginBottom: 4,
              minWidth: 80,
            }}>
              <li>
                <button
                  onClick={() => selectLevel(-1)}
                  style={itemStyle(selectedLevel === -1)}
                >
                  Auto
                </button>
              </li>
              {levels.map((level, i) => (
                <li key={i}>
                  <button
                    onClick={() => selectLevel(i)}
                    style={itemStyle(selectedLevel === i)}
                  >
                    {level.height}p
                  </button>
                </li>
              ))}
            </ul>
          )}
          <button
            onClick={() => setIsOpen((o) => !o)}
            style={{
              display: "block",
              width: "100%",
              padding: "4px 10px",
              background: "rgba(0,0,0,0.75)",
              color: "#fff",
              border: "1px solid rgba(255,255,255,0.3)",
              borderRadius: 4,
              cursor: "pointer",
              fontSize: 13,
              fontWeight: 600,
            }}
          >
            {activeLabel}
          </button>
        </div>
      )}
    </div>
  );
}

function itemStyle(active: boolean): React.CSSProperties {
  return {
    display: "block",
    width: "100%",
    padding: "6px 14px",
    background: "none",
    color: active ? "#fff" : "rgba(255,255,255,0.65)",
    fontWeight: active ? 700 : 400,
    border: "none",
    cursor: "pointer",
    fontSize: 13,
    textAlign: "left",
  };
}
