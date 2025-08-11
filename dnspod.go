package dnspod

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	dnspod "github.com/libdns/dnspod"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *dnspod.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.dnspod",
		New: func() caddy.Module { return &Provider{new(dnspod.Provider)} },
	}
}

// Provision resolves placeholders in the API token before use.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := ctx.Replacer() // 新版本推荐用 ctx.Replacer()
	p.Provider.APIToken = repl.ReplaceAll(p.Provider.APIToken, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens.
// Syntax:
//
// dnspod [<api_token>] {
//     api_token <api_token>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// 可选的第一个参数
		if d.NextArg() {
			p.Provider.APIToken = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr() // 位置参数超过1个
		}

		// 处理子块
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.Provider.APIToken != "" {
					return d.Err("API token already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.APIToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
