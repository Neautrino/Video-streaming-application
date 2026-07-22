import { useMutation } from "@tanstack/react-query";
import { useCallback, useState } from "react";

import { createVideo, uploadToPresignedUrl } from "../lib/api";

export interface UploadInput {
  file: File;
  title: string;
  description: string;
}

/**
 * Drives the two-step upload: create the record, then PUT the bytes to S3.
 *
 * The mutation owns the request lifecycle; progress lives in local state
 * because React Query has no concept of upload progress.
 */
export function useVideoUpload() {
  const [progress, setProgress] = useState(0);

  const mutation = useMutation<string, Error, UploadInput>({
    mutationFn: async ({ file, title, description }) => {
      setProgress(0);

      const { id, upload_url } = await createVideo({
        title: title.trim(),
        description: description.trim(),
        filename: file.name,
        size: file.size,
        content_type: file.type,
      });

      await uploadToPresignedUrl(upload_url, file, setProgress);

      return id;
    },
  });

  const reset = useCallback(() => {
    setProgress(0);
    mutation.reset();
  }, [mutation]);

  return {
    upload: mutation.mutate,
    reset,
    progress,
    isUploading: mutation.isPending,
    isSuccess: mutation.isSuccess,
    error: mutation.error,
    videoId: mutation.data ?? null,
  };
}
