import { useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { initUpload, confirmChunk, completeUpload, ApiError } from "../api/client";

const MAX_FILE_SIZE = 100 * 1024 * 1024; // 100MB
const CHUNK_SIZE = 10 * 1024 * 1024; // 10MB

export function UploadPage() {
  const navigate = useNavigate();
  const fileRef = useRef<HTMLInputElement>(null);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [secret, setSecret] = useState("");
  const [showSecret, setShowSecret] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState<{ current: number; total: number } | null>(null);

  function validateFile(file: File): string | null {
    if (file.size > MAX_FILE_SIZE) {
      return `File is too large (${(file.size / 1024 / 1024).toFixed(1)} MB). Maximum is 100 MB.`;
    }
    return null;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    const file = fileRef.current?.files?.[0];
    if (!file) { setError("Please select a video file."); return; }
    if (!title.trim()) { setError("Title is required."); return; }

    const fileError = validateFile(file);
    if (fileError) { setError(fileError); return; }

    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    setUploading(true);
    setProgress({ current: 0, total: totalChunks });

    try {
      // Init upload
      const initRes = await initUpload(
        { title: title.trim(), description: description.trim(), fileSize: file.size, totalChunks },
        secret
      );

      const { videoId, uploadId } = initRes;
      let { presignedUrl, nextPartNumber } = initRes;

      // Upload chunks
      for (let i = 0; i < totalChunks; i++) {
        const partNumber = nextPartNumber;
        const start = (partNumber - 1) * CHUNK_SIZE;
        const chunk = file.slice(start, start + CHUNK_SIZE);

        const putRes = await fetch(presignedUrl, { method: "PUT", body: chunk });
        if (!putRes.ok) throw new Error(`S3 upload failed for part ${partNumber}`);
        const eTag = putRes.headers.get("ETag") ?? "";

        setProgress({ current: partNumber, total: totalChunks });

        const confirmRes = await confirmChunk(videoId, { partNumber, eTag }, secret);
        if (!confirmRes.done) {
          presignedUrl = confirmRes.presignedUrl!;
          nextPartNumber = confirmRes.nextPartNumber!;
        }
      }

      // Complete
      await completeUpload(videoId, { uploadId }, secret);
      navigate("/");
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        setError("Invalid upload secret.");
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Upload failed.");
      }
    } finally {
      setUploading(false);
      setProgress(null);
    }
  }

  const inputStyle: React.CSSProperties = {
    padding: "10px 12px",
    borderRadius: 8,
    border: "1px solid #e5e5e5",
    fontSize: 14,
    color: "#0f0f0f",
    background: "#fff",
    outline: "none",
    width: "100%",
  };

  const labelStyle: React.CSSProperties = {
    display: "flex",
    flexDirection: "column",
    gap: 6,
    fontSize: 13,
    fontWeight: 500,
    color: "#606060",
  };

  return (
    <div style={{ maxWidth: 600, margin: "40px auto", padding: "0 24px 80px" }}>
      <h1 style={{ fontSize: 22, fontWeight: 600, color: "#0f0f0f", marginBottom: 32 }}>
        Upload video
      </h1>

      <div style={{ background: "#fff", borderRadius: 16, padding: "32px", border: "1px solid #e5e5e5" }}>
        <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          <label style={labelStyle}>
            Video file
            <input ref={fileRef} type="file" accept="video/*" disabled={uploading} style={{ fontSize: 13 }} />
          </label>

          <label style={labelStyle}>
            Title *
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              maxLength={200}
              placeholder="Add a title"
              disabled={uploading}
              style={inputStyle}
            />
          </label>

          <label style={labelStyle}>
            Description
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={4}
              maxLength={2000}
              placeholder="Tell viewers about your video"
              disabled={uploading}
              style={{ ...inputStyle, resize: "vertical", lineHeight: "22px" }}
            />
          </label>

          <label style={labelStyle}>
            Upload secret
            <div style={{ position: "relative" }}>
              <input
                type={showSecret ? "text" : "password"}
                value={secret}
                onChange={(e) => setSecret(e.target.value)}
                disabled={uploading}
                style={{ ...inputStyle, paddingRight: 40 }}
              />
              <button
                type="button"
                onClick={() => setShowSecret((v) => !v)}
                style={{
                  position: "absolute",
                  right: 10,
                  top: "50%",
                  transform: "translateY(-50%)",
                  background: "none",
                  border: "none",
                  cursor: "pointer",
                  padding: 4,
                  color: "#606060",
                  display: "flex",
                  alignItems: "center",
                }}
                aria-label={showSecret ? "Hide secret" : "Show secret"}
              >
                {showSecret ? (
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
                    <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
                    <line x1="1" y1="1" x2="23" y2="23" />
                  </svg>
                ) : (
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                )}
              </button>
            </div>
          </label>

          {error && (
            <p style={{ color: "#c00", fontSize: 13, margin: 0 }}>{error}</p>
          )}

          {progress && (
            <div>
              <p style={{ fontSize: 13, color: "#606060", marginBottom: 8 }}>
                Uploading part {progress.current} of {progress.total}…
              </p>
              <div style={{ background: "#f2f2f2", borderRadius: 4, height: 6, overflow: "hidden" }}>
                <div
                  style={{
                    background: "#ff0000",
                    height: "100%",
                    borderRadius: 4,
                    width: `${(progress.current / progress.total) * 100}%`,
                    transition: "width 0.3s ease",
                  }}
                />
              </div>
            </div>
          )}

          <button
            type="submit"
            disabled={uploading}
            style={{
              marginTop: 4,
              padding: "11px",
              background: uploading ? "#aaa" : "#ff0000",
              color: "#fff",
              border: "none",
              borderRadius: 8,
              fontSize: 14,
              fontWeight: 600,
              cursor: uploading ? "not-allowed" : "pointer",
              letterSpacing: "0.1px",
            }}
          >
            {uploading ? "Uploading…" : "Upload"}
          </button>
        </form>
      </div>
    </div>
  );
}
