package main

import (
	"io/ioutil"
)

func HtmlTemplate(htmlPath string) {

	html := "<html>\n<head>\n  <title>{TITLE}</title>\n  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, user-scalable=yes\">\n  <meta name=\"description\" content=\"\">\n  <meta name=\"keywords\" content=\"\">\n  <link type=\"text/css\" rel=\"stylesheet\" href=\"https://cdn.jsdelivr.net/gh/rastikerdar/samim-font@v4.0.5/dist/font-face.css\">\n  <style>\n    * {\n      margin: 0;\n      padding: 0;\n    }\n\n    a {\n      text-decoration: none;\n      color: #007fba;\n    }\n\n    html {\n      font-size: 16px;\n    }\n\n    body {\n      color: #000;\n      direction: rtl;\n      line-height: 2;\n      font-family: Samim,Tahoma,\"DejaVu Sans\",helvetica,arial,freesans,sans-serif;\n    }\n\n    header, .sample, footer {\n      margin: auto;\n      padding: 30px 50px;\n    }\n\n    header {\n      text-align: center;\n      background-color: #fff;\n      color: #222;\n    }\n\n    header h1 {\n      color: #6f6f91;\n      font-size: 3rem;\n      font-weight: 500;\n    }\n\n    header h1 span {\n      color: #6c5367;\n    }\n\n    header .version {\n      color: #000000;\n      font-weight: bold;\n    }\n\n    header .get {\n      margin: 10px 0;\n      font-size: 1.5rem;\n    }\n\n    header .get a {\n      display: inline-block;\n      border-bottom: 1px solid #d5d5d5;\n      color: #007fba;\n      margin: 0 5px;\n    }\n\n    header .info{\n      color: #999;\n    }\n\n    .infobox {\n      background-color: #f4f7ff;\n      border: 1px solid #e2e2e2;\n      border-radius: 5px;\n      color: #284559;\n      font-weight: 900;\n      line-height: 2;\n      margin: 0 auto;\n      max-width: 600px;\n      padding: 5px;\n    }\n    .sample {\n      background-color: #f5f8fa;\n      color: #181b1d;\n      border-top: 1px solid #eee;\n    }\n    .sample h2 {\n      font-weight: normal;\n      color: #00727b;\n    }\n    .sample p {\n      margin-bottom: 30px;\n    }\n    .sample a { color: #0b6db3; }\n    footer {\n      background-color: #1d2b2a;\n      /*background-image: url(\"paint.jpg\");*/\n      background-size: cover;\n      color: #fff;\n      text-align: center;\n    }\n    footer a {\n      color: #FEAF3A;\n    }\n    .clear { clear: both; }\n    @media all and (max-width: 480px){\n      html {\n        font-size: 15px;\n      }\n      header, .sample, footer {\n        padding: 30px 20px;\n      }\n\n      .sample {\n        text-align: right;\n      }\n\n      header h1 {\n        font-size: 2.5rem;\n      }\n    }\n    table {\n      font-family: Samim,Tahoma,\"DejaVu Sans\",helvetica,arial,freesans,sans-serif;\n      border-collapse: collapse;\n      /*width: 60%;*/\n      /*margin-right: 50px;*/\n    }\n\n    td, th {\n      border: 1px solid #dddddd;\n      padding: 8px;\n      text-align: center;\n    }\n    .center {\n      margin-left: auto;\n      margin-right: auto;\n    }\n    tr:nth-child(even) {\n      background-color: #dddddd;\n    }\n    .img-center {\n      display: block;\n      margin-left: auto;\n      margin-right: auto;\n      margin-top: 10px;\n    }\n    ul {\n      margin: 0px;\n      padding: 0px;\n    }\n    .footer-section {\n      background: #ffc000;\n      position: relative;\n    }\n    .footer-cta {\n      border-bottom: 1px solid #373636;\n    }\n    .single-cta i {\n      color: #ff5e14;\n      font-size: 30px;\n      float: left;\n      margin-top: 8px;\n    }\n    .cta-text {\n      padding-left: 15px;\n      display: inline-block;\n    }\n    .cta-text h4 {\n      color: #fff;\n      font-size: 20px;\n      font-weight: 600;\n      margin-bottom: 2px;\n    }\n    .cta-text span {\n      color: #757575;\n      font-size: 15px;\n    }\n    .footer-content {\n      position: relative;\n      z-index: 2;\n    }\n    .footer-pattern img {\n      position: absolute;\n      top: 0;\n      left: 0;\n      height: 330px;\n      background-size: cover;\n      background-position: 100% 100%;\n    }\n    .footer-logo {\n      margin-bottom: 30px;\n    }\n    .footer-logo img {\n      max-width: 200px;\n    }\n    .footer-text p {\n      margin-bottom: 14px;\n      font-size: 14px;\n      color: #7e7e7e;\n      line-height: 28px;\n    }\n    .footer-social-icon span {\n      color: #fff;\n      display: block;\n      font-size: 20px;\n      font-weight: 700;\n      font-family: 'Poppins', sans-serif;\n      margin-bottom: 20px;\n    }\n    .footer-social-icon a {\n      color: #fff;\n      font-size: 16px;\n      margin-right: 15px;\n    }\n    .footer-social-icon i {\n      height: 40px;\n      width: 40px;\n      text-align: center;\n      line-height: 38px;\n      border-radius: 50%;\n    }\n    .facebook-bg{\n      background: #3B5998;\n    }\n    .twitter-bg{\n      background: #55ACEE;\n    }\n    .google-bg{\n      background: #DD4B39;\n    }\n    .footer-widget-heading h3 {\n      color: #fff;\n      font-size: 20px;\n      font-weight: 600;\n      margin-bottom: 40px;\n      position: relative;\n    }\n    .footer-widget-heading h3::before {\n      content: \"\";\n      position: absolute;\n      left: 0;\n      bottom: -15px;\n      height: 2px;\n      width: 50px;\n      background: #ff5e14;\n    }\n    .footer-widget ul li {\n      display: inline-block;\n      float: left;\n      width: 50%;\n      margin-bottom: 12px;\n    }\n    .footer-widget ul li a:hover{\n      color: #ff5e14;\n    }\n    .footer-widget ul li a {\n      color: #878787;\n      text-transform: capitalize;\n    }\n    .subscribe-form {\n      position: relative;\n      overflow: hidden;\n    }\n    .subscribe-form input {\n      width: 100%;\n      padding: 14px 28px;\n      background: #2E2E2E;\n      border: 1px solid #2E2E2E;\n      color: #fff;\n    }\n    .subscribe-form button {\n      position: absolute;\n      right: 0;\n      background: #ff5e14;\n      padding: 13px 20px;\n      border: 1px solid #ff5e14;\n      top: 0;\n    }\n    .subscribe-form button i {\n      color: #fff;\n      font-size: 22px;\n      transform: rotate(-6deg);\n    }\n    .copyright-area{\n      background: #202020;\n      padding: 25px 0;\n    }\n    .copyright-text p {\n      margin: 0;\n      font-size: 14px;\n      color: #878787;\n    }\n    .copyright-text p a{\n      color: #ff5e14;\n    }\n    .footer-menu li {\n      display: inline-block;\n      margin-left: 20px;\n    }\n    .footer-menu li:hover a{\n      color: #ff5e14;\n    }\n    .footer-menu li a {\n      font-size: 14px;\n      color: #878787;\n    }\n    .card {\n      box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);\n      transition: 0.3s;\n      /*width: 100%;*/\n      border-radius: 5px;\n      margin-right: 50px;\n      margin-left: 50px;\n    }\n    .card:hover {\n      box-shadow: 0 8px 16px 0 rgba(0,0,0,0.2);\n    }\n    img {\n      border-radius: 5px 5px 0 0;\n    }\n    .container {\n      padding: 2px 16px;\n    }\n    .text-align {\n      text-align: center;\n      font-size: 20px;\n      margin-top: 20px;\n    }\n\n  </style>\n</head>\n<body>\n<img class=\"img-center\" width=\"400px\" height=\"200px\" src=\"https://www.ravro.ir/api/files/documents/image-1645346580017.png\" alt=\"image\" />\n<header>\n  <h1>{{.Title}}</h1>\n  <br />\n  <h3><a href=\"https://www.ravro.ir/company/{{.CompanyUserName}}\">{{.CompanyUserName}}</a></h3>\n</header>\n<br />\n<table class=\"center\">\n  <tr>\n    <th>شکارچی</th>\n    <th>شناسه گزارش</th>\n    <th>تاریخ ثبت</th>\n    <th>CVSS</th>\n    <th>پاداش پیشنهادی</th>\n    <th>تاریخ بررسی</th>\n    <th>ای پی شکارچی</th>\n  </tr>\n  <tr>\n    <td><a href=\"https://ravro.ir/hunter/{{.Hunter}}\">{{.Hunter}}</a></td>\n    <td><a href=\"https://www.ravro.ir/report/{{.ReportID}}\">{{.ReportID}}</a></td>\n    <td>{{.DateFrom}}</td>\n    <td>{{.CVSS}}</td>\n    <td>{{.Amount}} ریال</td>\n    <td>{{.DateTo}}</td>\n    <td>{{.Ips}}</td>\n  </tr>\n</table>\n<br />\n<hr>\n<h1 class=\"text-align\">شرح آسیب پذیری (شکارچی) </h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.PoC}}\n    </p>\n  </div>\n</div>\n<h1 class=\"text-align\">چگونه می‌توان آسیب‌پذیری را رفع کرد ؟ ( شکارچی ) </h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.Reproduce}}\n    </p>\n  </div>\n</div>\n<h1 class=\"text-align\">اطلاعات بیشتر ( شکارچی ) </h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.MoreInfo}}\n    </p>\n  </div>\n</div>\n<br />\n<h1 class=\"text-align\">بررسی تیم داوری ( راورو )</h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.JudgeInfo}}\n    </p>\n  </div>\n</div>\n<br />\n<h1 class=\"text-align\"> نوع آسیب پذیری ( راورو )</h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.VulType}}\n    </p>\n  </div>\n</div>\n<br />\n<h1 class=\"text-align\"> شرح آسیب پذیری ( راورو )</h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.VulDefine}}\n    </p>\n  </div>\n</div>\n<br />\n<h1 class=\"text-align\"> رفع آسیب پذیری ( راورو )</h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.VulFix}}\n    </p>\n  </div>\n</div>\n<br />\n<h1 class=\"text-align\"> مطالعه بیشتر از این آسیب پذیری ( راورو )</h1>\n<div class=\"card\">\n  <div class=\"container\">\n    <p>\n      {{.VulWriteup}}\n    </p>\n  </div>\n</div>\n<br />\n\n<footer class=\"footer-section\" style=\"border-radius: 45px;\">\n  <div class=\"container\">\n    <div class=\"footer-cta pt-5 pb-5\">\n      <div class=\"row\">\n        <div class=\"col-xl-4 col-md-4 mb-30\">\n          <div class=\"single-cta\">\n            <i class=\"fas fa-map-marker-alt\"></i>\n            <div class=\"cta-text\">\n              <h4>دفتر مرکزی</h4>\n              <span>  تهران، خيابان مطهری، نبش سهروردی، پلاك ۹۴ ، طبقه دوم، واحد ۲۵۰</span>\n            </div>\n          </div>\n        </div>\n        <div class=\"col-xl-4 col-md-4 mb-30\">\n          <div class=\"single-cta\">\n            <i class=\"fas fa-phone\"></i>\n            <div class=\"cta-text\">\n              <h4>شماره تماس</h4>\n              <span>021-9103-1553</span>\n            </div>\n          </div>\n        </div>\n        <div class=\"col-xl-4 col-md-4 mb-30\">\n          <div class=\"single-cta\">\n            <i class=\"far fa-envelope-open\"></i>\n            <div class=\"cta-text\">\n              <h4>آدرس ایمیل</h4>\n              <span>support[at]Ravro[dot]ir</span>\n            </div>\n          </div>\n        </div>\n        <div class=\"col-xl-4 col-md-4 mb-30\">\n          <div class=\"single-cta\">\n            <i class=\"far fa-envelope-open\"></i>\n            <div class=\"cta-text\">\n              <h4>شبکه اجتماعی</h4>\n              <a href=\"https://linkedin.com/company/ravro\"\n              ><img\n                      src=\"https://www.ravro.ir/api/files/documents/image-1643712266463.png\"\n                      height=\"30\"\n              /></a>\n              <a href=\"https://instagram.com/ravro_ir\"\n              ><img\n                      src=\"https://www.ravro.ir/api/files/documents/image-1643712340297.png\"\n                      height=\"30\"\n              /></a>\n              <a href=\"https://youtube.com/ch/ravro\"\n              ><img\n                      src=\"https://www.ravro.ir/api/files/documents/image-1643712346843.png\"\n                      height=\"30\"\n              /></a>\n              <a href=\"https://twitter.com/ravro_ir\"\n              ><img\n                      src=\"https://www.ravro.ir/api/files/documents/image-1643712216888.png\"\n                      height=\"30\"\n              /></a>\n              <a href=\"https://t.me/ravro_ir\"\n              ><img\n                      src=\"https://www.ravro.ir/api/files/documents/image-1643712255232.png\"\n                      height=\"30\"\n              /></a>\n            </div>\n          </div>\n        </div>\n\n      </div>\n    </div>\n  </div>\n</footer>\n</body>\n</html>"
	err := ioutil.WriteFile(htmlPath, []byte(html), 0644)
	if err != nil {
		panic(err)
	}

}
