// Màn soạn kịch bản (design v2, artboard Main): thanh tiêu đề, thanh công cụ, thư viện | danh sách bước | inspector, status bar
import { useEffect, useState } from 'react';
import { DebugButtons, RecordButton, StatusBar, Toolbar } from '../components/AppChrome';
import { MousePos, StopScript } from '../../wailsjs/go/main/App';
import { inWails, openOverlay } from '../windowMode';
import { removeStep, useEditor } from '../scriptStore';
import { closeRunForm, deletePreset, loadMemory, savePreset, startRun, submitRun, useRunForm } from '../runForm';
import { inputViews, runSections } from '../scriptInputs';
import RunDialog from '../components/inputs/RunDialog';
import { newScript, openScript, save, useFile } from '../scriptFile';
import SaveChangesDialog from '../components/SaveChangesDialog';
import TitleBar from '../components/TitleBar';
import Library from '../components/Library';
import StepList from '../components/StepList';
import StepForm from '../components/StepForm';
import RunSettings, { loopSummary } from '../components/RunSettings';
import { STOP_REASON, useEngine } from '../engineStore';
import { Point } from '../types';

const MOUSE_POLL_MS = 250;
const isMac = navigator.platform.startsWith('Mac');
const modKey = isMac ? '⌘' : 'Ctrl+';
const clock = (ms: number) => new Date(ms).toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit', second: '2-digit' });

export default function MainScreen() {
  const { script, selected, notice, dirty, path, countdown: editorCountdown } = useEditor();
  const { draftAt, error: fileError } = useFile();
  const { steps } = script;
  const cur = steps.find((s) => s.id === selected) ?? steps[0]; // undefined khi chưa có bước nào
  const idx = steps.indexOf(cur) + 1;
  const { running, lastError, stopReason, hotkeyError } = useEngine();
  const runForm = useRunForm();
  const autoRun = !!script.inputs?.length && !!runForm.memory.autoRun;
  useEffect(() => {
    loadMemory();
  }, [path, script.name]);

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

  // Cmd/Ctrl + S lưu, + Shift + S lưu thành, + O mở, + N kịch bản mới (cả khi đang gõ trong ô nhập).
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (!(isMac ? e.metaKey : e.ctrlKey) || e.altKey) return;
      const k = e.key.toLowerCase();
      if (k === 's') save(e.shiftKey);
      else if (k === 'o' && !e.shiftKey) openScript();
      else if (k === 'n' && !e.shiftKey) newScript();
      else return;
      e.preventDefault();
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  return (
    // h-screen + từng cột tự cuộn: thanh tiêu đề và status bar luôn nằm trong cửa sổ.
    <div className="flex h-screen flex-col overflow-hidden bg-bg text-ink">
      <TitleBar />

      <Toolbar name={script.name} nameTitle={path ?? 'Chưa lưu thành file'} unsaved={dirty}
        unsavedTitle={`Có thay đổi chưa lưu xuống file (${modKey}S) – bản nháp vẫn tự lưu trên máy`}
        loop={loopSummary(script)} loopPanel={<RunSettings />}
        running={running} onRun={(e) => (running ? StopScript() : startRun(e.shiftKey))}
        runDisabled={!inWails}
        runTitle={!inWails ? 'Chỉ chạy được trong app'
          : autoRun ? 'Chạy luôn với giá trị lần trước – Shift + bấm để sửa input' : undefined}
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
        {fileError && <span className="text-danger-text">{fileError}</span>}
        {notice && <span className="text-accent-text">{notice}</span>}
        {draftAt && dirty && <span title="Bản nháp lưu trên máy – mở lại app vẫn còn">Bản nháp tự lưu lúc <span className="font-mono">{clock(draftAt)}</span></span>}
      </StatusBar>
      <SaveChangesDialog />
      {runForm.open && (
        <RunDialog scriptName={script.name} inputs={inputViews(script)} sections={runSections(script)}
          presets={runForm.memory.presets ?? []} initialValues={runForm.memory.last} initialPreset={runForm.memory.lastPreset} initialRemember={!!runForm.memory.autoRun}
          countdown={editorCountdown} error={runForm.error} onClose={closeRunForm} onRun={submitRun}
          onSavePreset={savePreset} onDeletePreset={deletePreset} />
      )}
    </div>
  );
}
