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

  return (
    <div style={{ maxWidth: 560, margin: "48px auto", padding: "0 16px" }}>
      <h1 style={{ fontSize: 22, marginBottom: 24 }}>Upload Video</h1>

      <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: 16 }}>
        <label style={{ display: "flex", flexDirection: "column", gap: 4, fontSize: 14 }}>
          Video file
          <input ref={fileRef} type="file" accept="video/*" disabled={uploading} />
        </label>

        <label style={{ display: "flex", flexDirection: "column", gap: 4, fontSize: 14 }}>
          Title *
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            maxLength={200}
            disabled={uploading}
            style={{ padding: "6px 8px", borderRadius: 4, border: "1px solid #ccc" }}
          />
        </label>

        <label style={{ display: "flex", flexDirection: "column", gap: 4, fontSize: 14 }}>
          Description
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={3}
            maxLength={2000}
            disabled={uploading}
            style={{ padding: "6px 8px", borderRadius: 4, border: "1px solid #ccc", resize: "vertical" }}
          />
        </label>

        <label style={{ display: "flex", flexDirection: "column", gap: 4, fontSize: 14 }}>
          Upload secret
          <input
            type="password"
            value={secret}
            onChange={(e) => setSecret(e.target.value)}
            disabled={uploading}
            style={{ padding: "6px 8px", borderRadius: 4, border: "1px solid #ccc" }}
          />
        </label>

        {error && <p style={{ color: "red", margin: 0, fontSize: 13 }}>{error}</p>}

        {progress && (
          <div>
            <p style={{ margin: "0 0 4px", fontSize: 13 }}>
              Uploading chunk {progress.current} of {progress.total}…
            </p>
            <div style={{ background: "#eee", borderRadius: 4, height: 8 }}>
              <div
                style={{
                  background: "#e00",
                  height: "100%",
                  borderRadius: 4,
                  width: `${(progress.current / progress.total) * 100}%`,
                  transition: "width 0.2s",
                }}
              />
            </div>
          </div>
        )}

        <button
          type="submit"
          disabled={uploading}
          style={{
            padding: "10px",
            background: uploading ? "#ccc" : "#e00",
            color: "#fff",
            border: "none",
            borderRadius: 6,
            fontSize: 15,
            fontWeight: 600,
            cursor: uploading ? "not-allowed" : "pointer",
          }}
        >
          {uploading ? "Uploading…" : "Upload"}
        </button>
      </form>
    </div>
  );
}
