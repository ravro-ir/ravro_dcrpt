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
            line-height: 1.9;
            color: #333333;
            background: #f5f5f5;
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
            background: #ff5722;
            background: linear-gradient(135deg, #ff5722 0%, #e64a19 100%);
            color: white;
            padding: 40px 20px;
            text-align: center;
        }

        .logo {
            font-size: 32px;
            font-weight: bold;
            letter-spacing: 2px;
            margin-bottom: 15px;
            color: white;
        }

        .header h1 {
            font-size: 26px;
            font-weight: bold;
            margin: 10px 0;
        }

        .header .subtitle {
            font-size: 17px;
            margin-top: 10px;
            opacity: 0.95;
        }

        .header a {
            color: white;
            text-decoration: none;
        }

        /* Icon styling */
        .icon {
            display: inline-block;
            margin-left: 5px;
            font-size: 18px;
        }

        /* Info Table */
        .info-table {
            width: 100%;
            border-collapse: collapse;
            margin: 0;
            background: #f8f9fa;
        }

        .info-table td {
            padding: 15px 20px;
            border-bottom: 1px solid #e0e0e0;
            width: 50%;
        }

        .info-label {
            font-size: 11px;
            color: #666666;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            font-weight: bold;
            display: block;
            margin-bottom: 5px;
        }

        .info-value {
            font-size: 14px;
            font-weight: bold;
            color: #333333;
        }

        .info-value a {
            color: #2196f3;
            text-decoration: none;
        }

        /* Severity Badge */
        .severity-badge {
            display: inline-block;
            padding: 5px 12px;
            border-radius: 12px;
            font-weight: bold;
            font-size: 12px;
            text-transform: uppercase;
        }

        .severity-critical {
            background: #f44336;
            color: white;
        }

        .severity-high {
            background: #ff9800;
            color: white;
        }

        .severity-medium {
            background: #ffc107;
            color: #333;
        }

        .severity-low {
            background: #4caf50;
            color: white;
        }

        /* Content Section */
        .content-section {
                padding: 30px 20px;
            }

        .content-section-alt {
            padding: 30px 20px;
            background: #fafafa;
        }

        .section-title {
            font-size: 20px;
            font-weight: bold;
            color: #ff5722;
            margin-bottom: 20px;
            padding: 12px 15px;
            border-right: 5px solid #ff5722;
            background: #fff3f0;
            border-radius: 4px;
        }

        .section-title .icon {
            font-size: 22px;
            margin-left: 8px;
        }

        /* Card */
        .card {
            background: white;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
            border: 1px solid #e0e0e0;
        }

        .card-header {
            font-size: 15px;
            font-weight: bold;
            color: #333333;
            margin-bottom: 15px;
            padding: 10px 12px;
            background: #f8f9fa;
            border-right: 3px solid #ff5722;
            border-radius: 3px;
        }

        .card-header .icon {
            font-size: 16px;
            margin-left: 6px;
        }

        .card-body {
            color: #333333;
            line-height: 1.8;
        }

        .card-body p {
            margin-bottom: 12px;
        }

        .card-body strong {
            color: #ff5722;
        }

        /* Data Table */
        .data-table {
            width: 100%;
            border-collapse: collapse;
            background: white;
            margin: 20px 0;
        }

        .data-table th {
            background: #ff5722;
            color: white;
            padding: 12px 15px;
            text-align: right;
            font-weight: bold;
            font-size: 13px;
            border: 1px solid #e64a19;
        }

        .data-table td {
            padding: 12px 15px;
            border: 1px solid #e0e0e0;
            text-align: right;
        }

        .data-table tr:nth-child(even) {
            background: #f8f9fa;
        }

        /* Alert Boxes */
        .alert {
            padding: 15px 20px;
            border-radius: 6px;
            margin-bottom: 20px;
            border-right: 4px solid;
        }

        .alert-info {
            background: #e3f2fd;
            border-color: #2196f3;
            color: #0d47a1;
        }

        .alert strong {
            display: block;
            margin-bottom: 8px;
            font-size: 15px;
        }

        /* Footer */
        .footer {
            background: #1a1a2e;
            color: white;
            padding: 30px 20px;
        }

        .footer-table {
            width: 100%;
            border-collapse: collapse;
        }

        .footer-table td {
            padding: 10px 15px;
            vertical-align: top;
            width: 33.33%;
        }

        .footer h4 {
            color: #ffc107;
            margin-bottom: 10px;
            font-size: 14px;
        }

        .footer p {
            line-height: 1.6;
            opacity: 0.8;
            font-size: 12px;
            margin: 5px 0;
        }

        .footer-bottom {
            text-align: center;
            padding-top: 20px;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
            margin-top: 20px;
            opacity: 0.7;
            font-size: 11px;
        }

        .footer-bottom a {
            color: #ffc107;
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
            <div class="logo">🔐 RAVRO</div>
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
            <h2 class="section-title"><span class="icon">🎯</span>اطلاعات شکارچی</h2>
            
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
            <h2 class="section-title"><span class="icon">⚖️</span>قیمت پیشنهادی تیم داوری</h2>
            
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
            <h2 class="section-title"><span class="icon">📚</span>منابع آموزشی</h2>
            
            <div class="alert alert-info">
                <strong>برای اطلاعات بیشتر:</strong>
                <p>{{.LinkMoreInfo}}</p>
                        </div>
                    </div>
        {{end}}

        <!-- Attachments -->
        {{if .Attachment}}
        <div class="content-section">
            <h2 class="section-title"><span class="icon">📎</span>فایل‌های پیوست</h2>
            
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
