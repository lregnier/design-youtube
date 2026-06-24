import { Link } from "react-router-dom";

export function Navbar() {
  return (
    <header style={{
      position: "sticky",
      top: 0,
      zIndex: 100,
      background: "#fff",
      borderBottom: "1px solid #e5e5e5",
      height: 56,
      display: "flex",
      alignItems: "center",
      padding: "0 24px",
      justifyContent: "space-between",
    }}>
      <Link to="/" style={{ textDecoration: "none", display: "flex", alignItems: "center", gap: 8 }}>
        <svg width="28" height="20" viewBox="0 0 28 20" fill="none" aria-hidden="true">
          <rect width="28" height="20" rx="4" fill="#0d9488" />
          <polygon points="11,5 21,10 11,15" fill="#fff" />
        </svg>
        <span style={{ fontSize: 18, fontWeight: 700, letterSpacing: "-0.3px" }}>
          <span style={{ color: "#0d9488" }}>You</span><span style={{ color: "#0f0f0f" }}>Flick</span>
        </span>
      </Link>

      <Link
        to="/upload"
        style={{
          display: "flex",
          alignItems: "center",
          gap: 6,
          padding: "7px 16px",
          border: "1px solid #e5e5e5",
          borderRadius: 20,
          textDecoration: "none",
          color: "#0f0f0f",
          fontSize: 14,
          fontWeight: 500,
          background: "#fff",
        }}
      >
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
          <path d="M8 2v12M2 8h12" stroke="#0f0f0f" strokeWidth="1.8" strokeLinecap="round" />
        </svg>
        Upload
      </Link>
    </header>
  );
}
