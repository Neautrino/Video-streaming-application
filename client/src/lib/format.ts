const UNITS = ["B", "KB", "MB", "GB", "TB"] as const;

/** Formats a byte count into a short human-readable string, e.g. "34.5 MB". */
export function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return "0 B";

  const exponent = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    UNITS.length - 1,
  );
  const value = bytes / 1024 ** exponent;

  return `${value.toFixed(exponent === 0 ? 0 : 1)} ${UNITS[exponent]}`;
}
