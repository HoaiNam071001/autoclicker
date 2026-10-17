// Canvas flow của màn chính (bước 20 F4b): nối FlowCanvas (chỉ hiển thị) với store kịch bản.
// Dữ liệu thẻ tính ở flow/view.ts; mọi thao tác (nối, xoá, dời, thêm node) gọi hàm store – component không tự sửa script.
import { useEffect, useMemo, useState } from 'react';
import FlowCanvas from './FlowCanvas';
import AddStepPalette from '../AddStepPalette';
import { Region } from '../icons';
import { layoutFlow } from '../../flow/layout';
import { portsOf } from '../../flow/model';
import { nodeViews } from '../../flow/view';
import { RunState } from './model';
import {
  applyLayout, clearInsertAt, connectPort, copySelected, insertAfterSelected, moveNode, pasteClip, removeEdge, removeSteps, selectNodes, selectStep,
  getEditor, setInsertAt, setRegion, toggleMultiStep, useEditor,
} from '../../store';
import { stepInputCount } from '../../scriptInputs';
import { useEngine } from '../../engineStore';
import { openOverlay } from '../../windowMode';
import type { Script } from '../../types';

const isMac = navigator.platform.startsWith('Mac');

// Xếp tự động cho kịch bản (vị trí theo id node)
const layoutOf = (sc: Script) => layoutFlow(sc.steps.map((s) => ({ id: s.id, ports: portsOf(s) })), sc.edges);

export default function FlowEditor() {
  const { script, selected, picked, path } = useEditor();
  const { running, progress, failedStep } = useEngine();
  const [palette, setPalette] = useState(false);
  const [fitKey, setFitKey] = useState(0);

  // Trạng thái chạy trên thẻ: node đang chạy, node làm lần chạy trước dừng (nhật ký đầy đủ: F5)
  const run = useMemo(() => {
    const m = new Map<string, { run: RunState }>();
    if (failedStep) m.set(failedStep, { run: 'fail' });
    if (running && progress?.nodeId) m.set(progress.nodeId, { run: 'running' });
    return m;
  }, [failedStep, running, progress?.nodeId]);
  const views = useMemo(() => nodeViews(script, { inputs: (id) => stepInputCount(script, id), run }), [script, run]);
  const positions = useMemo(() => Object.fromEntries(script.steps.flatMap((s) => (s.pos ? [[s.id, s.pos]] : []))), [script.steps]);

  // Mở kịch bản khác: node chưa có vị trí (file viết tay / AI tạo) thì tự xếp – không tính là sửa; xong thì vừa màn hình.
  useEffect(() => {
    const sc = script;
    if (sc.steps.some((s) => !s.pos)) layoutOf(sc).then((pos) => { applyLayout(pos, true); setFitKey((k) => k + 1); });
    else setFitKey((k) => k + 1);
  }, [path]); // eslint-disable-line react-hooks/exhaustive-deps

  // Ctrl+K (macOS: Cmd+K) mở / đóng bảng thêm node – node mới nối sau cổng success của node đang chọn.
  // Ctrl+C / V chép / dán node đang chọn, Ctrl+A chọn mọi node – trừ khi đang gõ trong ô nhập.
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (!(isMac ? e.metaKey : e.ctrlKey) || e.altKey || e.shiftKey) return;
      const k = e.key.toLowerCase();
      if (k === 'k') {
        e.preventDefault();
        return setPalette((open) => {
          if (open) clearInsertAt();
          else insertAfterSelected();
          return !open;
        });
      }
      const typing = (e.target as HTMLElement).closest?.('input, textarea, select, [contenteditable="true"]');
      if (typing || !(k === 'c' || k === 'v' || k === 'a')) return;
      e.preventDefault();
      if (k === 'c') copySelected();
      else if (k === 'v') pasteClip();
      else selectNodes(getEditor().script.steps.map((s) => s.id));
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);
  const openPalette = () => {
    insertAfterSelected();
    setPalette(true);
  };

  const { region } = script;
  return (
    <main aria-label="Flow" className="flex min-w-0 flex-1 flex-col">
      <div className="flex flex-none flex-wrap items-center gap-2.5 border-b border-line px-4 py-2">
        <h1 className="text-base font-semibold">Flow</h1>
        <span className="whitespace-nowrap text-muted">{script.steps.length} node · {script.edges.length} cạnh</span>
        <span className="ml-auto flex items-center gap-1">
          <button onClick={() => openOverlay('region')} title="Chọn vùng trên màn hình – toạ độ điểm tính theo vùng"
            className="flex h-8 items-center gap-1.5 rounded-control border border-line bg-surface px-2.5 text-[13px] text-muted hover:bg-hover hover:text-ink">
            <Region className="h-3.5 w-3.5" />
            {region ? <span className="font-mono">Vùng {region.w} × {region.h}</span> : 'Toàn màn hình'}
          </button>
          {region && (
            <button onClick={() => setRegion(null)} aria-label="Về toàn màn hình" title="Về toàn màn hình"
              className="grid h-8 w-8 place-items-center rounded-control text-muted hover:bg-hover hover:text-ink">✕</button>
          )}
          <button onClick={openPalette} className="flex h-8 items-center gap-2 whitespace-nowrap rounded-control border border-line bg-surface px-3 hover:bg-hover">
            + Thêm node <span className="font-mono text-[11px] text-muted">{isMac ? '⌘K' : 'Ctrl+K'}</span>
          </button>
        </span>
      </div>

      <div className="relative min-h-0 flex-1">
        <FlowCanvas nodes={views} edges={script.edges} start={script.start} positions={positions} selected={selected || null} picked={picked} fitKey={fitKey}
          onSelect={(id, additive) => {
            if (id && additive) return toggleMultiStep(id);
            selectStep(id ?? '');
            selectNodes([]);
          }}
          onSelectBox={selectNodes}
          onMove={moveNode}
          onConnect={connectPort}
          onDeleteEdge={removeEdge}
          onDeleteNodes={removeSteps}
          onDropPort={(from, port, _client, pos) => {
            setInsertAt({ from, port, pos: { x: Math.round(pos.x), y: Math.round(pos.y - 28) } });
            setPalette(true);
          }}
          onLayout={() => layoutOf(script).then((pos) => { applyLayout(pos); setFitKey((k) => k + 1); })} />
        {script.steps.length === 0 && (
          <p className="pointer-events-none absolute inset-0 grid place-items-center text-center text-muted">
            Chưa có node nào – bấm “+ Thêm node” hoặc F6 để lưu vị trí chuột.
          </p>
        )}
      </div>

      {palette && <AddStepPalette onClose={() => setPalette(false)} onCancel={clearInsertAt} />}
    </main>
  );
}
