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

        .header-logo {
            display: block;
            margin: 0 auto 20px auto;
            width: 48px;
            height: 48px;
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
            text-align: center;
            display: block;
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

        .card-header svg,
        .section-title svg {
            width: 20px;
            height: 20px;
            margin-left: 10px;
            vertical-align: middle;
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
            <svg class="header-logo" width="48" height="48" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M20.871 4.84413C22.8023 3.71862 25.1977 3.71862 27.129 4.84413L38.9576 11.7373C40.8426 12.8358 42 14.8406 42 17.0072V30.9928C42 33.1594 40.8426 35.1642 38.9576 36.2627L27.129 43.1559C25.1977 44.2814 22.8023 44.2814 20.871 43.1559L9.04239 36.2627C7.1574 35.1642 6 33.1594 6 30.9928V17.0072C6 14.8406 7.1574 12.8358 9.04239 11.7373L20.871 4.84413Z" fill="#F7B500"/>
                <path d="M23.5309 9.28146C23.8196 9.10822 24.1804 9.10822 24.4691 9.28146L28.5573 11.7344C28.832 11.8992 29 12.196 29 12.5162V17.4838C29 17.804 28.832 18.1008 28.5573 18.2656L24.4691 20.7185C24.1804 20.8918 23.8196 20.8918 23.5309 20.7185L19.4427 18.2656C19.168 18.1008 19 17.804 19 17.4838V12.5162C19 12.196 19.168 11.8992 19.4427 11.7344L23.5309 9.28146Z" fill="white"/>
                <path d="M28.5573 23.7344C28.832 23.8992 29 24.196 29 24.5162V35.4838C29 35.804 28.832 36.1008 28.5573 36.2656L24.4691 38.7185C24.1804 38.8918 23.8196 38.8918 23.5309 38.7185L20.3018 36.7811C19.7118 36.4271 19.7115 35.5723 20.3011 35.2178L23.987 33.0019C23.9889 33.0007 23.9901 32.9986 23.9901 32.9963C23.9901 32.9941 23.9889 32.992 23.987 32.9908L19.4428 30.2656C19.1681 30.1008 19 29.804 19 29.4836V24.5162C19 24.196 19.168 23.8992 19.4427 23.7344L23.5309 21.2815C23.8196 21.1082 24.1804 21.1082 24.4691 21.2815L28.5573 23.7344Z" fill="white"/>
            </svg>
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
                    <p>{{.MoreInfo}}</p>
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
                <div class="card-header">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M21 15C21 15.5304 20.7893 16.0391 20.4142 16.4142C20.0391 16.7893 19.5304 17 19 17H7L3 21V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H19C19.5304 3 20.0391 3.21071 20.4142 3.58579C20.7893 3.96086 21 4.46957 21 5V15Z" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="#f59e0b" fill-opacity="0.1"/>
                    </svg>
                    نظر داور
                </div>
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
            <h2 class="section-title">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M13 2H6C5.46957 2 4.96086 2.21071 4.58579 2.58579C4.21071 2.96086 4 3.46957 4 4V20C4 20.5304 4.21071 21.0391 4.58579 21.4142C4.96086 21.7893 5.46957 22 6 22H18C18.5304 22 19.0391 21.7893 19.4142 21.4142C19.7893 21.0391 20 20.5304 20 20V9L13 2Z" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="#f59e0b" fill-opacity="0.1"/>
                    <path d="M13 2V9H20" stroke="#f59e0b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                فایل‌های پیوست
            </h2>
            
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
                        <p>021-9103-5315</p>
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
