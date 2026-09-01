package quotes

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

func formatRequiredBy(date *time.Time) string {
	if date == nil {
		return "Not specified"
	}
	return date.Format("02 Jan 2006")
}

func quoteEmailLayout(title string, subtitle string, content string, headerColor string) string {

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
		max-width:700px;
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
		font-size:25px;
		margin-bottom:14px;
	"
>
📄
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
		color:#e5e7eb;
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
This email was generated automatically from our quotation request system.
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

func quoteDetailRow(
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
		width:170px;
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

func quoteItemsTable(
	items []QuoteItem,
) string {

	if len(items) == 0 {
		return `
			<p
				style="
					color:#6b7280;
					font-size:14px;
				"
			>
				No products were provided.
			</p>
		`
	}

	var rows strings.Builder

	for i, item := range items {

		sku := optionalValue(item.SKU)

		imageHTML := `
			<div
				style="
					width:50px;
					height:50px;
					background-color:#f3f4f6;
					border-radius:6px;
					text-align:center;
					line-height:50px;
					font-size:10px;
					color:#9ca3af;
				"
			>
				No Image
			</div>
		`

		if item.Image != nil &&
			strings.TrimSpace(*item.Image) != "" {

			imageURL :=
				escape(
					strings.TrimSpace(
						*item.Image,
					),
				)

			imageHTML =
				fmt.Sprintf(`
					<img
						src="%s"
						alt="%s"
						width="50"
						height="50"
						style="
							width:50px;
							height:50px;
							object-fit:contain;
							border:1px solid #e5e7eb;
							border-radius:6px;
							background-color:#ffffff;
						"
					/>
				`,
					imageURL,
					escape(item.Name),
				)
		}

		rows.WriteString(
			fmt.Sprintf(`
<tr>

<td
	style="
		padding:12px;
		border-bottom:1px solid #e5e7eb;
		text-align:center;
		color:#6b7280;
		font-size:13px;
	"
>
%d
</td>

<td
	style="
		padding:12px;
		border-bottom:1px solid #e5e7eb;
	"
>
%s
</td>

<td
	style="
		padding:12px;
		border-bottom:1px solid #e5e7eb;
	"
>

<div
	style="
		color:#111827;
		font-size:14px;
		font-weight:700;
		line-height:20px;
		margin-bottom:5px;
	"
>
%s
</div>

<div
	style="
		color:#9ca3af;
		font-size:12px;
		line-height:18px;
	"
>
SKU: %s
</div>

</td>

<td
	align="center"
	style="
		padding:12px;
		border-bottom:1px solid #e5e7eb;
		color:#111827;
		font-size:15px;
		font-weight:700;
		white-space:nowrap;
	"
>
%d
</td>

</tr>
`,
				i+1,
				imageHTML,
				escape(item.Name),
				escape(sku),
				item.Quantity,
			),
		)
	}

	return fmt.Sprintf(`
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

<thead>

<tr style="background-color:#f9fafb;">

<th
	style="
		width:40px;
		padding:12px;
		text-align:center;
		color:#6b7280;
		font-size:12px;
		text-transform:uppercase;
	"
>
#
</th>

<th
	style="
		width:60px;
		padding:12px;
		text-align:left;
		color:#6b7280;
		font-size:12px;
		text-transform:uppercase;
	"
>
Image
</th>

<th
	style="
		padding:12px;
		text-align:left;
		color:#6b7280;
		font-size:12px;
		text-transform:uppercase;
	"
>
Product
</th>

<th
	style="
		width:80px;
		padding:12px;
		text-align:center;
		color:#6b7280;
		font-size:12px;
		text-transform:uppercase;
	"
>
Qty
</th>

</tr>

</thead>

<tbody>

%s

</tbody>

</table>
`,
		rows.String(),
	)
}

func CreateAdminQuoteEmail(
	quote Quote,
	items []QuoteItem,
) string {

	fullName :=
		strings.TrimSpace(
			quote.FullName,
		)

	totalQuantity := 0

	for _, item := range items {
		totalQuantity += item.Quantity
	}

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
A new Request for Quotation has been submitted through Smart Safety Hub.
</p>

</div>

<!-- RFQ INFORMATION -->

<h2
	style="
		margin:0 0 14px;
		color:#111827;
		font-size:17px;
	"
>
Quotation Information
</h2>

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

</table>

<!-- PRODUCTS -->

<div style="margin-top:30px;">

<h2
	style="
		margin:0 0 6px;
		color:#111827;
		font-size:17px;
	"
>
Requested Products
</h2>

<p
	style="
		margin:0 0 14px;
		color:#6b7280;
		font-size:13px;
	"
>
%d product line(s) • %d total unit(s)
</p>

%s

</div>

<!-- ADDITIONAL REQUIREMENTS -->

<div style="margin-top:30px;">

<p
	style="
		margin:0 0 10px;
		color:#374151;
		font-size:14px;
		font-weight:700;
	"
>
Additional Requirements
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
		quoteDetailRow(
			"RFQ Number",
			quote.QuoteNumber,
		),

		quoteDetailRow(
			"Customer Name",
			fullName,
		),

		quoteDetailRow(
			"Company",
			quote.CompanyName,
		),

		quoteDetailRow(
			"Email",
			quote.Email,
		),

		quoteDetailRow(
			"Phone",
			quote.Phone,
		),

		quoteDetailRow(
			"Delivery Location",
			quote.DeliveryLocation,
		),

		quoteDetailRow(
			"Required By",
			formatRequiredBy(
				quote.RequiredBy,
			),
		),

		quoteDetailRow(
			"Status",
			quote.Status,
		),

		quoteDetailRow(
			"Submitted At",
			formatEmailDate(
				quote.CreatedAt,
			),
		),

		quoteDetailRow(
			"Internal ID",
			quote.ID,
		),

		len(items),

		totalQuantity,

		quoteItemsTable(items),

		escape(
			optionalValue(
				quote.AdditionalRequirements,
			),
		),
	)

	return quoteEmailLayout(
		"New Request for Quotation",
		fmt.Sprintf(
			"%s • %s",
			quote.QuoteNumber,
			quote.CompanyName,
		),
		content,
		"#111827",
	)
}

func CreateCustomerQuoteEmail(
	quote Quote,
	items []QuoteItem,
) string {

	fullName :=
		strings.TrimSpace(
			quote.FullName,
		)

	totalQuantity := 0

	for _, item := range items {
		totalQuantity += item.Quantity
	}

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
Thank you for requesting a quotation from Smart Safety Hub.
We have successfully received your requirements.
</p>

<!-- REFERENCE -->

<div
	style="
		padding:20px;
		margin-bottom:26px;
		background-color:#fff7ed;
		border-left:4px solid #f97316;
		border-radius:8px;
	"
>

<p
	style="
		margin:0 0 5px;
		color:#9a3412;
		font-size:12px;
		font-weight:700;
		text-transform:uppercase;
		letter-spacing:0.5px;
	"
>
Your RFQ Reference
</p>

<p
	style="
		margin:0;
		color:#111827;
		font-size:22px;
		line-height:30px;
		font-weight:800;
	"
>
%s
</p>

</div>

<p
	style="
		margin:0 0 24px;
		color:#4b5563;
		font-size:14px;
		line-height:24px;
	"
>
Our sales team will review pricing, availability, delivery requirements and the requested product quantities before contacting you.
</p>

<!-- DETAILS -->

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

<!-- PRODUCTS -->

<div style="margin-top:30px;">

<h2
	style="
		margin:0 0 6px;
		color:#111827;
		font-size:17px;
	"
>
Your Requested Products
</h2>

<p
	style="
		margin:0 0 14px;
		color:#6b7280;
		font-size:13px;
	"
>
%d product line(s) • %d total unit(s)
</p>

%s

</div>

<!-- REQUIREMENTS -->

<div style="margin-top:28px;">

<p
	style="
		margin:0 0 10px;
		color:#374151;
		font-size:14px;
		font-weight:700;
	"
>
Additional Requirements
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

<div
	style="
		margin-top:28px;
		padding:16px;
		background-color:#f9fafb;
		border-radius:8px;
	"
>

<p
	style="
		margin:0;
		color:#6b7280;
		font-size:13px;
		line-height:21px;
	"
>
Please keep your RFQ reference number
<strong>%s</strong>
for future communication regarding this request.
You do not need to submit the quotation request again.
</p>

</div>
`,
		escape(fullName),

		escape(
			quote.QuoteNumber,
		),

		quoteDetailRow(
			"Company",
			quote.CompanyName,
		),

		quoteDetailRow(
			"Delivery Location",
			quote.DeliveryLocation,
		),

		quoteDetailRow(
			"Required By",
			formatRequiredBy(
				quote.RequiredBy,
			),
		),

		quoteDetailRow(
			"Submitted At",
			formatEmailDate(
				quote.CreatedAt,
			),
		),

		quoteDetailRow(
			"RFQ Status",
			quote.Status,
		),

		len(items),

		totalQuantity,

		quoteItemsTable(items),

		escape(
			optionalValue(
				quote.AdditionalRequirements,
			),
		),

		escape(
			quote.QuoteNumber,
		),
	)

	return quoteEmailLayout(
		"We Received Your Quote Request",
		fmt.Sprintf(
			"Reference: %s",
			quote.QuoteNumber,
		),
		content,
		"#f97316",
	)
}
