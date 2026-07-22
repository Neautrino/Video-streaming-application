/**
 * Application-wide constants.
 *
 * Anything that must stay in sync with the Go API lives here so there is a
 * single place to change it (notably MAX_VIDEO_BYTES, which mirrors the
 * server's `maxVideoBytes`).
 */

/** Base URL of the Go API server. Override with VITE_API_BASE_URL in .env. */
export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

/** Creates the video record and returns a presigned S3 upload URL. */
export const CREATE_VIDEO_ENDPOINT = `${API_BASE_URL}/videos/upload`;

/**
 * Largest video we accept, in bytes (1 GiB).
 * MUST match the server's `maxVideoBytes` — the server rejects anything
 * larger, so checking here just saves the user a pointless round trip.
 */
export const MAX_VIDEO_BYTES = 1 * 1024 * 1024 * 1024;

/** The server only accepts `video/*` content types. */
export const ACCEPTED_MIME_PREFIX = "video/";

/** `accept` attribute for the native file picker. */
export const FILE_INPUT_ACCEPT = "video/*";

/** Keeps the JSON metadata body small and the UI tidy. */
export const MAX_TITLE_LENGTH = 120;
export const MAX_DESCRIPTION_LENGTH = 500;
