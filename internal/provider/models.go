package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type HttpDownloadFile struct {
	Checksum types.String `tfsdk:"checksum"`
	Url      types.String `tfsdk:"url"`
	Dest     types.String `tfsdk:"dest"`
}
