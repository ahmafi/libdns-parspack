ParsPack for [`libdns`](https://github.com/libdns/libdns)
=======================

[![Go Reference](https://pkg.go.dev/badge/test.svg)](https://pkg.go.dev/github.com/libdns/parspack)

This package implements the [libdns interfaces](https://github.com/libdns/libdns) for ParsPack, allowing you to manage DNS records.

## Authentication

1. You need to a zone (domain) in your ParsPack panel.
2. Go to CDN Sections and you can see the button for creating a new token.

![ParsPack CDN Panel](assets/images/screenshot-1.webp)

3. Create a new API token with the following scopes:

![Create API Token Dialog](assets/images/screenshot-2.webp)

- لیست سرویس‌ها
- لیست رکوردهای DNS
- ایجاد رکورد DNS
- آپدیت رکورد DNS
- حذف رکورد DNS


## Example Configuration

```golang
p := parspack.Provider{
    APIToken: "your-apitoken-here",
}
```
