import { useRef, useState, type ChangeEvent, type DragEvent } from "react";

import { FILE_INPUT_ACCEPT, MAX_VIDEO_BYTES } from "../config/constants";
import { formatBytes } from "../lib/format";

interface FileDropzoneProps {
  onSelect: (file: File) => void;
  disabled?: boolean;
}

export function FileDropzone({ onSelect, disabled = false }: FileDropzoneProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleDragOver = (event: DragEvent<HTMLButtonElement>) => {
    event.preventDefault();
    if (!disabled) setIsDragging(true);
  };

  const handleDragLeave = (event: DragEvent<HTMLButtonElement>) => {
    event.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (event: DragEvent<HTMLButtonElement>) => {
    event.preventDefault();
    setIsDragging(false);
    if (disabled) return;

    const file = event.dataTransfer.files?.[0];
    if (file) onSelect(file);
  };

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) onSelect(file);
    // Reset so picking the same file twice still fires onChange.
    event.target.value = "";
  };

  return (
    <>
      <button
        type="button"
        disabled={disabled}
        onClick={() => inputRef.current?.click()}
        onDragOver={handleDragOver}
        onDragEnter={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={[
          "flex w-full flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed px-6 py-12 text-center transition",
          "focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-500/40",
          disabled
            ? "cursor-not-allowed border-slate-200 opacity-60 dark:border-slate-800"
            : "cursor-pointer",
          isDragging
            ? "border-violet-500 bg-violet-50 dark:bg-violet-500/10"
            : "border-slate-300 bg-slate-50/60 hover:border-violet-400 hover:bg-violet-50/50 dark:border-slate-700 dark:bg-slate-950/40 dark:hover:border-violet-500 dark:hover:bg-violet-500/5",
        ].join(" ")}
      >
        <span className="flex h-12 w-12 items-center justify-center rounded-full bg-violet-100 text-violet-600 dark:bg-violet-500/15 dark:text-violet-400">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth={1.8}
            strokeLinecap="round"
            strokeLinejoin="round"
            className="h-6 w-6"
            aria-hidden="true"
          >
            <path d="M12 16V4" />
            <path d="m7 9 5-5 5 5" />
            <path d="M4 16v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2" />
          </svg>
        </span>

        <span className="space-y-1">
          <span className="block text-sm font-medium text-slate-900 dark:text-slate-100">
            Drop a video here, or{" "}
            <span className="text-violet-600 dark:text-violet-400">browse</span>
          </span>
          <span className="block text-xs text-slate-500 dark:text-slate-400">
            MP4, MOV, or WebM · up to {formatBytes(MAX_VIDEO_BYTES)}
          </span>
        </span>
      </button>

      <input
        ref={inputRef}
        type="file"
        accept={FILE_INPUT_ACCEPT}
        onChange={handleChange}
        className="sr-only"
        tabIndex={-1}
      />
    </>
  );
}
