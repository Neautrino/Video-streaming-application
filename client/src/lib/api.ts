import { CREATE_VIDEO_ENDPOINT } from "../config/constants";

export interface CreateVideoRequest {
  title: string;
  description: string;
  filename: string;
  size: number;
  content_type: string;
}

export interface CreateVideoResponse {
  id: string;
  upload_url: string;
}

/**
 * Step 1 — register the video and get a presigned S3 URL back.
 * Only small JSON travels through our API; the bytes never do.
 */
export async function createVideo(
  payload: CreateVideoRequest,
): Promise<CreateVideoResponse> {
  const response = await fetch(CREATE_VIDEO_ENDPOINT, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const message = (await response.text()).trim();
    throw new Error(message || `Could not start the upload (${response.status})`);
  }

  return response.json();
}

/**
 * Step 2 — PUT the file straight to S3 using the presigned URL.
 *
 * Deliberately XMLHttpRequest rather than fetch: fetch cannot report upload
 * progress. No extra headers are set either — the presigned signature only
 * covers what the server signed, and adding headers can invalidate it.
 */
export function uploadToPresignedUrl(
  uploadUrl: string,
  file: File,
  onProgress: (percent: number) => void,
): Promise<void> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", uploadUrl, true);

    xhr.upload.onprogress = (event) => {
      if (event.lengthComputable) {
        onProgress(Math.round((event.loaded / event.total) * 100));
      }
    };

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        onProgress(100);
        resolve();
      } else {
        reject(new Error(`S3 rejected the upload (${xhr.status})`));
      }
    };

    xhr.onerror = () =>
      reject(new Error("Network error while uploading. Check your connection."));
    xhr.onabort = () => reject(new Error("Upload cancelled."));

    xhr.send(file);
  });
}
