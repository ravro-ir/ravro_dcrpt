package pdfgen

const reportHTMLStyles = `
    :root {
      --bg: #0f1117;
      --surface: #171b24;
      --surface-2: #1e2330;
      --border: #262c3d;
      --border-subtle: #1e2330;
      --text-primary: #e8eaf0;
      --text-secondary: #8892a4;
      --text-tertiary: #555f72;
      --accent: #6366f1;
      --accent-dim: rgba(99,102,241,0.12);
      --sev-critical-bg: rgba(239,68,68,0.12);
      --sev-critical-text: #f87171;
      --sev-critical-border: rgba(239,68,68,0.25);
      --sev-high-bg: rgba(249,115,22,0.12);
      --sev-high-text: #fb923c;
      --sev-high-border: rgba(249,115,22,0.25);
      --sev-medium-bg: rgba(234,179,8,0.12);
      --sev-medium-text: #facc15;
      --sev-medium-border: rgba(234,179,8,0.25);
      --sev-low-bg: rgba(34,197,94,0.12);
      --sev-low-text: #4ade80;
      --sev-low-border: rgba(34,197,94,0.25);
      --sev-info-bg: rgba(99,102,241,0.12);
      --sev-info-text: #818cf8;
      --sev-info-border: rgba(99,102,241,0.25);
      --status-bg: rgba(16,185,129,0.1);
      --status-text: #34d399;
      --status-border: rgba(16,185,129,0.2);
      --mono: 'JetBrains Mono', 'Courier New', monospace;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: 'Vazirmatn', Tahoma, system-ui, sans-serif;
      font-size: 14px;
      background: var(--bg);
      color: var(--text-primary);
      line-height: 1.75;
      min-height: 100vh;
      padding: 24px 16px 48px;
    }
    .topbar {
      max-width: 820px;
      margin: 0 auto 16px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
    }
    .topbar-left {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 12px;
      color: var(--text-tertiary);
    }
    .topbar-left svg { flex-shrink: 0; }
    .topbar-actions { display: flex; gap: 8px; }
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 7px 14px;
      border-radius: 7px;
      font-family: inherit;
      font-size: 12px;
      font-weight: 500;
      cursor: pointer;
      border: 1px solid var(--border);
      background: var(--surface);
      color: var(--text-secondary);
      transition: color .15s, background .15s, border-color .15s;
      text-decoration: none;
    }
    .btn:hover { background: var(--surface-2); color: var(--text-primary); border-color: var(--accent); }
    .btn svg { flex-shrink: 0; }
    .card {
      max-width: 820px;
      margin: 0 auto 16px;
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: 14px;
      overflow: hidden;
    }
    .hero {
      padding: 28px 28px 24px;
      border-bottom: 1px solid var(--border);
      position: relative;
    }
    .hero::before {
      content: '';
      position: absolute;
      top: 0; left: 0; right: 0;
      height: 3px;
      background: linear-gradient(90deg, #f87171 0%, #fb923c 40%, #6366f1 100%);
    }
    .hero-top {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 16px;
      flex-wrap: wrap;
      margin-bottom: 16px;
    }
    .hero-id {
      font-family: var(--mono);
      font-size: 11px;
      color: var(--text-tertiary);
      letter-spacing: .04em;
      margin-bottom: 8px;
    }
    .hero-title {
      font-size: 20px;
      font-weight: 700;
      line-height: 1.4;
      color: var(--text-primary);
    }
    .hero-right {
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      gap: 8px;
      flex-shrink: 0;
    }
    .severity-block {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 4px;
      padding: 5px;
      border-radius: 16px;
    }
    .severity-score {
      font-family: var(--mono);
      font-size: 28px;
      font-weight: 700;
      line-height: 1;
    }
    .severity-badge {
      display: inline-flex;
      align-items: center;
      gap: 5px;
      padding: 3px 10px;
      border-radius: 999px;
      font-size: 11px;
      font-weight: 600;
      letter-spacing: .03em;
      border: 1px solid;
    }
    .sev-critical { background: var(--sev-critical-bg); color: var(--sev-critical-text); border-color: var(--sev-critical-border); }
    .sev-high { background: var(--sev-high-bg); color: var(--sev-high-text); border-color: var(--sev-high-border); }
    .sev-medium { background: var(--sev-medium-bg); color: var(--sev-medium-text); border-color: var(--sev-medium-border); }
    .sev-low { background: var(--sev-low-bg); color: var(--sev-low-text); border-color: var(--sev-low-border); }
    .sev-info { background: var(--sev-info-bg); color: var(--sev-info-text); border-color: var(--sev-info-border); }
    .sev-critical .severity-score { color: var(--sev-critical-text); }
    .sev-high .severity-score { color: var(--sev-high-text); }
    .sev-medium .severity-score { color: var(--sev-medium-text); }
    .sev-low .severity-score { color: var(--sev-low-text); }
    .sev-info .severity-score { color: var(--sev-info-text); }
    .status-pill {
      display: inline-flex;
      align-items: center;
      gap: 5px;
      padding: 3px 10px;
      border-radius: 999px;
      font-size: 11px;
      font-weight: 600;
      border: 1px solid;
    }
    .status-company { background: var(--status-bg); color: var(--status-text); border-color: var(--status-border); }
    .status-judge { background: rgba(234,179,8,0.1); color: #facc15; border-color: rgba(234,179,8,0.2); }
    .status-resolved { background: var(--status-bg); color: var(--status-text); border-color: var(--status-border); }
    .status-rejected { background: rgba(239,68,68,0.1); color: #f87171; border-color: rgba(239,68,68,0.2); }
    .status-pending { background: rgba(148,163,184,0.1); color: #94a3b8; border-color: rgba(148,163,184,0.2); }
    .status-dot {
      width: 6px; height: 6px;
      border-radius: 50%;
      background: currentColor;
    }
    .status-pill.pulse .status-dot { animation: pulse 2s ease-in-out infinite; }
    @keyframes pulse { 0%,100% { opacity:1; } 50% { opacity:.4; } }
    .meta-row {
      display: flex;
      gap: 20px;
      flex-wrap: wrap;
      font-size: 12px;
      color: var(--text-secondary);
    }
    .meta-item { display: flex; align-items: center; gap: 5px; }
    .meta-item svg { color: var(--text-tertiary); flex-shrink: 0; }
    .info-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 0;
    }
    .info-cell { padding: 16px 20px; }
    .info-cell:last-child { border-left: none; }
    .info-cell.span-3 { grid-column: span 3; border-left: none; }
    .cell-label {
      font-size: 10px;
      font-weight: 600;
      letter-spacing: .08em;
      text-transform: uppercase;
      color: var(--text-tertiary);
      margin-bottom: 6px;
    }
    .cell-value {
      font-size: 13px;
      font-weight: 500;
      color: var(--text-primary);
      word-break: break-word;
    }
    .cell-value.mono {
      font-family: var(--mono);
      font-size: 12px;
      color: var(--text-secondary);
    }
    .cell-value a {
      color: #818cf8;
      text-decoration: none;
      font-family: var(--mono);
      font-size: 12px;
    }
    .cell-value a:hover { text-decoration: underline; }
    .section { border-bottom: 1px solid var(--border); }
    .section:last-child { border-bottom: none; }
    .section-head {
      padding: 14px 20px;
      display: flex;
      align-items: center;
      gap: 10px;
      border-bottom: 1px solid var(--border-subtle);
      background: rgba(255,255,255,.015);
    }
    summary.section-head {
      list-style: none;
      cursor: pointer;
      user-select: none;
    }
    summary.section-head::-webkit-details-marker,
    summary.section-head::marker { display: none; content: none; }
    details.section:not([open]) > summary.section-head { border-bottom: none; }
    .section-chevron {
      flex-shrink: 0;
      margin-inline-start: 8px;
      color: var(--text-tertiary);
      transition: transform .15s ease;
    }
    details.section[open] > summary .section-chevron { transform: rotate(180deg); }
    .section-icon {
      width: 28px; height: 28px;
      border-radius: 7px;
      background: var(--accent-dim);
      color: var(--accent);
      display: flex; align-items: center; justify-content: center;
      flex-shrink: 0;
    }
    .section-title {
      font-size: 13px;
      font-weight: 600;
      color: var(--text-primary);
      flex: 1;
    }
    .section-count {
      margin-inline-start: auto;
      font-size: 11px;
      color: var(--text-tertiary);
      font-family: var(--mono);
    }
    .section-body { padding: 20px; }
    .cvss-table {
      width: 100%;
      border-collapse: collapse;
      border: 1px solid var(--border);
      border-radius: 10px;
      overflow: hidden;
      font-size: 12px;
    }
    .cvss-table th {
      background: var(--surface-2);
      color: var(--text-tertiary);
      font-weight: 600;
      font-size: 10px;
      letter-spacing: .07em;
      text-transform: uppercase;
      padding: 10px 14px;
      text-align: start;
      border-bottom: 1px solid var(--border);
    }
    .cvss-table td {
      padding: 12px 14px;
      border-bottom: 1px solid var(--border);
      vertical-align: middle;
    }
    .cvss-table tr:last-child td { border-bottom: none; }
    .cvss-table .role-label { font-weight: 600; color: var(--text-secondary); }
    .cvss-vector {
      font-family: var(--mono);
      font-size: 11px;
      color: var(--text-secondary);
      word-break: break-all;
      direction: ltr;
    }
    .null-value {
      font-size: 11px;
      color: var(--text-tertiary);
      font-style: italic;
    }
    .cvss-score-cell {
      font-family: var(--mono);
      font-size: 14px;
      font-weight: 700;
      background: transparent;
    }
    .text-block {
      font-size: 14px;
      color: var(--text-secondary);
      line-height: 1.7;
      word-break: break-word;
    }
    .md-content { overflow-wrap: break-word; }
    .md-content > :first-child { margin-top: 0; }
    .md-content > :last-child { margin-bottom: 0; }
    .md-content p { margin: 0 0 0.65em; }
    .md-content strong, .md-content b { font-weight: 700; }
    .md-content h1, .md-content h2, .md-content h3, .md-content h4, .md-content h5, .md-content h6 {
      color: var(--text-primary);
      font-weight: 700;
      margin: 1em 0 0.5em;
      line-height: 1.35;
    }
    .md-content h1 { font-size: 1.25em; }
    .md-content h2 { font-size: 1.15em; }
    .md-content h3 { font-size: 1.05em; }
    .md-content h4, .md-content h5, .md-content h6 { font-size: 1em; }
    .md-content ul, .md-content ol {
      margin: 0 0 0.75em;
      padding-inline-start: 1.4em;
    }
    .md-content li { margin: 0.25em 0; }
    .md-content blockquote {
      margin: 0.75em 0;
      padding: 0.5em 1em;
      border-inline-start: 3px solid var(--accent);
      background: var(--surface-2);
      color: var(--text-secondary);
    }
    .md-content hr {
      border: none;
      border-top: 1px solid var(--border);
      margin: 1em 0;
    }
    .md-content a {
      color: #818cf8;
      text-decoration: none;
      word-break: break-all;
    }
    .md-content a:hover { text-decoration: underline; }
    .md-content code {
      font-family: var(--mono);
      font-size: 0.9em;
      background: var(--surface-2);
      color: var(--text-primary);
      padding: 0.1em 0.35em;
      border-radius: 4px;
      direction: ltr;
      unicode-bidi: isolate;
    }
    .md-content pre {
      margin: 0.75em 0;
      padding: 12px 14px;
      border-radius: 8px;
      overflow-x: auto;
      background: var(--surface-2);
      border: 1px solid var(--border);
      direction: ltr;
      text-align: left;
    }
    .md-content pre code {
      background: transparent;
      padding: 0;
      border-radius: 0;
      font-size: 12px;
      line-height: 1.65;
    }
    .md-content table {
      display: block;
      width: 100%;
      overflow-x: auto;
      border-collapse: collapse;
      margin: 0.75em 0;
      font-size: 12px;
    }
    .md-content th, .md-content td {
      border: 1px solid var(--border);
      padding: 6px 10px;
    }
    .md-content th {
      background: var(--surface-2);
      font-weight: 600;
      color: var(--text-primary);
    }
    .md-content img {
      max-width: 100%;
      height: auto;
      border-radius: 8px;
    }
    .attachments-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
      gap: 10px;
    }
    .attachment-card {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px 14px;
      background: var(--surface-2);
      border: 1px solid var(--border);
      border-radius: 10px;
      transition: border-color .15s;
      text-decoration: none;
      color: inherit;
    }
    .attachment-card:hover { border-color: var(--accent); }
    .file-icon {
      width: 36px; height: 36px;
      border-radius: 8px;
      display: flex; align-items: center; justify-content: center;
      flex-shrink: 0;
      font-size: 11px;
      font-weight: 700;
      letter-spacing: .03em;
    }
    .file-icon-pdf { background: rgba(239,68,68,.15); color: #f87171; }
    .file-icon-img { background: rgba(99,102,241,.15); color: #818cf8; }
    .file-icon-zip { background: rgba(234,179,8,.15); color: #facc15; }
    .file-icon-vid { background: rgba(16,185,129,.15); color: #34d399; }
    .file-icon-gen { background: rgba(148,163,184,.1); color: #94a3b8; }
    .att-meta { min-width: 0; }
    .att-name {
      font-family: var(--mono);
      font-size: 11px;
      color: var(--text-primary);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      direction: ltr;
      text-align: right;
    }
    .att-size { font-size: 11px; color: var(--text-tertiary); margin-top: 2px; }
    .attachments-note {
      margin-top: 14px;
      padding: 10px 14px;
      background: var(--accent-dim);
      border: 1px solid rgba(99,102,241,.2);
      border-radius: 8px;
      font-size: 12px;
      color: var(--text-secondary);
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .comment-item { display: flex; gap: 12px; margin-bottom: 16px; }
    .comment-item:last-child { margin-bottom: 0; }
    .avatar {
      width: 32px; height: 32px;
      border-radius: 8px;
      background: var(--accent-dim);
      border: 1px solid rgba(99,102,241,.25);
      color: var(--accent);
      display: flex; align-items: center; justify-content: center;
      font-size: 12px;
      font-weight: 700;
      flex-shrink: 0;
    }
    .comment-body { flex: 1; }
    .comment-meta {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 6px;
      flex-wrap: wrap;
    }
    .comment-author { font-size: 12px; font-weight: 600; color: var(--text-primary); }
    .comment-role {
      font-size: 10px;
      padding: 1px 7px;
      border-radius: 999px;
      font-weight: 600;
      letter-spacing: .03em;
      border: 1px solid;
    }
    .comment-role-bughunter { background: var(--accent-dim); color: #818cf8; border-color: rgba(99,102,241,.2); }
    .comment-role-company { background: rgba(16,185,129,0.1); color: #34d399; border-color: rgba(16,185,129,0.2); }
    .comment-role-judge { background: rgba(234,179,8,0.1); color: #facc15; border-color: rgba(234,179,8,0.2); }
    .comment-role-support { background: rgba(56,189,248,0.1); color: #38bdf8; border-color: rgba(56,189,248,0.2); }
    .party-chip { display: inline-flex; align-items: center; gap: 8px; }
    .party-chip .avatar,
    .party-chip .avatar-img {
      width: 24px; height: 24px; font-size: 11px; border-radius: 6px;
    }
    .party-chip .avatar-img { object-fit: cover; display: block; }
    .party-chip .avatar-wrap { display: inline-flex; flex-shrink: 0; }
    .party-link {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      color: #818cf8;
      text-decoration: none;
    }
    .party-link:hover .party-name { text-decoration: underline; }
    .timeline-list { display: flex; flex-direction: column; gap: 10px; }
    .timeline-row { display: flex; gap: 12px; align-items: baseline; flex-wrap: wrap; }
    .timeline-time { font-family: var(--mono); font-size: 11px; color: var(--text-tertiary); }
    .timeline-stages { font-size: 13px; color: var(--text-primary); }
    .appeal-grid {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 0;
    }
    .plain-text { white-space: pre-wrap; word-break: break-word; }
    .comment-time { font-size: 11px; color: var(--text-tertiary); font-family: var(--mono); }
    .comment-text {
      font-size: 13px;
      color: var(--text-secondary);
      background: var(--surface-2);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: 10px 14px;
      line-height: 1.7;
      word-break: break-word;
    }
    .amendment-item {
      padding: 16px;
      background: var(--surface-2);
      border: 1px solid var(--border);
      border-radius: 10px;
      margin-bottom: 12px;
    }
    .amendment-item:last-child { margin-bottom: 0; }
    .amendment-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 8px;
      flex-wrap: wrap;
    }
    .amendment-number { font-size: 12px; font-weight: 600; color: #c084fc; }
    .amendment-time { font-size: 11px; color: var(--text-tertiary); font-family: var(--mono); }
    .amendment-stage {
      font-size: 11px;
      color: var(--text-tertiary);
      margin-bottom: 10px;
      line-height: 1.5;
    }
    .amendment-stage-label {
      display: block;
      font-size: 10px;
      font-weight: 600;
      letter-spacing: .06em;
      text-transform: uppercase;
      color: var(--text-tertiary);
      margin-bottom: 4px;
    }
    .amendment-explanation-label {
      font-size: 10px;
      font-weight: 600;
      letter-spacing: .06em;
      text-transform: uppercase;
      color: var(--text-tertiary);
      margin-bottom: 6px;
    }
    .amendment-explanation {
      font-size: 13px;
      color: var(--text-secondary);
      line-height: 1.7;
      word-break: break-word;
    }
    .amendment-attachments { margin-top: 14px; }
    .amendment-attachments-label {
      font-size: 11px;
      color: var(--text-tertiary);
      margin-bottom: 8px;
    }
    .report-footer {
      max-width: 820px;
      margin: 0 auto;
      padding: 20px 0 0;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      flex-wrap: wrap;
    }
    .footer-brand {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 12px;
      color: var(--text-tertiary);
    }
    .footer-brand strong { color: var(--text-secondary); }
    [data-theme="light"] {
      --bg: #fff;
      --surface: #fff;
      --surface-2: #f5f5f5;
      --border: #e0e0e0;
      --border-subtle: #eee;
      --text-primary: #111;
      --text-secondary: #444;
      --text-tertiary: #888;
      --accent: #4f46e5;
      --accent-dim: rgba(79,70,229,.06);
    }
    @media print {
      .no-print { display: none !important; }
      details.section { display: block; }
      details.section > .section-body { display: block !important; }
      .section-chevron { display: none; }
      body { background: #fff; color: #000; padding: 0; font-size: 12px; }
      .card { border: 1px solid #ddd; border-radius: 0; margin: 0; max-width: 100%; }
      :root {
        --bg: #fff;
        --surface: #fff;
        --surface-2: #f5f5f5;
        --border: #e0e0e0;
        --border-subtle: #eee;
        --text-primary: #111;
        --text-secondary: #444;
        --text-tertiary: #888;
        --accent: #4f46e5;
        --accent-dim: rgba(79,70,229,.06);
      }
    }
    @media (max-width: 600px) {
      .info-grid { grid-template-columns: 1fr 1fr; }
      .appeal-grid { grid-template-columns: 1fr 1fr; }
      .info-cell.span-3 { grid-column: span 2; }
      .hero-top { flex-direction: column-reverse; }
      .hero-right { align-items: flex-start; flex-direction: row; }
    }
`
