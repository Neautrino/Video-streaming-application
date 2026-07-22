import { VideoUploader } from "./components/VideoUploader";

export default function App() {
  return (
    <main className="min-h-screen bg-slate-50 px-4 py-12 dark:bg-slate-950">
      <div className="mx-auto w-full max-w-xl">
        <header className="mb-8">
          <h1 className="text-2xl font-semibold tracking-tight text-slate-900 dark:text-white">
            Upload a video
          </h1>
          <p className="mt-1.5 text-sm text-slate-500 dark:text-slate-400">
            Files are uploaded straight to storage — they never pass through the
            API server.
          </p>
        </header>

        <VideoUploader />
      </div>
    </main>
  );
}
