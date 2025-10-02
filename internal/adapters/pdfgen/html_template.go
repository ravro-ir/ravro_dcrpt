package pdfgen

const HTMLTemplate = `<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
    <meta charset="UTF-8">
    <title>گزارش آسیب‌پذیری - {{.ReportID}}</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Vazirmatn:wght@300;400;500;600;700&display=swap');

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Vazirmatn', 'B Nazanin', Tahoma, Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
            background: #ffffff;
            padding: 0;
            margin: 0;
            font-size: 14px;
        }

        .container {
            max-width: 900px;
            margin: 0 auto;
            background: white;
        }

        /* Header Section */
        .header {
            background: #1a1a1a;
            color: white;
            padding: 60px 40px;
            text-align: center;
            position: relative;
        }

        .logo {
            font-size: 18px;
            font-weight: 400;
            letter-spacing: 1px;
            margin-bottom: 20px;
            color: white;
        }

        .header h1 {
            font-size: 24px;
            font-weight: 400;
            margin: 15px 0;
            line-height: 1.4;
        }

        .header .subtitle {
            font-size: 16px;
            margin-top: 15px;
            opacity: 0.9;
            font-weight: 300;
        }

        .header a {
            color: white;
            text-decoration: none;
        }

        /* Icon styling */
        .icon {
            display: inline-block;
            margin-left: 8px;
            font-size: 16px;
        }

        /* Info Table */
        .info-table {
            width: 100%;
            border-collapse: collapse;
            margin: 0;
            background: #ffffff;
        }

        .info-table td {
            padding: 20px 30px;
            border-bottom: 1px solid #f0f0f0;
            width: 50%;
            vertical-align: top;
        }

        .info-label {
            font-size: 12px;
            color: #888888;
            font-weight: 500;
            display: block;
            margin-bottom: 8px;
        }

        .info-value {
            font-size: 16px;
            font-weight: 600;
            color: #333333;
        }

        .info-value a {
            color: #6366f1;
            text-decoration: none;
        }

        .info-value a:hover {
            text-decoration: underline;
        }

        /* Severity Badge */
        .severity-badge {
            display: inline-block;
            padding: 6px 16px;
            border-radius: 20px;
            font-weight: 600;
            font-size: 11px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .severity-critical {
            background: #dc2626;
            color: white;
        }

        .severity-high {
            background: #ea580c;
            color: white;
        }

        .severity-medium {
            background: #d97706;
            color: white;
        }

        .severity-low {
            background: #16a34a;
            color: white;
        }

        /* Content Section */
        .content-section {
            padding: 40px 30px;
        }

        .content-section-alt {
            padding: 40px 30px;
            background: #fafbfc;
        }

        .section-title {
            font-size: 18px;
            font-weight: 600;
            color: #1a1a1a;
            margin-bottom: 30px;
            padding: 16px 20px;
            border-right: 4px solid #f59e0b;
            background: #fffbeb;
            border-radius: 8px;
            display: flex;
            align-items: center;
        }

        .section-title .icon {
            font-size: 20px;
            margin-left: 12px;
            color: #f59e0b;
        }

        /* Card */
        .card {
            background: white;
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 24px;
            border: 1px solid #e5e7eb;
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        }

        .card-header {
            font-size: 16px;
            font-weight: 600;
            color: #1a1a1a;
            margin-bottom: 20px;
            padding: 12px 16px;
            background: #f9fafb;
            border-right: 3px solid #f59e0b;
            border-radius: 6px;
            display: flex;
            align-items: center;
        }

        .card-header .icon {
            font-size: 18px;
            margin-left: 10px;
            color: #f59e0b;
        }

        .card-body {
            color: #374151;
            line-height: 1.7;
            font-size: 14px;
        }

        .card-body p {
            margin-bottom: 16px;
        }

        .card-body strong {
            color: #1a1a1a;
            font-weight: 600;
        }

        /* Data Table */
        .data-table {
            width: 100%;
            border-collapse: collapse;
            background: white;
            margin: 24px 0;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        }

        .data-table th {
            background: #1f2937;
            color: white;
            padding: 16px 20px;
            text-align: right;
            font-weight: 600;
            font-size: 14px;
            border-bottom: 1px solid #374151;
        }

        .data-table td {
            padding: 16px 20px;
            border-bottom: 1px solid #f3f4f6;
            text-align: right;
            font-size: 14px;
        }

        .data-table tr:nth-child(even) {
            background: #f9fafb;
        }

        .data-table tr:hover {
            background: #f3f4f6;
        }

        /* Alert Boxes */
        .alert {
            padding: 20px 24px;
            border-radius: 8px;
            margin-bottom: 24px;
            border-right: 4px solid;
        }

        .alert-info {
            background: #eff6ff;
            border-color: #3b82f6;
            color: #1e40af;
        }

        .alert strong {
            display: block;
            margin-bottom: 12px;
            font-size: 16px;
            font-weight: 600;
        }

        /* Footer */
        .footer {
            background: #1a1a1a;
            color: white;
            padding: 40px 30px;
        }

        .footer-table {
            width: 100%;
            border-collapse: collapse;
        }

        .footer-table td {
            padding: 15px 20px;
            vertical-align: top;
            width: 33.33%;
        }

        .footer h4 {
            color: #f59e0b;
            margin-bottom: 15px;
            font-size: 16px;
            font-weight: 600;
        }

        .footer p {
            line-height: 1.6;
            opacity: 0.9;
            font-size: 14px;
            margin: 8px 0;
        }

        .footer-bottom {
            text-align: center;
            padding-top: 30px;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
            margin-top: 30px;
            opacity: 0.8;
            font-size: 12px;
        }

        .footer-bottom a {
            color: #f59e0b;
            text-decoration: none;
        }

        /* Links */
        a {
            color: #2196f3;
            text-decoration: none;
        }

        a:hover {
            text-decoration: underline;
        }

        /* Page break control */
        .no-break {
            page-break-inside: avoid;
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- Header -->
        <div class="header">
            <h1>{{.Title}}</h1>
            <div class="subtitle">
                <span class="icon">🏢</span>
                <a href="https://www.ravro.ir/company/{{.CompanyUserName}}">
                    {{.CompanyUserName}}
                </a>
            </div>
        </div>

        <!-- Info Grid -->
        <table class="info-table">
            <tr>
                <td>
                    <span class="info-label">شناسه گزارش</span>
                    <div class="info-value">
                        <a href="https://www.ravro.ir/report/{{.ReportID}}">{{.ReportID}}</a>
                    </div>
                </td>
                <td>
                    <span class="info-label">شکارچی</span>
                    <div class="info-value">
                        <a href="https://ravro.ir/hunter/{{.Hunter}}">{{.Hunter}}</a>
                    </div>
                </td>
        </tr>
        <tr>
                <td>
                    <span class="info-label">وضعیت</span>
                    <div class="info-value">{{.Status}}</div>
                </td>
                <td>
                    <span class="info-label">تاریخ ثبت</span>
                    <div class="info-value">{{.DateFrom}}</div>
                </td>
    </tr>
    <tr>
                <td>
                    <span class="info-label">عنوان هدف</span>
                    <div class="info-value">{{.Targets}}</div>
                </td>
                <td>
                    <span class="info-label">بازه زمانی فعالیت</span>
                    <div class="info-value">{{.RangeDate}}</div>
                </td>
    </tr>
</table>

        <!-- Hunter Section -->
        <div class="content-section-alt">
            <h2 class="section-title"><span class="icon">📊</span>اطلاعات شکارچی</h2>
            
            <table class="data-table no-break">
                <thead>
                    <tr>
                        <th>CVSS شکارچی</th>
                        <th>درجه حساسیت</th>
                        <th>IP شکارچی</th>
    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>{{.CVSSHunter}}</td>
                        <td><span class="severity-badge severity-{{.ScoreHunter}}">{{.ScoreHunter}}</span></td>
                        <td>{{.Ips}}</td>
    </tr>
                </tbody>
</table>

            <div class="card no-break">
                <div class="card-header"><span class="icon">📋</span>سناریو آسیب‌پذیری</div>
                <div class="card-body">
                    <p>{{.Scenario}}</p>
    </div>
</div>

            <div class="card no-break">
                <div class="card-header"><span class="icon">🔍</span>شرح آسیب‌پذیری</div>
                <div class="card-body">
                    <p><strong>هدف:</strong> {{.UrlTarget}}</p>
                    <br>
                    <p>{{.PoC}}</p>
    </div>
</div>

            {{if .MoreInfo}}
            <div class="card no-break">
                <div class="card-header"><span class="icon">ℹ️</span>نیاز به اطلاعات بیشتر</div>
                <div class="card-body">
                    <div class="alert alert-info">
                        <p>{{.MoreInfo}}</p>
                        </div>
                    </div>
                </div>
            {{end}}
</div>

        <!-- Judge Section -->
        {{if .JudgeUser}}
        <div class="content-section">
            <h2 class="section-title"><span class="icon">💰</span>اطلاعات تیم داوری</h2>
            
            <table class="data-table no-break">
                <thead>
                    <tr>
                        <th>داور</th>
                        <th>پاداش پیشنهادی</th>
                        <th>CVSS داور</th>
                        <th>درجه حساسیت</th>
                        <th>تاریخ بررسی</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>{{.JudgeUser}}</td>
                        <td>{{.Amount}}</td>
                        <td>{{.CVSSJudge}}</td>
                        <td><span class="severity-badge severity-{{.ScoreJudge}}">{{.ScoreJudge}}</span></td>
                        <td>{{.DateTo}}</td>
    </tr>
                </tbody>
</table>

            {{if .JudgeInfo}}
            <div class="card no-break">
                <div class="card-header"><span class="icon">💬</span>نظر داور</div>
                <div class="card-body">
                    <p>{{.JudgeInfo}}</p>
                        </div>
                    </div>
            {{end}}
                </div>
        {{end}}

        <!-- Additional Info -->
        {{if .LinkMoreInfo}}
        <div class="content-section-alt">
            <h2 class="section-title"><span class="icon">🔗</span>شرح آسیب‌پذیری</h2>
            
            <div class="alert alert-info">
                <strong>برای اطلاعات بیشتر:</strong>
                <p>{{.LinkMoreInfo}}</p>
                        </div>
                    </div>
        {{end}}

        <!-- Attachments -->
        {{if .Attachment}}
        <div class="content-section">
            <h2 class="section-title"><span class="icon">📁</span>فایل‌های پیوست</h2>
            
            {{.Attachment}}
                </div>
        {{end}}

        <!-- Footer -->
        <div class="footer">
            <table class="footer-table">
                <tr>
                    <td>
                        <h4>📍 آدرس</h4>
                        <p>تهران، خیابان مطهری، نبش سهروردی، پلاک ۹۴، طبقه دوم، واحد ۲۵۰</p>
                    </td>
                    <td>
                        <h4>📞 تماس</h4>
                        <p>021-9103-1553</p>
                        <p>support@ravro.ir</p>
                    </td>
                    <td>
                        <h4>🌐 وب‌سایت</h4>
                        <p><a href="https://www.ravro.ir" style="color: #ffc107;">www.ravro.ir</a></p>
                    </td>
                </tr>
            </table>

            <div class="footer-bottom">
                <p>این گزارش با نسخه {{.RavroVer}} ابزار 
                    <a href="https://github.com/ravro-ir/ravro_dcrpt">ravro_dcrpt</a> رمزگشایی شده است.
                </p>
            </div>
        </div>
    </div>
</body>
</html>
`
