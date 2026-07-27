package contacts

import (
	"fmt"
	"html"
	"strings"
	"time"
)

func escape(value string) string {
	return html.EscapeString(value)
}

func optionalValue(value *string) string {
	if value == nil {
		return "Not provided"
	}

	trimmed := strings.TrimSpace(*value)

	if trimmed == "" {
		return "Not provided"
	}

	return trimmed
}

func formatEmailDate(date time.Time) string {
	return date.Format("02 Jan 2006, 03:04 PM")
}

func emailLayout(title string, subtitle string, content string, headerColor string) string {
	return fmt.Sprintf(`
	<!doctype html>

<html lang="en">

<head>

<meta charset="UTF-8">

<meta
	name="viewport"
	content="width=device-width, initial-scale=1.0"
/>

<title>%s</title>

</head>

<body
	style="
		margin:0;
		padding:0;
		background-color:#f3f4f6;
		font-family:Arial,Helvetica,sans-serif;
		color:#111827;
	"
>

<table
	role="presentation"
	width="100%%"
	cellspacing="0"
	cellpadding="0"
	border="0"
	style="
		background-color:#f3f4f6;
		padding:32px 16px;
	"
>

<tr>

<td align="center">

<table
	role="presentation"
	width="100%%"
	cellspacing="0"
	cellpadding="0"
	border="0"
	style="
		max-width:640px;
		background-color:#ffffff;
		border-radius:14px;
		overflow:hidden;
		box-shadow:0 8px 30px rgba(15,23,42,0.08);
	"
>

<!-- HEADER -->

<tr>

<td
	style="
		background-color:%s;
		padding:32px 30px;
		text-align:center;
	"
>

<div
	style="
		display:inline-block;
		width:54px;
		height:54px;
		line-height:54px;
		background-color:#ffffff;
		border-radius:50%%;
		font-size:26px;
		margin-bottom:14px;
	"
>
✉️
</div>

<h1
	style="
		margin:0;
		color:#ffffff;
		font-size:25px;
		line-height:34px;
	"
>
%s
</h1>

<p
	style="
		margin:10px 0 0;
		color:#d1d5db;
		font-size:14px;
		line-height:22px;
	"
>
%s
</p>

</td>

</tr>

<!-- CONTENT -->

<tr>

<td style="padding:34px 30px;">

%s

</td>

</tr>

<!-- FOOTER -->

<tr>

<td
	style="
		padding:20px 30px;
		background-color:#f9fafb;
		border-top:1px solid #e5e7eb;
		text-align:center;
	"
>

<p
	style="
		margin:0;
		color:#9ca3af;
		font-size:12px;
		line-height:18px;
	"
>
Smart Safety Hub
<br>
This email was generated automatically from our website enquiry system.
</p>

</td>

</tr>

</table>

</td>

</tr>

</table>

</body>

</html>
`,
		escape(title),
		headerColor,
		escape(title),
		escape(subtitle),
		content,
	)
}

func detailRow(
	label string,
	value string,
) string {
	if strings.TrimSpace(value) == "" {
		value = "Not provided"
	}

	return fmt.Sprintf(`
<tr>

<td
	style="
		width:150px;
		padding:13px 14px;
		background-color:#f9fafb;
		border-bottom:1px solid #e5e7eb;
		color:#6b7280;
		font-size:14px;
		font-weight:600;
		vertical-align:top;
	"
>
%s
</td>

<td
	style="
		padding:13px 14px;
		border-bottom:1px solid #e5e7eb;
		color:#111827;
		font-size:14px;
		line-height:22px;
		vertical-align:top;
		word-break:break-word;
	"
>
%s
</td>

</tr>
`,
		escape(label),
		escape(value),
	)
}

func CreateAdminContactEmail(
	contact Contact,
) string {

	content := fmt.Sprintf(`
<div
	style="
		padding:16px;
		margin-bottom:24px;
		background-color:#fff7ed;
		border-left:4px solid #f97316;
		border-radius:8px;
	"
>

<p
	style="
		margin:0;
		color:#9a3412;
		font-size:14px;
		line-height:22px;
	"
>
A new customer enquiry has been submitted through Smart Safety Hub.
</p>

</div>

<table
	role="presentation"
	width="100%%"
	cellspacing="0"
	cellpadding="0"
	border="0"
	style="
		border:1px solid #e5e7eb;
		border-radius:8px;
		overflow:hidden;
		border-collapse:separate;
		border-spacing:0;
	"
>

%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s

</table>

<div style="margin-top:26px;">

<p
	style="
		margin:0 0 10px;
		color:#374151;
		font-size:14px;
		font-weight:700;
	"
>
Customer Message
</p>

<div
	style="
		padding:18px;
		background-color:#f9fafb;
		border:1px solid #e5e7eb;
		border-radius:8px;
		color:#374151;
		font-size:14px;
		line-height:24px;
		white-space:pre-wrap;
		word-break:break-word;
	"
>
%s
</div>

</div>
`,
		detailRow(
			"Inquiry ID",
			contact.ID,
		),

		detailRow(
			"Full Name",
			contact.Name,
		),

		detailRow(
			"Company",
			optionalValue(
				contact.CompanyName,
			),
		),

		detailRow(
			"Email",
			optionalValue(
				contact.Email,
			),
		),

		detailRow(
			"Phone",
			contact.Phone,
		),

		detailRow(
			"Country",
			optionalValue(
				contact.Country,
			),
		),

		detailRow(
			"Inquiry Type",
			optionalValue(
				contact.InquiryType,
			),
		),

		detailRow(
			"Product / Brand",
			optionalValue(
				contact.ProductName,
			),
		),

		detailRow(
			"Quantity",
			optionalValue(
				contact.Quantity,
			),
		),

		detailRow(
			"Source",
			contact.Source,
		),

		detailRow(
			"Submitted At",
			formatEmailDate(
				contact.CreatedAt,
			),
		),

		escape(
			optionalValue(
				contact.Message,
			),
		),
	)

	return emailLayout(
		"New Website Enquiry",
		fmt.Sprintf(
			"Submitted by %s",
			contact.Name,
		),
		content,
		"#111827",
	)
}

func CreateCustomerContactEmail(
	contact Contact,
) string {

	content := fmt.Sprintf(`
<p
	style="
		margin:0 0 18px;
		font-size:16px;
		line-height:26px;
	"
>
Hello <strong>%s</strong>,
</p>

<p
	style="
		margin:0 0 20px;
		color:#4b5563;
		font-size:15px;
		line-height:25px;
	"
>
Thank you for contacting Smart Safety Hub.
We have successfully received your enquiry.
</p>

<div
	style="
		padding:16px;
		margin-bottom:24px;
		background-color:#fff7ed;
		border-left:4px solid #f97316;
		border-radius:8px;
	"
>

<p
	style="
		margin:0;
		color:#9a3412;
		font-size:14px;
		line-height:22px;
	"
>
Our sales team will review your requirement and contact you as soon as possible.
</p>

</div>

<table
	role="presentation"
	width="100%%"
	cellspacing="0"
	cellpadding="0"
	border="0"
	style="
		border:1px solid #e5e7eb;
		border-radius:8px;
		overflow:hidden;
		border-collapse:separate;
		border-spacing:0;
	"
>

%s
%s
%s
%s
%s

</table>

<div style="margin-top:24px;">

<p
	style="
		margin:0 0 10px;
		color:#374151;
		font-size:14px;
		font-weight:700;
	"
>
Your Message
</p>

<div
	style="
		padding:18px;
		background-color:#f9fafb;
		border:1px solid #e5e7eb;
		border-radius:8px;
		color:#4b5563;
		font-size:14px;
		line-height:24px;
		white-space:pre-wrap;
		word-break:break-word;
	"
>
%s
</div>

</div>

<p
	style="
		margin:26px 0 0;
		color:#6b7280;
		font-size:14px;
		line-height:23px;
	"
>
You do not need to submit the form again.
Please keep your enquiry reference number for future communication.
</p>
`,
		escape(
			contact.Name,
		),

		detailRow(
			"Reference ID",
			contact.ID,
		),

		detailRow(
			"Inquiry Type",
			optionalValue(
				contact.InquiryType,
			),
		),

		detailRow(
			"Product / Brand",
			optionalValue(
				contact.ProductName,
			),
		),

		detailRow(
			"Quantity",
			optionalValue(
				contact.Quantity,
			),
		),

		detailRow(
			"Submitted At",
			formatEmailDate(
				contact.CreatedAt,
			),
		),

		escape(
			optionalValue(
				contact.Message,
			),
		),
	)

	return emailLayout(
		"We Received Your Enquiry",
		"Thank you for contacting Smart Safety Hub",
		content,
		"#f97316",
	)
}
