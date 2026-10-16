// Màn soạn kịch bản (design v2, artboard Main): thanh tiêu đề, thanh công cụ, thư viện | danh sách bước | inspector, status bar
import { useEffect, useState } from 'react';
import { DebugButtons, RecordButton, StatusBar, Toolbar } from '../components/AppChrome';
import { MousePos, StopScript } from '../../wailsjs/go/main/App';
import { inWails, openOverlay } from '../windowMode';
import { removeStep, runCurrent, useEditor } from '../scriptStore';
import TitleBar from '../components/TitleBar';
import Library from '../components/Library';
import StepList from '../components/StepList';
import StepForm from '../components/StepForm';
import RunSettings, { loopSummary } from '../components/RunSettings';
import { STOP_REASON, useEngine } from '../engineStore';
import { Point } from '../types';

const MOUSE_POLL_MS = 250;

export default function MainScreen() {
  const { script, selected, notice } = useEditor();
  const { steps } = script;
  const cur = steps.find((s) => s.id === selected) ?? steps[0]; // undefined khi chưa có bước nào
  const idx = steps.indexOf(cur) + 1;
  const { running, lastError, stopReason, hotkeyError } = useEngine();

  // Toạ độ chuột thật ở status bar – để đọc nhanh toạ độ khi nhập tay.
  const [mouse, setMouse] = useState<Point | null>(null);
  useEffect(() => {
    if (!inWails) return;
    const id = setInterval(() => MousePos().then(setMouse), MOUSE_POLL_MS);
    return () => clearInterval(id);
  }, []);

  // Delete / Backspace xoá bước đang chọn (trừ khi đang gõ trong ô nhập).
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const el = e.target as HTMLElement;
      if ((e.key !== 'Delete' && e.key !== 'Backspace') || el.closest('input, textarea, select')) return;
      if (cur) removeStep(cur.id);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [cur]);

  return (
    // h-screen + từng cột tự cuộn: thanh tiêu đề và status bar luôn nằm trong cửa sổ.
    <div className="flex h-screen flex-col overflow-hidden bg-bg text-ink">
      <TitleBar />

      <Toolbar name={script.name} unsaved unsavedTitle="Chưa lưu xuống file – Lưu / Mở sắp có (bước 21)"
        loop={loopSummary(script)} loopPanel={<RunSettings />}
        running={running} onRun={() => (running ? StopScript() : runCurrent())}
        runDisabled={!inWails} runTitle={inWails ? undefined : 'Chỉ chạy được trong app'}
        onPreview={() => openOverlay('preview')} extra={<><RecordButton soon /><DebugButtons soon /></>} />

      <div className="flex min-h-0 flex-1 border-t border-line">
        <Library />
        <StepList />
        {cur ? (
          <StepForm step={cur} index={idx} />
        ) : (
          <aside className="grid w-[340px] flex-none place-items-center border-l border-line bg-surface p-4 text-center text-muted">
            Chọn hoặc thêm một bước để cấu hình
          </aside>
        )}
      </div>

      <StatusBar running={running} mouse={mouse ? `X ${mouse.x} Y ${mouse.y}` : 'X — Y —'}>
        {hotkeyError && <span className="text-danger-text">{hotkeyError}</span>}
        {lastError && <span className="text-danger-text">Lỗi lần chạy trước: {lastError}</span>}
        {stopReason && <span>{STOP_REASON[stopReason] ?? stopReason}</span>}
        {notice && <span className="text-accent-text">{notice}</span>}
      </StatusBar>
    </div>
  );
}
