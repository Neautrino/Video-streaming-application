import { useState, type FormEvent, type ReactNode } from "react";

import {
  MAX_DESCRIPTION_LENGTH,
  MAX_TITLE_LENGTH,
} from "../config/constants";
import { useVideoUpload } from "../hooks/useVideoUpload";
import { formatBytes } from "../lib/format";
import { validateVideoFile } from "../lib/validation";
import { FileDropzone } from "./FileDropzone";
import { ProgressBar } from "./ProgressBar";

function stripExtension(filename: string): string {
  const dot = filename.lastIndexOf(".");
  return dot > 0 ? filename.slice(0, dot) : filename;
}

export function VideoUploader() {
  const [file, setFile] = useState<File | null>(null);
  const [fileError, setFileError] = useState<string | null>(null);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  const { upload, reset, progress, isUploading, isSuccess, error, videoId } =
    useVideoUpload();

  const handleSelect = (selected: File) => {
    const message = validateVideoFile(selected);

    if (message) {
      setFile(null);
      setFileError(message);
      return;
    }

    setFileError(null);
    setFile(selected);

    // Small nicety: seed the title from the filename if it's still empty.
    if (!title.trim()) setTitle(stripExtension(selected.name));
  };

  const startOver = () => {
    reset();
    setFile(null);
    setFileError(null);
    setTitle("");
    setDescription("");
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    if (!file || !title.trim() || isUploading) return;
    upload({ file, title, description });
  };

  const canSubmit = Boolean(file) && title.trim().length > 0 && !isUploading;

  if (isSuccess && videoId) {
    return <SuccessCard videoId={videoId} onUploadAnother={startOver} />;
  }

  return (
    <form onSubmit={handleSubmit} className={cardClass}>
      {file ? (
        <SelectedFile
          file={file}
          disabled={isUploading}
          onClear={() => setFile(null)}
        />
      ) : (
        <FileDropzone onSelect={handleSelect} disabled={isUploading} />
      )}

      {fileError && <Alert tone="error">{fileError}</Alert>}

      <div className="mt-6 space-y-4">
        <Field label="Title" required htmlFor="title">
          <input
            id="title"
            type="text"
            value={title}
            maxLength={MAX_TITLE_LENGTH}
            disabled={isUploading}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Give your video a name"
            className={inputClass}
          />
        </Field>

        <Field label="Description" htmlFor="description">
          <textarea
            id="description"
            rows={3}
            value={description}
            maxLength={MAX_DESCRIPTION_LENGTH}
            disabled={isUploading}
            onChange={(event) => setDescription(event.target.value)}
            placeholder="What's this video about? (optional)"
            className={`${inputClass} resize-none`}
          />
          <span className="mt-1 block text-right text-xs text-slate-400">
            {description.length}/{MAX_DESCRIPTION_LENGTH}
          </span>
        </Field>
      </div>

      {error && <Alert tone="error">{error.message}</Alert>}

      {isUploading && (
        <div className="mt-6 space-y-2">
          <div className="flex items-center justify-between text-xs font-medium">
            <span className="text-slate-600 dark:text-slate-300">
              Uploading to storage…
            </span>
            <span className="tabular-nums text-slate-500 dark:text-slate-400">
              {progress}%
            </span>
          </div>
          <ProgressBar value={progress} />
        </div>
      )}

      <button
        type="submit"
        disabled={!canSubmit}
        className="mt-6 flex w-full items-center justify-center gap-2 rounded-lg bg-violet-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-violet-500 focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-500/50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-violet-600 dark:hover:bg-violet-500"
      >
        {isUploading ? (
          <>
            <Spinner />
            Uploading…
          </>
        ) : (
          "Upload video"
        )}
      </button>
    </form>
  );
}

/* ---------- pieces ---------- */

const cardClass =
  "rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900";

const inputClass =
  "w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 disabled:opacity-60 dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100 dark:placeholder:text-slate-500";

function Field({
  label,
  htmlFor,
  required = false,
  children,
}: {
  label: string;
  htmlFor: string;
  required?: boolean;
  children: ReactNode;
}) {
  return (
    <div>
      <label
        htmlFor={htmlFor}
        className="mb-1.5 block text-sm font-medium text-slate-700 dark:text-slate-300"
      >
        {label}
        {required && <span className="ml-0.5 text-violet-600">*</span>}
      </label>
      {children}
    </div>
  );
}

function SelectedFile({
  file,
  disabled,
  onClear,
}: {
  file: File;
  disabled: boolean;
  onClear: () => void;
}) {
  return (
    <div className="flex items-center gap-3 rounded-xl border border-slate-200 bg-slate-50 p-4 dark:border-slate-800 dark:bg-slate-950/50">
      <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-violet-100 text-violet-600 dark:bg-violet-500/15 dark:text-violet-400">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth={1.8}
          strokeLinecap="round"
          strokeLinejoin="round"
          className="h-5 w-5"
          aria-hidden="true"
        >
          <path d="m10 8 6 4-6 4V8Z" />
          <rect x="2" y="4" width="20" height="16" rx="3" />
        </svg>
      </span>

      <div className="min-w-0 flex-1">
        <p className="truncate text-sm font-medium text-slate-900 dark:text-slate-100">
          {file.name}
        </p>
        <p className="text-xs text-slate-500 dark:text-slate-400">
          {formatBytes(file.size)} · {file.type || "unknown type"}
        </p>
      </div>

      {!disabled && (
        <button
          type="button"
          onClick={onClear}
          aria-label="Remove file"
          className="shrink-0 rounded-md p-1.5 text-slate-400 transition hover:bg-slate-200 hover:text-slate-700 dark:hover:bg-slate-800 dark:hover:text-slate-200"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth={2}
            strokeLinecap="round"
            className="h-4 w-4"
            aria-hidden="true"
          >
            <path d="M18 6 6 18M6 6l12 12" />
          </svg>
        </button>
      )}
    </div>
  );
}

function SuccessCard({
  videoId,
  onUploadAnother,
}: {
  videoId: string;
  onUploadAnother: () => void;
}) {
  return (
    <div className={`${cardClass} text-center`}>
      <span className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth={2}
          strokeLinecap="round"
          strokeLinejoin="round"
          className="h-6 w-6"
          aria-hidden="true"
        >
          <path d="M20 6 9 17l-5-5" />
        </svg>
      </span>

      <h2 className="mt-4 text-lg font-semibold text-slate-900 dark:text-white">
        Upload complete
      </h2>
      <p className="mt-1 text-sm text-slate-500 dark:text-slate-400">
        Your video is queued for processing.
      </p>

      <code className="mt-4 inline-block max-w-full truncate rounded-lg bg-slate-100 px-3 py-1.5 font-mono text-xs text-slate-600 dark:bg-slate-950 dark:text-slate-300">
        {videoId}
      </code>

      <button
        type="button"
        onClick={onUploadAnother}
        className="mt-6 w-full rounded-lg border border-slate-300 px-4 py-2.5 text-sm font-medium text-slate-700 transition hover:bg-slate-50 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800"
      >
        Upload another
      </button>
    </div>
  );
}

function Alert({
  tone,
  children,
}: {
  tone: "error";
  children: ReactNode;
}) {
  return (
    <p
      role="alert"
      className={
        tone === "error"
          ? "mt-4 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-300"
          : ""
      }
    >
      {children}
    </p>
  );
}

function Spinner() {
  return (
    <svg
      className="h-4 w-4 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <circle
        className="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="4"
      />
      <path
        className="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 0 1 8-8v4a4 4 0 0 0-4 4H4Z"
      />
    </svg>
  );
}
