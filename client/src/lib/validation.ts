import { ACCEPTED_MIME_PREFIX, MAX_VIDEO_BYTES } from "../config/constants";
import { formatBytes } from "./format";

/**
 * Client-side mirror of the server's validation rules.
 *
 * This is a convenience, not a security boundary — the server re-checks
 * everything. It just avoids uploading a gigabyte before finding out it
 * was rejected.
 *
 * @returns an error message, or null when the file is acceptable.
 */
export function validateVideoFile(file: File): string | null {
  if (!file.type.startsWith(ACCEPTED_MIME_PREFIX)) {
    return "That doesn't look like a video file. Try an MP4, MOV, or WebM.";
  }

  if (file.size <= 0) {
    return "That file is empty.";
  }

  if (file.size > MAX_VIDEO_BYTES) {
    return `That video is ${formatBytes(file.size)}. The limit is ${formatBytes(
      MAX_VIDEO_BYTES,
    )}.`;
  }

  return null;
}
